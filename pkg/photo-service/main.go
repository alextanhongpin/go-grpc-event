package main

import (
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	jaeger "github.com/alextanhongpin/go-grpc-event/internal/jaeger"
	pb "github.com/alextanhongpin/go-grpc-event/proto/photo"
)

func main() {
	//
	// TCP
	//
	lis, err := net.Listen("tcp", viper.GetString("port"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//
	// TRACER
	//
	trc, closer := jaeger.New(viper.GetString("tracer"), viper.GetString("tracer_sampler_url"), viper.GetString("tracer_reporter_url"))
	defer closer.Close()

	tracerOpts := []grpc_opentracing.Option{
		grpc_opentracing.WithTracer(trc),
	}

	//
	// DATABASE
	//
	db, err := database.New(
		database.Host(viper.GetString("mgo_host")),
		database.Name(viper.GetString("mgo_db")),
		database.Username(viper.GetString("mgo_usr")),
		database.Password(viper.GetString("mgo_pwd")),
	)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	//
	// GRPC
	//
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_opentracing.StreamServerInterceptor(tracerOpts...),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_opentracing.UnaryServerInterceptor(tracerOpts...),
		)),
	)
	pb.RegisterPhotoServiceServer(grpcServer, &photoServer{
		db: db,
	})
	log.Printf("listening to port *%v", viper.GetString("port"))
	grpcServer.Serve(lis)
}
