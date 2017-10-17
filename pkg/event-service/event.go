package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gopkg.in/mgo.v2/bson"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	jaeger "github.com/alextanhongpin/go-grpc-event/internal/jaeger"
	"github.com/alextanhongpin/go-grpc-event/internal/slack"
	pb "github.com/alextanhongpin/go-grpc-event/proto/event"
	pbUser "github.com/alextanhongpin/go-grpc-event/proto/user"
)

type eventServer struct {
	db    *database.DB
	trc   opentracing.Tracer
	slack *slack.SlackWebhook
}

func (s eventServer) GetEvents(ctx context.Context, msg *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	span := jaeger.NewSpanFromContext(ctx, "get_events")
	defer span.Finish()

	span.LogEvent("create_session")
	sess := s.db.Copy()
	defer sess.Close()
	c := s.db.Collection(sess, "events")

	query := bson.M{}
	if msg.Filter != "" && strings.Contains(msg.Filter, "published") {
		query["is_published"] = !strings.Contains(msg.Filter, "-")
	}

	span.LogEvent("query")
	var events []*pb.Event
	if err := c.Find(query).All(&events); err != nil {
		msg := fmt.Sprintf("Error performing query: %s", err.Error())
		span.SetTag("error", msg)
		return nil, err
	}

	span.LogEvent("parse")

	for _, event := range events {
		event.Id = event.Mgoid.Hex()
	}
	return &pb.GetEventsResponse{
		Data:  events,
		Count: int64(len(events)),
	}, nil
}

func (s eventServer) GetEvent(ctx context.Context, msg *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	span := jaeger.NewSpanFromContext(ctx, "get_event")
	defer span.Finish()

	span.LogEvent("validate_id")
	if !bson.IsObjectIdHex(msg.Id) {
		span.SetTag("error", "Event does not exist or has been deleted")
		return nil, grpc.Errorf(codes.FailedPrecondition, "Event does not exist or has been deleted")
	}
	span.LogEvent("create_session")
	sess := s.db.Copy()
	defer sess.Close()
	c := s.db.Collection(sess, "events")

	var event *pb.Event
	span.LogEvent("query")
	if err := c.FindId(bson.ObjectIdHex(msg.Id)).One(&event); err != nil {
		span.SetTag("error", err.Error())
		return nil, err
	}
	return &pb.GetEventResponse{
		Data: event,
	}, nil
}

func (s eventServer) CreateEvent(ctx context.Context, msg *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	span := jaeger.NewSpanFromContext(ctx, "create_event")
	defer span.Finish()

	span.LogEvent("create_session")
	sess := s.db.Copy()
	defer sess.Close()
	c := s.db.Collection(sess, "events")

	span.LogEvent("get_metadata")
	var usr UserInfo
	usr.Extract(ctx)
	span.LogEventWithPayload("metadata", usr)

	// Create a new id, because we want to return it after creating
	id := bson.NewObjectId()

	msg.Data.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	msg.Data.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	msg.Data.IsPublished = usr.IsAdmin() // Events created by admin defaults to true
	msg.Data.Mgoid = id
	msg.Data.Id = id.Hex()

	if usr.IsAuthorized() {
		msg.Data.User = &pbUser.User{
			Id:       usr.Sub,
			Email:    usr.Email,
			Name:     usr.Name,
			Picture:  usr.Picture,
			Nickname: usr.Nickname,
			Sub:      usr.Sub,
		}
	}

	// Carry out validation
	span.LogEvent("validate_struct")
	ok, err := govalidator.ValidateStruct(msg.Data)
	if err != nil || !ok {
		span.SetTag("error", err.Error())
		return nil, err
	}

	span.LogEvent("insert")
	if err := c.Insert(msg.Data); err != nil {
		msg := fmt.Sprintf("Error inserting data: %s", err.Error())
		span.SetTag("error", msg)
		return nil, err
	}

	if err := s.slack.Send(fmt.Sprintf("A new event *%s* is pending for approval. View it <https://events.engineers.my/admin|here>.", msg.Data.Name)); err != nil {
		msg := fmt.Sprintf("Error sending slack notification: %s", err.Error())
		span.SetTag("error", msg)
	}

	span.LogKV("event_id", id.Hex())
	return &pb.CreateEventResponse{
		Id: id.Hex(),
	}, nil
}

func (s eventServer) UpdateEvent(ctx context.Context, msg *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	span := jaeger.NewSpanFromContext(ctx, "update_event")
	defer span.Finish()

	span.LogEvent("get_metadata")
	var usr UserInfo
	usr.Extract(ctx)

	if !usr.IsAdmin() {
		span.SetTag("error", "User is not authorized to perform this action")
		return nil, grpc.Errorf(codes.Unauthenticated, "User is not authorized to perform this action")
	}

	span.LogEvent("validate_id")
	if !bson.IsObjectIdHex(msg.Data.Id) {
		span.SetTag("error", "Event does not exist or has been deleted")
		return nil, grpc.Errorf(codes.FailedPrecondition, "Event does not exist or has been deleted")
	}

	span.LogEvent("create_session")
	sess := s.db.Copy()
	defer sess.Close()
	c := s.db.Collection(sess, "events")

	// Perform partial update

	msg.Data.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	var i map[string]interface{}
	o, err := json.Marshal(msg.Data)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(o, &i)

	span.LogEvent("update")
	if err := c.UpdateId(bson.ObjectIdHex(msg.Data.Id), bson.M{
		"$set": i,
	}); err != nil {
		msg := fmt.Sprintf("Error updating db: %s", err.Error())
		span.SetTag("error", msg)
		return nil, err
	}

	return &pb.UpdateEventResponse{
		Ok: true,
	}, nil
}

func (s eventServer) DeleteEvent(ctx context.Context, msg *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	span := jaeger.NewSpanFromContext(ctx, "delete_event")
	defer span.Finish()

	var usr UserInfo
	span.LogEvent("get_metadata")
	usr.Extract(ctx)
	if !usr.IsAdmin() {
		span.SetTag("error", "User is not authorized to perform this action")
		return nil, grpc.Errorf(codes.Unauthenticated, "User is not authorized to perform this action")
	}

	span.LogEvent("validate")
	if !bson.IsObjectIdHex(msg.Id) {
		span.SetTag("error", "Event does not exist or has been deleted")
		return nil, grpc.Errorf(codes.FailedPrecondition, "Event does not exist or has been deleted")
	}

	span.LogEvent("create_session")
	sess := s.db.Copy()
	defer sess.Close()
	c := s.db.Collection(sess, "events")

	span.LogEvent("delete")
	if err := c.RemoveId(bson.ObjectIdHex(msg.Id)); err != nil {
		msg := fmt.Sprintf("Error deleting from db: %s", err.Error())
		span.SetTag("error", msg)
		return nil, err
	}

	return &pb.DeleteEventResponse{
		Ok: true,
	}, nil
}
