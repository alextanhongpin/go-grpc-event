package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	"github.com/alextanhongpin/go-grpc-event/internal/tracer"
	pb "github.com/alextanhongpin/go-grpc-event/proto/event-public"
)

type eventServer struct {
	db     *database.Database
	tracer opentracing.Tracer
}

// Event shadows the ID field to map the bson.ObjectId that is not available
// in .proto
type Event struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	pb.Event `bson:",inline"`
}

func (s eventServer) GetEvents(ctx context.Context, msg *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {

	// Receive metadata server-side
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Printf("got metadata: %#v", md)
	}
	sess := s.db.Copy()
	defer sess.Close()

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetEvents")
	defer span.Finish()

	c := s.db.Collection(sess, "events")

	var tmpEvents []Event
	if err := c.Find(bson.M{
		"is_published": true, // Public users can only view published events
	}).All(&tmpEvents); err != nil {
		span.SetTag("error", err.Error())
		return nil, err
	}

	var events []*pb.Event
	for _, event := range tmpEvents {
		// Convert the objectId to string id
		event.Id = event.ID.Hex()
		cvt := pb.Event(event.Event)
		events = append(events, &cvt)
	}
	return &pb.GetEventsResponse{
		Data: events,
	}, nil
}

func (s eventServer) CreateEvent(ctx context.Context, msg *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateEvent")
	defer span.Finish()

	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "events")

	msg.Data.CreatedAt = time.Now().UnixNano() / 1000000
	msg.Data.UpdatedAt = time.Now().UnixNano() / 1000000
	msg.Data.IsPublished = false // Public users can only suggest, moderator is required to approve

	if err := c.Insert(msg.Data); err != nil {
		span.SetTag("error", err.Error())
		return nil, err
	}

	return &pb.CreateEventResponse{
		Ok: true,
	}, nil
}

func main() {
	var (
		port       = flag.String("port", ":8090", "TCP port to listen on")
		mgoHost    = flag.String("mgo_host", "mongodb://localhost:27017", "MongoDB uri string")
		tracerHost = flag.String("tracer_host", "http://localhost:9411/api/v1/spans", "The jaeger host for opentracing")
		tracerKind = flag.String("tracer_kind", "grpc_event_public", "The namespace of the tracer we are running")
	)
	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	trc, err := tracer.New(
		tracer.Name(*tracerKind),
		tracer.Host(*tracerHost),
	)

	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}

	tracerOpts := []grpc_opentracing.Option{
		grpc_opentracing.WithTracer(trc),
	}

	db, err := database.New(database.Host(*mgoHost))
	if err != nil {
		log.Fatalf("error connecting to db: %v\n", err)
	}
	defer db.Close()

	log.Printf("connected to mongo=%s", *mgoHost)
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
		db: db,
	})

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", *port)
	grpcServer.Serve(lis)
}
