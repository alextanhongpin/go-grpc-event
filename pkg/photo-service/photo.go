package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	pb "github.com/alextanhongpin/go-grpc-event/proto/photo"
)

type photoServer struct {
	db *database.DB
}

func (s photoServer) GetPhotos(ctx context.Context, msg *pb.GetPhotosRequest) (*pb.GetPhotosResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) GetPhoto(ctx context.Context, msg *pb.GetPhotoRequest) (*pb.GetPhotoResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) UpdatePhoto(ctx context.Context, msg *pb.UpdatePhotoRequest) (*pb.UpdatePhotoResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) DeletePhoto(ctx context.Context, msg *pb.DeletePhotoRequest) (*pb.DeletePhotoResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) CreatePhoto(ctx context.Context, msg *pb.CreatePhotoRequest) (*pb.CreatePhotoResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}
