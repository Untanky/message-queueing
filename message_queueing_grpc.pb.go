// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: message_queueing.proto

package queueing

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// QueueServiceClient is the client API for QueueService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueueServiceClient interface {
	WriteMessages(ctx context.Context, opts ...grpc.CallOption) (QueueService_WriteMessagesClient, error)
	SubmitMessageBatch(ctx context.Context, in *SubmitMessagesRequest, opts ...grpc.CallOption) (*SubmitMessagesResponse, error)
	RetrieveMessageBatch(ctx context.Context, in *RetrieveMessagesRequest, opts ...grpc.CallOption) (*RetrieveMessagesResponse, error)
	AcknowledgeMessageBatch(ctx context.Context, in *AcknowledgeMessagesRequest, opts ...grpc.CallOption) (*AcknowledgeMessagesResponse, error)
}

type queueServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQueueServiceClient(cc grpc.ClientConnInterface) QueueServiceClient {
	return &queueServiceClient{cc}
}

func (c *queueServiceClient) WriteMessages(ctx context.Context, opts ...grpc.CallOption) (QueueService_WriteMessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &QueueService_ServiceDesc.Streams[0], "/message_queueing.QueueService/WriteMessages", opts...)
	if err != nil {
		return nil, err
	}
	x := &queueServiceWriteMessagesClient{stream}
	return x, nil
}

type QueueService_WriteMessagesClient interface {
	Send(*RawQueueMessage) error
	Recv() (*SubmitReceipt, error)
	grpc.ClientStream
}

type queueServiceWriteMessagesClient struct {
	grpc.ClientStream
}

func (x *queueServiceWriteMessagesClient) Send(m *RawQueueMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *queueServiceWriteMessagesClient) Recv() (*SubmitReceipt, error) {
	m := new(SubmitReceipt)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *queueServiceClient) SubmitMessageBatch(ctx context.Context, in *SubmitMessagesRequest, opts ...grpc.CallOption) (*SubmitMessagesResponse, error) {
	out := new(SubmitMessagesResponse)
	err := c.cc.Invoke(ctx, "/message_queueing.QueueService/SubmitMessageBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueServiceClient) RetrieveMessageBatch(ctx context.Context, in *RetrieveMessagesRequest, opts ...grpc.CallOption) (*RetrieveMessagesResponse, error) {
	out := new(RetrieveMessagesResponse)
	err := c.cc.Invoke(ctx, "/message_queueing.QueueService/RetrieveMessageBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueServiceClient) AcknowledgeMessageBatch(ctx context.Context, in *AcknowledgeMessagesRequest, opts ...grpc.CallOption) (*AcknowledgeMessagesResponse, error) {
	out := new(AcknowledgeMessagesResponse)
	err := c.cc.Invoke(ctx, "/message_queueing.QueueService/AcknowledgeMessageBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueueServiceServer is the server API for QueueService service.
// All implementations must embed UnimplementedQueueServiceServer
// for forward compatibility
type QueueServiceServer interface {
	WriteMessages(QueueService_WriteMessagesServer) error
	SubmitMessageBatch(context.Context, *SubmitMessagesRequest) (*SubmitMessagesResponse, error)
	RetrieveMessageBatch(context.Context, *RetrieveMessagesRequest) (*RetrieveMessagesResponse, error)
	AcknowledgeMessageBatch(context.Context, *AcknowledgeMessagesRequest) (*AcknowledgeMessagesResponse, error)
	mustEmbedUnimplementedQueueServiceServer()
}

// UnimplementedQueueServiceServer must be embedded to have forward compatible implementations.
type UnimplementedQueueServiceServer struct {
}

func (UnimplementedQueueServiceServer) WriteMessages(QueueService_WriteMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method WriteMessages not implemented")
}
func (UnimplementedQueueServiceServer) SubmitMessageBatch(context.Context, *SubmitMessagesRequest) (*SubmitMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitMessageBatch not implemented")
}
func (UnimplementedQueueServiceServer) RetrieveMessageBatch(context.Context, *RetrieveMessagesRequest) (*RetrieveMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveMessageBatch not implemented")
}
func (UnimplementedQueueServiceServer) AcknowledgeMessageBatch(context.Context, *AcknowledgeMessagesRequest) (*AcknowledgeMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcknowledgeMessageBatch not implemented")
}
func (UnimplementedQueueServiceServer) mustEmbedUnimplementedQueueServiceServer() {}

// UnsafeQueueServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueueServiceServer will
// result in compilation errors.
type UnsafeQueueServiceServer interface {
	mustEmbedUnimplementedQueueServiceServer()
}

func RegisterQueueServiceServer(s grpc.ServiceRegistrar, srv QueueServiceServer) {
	s.RegisterService(&QueueService_ServiceDesc, srv)
}

func _QueueService_WriteMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(QueueServiceServer).WriteMessages(&queueServiceWriteMessagesServer{stream})
}

type QueueService_WriteMessagesServer interface {
	Send(*SubmitReceipt) error
	Recv() (*RawQueueMessage, error)
	grpc.ServerStream
}

type queueServiceWriteMessagesServer struct {
	grpc.ServerStream
}

func (x *queueServiceWriteMessagesServer) Send(m *SubmitReceipt) error {
	return x.ServerStream.SendMsg(m)
}

func (x *queueServiceWriteMessagesServer) Recv() (*RawQueueMessage, error) {
	m := new(RawQueueMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _QueueService_SubmitMessageBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServiceServer).SubmitMessageBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message_queueing.QueueService/SubmitMessageBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServiceServer).SubmitMessageBatch(ctx, req.(*SubmitMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueueService_RetrieveMessageBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServiceServer).RetrieveMessageBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message_queueing.QueueService/RetrieveMessageBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServiceServer).RetrieveMessageBatch(ctx, req.(*RetrieveMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueueService_AcknowledgeMessageBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcknowledgeMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServiceServer).AcknowledgeMessageBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message_queueing.QueueService/AcknowledgeMessageBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServiceServer).AcknowledgeMessageBatch(ctx, req.(*AcknowledgeMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// QueueService_ServiceDesc is the grpc.ServiceDesc for QueueService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QueueService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "message_queueing.QueueService",
	HandlerType: (*QueueServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitMessageBatch",
			Handler:    _QueueService_SubmitMessageBatch_Handler,
		},
		{
			MethodName: "RetrieveMessageBatch",
			Handler:    _QueueService_RetrieveMessageBatch_Handler,
		},
		{
			MethodName: "AcknowledgeMessageBatch",
			Handler:    _QueueService_AcknowledgeMessageBatch_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WriteMessages",
			Handler:       _QueueService_WriteMessages_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "message_queueing.proto",
}
