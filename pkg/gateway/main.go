package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	// auth0 "github.com/auth0-community/go-auth0"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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

var TraceKey string = "tracer"

func run() error {
	var (
		addr              = flag.String("addr", "localhost:8081", "Address of the private event GRPC service")
		port              = flag.String("port", ":9081", "TCP address to listen on")
		jwksURI           = flag.String("jwks_uri", "", "Auth0 jwks uri available at auth0 dashboard")
		auth0APIIssuer    = flag.String("auth0_iss", "", "Auth0 api issuer available at auth0 dashboard")
		auth0APIAudience  = flag.String("auth0_aud", "", "Auth0 api audience available at auth0 dashboard")
		whitelistedEmails = flag.String("whitelisted_emails", "", "A list of admin emails that are whitelisted")
		tracerKind        = flag.String("tracker_kind", "grpc_gateway_event", "Namespace for the opentracing")
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
			GetUnaryClientInterceptor(),
			grpc_opentracing.UnaryClientInterceptor(tracerOpts...),
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

func GetUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		r, _ := http.NewRequest("", "http://localhost", nil)
		authHeader, ok := md["authorization"]
		if ok && len(authHeader) > 0 {
			r.Header.Add("Authorization", authHeader[0])
			_, err := auth0Validator.Validate(r)
			if err != nil {
				log.Println("error validating lo")
			}

			md = metadata.Pairs("holla!", "Bearer XXXX")
			ctx = metadata.NewContext(ctx, md)
			ctx, err = fetchUserDetails(ctx, authHeader[0])
			if err != nil {
				log.Println("error fetching user details", err)
			}
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func fetchUserDetails(ctx context.Context, auth string) (context.Context, error) {
	// Start a new span
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://engineersmy.auth0.com/userinfo", nil)
	if err != nil {
		return ctx, err
	}

	req.Header.Add("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		return ctx, err
	}
	defer resp.Body.Close()

	userinfo := make(map[string]string) // UserInfo
	if err = json.NewDecoder(resp.Body).Decode(&userinfo); err != nil {
		return ctx, err
	}

	if email, ok := userinfo["email"]; ok {
		if len(whitelist) > 0 {
			for _, v := range whitelist {
				if v == email {
					userinfo["admin"] = "true"
				}
			}
		}
	}

	md := metadata.New(userinfo)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Do a checking to see the user emails which are whitelisted
	return ctx, nil
}
