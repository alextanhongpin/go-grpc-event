package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	gw "github.com/alextanhongpin/go-grpc-event/proto/event-public"
)

func run() error {
	var (
		addr   = flag.String("addr", "localhost:8090", "Address of the public event GRPC service")
		port   = flag.String("port", ":9090", "TCP address to listen on")
		origin = flag.String("origin", "*", "The origin allowed")
	)
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		// Add an interceptor for the grpc-gateway
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(GetUnaryClientInterceptor())),
	}

	mux := runtime.NewServeMux()
	if err := gw.RegisterEventServiceHandlerFromEndpoint(ctx, mux, *addr, opts); err != nil {
		return err
	}
	log.Printf("listening to service=public_event at endpoint=%s\n", *addr)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{*origin},
		AllowedHeaders:   []string{"*"}, //[]string{"Authorization", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowCredentials: true,
	}).Handler(mux)

	log.Printf("listening to port *%s\n", *port)
	return http.ListenAndServe(*port, authMiddleware(handler))
}

// GetUnaryClientInterceptor is responsible for intercepting the grpc request and decide the ACL
func GetUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		log.Println("in GetUnaryClientInterceptor")
		log.Println(ctx)
		// Note that this metadata also receives the `Grpc-Metadata-<field>` set from the headers in
		// a curl request
		md, ok := metadata.FromOutgoingContext(ctx)
		if ok {
			// Override context if the user is in the whitelist
			roles := md["role"]
			if len(roles) > 0 {
				// Set metadata to send to grpc-server
				md = metadata.Pairs(
					"role", "admin",
					"can-edit", "true",
				)
				ctx = metadata.NewOutgoingContext(ctx, md)
			}
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate JWT, then set the role base on the conditions
		r.Header.Set("Grpc-Metadata-Role", "Admin")
		h.ServeHTTP(w, r)
		return
	})
}
func main() {
	defer glog.Flush()
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

type Response struct {
	Message string
}
