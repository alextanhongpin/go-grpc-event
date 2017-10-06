package main

import (
	"context"

	photopb "github.com/alextanhongpin/proto/photo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type photoServer struct{}

func (s photoServer) GetPhotos(ctx context.Context, msg *photopb.GetPhotosRequest) (*photopb.GetPhotosResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) GetPhoto(ctx context.Context, msg *photopb.GetPhotoRequest) (*photopb.GetPhotoResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) UpdatePhoto(ctx context.Context, msg *photopb.UpdatePhotoRequest) (*photopb.UpdatePhotoResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) DeletePhoto(ctx context.Context, msg *photopb.DeletePhotoRequest) (*photopb.DeletePhotoResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func (s photoServer) CreatePhoto(ctx context.Context, msg *photopb.CreatePhotoRequest) (*photopb.CreatePhotoResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}
