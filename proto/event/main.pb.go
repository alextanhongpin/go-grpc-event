// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/event/main.proto

/*
Package event is a generated protocol buffer package.

It is generated from these files:
	proto/event/main.proto

It has these top-level messages:
	Event
	GetEventsRequest
	GetEventsResponse
	GetEventRequest
	GetEventResponse
	CreateEventRequest
	CreateEventResponse
	UpdateEventRequest
	UpdateEventResponse
	DeleteEventRequest
	DeleteEventResponse
*/
package event

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import grpc_user "github.com/alextanhongpin/go-grpc-event/proto/user"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Event struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// @inject_tag: bson:"created_at"
	CreatedAt int64 `protobuf:"varint,2,opt,name=created_at,json=createdAt" json:"created_at,omitempty" bson:"created_at"`
	// @inject_tag: bson:"updated_at"
	UpdatedAt int64 `protobuf:"varint,3,opt,name=updated_at,json=updatedAt" json:"updated_at,omitempty" bson:"updated_at"`
	// @inject_tag: bson:"start_date"
	StartDate int64 `protobuf:"varint,4,opt,name=start_date,json=startDate" json:"start_date,omitempty" bson:"start_date"`
	// @inject_tag: bson:"name"
	Name string `protobuf:"bytes,5,opt,name=name" json:"name,omitempty" bson:"name"`
	// @inject_tag: bson:"uri"
	Uri string `protobuf:"bytes,6,opt,name=uri" json:"uri,omitempty" bson:"uri"`
	// @inject_tag: bson:"tags"
	Tags []string `protobuf:"bytes,7,rep,name=tags" json:"tags,omitempty" bson:"tags"`
	// @inject_tag: bson:"is_published"
	IsPublished bool            `protobuf:"varint,8,opt,name=is_published,json=isPublished" json:"is_published,omitempty" bson:"is_published"`
	User        *grpc_user.User `protobuf:"bytes,10,opt,name=user" json:"user,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Event) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Event) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

func (m *Event) GetUpdatedAt() int64 {
	if m != nil {
		return m.UpdatedAt
	}
	return 0
}

func (m *Event) GetStartDate() int64 {
	if m != nil {
		return m.StartDate
	}
	return 0
}

func (m *Event) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Event) GetUri() string {
	if m != nil {
		return m.Uri
	}
	return ""
}

func (m *Event) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *Event) GetIsPublished() bool {
	if m != nil {
		return m.IsPublished
	}
	return false
}

func (m *Event) GetUser() *grpc_user.User {
	if m != nil {
		return m.User
	}
	return nil
}

type GetEventsRequest struct {
	Query  string `protobuf:"bytes,1,opt,name=query" json:"query,omitempty"`
	Filter string `protobuf:"bytes,2,opt,name=filter" json:"filter,omitempty"`
}

func (m *GetEventsRequest) Reset()                    { *m = GetEventsRequest{} }
func (m *GetEventsRequest) String() string            { return proto.CompactTextString(m) }
func (*GetEventsRequest) ProtoMessage()               {}
func (*GetEventsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetEventsRequest) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

func (m *GetEventsRequest) GetFilter() string {
	if m != nil {
		return m.Filter
	}
	return ""
}

type GetEventsResponse struct {
	Data  []*Event `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
	Count int64    `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
}

func (m *GetEventsResponse) Reset()                    { *m = GetEventsResponse{} }
func (m *GetEventsResponse) String() string            { return proto.CompactTextString(m) }
func (*GetEventsResponse) ProtoMessage()               {}
func (*GetEventsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetEventsResponse) GetData() []*Event {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *GetEventsResponse) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type GetEventRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *GetEventRequest) Reset()                    { *m = GetEventRequest{} }
func (m *GetEventRequest) String() string            { return proto.CompactTextString(m) }
func (*GetEventRequest) ProtoMessage()               {}
func (*GetEventRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetEventRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GetEventResponse struct {
	Data *Event `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *GetEventResponse) Reset()                    { *m = GetEventResponse{} }
func (m *GetEventResponse) String() string            { return proto.CompactTextString(m) }
func (*GetEventResponse) ProtoMessage()               {}
func (*GetEventResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GetEventResponse) GetData() *Event {
	if m != nil {
		return m.Data
	}
	return nil
}

