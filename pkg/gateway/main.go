package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"google.golang.org/grpc/codes"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/alextanhongpin/go-grpc-event/internal/auth0"
	"github.com/alextanhongpin/go-grpc-event/internal/cors"
	jaeger "github.com/alextanhongpin/go-grpc-event/internal/jaeger"
	egw "github.com/alextanhongpin/go-grpc-event/proto/event"
	pgw "github.com/alextanhongpin/go-grpc-event/proto/photo"
)

//  -pkg asset
//go:generate go-bindata-assetfs assets
// Response represents the payload that is returned on error
type Response struct {
	Message string `json:"message"`
}

func run() error {
	trc, closer := jaeger.New(viper.GetString("tracer"), viper.GetString("tracer_sampler_url"), viper.GetString("tracer_reporter_url"))
	defer closer.Close()

	tracerOpts := []grpc_opentracing.Option{
		grpc_opentracing.WithTracer(trc),
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(tracerOpts...),
			AuthClientInterceptor(),
		)),
	}

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{EnumsAsInts: true, OrigName: true, EmitDefaults: false}),
	)

	// Register the event gateway
	if err := egw.RegisterEventServiceHandlerFromEndpoint(ctx, mux, viper.GetString("event_addr"), opts); err != nil {
		return err
	}

	// Register the photo gateway
	if err := pgw.RegisterPhotoServiceHandlerFromEndpoint(ctx, mux, viper.GetString("photo_addr"), opts); err != nil {
		return err
	}

	log.Printf("event_service = %s\n", viper.GetString("event_addr"))
	log.Printf("photo_service = %s\n", viper.GetString("photo_addr"))
	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", viper.GetString("port"))

	httpMux := http.NewServeMux()
	httpMux.Handle("/swagger/", http.StripPrefix("/swagger", http.FileServer(assetFS())))
	httpMux.Handle("/", mux)

	return http.ListenAndServe(viper.GetString("port"), cors.New(httpMux))
}

func main() {
	defer glog.Flush()
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

// AuthClientInterceptor is a middleware to carry out authorization
func AuthClientInterceptor() grpc.UnaryClientInterceptor {
	var auth0Validator *auth0.Auth0
	// Create a global validator that is only initialized once
	auth0Validator = auth0.New(
		auth0.Audience(viper.GetString("auth0_aud")),
		auth0.JWKSURI(viper.GetString("auth0_jwk")),
		auth0.Issuer(viper.GetString("auth0_iss")),
	)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		span := jaeger.NewSpanFromContext(ctx, "jwt")
		defer span.Finish()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		if authHeader, ok := md["authorization"]; ok && len(authHeader) > 0 {
			// Create a new request object
			r, _ := http.NewRequest("", "http://localhost", nil)

			// Add the authorization header
			r.Header.Add("Authorization", authHeader[0])
			span.LogEvent("validate")

			if _, err := auth0Validator.Validate(r); err != nil {
				span.SetTag("error", err.Error())
				return grpc.Errorf(codes.Unauthenticated, "User is unauthorized")
			}
			span.LogKV("guest", "true")
			span.LogEvent("fetch_user")
			newMD, err := fetchUserDetails(span, authHeader[0])

			if err != nil {
				span.SetTag("error", fmt.Sprintf("Unable to fetch user details: %#v", err.Error()))
				return err
			}
			// This hack allows us to send the metadata, at the same time connect the spans between the server and gateway
			ctx = metadata.NewOutgoingContext(ctx, metadata.Join(newMD, md))

		} else {
			span.LogKV("guest", "false")
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func fetchUserDetails(parentSpan opentracing.Span, auth string) (metadata.MD, error) {
	var md metadata.MD
	span := opentracing.StartSpan("userinfo", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	// Start a new span
	client := &http.Client{}
	span.LogEvent("create_request")
	req, err := http.NewRequest("GET", "https://engineersmy.auth0.com/userinfo", nil)
	if err != nil {
		msg := fmt.Sprintf("Error creating new userinfor request: %s", err.Error())
		span.SetTag("error", msg)
		return md, err
	}

	req.Header.Add("Authorization", auth)
	span.LogEvent("make_request")
	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("Error getting user info: %s", err.Error())
		span.SetTag("error", msg)
		return md, err
	}
	defer resp.Body.Close()

	userinfo := make(map[string]interface{}) // UserInfo
	span.LogEvent("decode_payload")
	if err = json.NewDecoder(resp.Body).Decode(&userinfo); err != nil {
		msg := fmt.Sprintf("Error decoding userinfo: %s", err.Error())
		span.SetTag("error", msg)
		return md, err
	}

	newUserinfo := stringify(userinfo)

	span.LogEvent("extract_metadata")
	if email, ok := newUserinfo["email"]; ok {
		whitelist := viper.GetStringSlice("auth0_whitelist")
		if len(whitelist) > 0 {
			for _, v := range whitelist {
				if v == email {
					newUserinfo["admin"] = "true"
				}
			}
		}
	}
	return metadata.New(newUserinfo), nil
}

func stringify(in map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for k, v := range in {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}

// func serveSwagger(w http.ResponseWriter, r *http.Request) {
// 	if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
// 		glog.Errorf("Not Found: %s", r.URL.Path)
// 		http.NotFound(w, r)
// 		return
// 	}
// 	glog.Infof("serving: %s", r.URL.Path)
// 	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
// 	p = path.Join(viper.GetString("swagger_dir"), p)
// 	http.ServeFile(w, r, p)

// }
