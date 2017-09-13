package main

import (
	"flag"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/alextanhongpin/go-grpc-event/app/database"
	pb "github.com/alextanhongpin/go-grpc-event/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type eventServer struct {
	db *database.Database
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

	c := s.db.Collection(sess, "events")

	var tmpEvents []Event
	if err := c.Find(bson.M{}).All(&tmpEvents); err != nil {
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
	if !bson.IsObjectIdHex(msg.Id) {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Event does not exist or has been deleted")
	}
	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "events")

	var tmpEvt Event
	if err := c.FindId(bson.ObjectIdHex(msg.Id)).One(&tmpEvt); err != nil {
		return nil, err
	}
	tmpEvt.Id = tmpEvt.ID.Hex()
	evt := &pb.GetEventResponse{
		Data: &tmpEvt.Event,
	}
	return evt, nil
}

func (s eventServer) CreateEvent(ctx context.Context, msg *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "events")

	msg.Data.CreatedAt = time.Now().UnixNano() / 1000000
	msg.Data.UpdatedAt = time.Now().UnixNano() / 1000000

	if err := c.Insert(msg.Data); err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{
		Ok: true,
	}, nil
}

func (s eventServer) UpdateEvent(ctx context.Context, msg *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	if !bson.IsObjectIdHex(msg.Data.Id) {
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
	// Optional, use reflect to check the type
	// v := reflect.ValueOf(x)
	// switch v.Kind() {
	// case reflect.Bool:
	//     fmt.Printf("bool: %v\n", v.Bool())
	// case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
	//     fmt.Printf("int: %v\n", v.Int())
	// case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
	//     fmt.Printf("int: %v\n", v.Uint())
	// case reflect.Float32, reflect.Float64:
	//     fmt.Printf("float: %v\n", v.Float())
	// case reflect.String:
	//     fmt.Printf("string: %v\n", v.String())
	// case reflect.Slice:
	//     fmt.Printf("slice: len=%d, %v\n", v.Len(), v.Interface())
	// case reflect.Map:
	//     fmt.Printf("map: %v\n", v.Interface())
	// case reflect.Chan:
	//     fmt.Printf("chan %v\n", v.Interface())
	// default:
	//     fmt.Println(x)
	// }
	if err := c.UpdateId(bson.ObjectIdHex(msg.Data.Id), bson.M{
		"$set": m,
	}); err != nil {
		return nil, err
	}

	return &pb.UpdateEventResponse{
		Ok: true,
	}, nil
}

func (s eventServer) DeleteEvent(ctx context.Context, msg *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	if !bson.IsObjectIdHex(msg.Id) {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Event does not exist or has been deleted")
	}
	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "events")
	if err := c.RemoveId(bson.ObjectIdHex(msg.Id)); err != nil {
		return nil, err
	}
	return &pb.DeleteEventResponse{
		Ok: true,
	}, nil
}

func main() {
	var (
		port = flag.String("port", ":8080", "TCP port to listen on")
	)
	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := database.New()
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

	grpcServer := grpc.NewServer()
	pb.RegisterEventServiceServer(grpcServer, &eventServer{
		db: db,
	})

	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", *port)
	grpcServer.Serve(lis)
}
