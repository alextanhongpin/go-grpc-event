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

	gw "github.com/alextanhongpin/go-engineersmy-event/proto"
)

func run() error {
	var (
		addr = flag.String("addr", "localhost:8080", "Address of the GRPC service")
		port = flag.String("port", ":9090", "TCP address to listen on")
	)
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := gw.RegisterEventServiceHandlerFromEndpoint(ctx, mux, *addr, opts); err != nil {
		return err
	}

	handler := cors.Default().Handler(mux)
	log.Printf("listening to port *%s\n", *port)
	return http.ListenAndServe(*port, handler)
}

func main() {
	defer glog.Flush()
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
