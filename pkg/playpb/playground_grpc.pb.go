// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: pkg/playpb/playground.proto

package playpb

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

// PlaygroundApiClient is the client API for PlaygroundApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlaygroundApiClient interface {
	AddDeviceTrait(ctx context.Context, in *AddDeviceTraitRequest, opts ...grpc.CallOption) (*AddDeviceTraitResponse, error)
	ListSupportedTraits(ctx context.Context, in *ListSupportedTraitsRequest, opts ...grpc.CallOption) (*ListSupportedTraitsResponse, error)
	AddRemoteDevice(ctx context.Context, in *AddRemoteDeviceRequest, opts ...grpc.CallOption) (*AddRemoteDeviceResponse, error)
	// PullPerformance returns the current performance metrics for the playground.
	PullPerformance(ctx context.Context, in *PullPerformanceRequest, opts ...grpc.CallOption) (PlaygroundApi_PullPerformanceClient, error)
}

type playgroundApiClient struct {
	cc grpc.ClientConnInterface
}

func NewPlaygroundApiClient(cc grpc.ClientConnInterface) PlaygroundApiClient {
	return &playgroundApiClient{cc}
}

func (c *playgroundApiClient) AddDeviceTrait(ctx context.Context, in *AddDeviceTraitRequest, opts ...grpc.CallOption) (*AddDeviceTraitResponse, error) {
	out := new(AddDeviceTraitResponse)
	err := c.cc.Invoke(ctx, "/smartcore.playground.api.PlaygroundApi/AddDeviceTrait", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playgroundApiClient) ListSupportedTraits(ctx context.Context, in *ListSupportedTraitsRequest, opts ...grpc.CallOption) (*ListSupportedTraitsResponse, error) {
	out := new(ListSupportedTraitsResponse)
	err := c.cc.Invoke(ctx, "/smartcore.playground.api.PlaygroundApi/ListSupportedTraits", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playgroundApiClient) AddRemoteDevice(ctx context.Context, in *AddRemoteDeviceRequest, opts ...grpc.CallOption) (*AddRemoteDeviceResponse, error) {
	out := new(AddRemoteDeviceResponse)
	err := c.cc.Invoke(ctx, "/smartcore.playground.api.PlaygroundApi/AddRemoteDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playgroundApiClient) PullPerformance(ctx context.Context, in *PullPerformanceRequest, opts ...grpc.CallOption) (PlaygroundApi_PullPerformanceClient, error) {
	stream, err := c.cc.NewStream(ctx, &PlaygroundApi_ServiceDesc.Streams[0], "/smartcore.playground.api.PlaygroundApi/PullPerformance", opts...)
	if err != nil {
		return nil, err
	}
	x := &playgroundApiPullPerformanceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PlaygroundApi_PullPerformanceClient interface {
	Recv() (*PullPerformanceResponse, error)
	grpc.ClientStream
}

type playgroundApiPullPerformanceClient struct {
	grpc.ClientStream
}

func (x *playgroundApiPullPerformanceClient) Recv() (*PullPerformanceResponse, error) {
	m := new(PullPerformanceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PlaygroundApiServer is the server API for PlaygroundApi service.
// All implementations must embed UnimplementedPlaygroundApiServer
// for forward compatibility
type PlaygroundApiServer interface {
	AddDeviceTrait(context.Context, *AddDeviceTraitRequest) (*AddDeviceTraitResponse, error)
	ListSupportedTraits(context.Context, *ListSupportedTraitsRequest) (*ListSupportedTraitsResponse, error)
	AddRemoteDevice(context.Context, *AddRemoteDeviceRequest) (*AddRemoteDeviceResponse, error)
	// PullPerformance returns the current performance metrics for the playground.
	PullPerformance(*PullPerformanceRequest, PlaygroundApi_PullPerformanceServer) error
	mustEmbedUnimplementedPlaygroundApiServer()
}

// UnimplementedPlaygroundApiServer must be embedded to have forward compatible implementations.
type UnimplementedPlaygroundApiServer struct {
}

func (UnimplementedPlaygroundApiServer) AddDeviceTrait(context.Context, *AddDeviceTraitRequest) (*AddDeviceTraitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddDeviceTrait not implemented")
}
func (UnimplementedPlaygroundApiServer) ListSupportedTraits(context.Context, *ListSupportedTraitsRequest) (*ListSupportedTraitsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSupportedTraits not implemented")
}
func (UnimplementedPlaygroundApiServer) AddRemoteDevice(context.Context, *AddRemoteDeviceRequest) (*AddRemoteDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRemoteDevice not implemented")
}
func (UnimplementedPlaygroundApiServer) PullPerformance(*PullPerformanceRequest, PlaygroundApi_PullPerformanceServer) error {
	return status.Errorf(codes.Unimplemented, "method PullPerformance not implemented")
}
func (UnimplementedPlaygroundApiServer) mustEmbedUnimplementedPlaygroundApiServer() {}

// UnsafePlaygroundApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlaygroundApiServer will
// result in compilation errors.
type UnsafePlaygroundApiServer interface {
	mustEmbedUnimplementedPlaygroundApiServer()
}

func RegisterPlaygroundApiServer(s grpc.ServiceRegistrar, srv PlaygroundApiServer) {
	s.RegisterService(&PlaygroundApi_ServiceDesc, srv)
}

func _PlaygroundApi_AddDeviceTrait_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddDeviceTraitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaygroundApiServer).AddDeviceTrait(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.playground.api.PlaygroundApi/AddDeviceTrait",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaygroundApiServer).AddDeviceTrait(ctx, req.(*AddDeviceTraitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaygroundApi_ListSupportedTraits_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSupportedTraitsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaygroundApiServer).ListSupportedTraits(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.playground.api.PlaygroundApi/ListSupportedTraits",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaygroundApiServer).ListSupportedTraits(ctx, req.(*ListSupportedTraitsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaygroundApi_AddRemoteDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRemoteDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaygroundApiServer).AddRemoteDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smartcore.playground.api.PlaygroundApi/AddRemoteDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaygroundApiServer).AddRemoteDevice(ctx, req.(*AddRemoteDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaygroundApi_PullPerformance_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullPerformanceRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PlaygroundApiServer).PullPerformance(m, &playgroundApiPullPerformanceServer{stream})
}

type PlaygroundApi_PullPerformanceServer interface {
	Send(*PullPerformanceResponse) error
	grpc.ServerStream
}

type playgroundApiPullPerformanceServer struct {
	grpc.ServerStream
}

func (x *playgroundApiPullPerformanceServer) Send(m *PullPerformanceResponse) error {
	return x.ServerStream.SendMsg(m)
}

// PlaygroundApi_ServiceDesc is the grpc.ServiceDesc for PlaygroundApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PlaygroundApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smartcore.playground.api.PlaygroundApi",
	HandlerType: (*PlaygroundApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddDeviceTrait",
			Handler:    _PlaygroundApi_AddDeviceTrait_Handler,
		},
		{
			MethodName: "ListSupportedTraits",
			Handler:    _PlaygroundApi_ListSupportedTraits_Handler,
		},
		{
			MethodName: "AddRemoteDevice",
			Handler:    _PlaygroundApi_AddRemoteDevice_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PullPerformance",
			Handler:       _PlaygroundApi_PullPerformance_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/playpb/playground.proto",
}
