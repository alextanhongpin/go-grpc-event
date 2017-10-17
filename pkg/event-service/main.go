package main

import (
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	jaeger "github.com/alextanhongpin/go-grpc-event/internal/jaeger"
	"github.com/alextanhongpin/go-grpc-event/internal/slack"
	pb "github.com/alextanhongpin/go-grpc-event/proto/event"
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
	slk := slack.New(
		slack.WebhookURL(viper.GetString("slack_webhook")),
		slack.IconEmoji(viper.GetString("slack_icon")),
		slack.Username(viper.GetString("slack_username")),
		slack.Channel(viper.GetString("slack_channel")),
	)

	pb.RegisterEventServiceServer(grpcServer, &eventServer{
		db:    db,
		trc:   trc,
		slack: slk,
	})

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", viper.GetString("port"))
	grpcServer.Serve(lis)
}
