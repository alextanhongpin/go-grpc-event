package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	gw "github.com/alextanhongpin/go-grpc-event/proto/event-public"
)

func run() error {
	var (
		addr = flag.String("addr", "localhost:8090", "Address of the public event GRPC service")
		port = flag.String("port", ":9090", "TCP address to listen on")
	)
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	mux := runtime.NewServeMux()
	if err := gw.RegisterEventServiceHandlerFromEndpoint(ctx, mux, *addr, opts); err != nil {
		return err
	}
	log.Printf("listening to service=public_event at endpoint=%s\n", *addr)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:9090/v1/events", "http://localhost:8080"},
		AllowedHeaders:   []string{"Authorization", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}).Handler(mux)

	log.Printf("listening to port *%s\n", *port)
	return http.ListenAndServe(*port, handler)
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