type CreateEventRequest struct {
	Data *Event `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *CreateEventRequest) Reset()                    { *m = CreateEventRequest{} }
func (m *CreateEventRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateEventRequest) ProtoMessage()               {}
func (*CreateEventRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *CreateEventRequest) GetData() *Event {
	if m != nil {
		return m.Data
	}
	return nil
}

type CreateEventResponse struct {
	Error string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	Ok    bool   `protobuf:"varint,2,opt,name=ok" json:"ok,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *CreateEventResponse) Reset()                    { *m = CreateEventResponse{} }
func (m *CreateEventResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateEventResponse) ProtoMessage()               {}
func (*CreateEventResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CreateEventResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *CreateEventResponse) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func (m *CreateEventResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type UpdateEventRequest struct {
	Data *Event `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *UpdateEventRequest) Reset()                    { *m = UpdateEventRequest{} }
func (m *UpdateEventRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateEventRequest) ProtoMessage()               {}
func (*UpdateEventRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *UpdateEventRequest) GetData() *Event {
	if m != nil {
		return m.Data
	}
	return nil
}

type UpdateEventResponse struct {
	Error string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	Ok    bool   `protobuf:"varint,2,opt,name=ok" json:"ok,omitempty"`
}

func (m *UpdateEventResponse) Reset()                    { *m = UpdateEventResponse{} }
func (m *UpdateEventResponse) String() string            { return proto.CompactTextString(m) }
func (*UpdateEventResponse) ProtoMessage()               {}
func (*UpdateEventResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *UpdateEventResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *UpdateEventResponse) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

type DeleteEventRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteEventRequest) Reset()                    { *m = DeleteEventRequest{} }
func (m *DeleteEventRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteEventRequest) ProtoMessage()               {}
func (*DeleteEventRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *DeleteEventRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type DeleteEventResponse struct {
	Error string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	Ok    bool   `protobuf:"varint,2,opt,name=ok" json:"ok,omitempty"`
}

func (m *DeleteEventResponse) Reset()                    { *m = DeleteEventResponse{} }
func (m *DeleteEventResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteEventResponse) ProtoMessage()               {}
func (*DeleteEventResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *DeleteEventResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *DeleteEventResponse) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func init() {
	proto.RegisterType((*Event)(nil), "event.Event")
	proto.RegisterType((*GetEventsRequest)(nil), "event.GetEventsRequest")
	proto.RegisterType((*GetEventsResponse)(nil), "event.GetEventsResponse")
	proto.RegisterType((*GetEventRequest)(nil), "event.GetEventRequest")
	proto.RegisterType((*GetEventResponse)(nil), "event.GetEventResponse")
	proto.RegisterType((*CreateEventRequest)(nil), "event.CreateEventRequest")
	proto.RegisterType((*CreateEventResponse)(nil), "event.CreateEventResponse")
	proto.RegisterType((*UpdateEventRequest)(nil), "event.UpdateEventRequest")
	proto.RegisterType((*UpdateEventResponse)(nil), "event.UpdateEventResponse")
	proto.RegisterType((*DeleteEventRequest)(nil), "event.DeleteEventRequest")
	proto.RegisterType((*DeleteEventResponse)(nil), "event.DeleteEventResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for EventService service

type EventServiceClient interface {
	GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error)
	GetEvent(ctx context.Context, in *GetEventRequest, opts ...grpc.CallOption) (*GetEventResponse, error)
	CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error)
	UpdateEvent(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*UpdateEventResponse, error)
	DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*DeleteEventResponse, error)
}

type eventServiceClient struct {
	cc *grpc.ClientConn
}

func NewEventServiceClient(cc *grpc.ClientConn) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error) {
	out := new(GetEventsResponse)
	err := grpc.Invoke(ctx, "/event.EventService/GetEvents", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) GetEvent(ctx context.Context, in *GetEventRequest, opts ...grpc.CallOption) (*GetEventResponse, error) {
	out := new(GetEventResponse)
	err := grpc.Invoke(ctx, "/event.EventService/GetEvent", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error) {
	out := new(CreateEventResponse)
	err := grpc.Invoke(ctx, "/event.EventService/CreateEvent", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) UpdateEvent(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*UpdateEventResponse, error) {
	out := new(UpdateEventResponse)
	err := grpc.Invoke(ctx, "/event.EventService/UpdateEvent", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*DeleteEventResponse, error) {
	out := new(DeleteEventResponse)
	err := grpc.Invoke(ctx, "/event.EventService/DeleteEvent", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for EventService service

type EventServiceServer interface {
	GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error)
	GetEvent(context.Context, *GetEventRequest) (*GetEventResponse, error)
	CreateEvent(context.Context, *CreateEventRequest) (*CreateEventResponse, error)
	UpdateEvent(context.Context, *UpdateEventRequest) (*UpdateEventResponse, error)
	DeleteEvent(context.Context, *DeleteEventRequest) (*DeleteEventResponse, error)
}

func RegisterEventServiceServer(s *grpc.Server, srv EventServiceServer) {
	s.RegisterService(&_EventService_serviceDesc, srv)
}

func _EventService_GetEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/GetEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetEvents(ctx, req.(*GetEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_GetEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/GetEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetEvent(ctx, req.(*GetEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/CreateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).CreateEvent(ctx, req.(*CreateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).UpdateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/UpdateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).UpdateEvent(ctx, req.(*UpdateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/DeleteEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).DeleteEvent(ctx, req.(*DeleteEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EventService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "event.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEvents",
			Handler:    _EventService_GetEvents_Handler,
		},
		{
			MethodName: "GetEvent",
			Handler:    _EventService_GetEvent_Handler,
		},
		{
			MethodName: "CreateEvent",
			Handler:    _EventService_CreateEvent_Handler,
		},
		{
			MethodName: "UpdateEvent",
			Handler:    _EventService_UpdateEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _EventService_DeleteEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/event/main.proto",
}

func init() { proto.RegisterFile("proto/event/main.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 590 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x4f, 0x6f, 0xd3, 0x4e,
	0x10, 0x95, 0xe3, 0xa4, 0xbf, 0x64, 0x1c, 0xfd, 0xd2, 0x4e, 0xd2, 0xc4, 0x8d, 0x40, 0xb8, 0x86,
	0x43, 0xd4, 0x83, 0x2d, 0x02, 0xea, 0x81, 0x5e, 0x40, 0x14, 0x71, 0xe8, 0x05, 0x19, 0x7a, 0x42,
	0x28, 0xda, 0xc6, 0x4b, 0x58, 0x35, 0xb5, 0xdd, 0xdd, 0x75, 0x24, 0x84, 0x7a, 0xe1, 0x2b, 0xf0,
	0xd1, 0x38, 0x72, 0xe5, 0x33, 0x70, 0x46, 0xbb, 0xfe, 0x53, 0xc7, 0x89, 0x84, 0x72, 0xdb, 0x7d,
	0x33, 0xf3, 0xde, 0xec, 0xbc, 0xb1, 0x61, 0x98, 0xf0, 0x58, 0xc6, 0x3e, 0x5d, 0xd1, 0x48, 0xfa,
	0x37, 0x84, 0x45, 0x9e, 0x06, 0xb0, 0xa5, 0x91, 0xf1, 0x83, 0x45, 0x1c, 0x2f, 0x96, 0xd4, 0x27,
	0x09, 0xf3, 0x49, 0x14, 0xc5, 0x92, 0x48, 0x16, 0x47, 0x22, 0x4b, 0x1a, 0x1f, 0x66, 0xc5, 0xa9,
	0xa0, 0xbc, 0x52, 0xeb, 0xfe, 0x31, 0xa0, 0xf5, 0x46, 0x95, 0xe3, 0xff, 0xd0, 0x60, 0xa1, 0x6d,
	0x38, 0xc6, 0xa4, 0x13, 0x34, 0x58, 0x88, 0x0f, 0x01, 0xe6, 0x9c, 0x12, 0x49, 0xc3, 0x19, 0x91,
	0x76, 0xc3, 0x31, 0x26, 0x66, 0xd0, 0xc9, 0x91, 0x57, 0x52, 0x85, 0xd3, 0x24, 0x2c, 0xc2, 0x66,
	0x16, 0xce, 0x91, 0x2c, 0x2c, 0x24, 0xe1, 0x72, 0xa6, 0x00, 0xbb, 0x99, 0x85, 0x35, 0x72, 0x4e,
	0x24, 0x45, 0x84, 0x66, 0x44, 0x6e, 0xa8, 0xdd, 0xd2, 0x72, 0xfa, 0x8c, 0xfb, 0x60, 0xa6, 0x9c,
	0xd9, 0x7b, 0x1a, 0x52, 0x47, 0x95, 0x25, 0xc9, 0x42, 0xd8, 0xff, 0x39, 0xa6, 0xca, 0x52, 0x67,
	0x3c, 0x86, 0x2e, 0x13, 0xb3, 0x24, 0xbd, 0x5a, 0x32, 0xf1, 0x85, 0x86, 0x76, 0xdb, 0x31, 0x26,
	0xed, 0xc0, 0x62, 0xe2, 0x5d, 0x01, 0xe1, 0x63, 0x68, 0xaa, 0x67, 0xda, 0xe0, 0x18, 0x13, 0x6b,
	0xda, 0xf3, 0x16, 0x3c, 0x99, 0x7b, 0x0a, 0xf1, 0x2e, 0x05, 0xe5, 0x81, 0x0e, 0xba, 0x2f, 0x61,
	0xff, 0x2d, 0x95, 0xfa, 0xe9, 0x22, 0xa0, 0xb7, 0x29, 0x15, 0x12, 0x07, 0xd0, 0xba, 0x4d, 0x29,
	0xff, 0x9a, 0x4f, 0x21, 0xbb, 0xe0, 0x10, 0xf6, 0x3e, 0xb3, 0xa5, 0xa4, 0x5c, 0x0f, 0xa1, 0x13,
	0xe4, 0x37, 0xf7, 0x02, 0x0e, 0x2a, 0x0c, 0x22, 0x89, 0x23, 0x41, 0xd1, 0x81, 0x66, 0x48, 0x24,
	0xb1, 0x0d, 0xc7, 0x9c, 0x58, 0xd3, 0xae, 0xa7, 0xad, 0xf1, 0x74, 0x52, 0xa0, 0x23, 0x4a, 0x64,
	0x1e, 0xa7, 0x51, 0x31, 0xd2, 0xec, 0xe2, 0x1e, 0x43, 0xaf, 0x20, 0x2b, 0xba, 0xa9, 0x19, 0xe2,
	0x3e, 0xbf, 0xef, 0x78, 0x8b, 0x9c, 0xb1, 0x5d, 0xce, 0x3d, 0x05, 0x7c, 0xad, 0x4d, 0x5b, 0xe3,
	0xfe, 0x77, 0xdd, 0x05, 0xf4, 0xd7, 0xea, 0x72, 0xc1, 0x01, 0xb4, 0x28, 0xe7, 0x31, 0x2f, 0x46,
	0xa4, 0x2f, 0xaa, 0xd5, 0xf8, 0x5a, 0x3f, 0xa8, 0x1d, 0x34, 0xe2, 0xeb, 0xbc, 0x75, 0xb3, 0x6c,
	0xfd, 0x14, 0xf0, 0x52, 0xaf, 0xc6, 0x8e, 0x4d, 0x9c, 0x41, 0x7f, 0xad, 0x6e, 0x97, 0x26, 0xdc,
	0x27, 0x80, 0xe7, 0x74, 0x49, 0x6b, 0xa2, 0xf5, 0xa9, 0x9e, 0x41, 0x7f, 0x2d, 0x6b, 0x17, 0x89,
	0xe9, 0x2f, 0x13, 0xba, 0xba, 0xee, 0x3d, 0xe5, 0x2b, 0x36, 0xa7, 0x18, 0x40, 0xa7, 0xdc, 0x09,
	0x1c, 0xe5, 0x2f, 0xaa, 0xef, 0xd9, 0xd8, 0xde, 0x0c, 0x64, 0xb2, 0x2e, 0x7e, 0xff, 0xf9, 0xfb,
	0x47, 0xa3, 0x8b, 0xe0, 0xaf, 0x9e, 0x66, 0x1f, 0xba, 0xc0, 0x0f, 0xd0, 0x2e, 0x12, 0x71, 0x58,
	0xab, 0x2c, 0x18, 0x47, 0x1b, 0x78, 0x4e, 0x38, 0xd2, 0x84, 0x07, 0xd8, 0xbb, 0x27, 0xf4, 0xbf,
	0xb1, 0xf0, 0x0e, 0x3f, 0x82, 0x55, 0xf1, 0x17, 0x8f, 0x72, 0x82, 0xcd, 0x5d, 0x19, 0x8f, 0xb7,
	0x85, 0x72, 0xfa, 0x43, 0x4d, 0xdf, 0x73, 0x2b, 0xfd, 0xbe, 0x30, 0x4e, 0x90, 0x82, 0x55, 0xf1,
	0xad, 0x24, 0xdf, 0xdc, 0x81, 0x92, 0x7c, 0x8b, 0xcd, 0xee, 0x23, 0x4d, 0x7e, 0x34, 0x1d, 0x54,
	0x7b, 0x57, 0x7b, 0xe1, 0xb1, 0xf0, 0x4e, 0xc9, 0x7c, 0x02, 0xab, 0xe2, 0x5d, 0x29, 0xb3, 0xe9,
	0x7a, 0x29, 0xb3, 0xc5, 0xea, 0x62, 0x44, 0x27, 0xf5, 0x11, 0x5d, 0xed, 0xe9, 0x5f, 0xe4, 0xb3,
	0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xf4, 0xf3, 0xa6, 0xeb, 0x78, 0x05, 0x00, 0x00,
}
