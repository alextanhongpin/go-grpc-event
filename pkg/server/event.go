package main

import (
	"strings"
	"time"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	pb "github.com/alextanhongpin/go-grpc-event/proto/event"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	pb.Event `bson:",inline"`
}

type eventServer struct {
	db  *database.Database
	trc opentracing.Tracer
}

func (s eventServer) GetEvents(ctx context.Context, msg *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	sess := s.db.Copy()
	defer sess.Close()

	// usr := getUserInfoFromCtx(ctx)
	// log.Println("got users", usr)
	var parentCtx opentracing.SpanContext
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		parentCtx = parentSpan.Context()
	}
	span := s.trc.StartSpan("get_events", opentracing.ChildOf(parentCtx))

	// span, ctx := opentracing.StartSpanFromContext(ctx, "GetEvents")
	defer span.Finish()

	span.LogFields(otlog.String("hello", "World"))
	c := s.db.Collection(sess, "events")

	var tmpEvents []Event

	query := bson.M{}
	if msg.Filter != "" && strings.Contains(msg.Filter, "published") {
		query["is_published"] = !strings.Contains(msg.Filter, "-")
	}

	if err := c.Find(query).All(&tmpEvents); err != nil {
		span.SetTag("error", err.Error())
		return nil, err
	}

	var events []*pb.Event
	for _, event := range tmpEvents {
		// Convert the objectId to string id
		event.Id = event.ID.Hex()
		cvt := pb.Event(event.Event)
		// Delete the user sub
		if cvt.User != nil {
			// if cvt.User.isAnonymous, remove all users object
			cvt.User.UserId = ""
			cvt.User.Sub = ""
		}

		events = append(events, &cvt)
	}
	return &pb.GetEventsResponse{
		Data:  events,
		Count: int64(len(events)),
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

	span, ctx := opentracing.StartSpanFromContext(ctx, "create-event")
	defer span.Finish()

	var usr UserInfo
	usr.Extract(ctx)
	// if !usr.IsAdmin() {
	// 	span.SetTag("error", "User is not authorized to perform this action")
	// 	return nil, grpc.Errorf(codes.Unauthenticated, "User is not authorized to perform this action")
	// }

	sess := s.db.Copy()
	defer sess.Close()

	// Create a new id, because we want to return it after creating
	id := bson.NewObjectId()
	c := s.db.Collection(sess, "events")

	msg.Data.CreatedAt = time.Now().UnixNano() / 1000000
	msg.Data.UpdatedAt = time.Now().UnixNano() / 1000000
	msg.Data.IsPublished = usr.IsAdmin() // Events created by admin defaults to true

	if usr.IsAdmin() {
		msg.Data.User = &pb.User{
			UserId:   usr.Sub,
			Email:    usr.Email,
			Name:     usr.Name,
			Picture:  usr.Picture,
			Nickname: usr.Nickname,
			Sub:      usr.Sub,
		}
	}
	evt := Event{
		id,
		*msg.Data,
	}
	// Set user
	// msg.Data.User = usr

	if err := c.Insert(evt); err != nil {
		span.SetTag("error", err.Error())
		return nil, err
	}

	return &pb.CreateEventResponse{
		Ok: true,
		Id: id.Hex(),
	}, nil
}

func (s eventServer) UpdateEvent(ctx context.Context, msg *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateEvent")
	defer span.Finish()

	var usr UserInfo
	usr.Extract(ctx)
	if !usr.IsAdmin() {
		span.SetTag("error", "User is not authorized to perform this action")
		return nil, grpc.Errorf(codes.Unauthenticated, "User is not authorized to perform this action")
	}

	if !bson.IsObjectIdHex(msg.Data.Id) {
		span.SetTag("error", "Event does not exist or has been deleted")
		return nil, grpc.Errorf(codes.FailedPrecondition, "Event does not exist or has been deleted")
	}
	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "events")

	// Perform partial update
	m := bson.M{
		"name":         msg.Data.Name,
		"uri":          msg.Data.Uri,
		"start_date":   msg.Data.StartDate,
		"updated_at":   time.Now().UnixNano() / 1000000,
		"is_published": msg.Data.IsPublished,
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

	var usr UserInfo
	usr.Extract(ctx)
	if !usr.IsAdmin() {
		span.SetTag("error", "User is not authorized to perform this action")
		return nil, grpc.Errorf(codes.Unauthenticated, "User is not authorized to perform this action")
	}

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
