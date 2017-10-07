package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	pb "github.com/alextanhongpin/go-grpc-event/proto/photo"
)

type photoServer struct {
	db *database.DB
}

func (s photoServer) GetPhotos(ctx context.Context, msg *pb.GetPhotosRequest) (*pb.GetPhotosResponse, error) {
	log.Println("get_photos", msg)
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) GetPhoto(ctx context.Context, msg *pb.GetPhotoRequest) (*pb.GetPhotoResponse, error) {
	log.Println("get_photo", msg)
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) UpdatePhoto(ctx context.Context, msg *pb.UpdatePhotoRequest) (*pb.UpdatePhotoResponse, error) {
	log.Println("update_photo", msg)
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) DeletePhoto(ctx context.Context, msg *pb.DeletePhotoRequest) (*pb.DeletePhotoResponse, error) {
	log.Println("delete_photo", msg)
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) CreatePhoto(ctx context.Context, msg *pb.CreatePhotoRequest) (*pb.CreatePhotoResponse, error) {
	log.Println("create_photo", msg)
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}
