package main

import (
	"flag"
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	jaeger "github.com/alextanhongpin/go-grpc-event/internal/jaeger"
	pb "github.com/alextanhongpin/go-grpc-event/proto/event"
)

func main() {
	var (
		port       = flag.String("port", ":8080", "TCP port to listen on")
		mgoHost    = flag.String("mgo_host", "mongodb://localhost:27017", "MongoDB uri string")
		tracerKind = flag.String("tracer_kind", "event", "The namespace of the tracer we are running")
	)

	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	trc, closer := jaeger.New(*tracerKind)
	defer closer.Close()

	tracerOpts := []grpc_opentracing.Option{
		grpc_opentracing.WithTracer(trc),
	}

	// TODO: Setup database in `internals`` folder
	db, err := database.New(database.Host(*mgoHost))
	if err != nil {
		log.Fatalf("error connecting to db: %v\n", err)
	}
	defer db.Close()

	db.Ref.SetMode(mgo.Monotonic, true)
	c := db.Collection(db.Ref, "events")

	if err := c.EnsureIndex(mgo.Index{
		Key: []string{"$text:name"},
	}); err != nil {
		log.Printf("error creating index: %v\n", err)
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_opentracing.StreamServerInterceptor(tracerOpts...),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_opentracing.UnaryServerInterceptor(tracerOpts...),
		)),
	)

	pb.RegisterEventServiceServer(grpcServer, &eventServer{
		db:  db,
		trc: trc,
	})

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", *port)
	grpcServer.Serve(lis)
}
