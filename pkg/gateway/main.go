package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	gw "github.com/alextanhongpin/go-grpc-event/proto/event"
)

// Response represents the payload that is returned on error
type Response struct {
	Message string `json:"message"`
}

var auth0Validator *auth0.Auth0
var whitelist []string

// var tracer opentracing.Tracer

func run() error {
	var (
		addr              = flag.String("addr", "localhost:8081", "Address of the private event GRPC service")
		port              = flag.String("port", ":9081", "TCP address to listen on")
		jwksURI           = flag.String("jwks_uri", "", "Auth0 jwks uri available at auth0 dashboard")
		auth0APIIssuer    = flag.String("auth0_iss", "", "Auth0 api issuer available at auth0 dashboard")
		auth0APIAudience  = flag.String("auth0_aud", "", "Auth0 api audience available at auth0 dashboard")
		whitelistedEmails = flag.String("whitelisted_emails", "", "A list of admin emails that are whitelisted")
		tracerKind        = flag.String("tracker_kind", "gateway", "Namespace for the opentracing")
		// https://engineersmy.auth0.com/userinfo
	)
	flag.Parse()

	if *whitelistedEmails != "" {
		tmp := strings.Split(*whitelistedEmails, ",")
		for _, v := range tmp {
			whitelist = append(whitelist, strings.TrimSpace(v))
		}
	}

	trc, closer := jaeger.New(*tracerKind)
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

	if err := gw.RegisterEventServiceHandlerFromEndpoint(ctx, mux, *addr, opts); err != nil {
		return err
	}

	// Create a global validator that is only initialized once
	auth0Validator = auth0.New(
		auth0.Audience(*auth0APIAudience),
		auth0.JWKSURI(*jwksURI),
		auth0.Issuer(*auth0APIIssuer),
	)

	log.Printf("grpc_server = %s\n", *addr)
	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", *port)
	return http.ListenAndServe(*port, cors.New(mux))
}

func main() {
	defer glog.Flush()
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

// AuthClientInterceptor is a middleware to carry out authorization
func AuthClientInterceptor() grpc.UnaryClientInterceptor {
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

	userinfo := make(map[string]string) // UserInfo
	span.LogEvent("decode_payload")
	if err = json.NewDecoder(resp.Body).Decode(&userinfo); err != nil {
		msg := fmt.Sprintf("Error decoding userinfo: %s", err.Error())
		span.SetTag("error", msg)
		return md, err
	}
	span.LogEvent("extract_metadata")
	if email, ok := userinfo["email"]; ok {
		if len(whitelist) > 0 {
			for _, v := range whitelist {
				if v == email {
					userinfo["admin"] = "true"
				}
			}
		}
	}
	return metadata.New(userinfo), nil
}
