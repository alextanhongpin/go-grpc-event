package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gopkg.in/mgo.v2/bson"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	jaeger "github.com/alextanhongpin/go-grpc-event/internal/jaeger"
	pb "github.com/alextanhongpin/go-grpc-event/proto/photo"
	pbUser "github.com/alextanhongpin/go-grpc-event/proto/user"
)

type photoServer struct {
	db *database.DB
}

func (s photoServer) GetPhotos(ctx context.Context, msg *pb.GetPhotosRequest) (*pb.GetPhotosResponse, error) {
	span := jaeger.NewSpanFromContext(ctx, "get_photos")
	defer span.Finish()

	span.LogEvent("create_session")
	sess := s.db.Copy()
	defer sess.Close()
	c := s.db.Collection(sess, "photos")

	span.LogEvent("query")
	var photos []*pb.Photo
	if err := c.Find(nil).All(&photos); err != nil {
		msg := fmt.Sprintf("Error performing query: %s", err.Error())
		span.SetTag("error", msg)
		return nil, err
	}

	return &pb.GetPhotosResponse{
		Data:       photos,
		TotalCount: int64(len(photos)),
	}, nil
}

func (s photoServer) GetPhoto(ctx context.Context, msg *pb.GetPhotoRequest) (*pb.GetPhotoResponse, error) {
	span := jaeger.NewSpanFromContext(ctx, "get_photo")
	defer span.Finish()

	span.LogEvent("validate_id")
	if !bson.IsObjectIdHex(msg.Id) {
		msg := fmt.Sprintf("Id provided is invalid: %s", msg.Id)
		span.SetTag("error", msg)
		return nil, grpc.Errorf(codes.InvalidArgument, "Photo does not exist or has been deleted")
	}

	span.LogEvent("create_session")
	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "photos")

	span.LogEvent("query")
	var photo *pb.Photo
	if err := c.FindId(bson.ObjectIdHex(msg.Id)).One(&photo); err != nil {
		msg := fmt.Sprintf("Error performing query: %s", err.Error())
		span.SetTag("error", msg)
		return nil, err
	}

	return &pb.GetPhotoResponse{
		Data: photo,
	}, nil
}

func (s photoServer) CreatePhoto(ctx context.Context, msg *pb.CreatePhotoRequest) (*pb.CreatePhotoResponse, error) {
	span := jaeger.NewSpanFromContext(ctx, "create_photo")
	defer span.Finish()

	span.LogEvent("create_session")
	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "photos")

	var usr UserInfo
	usr.Extract(ctx)
	span.LogEventWithPayload("metadata", usr)

	span.LogEvent("insert")
	var ids []string
	var photos []interface{}
	for _, p := range msg.Data {
		id := bson.NewObjectId()
		p.Mgoid = id
		p.Id = id.Hex()
		p.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		p.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

		if usr.IsAuthorized() {
			p.User = &pbUser.User{
				Id:       usr.Sub,
				Email:    usr.Email,
				Name:     usr.Name,
				Picture:  usr.Picture,
				Nickname: usr.Nickname,
				Sub:      usr.Sub,
			}
		}
		ids = append(ids, id.Hex())
		photos = append(photos, p)
	}

	if err := c.Insert(photos...); err != nil {
		msg := fmt.Sprintf("Error inserting: %s", err.Error())
		span.SetTag("error", msg)
		return nil, err
	}

	span.LogKV("event_ids", strings.Join(ids, ","))

	return &pb.CreatePhotoResponse{
		Ids: ids,
	}, nil
}

func (s photoServer) UpdatePhoto(ctx context.Context, msg *pb.UpdatePhotoRequest) (*pb.UpdatePhotoResponse, error) {
	span := jaeger.NewSpanFromContext(ctx, "update_photo")
	defer span.Finish()

	span.LogEvent("validate_id")
	if !bson.IsObjectIdHex(msg.Data.Id) {
		msg := fmt.Sprintf("Id provided is invalid: %s", msg.Data.Id)
		span.SetTag("error", msg)
		return nil, grpc.Errorf(codes.InvalidArgument, "Photo does not exist or has been deleted")
	}

	span.LogEvent("create_session")
	sess := s.db.Copy()
	defer sess.Close()
	c := s.db.Collection(sess, "photos")

	span.LogEvent("partial_update")
	var i map[string]interface{}
	msg.Data.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	o, err := json.Marshal(msg.Data)
	if err != nil {
		msg := fmt.Sprintf("Error marshalling field: %s", err.Error())
		span.SetTag("error", msg)
		return nil, err
	}
	json.Unmarshal(o, &i)

	span.LogEvent("update")
	if err := c.UpdateId(bson.ObjectIdHex(msg.Data.Id), bson.M{
		"$set": i,
	}); err != nil {
		msg := fmt.Sprintf("Error updating document: %s", err.Error())
		span.SetTag("error", msg)
		return nil, err
	}

	return &pb.UpdatePhotoResponse{
		Ok: true,
	}, nil
}

func (s photoServer) DeletePhoto(ctx context.Context, msg *pb.DeletePhotoRequest) (*pb.DeletePhotoResponse, error) {
	span := jaeger.NewSpanFromContext(ctx, "delete_photo")
	defer span.Finish()

	var usr UserInfo
	usr.Extract(ctx)

	if !usr.IsAuthorized() {
		span.SetTag("error", "User is not authorized to perform this action")
		return nil, grpc.Errorf(codes.Unauthenticated, "User is not authorized to perform this action")
	}

	span.LogEvent("validate_id")
	if !bson.IsObjectIdHex(msg.Id) {
		msg := fmt.Sprintf("Id provided is invalid: %s", msg.Id)
		span.SetTag("error", msg)
		return nil, grpc.Errorf(codes.InvalidArgument, "Photo does not exist or has been deleted")
	}

	span.LogEvent("create_session")
	sess := s.db.Copy()
	defer sess.Close()

	c := s.db.Collection(sess, "photos")

	span.LogEvent("delete")

	if usr.IsAdmin() {
		// Admin has full access to delete everything
		if err := c.RemoveId(bson.ObjectIdHex(msg.Id)); err != nil {
			msg := fmt.Sprintf("Error deleting document from database: %s", err.Error())
			span.SetTag("error", msg)
			return nil, err
		}
	} else {
		// Authorized user can only delete the image that belongs to them
		if err := c.Remove(bson.M{
			"_id":     bson.ObjectIdHex(msg.Id),
			"user.id": usr.UserID,
		}); err != nil {
			msg := fmt.Sprintf("Error deleting document from database: %s", err.Error())
			span.SetTag("error", msg)
			return nil, err
		}
	}

	return &pb.DeletePhotoResponse{
		Ok: true,
	}, nil
}
