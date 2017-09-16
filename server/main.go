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
	"google.golang.org/grpc/codes"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	"github.com/alextanhongpin/go-grpc-event/internal/tracer"
	pb "github.com/alextanhongpin/go-grpc-event/proto"
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
	sess := s.db.Copy()
	defer sess.Close()

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetEvents")
	defer span.Finish()

	c := s.db.Collection(sess, "events")

	var tmpEvents []Event
	if err := c.Find(bson.M{}).All(&tmpEvents); err != nil {
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

func (s eventServer) GetEvent(ctx context.Context, msg *pb.GetEventRequest) (*pb.GetEventResponse, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "GetEvent")
	defer span.Finish()

	if !bson.IsObjectIdHex(msg.Id) {
		span.SetTag("error", "Event does not exist or has been deleted")
		return nil, grpc.Errorf(codes.FailedPrecondition, "Event does not exist or has been deleted")
	}
	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "events")

	var tmpEvt Event
	if err := c.FindId(bson.ObjectIdHex(msg.Id)).One(&tmpEvt); err != nil {
		span.SetTag("error", err.Error())
		return nil, err
	}
	tmpEvt.Id = tmpEvt.ID.Hex()
	evt := &pb.GetEventResponse{
		Data: &tmpEvt.Event,
	}
	return evt, nil
}

func (s eventServer) CreateEvent(ctx context.Context, msg *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateEvent")
	defer span.Finish()

	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "events")

	msg.Data.CreatedAt = time.Now().UnixNano() / 1000000
	msg.Data.UpdatedAt = time.Now().UnixNano() / 1000000

	if err := c.Insert(msg.Data); err != nil {
		span.SetTag("error", err.Error())
		return nil, err
	}

	return &pb.CreateEventResponse{
		Ok: true,
	}, nil
}

func (s eventServer) UpdateEvent(ctx context.Context, msg *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateEvent")
	defer span.Finish()

	if !bson.IsObjectIdHex(msg.Data.Id) {
		span.SetTag("error", "Event does not exist or has been deleted")
		return nil, grpc.Errorf(codes.FailedPrecondition, "Event does not exist or has been deleted")
	}
	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "events")

	// Perform partial update
	m := bson.M{
		"name":       msg.Data.Name,
		"uri":        msg.Data.Uri,
		"start_date": msg.Data.StartDate,
		"updated_at": time.Now().UnixNano() / 1000000,
	}

	if len(msg.Data.Tags) != 0 {
		m["tags"] = msg.Data.Tags
	}

	// Remove unused fields
	for k, v := range m {
		switch i := v.(type) {
		case int:
			if i == 0 {
				delete(m, k)
			}
		case string:
			if i == "" {
				delete(m, k)
			}
		}
	}
	if err := c.UpdateId(bson.ObjectIdHex(msg.Data.Id), bson.M{
		"$set": m,
	}); err != nil {
		span.SetTag("error", err.Error())
		return nil, err
	}

	return &pb.UpdateEventResponse{
		Ok: true,
	}, nil
}

func (s eventServer) DeleteEvent(ctx context.Context, msg *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteEvent")
	defer span.Finish()

	if !bson.IsObjectIdHex(msg.Id) {
		span.SetTag("error", "Event does not exist or has been deleted")
		return nil, grpc.Errorf(codes.FailedPrecondition, "Event does not exist or has been deleted")
	}
	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "events")
	if err := c.RemoveId(bson.ObjectIdHex(msg.Id)); err != nil {
		span.SetTag("error", err.Error())
		return nil, err
	}
	return &pb.DeleteEventResponse{
		Ok: true,
	}, nil
}

func main() {
	var (
		port       = flag.String("port", ":8080", "TCP port to listen on")
		mgoHost    = flag.String("mgo_host", "mongodb://localhost:27017", "MongoDB uri string")
		tracerHost = flag.String("tracer_host", "http://localhost:9411/api/v1/spans", "The jaeger host for opentracing")
		tracerKind = flag.String("tracer_kind", "grpc_event", "The namespace of the tracer we are running")
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
	// trc, closer := hunter.New()
	// defer closer.Close()

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
		db: db,
	})

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", *port)
	grpcServer.Serve(lis)
}
