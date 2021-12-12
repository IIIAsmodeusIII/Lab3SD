// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package Proto

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

// BrokerClient is the client API for Broker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrokerClient interface {
	GetServer(ctx context.Context, in *GetServerReq, opts ...grpc.CallOption) (*GetServerResp, error)
	GetRebelds(ctx context.Context, in *GetRebeldsReq, opts ...grpc.CallOption) (*GetRebeldsResp, error)
}

type brokerClient struct {
	cc grpc.ClientConnInterface
}

func NewBrokerClient(cc grpc.ClientConnInterface) BrokerClient {
	return &brokerClient{cc}
}

func (c *brokerClient) GetServer(ctx context.Context, in *GetServerReq, opts ...grpc.CallOption) (*GetServerResp, error) {
	out := new(GetServerResp)
	err := c.cc.Invoke(ctx, "/grpc.Broker/GetServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerClient) GetRebelds(ctx context.Context, in *GetRebeldsReq, opts ...grpc.CallOption) (*GetRebeldsResp, error) {
	out := new(GetRebeldsResp)
	err := c.cc.Invoke(ctx, "/grpc.Broker/GetRebelds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BrokerServer is the server API for Broker service.
// All implementations must embed UnimplementedBrokerServer
// for forward compatibility
type BrokerServer interface {
	GetServer(context.Context, *GetServerReq) (*GetServerResp, error)
	GetRebelds(context.Context, *GetRebeldsReq) (*GetRebeldsResp, error)
	mustEmbedUnimplementedBrokerServer()
}

// UnimplementedBrokerServer must be embedded to have forward compatible implementations.
type UnimplementedBrokerServer struct {
}

func (UnimplementedBrokerServer) GetServer(context.Context, *GetServerReq) (*GetServerResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServer not implemented")
}
func (UnimplementedBrokerServer) GetRebelds(context.Context, *GetRebeldsReq) (*GetRebeldsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRebelds not implemented")
}
func (UnimplementedBrokerServer) mustEmbedUnimplementedBrokerServer() {}

// UnsafeBrokerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrokerServer will
// result in compilation errors.
type UnsafeBrokerServer interface {
	mustEmbedUnimplementedBrokerServer()
}

func RegisterBrokerServer(s grpc.ServiceRegistrar, srv BrokerServer) {
	s.RegisterService(&Broker_ServiceDesc, srv)
}

func _Broker_GetServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).GetServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Broker/GetServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).GetServer(ctx, req.(*GetServerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broker_GetRebelds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRebeldsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServer).GetRebelds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Broker/GetRebelds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServer).GetRebelds(ctx, req.(*GetRebeldsReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Broker_ServiceDesc is the grpc.ServiceDesc for Broker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Broker",
	HandlerType: (*BrokerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetServer",
			Handler:    _Broker_GetServer_Handler,
		},
		{
			MethodName: "GetRebelds",
			Handler:    _Broker_GetRebelds_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/services.proto",
}

// FulcrumClient is the client API for Fulcrum service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FulcrumClient interface {
	CRUD(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Data, error)
	GetRebelds(ctx context.Context, in *GetRebeldsReq, opts ...grpc.CallOption) (*GetRebeldsResp, error)
}

type fulcrumClient struct {
	cc grpc.ClientConnInterface
}

func NewFulcrumClient(cc grpc.ClientConnInterface) FulcrumClient {
	return &fulcrumClient{cc}
}

func (c *fulcrumClient) CRUD(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Data, error) {
	out := new(Data)
	err := c.cc.Invoke(ctx, "/grpc.Fulcrum/CRUD", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fulcrumClient) GetRebelds(ctx context.Context, in *GetRebeldsReq, opts ...grpc.CallOption) (*GetRebeldsResp, error) {
	out := new(GetRebeldsResp)
	err := c.cc.Invoke(ctx, "/grpc.Fulcrum/GetRebelds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FulcrumServer is the server API for Fulcrum service.
// All implementations must embed UnimplementedFulcrumServer
// for forward compatibility
type FulcrumServer interface {
	CRUD(context.Context, *Command) (*Data, error)
	GetRebelds(context.Context, *GetRebeldsReq) (*GetRebeldsResp, error)
	mustEmbedUnimplementedFulcrumServer()
}

// UnimplementedFulcrumServer must be embedded to have forward compatible implementations.
type UnimplementedFulcrumServer struct {
}

func (UnimplementedFulcrumServer) CRUD(context.Context, *Command) (*Data, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CRUD not implemented")
}
func (UnimplementedFulcrumServer) GetRebelds(context.Context, *GetRebeldsReq) (*GetRebeldsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRebelds not implemented")
}
func (UnimplementedFulcrumServer) mustEmbedUnimplementedFulcrumServer() {}

// UnsafeFulcrumServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FulcrumServer will
// result in compilation errors.
type UnsafeFulcrumServer interface {
	mustEmbedUnimplementedFulcrumServer()
}

func RegisterFulcrumServer(s grpc.ServiceRegistrar, srv FulcrumServer) {
	s.RegisterService(&Fulcrum_ServiceDesc, srv)
}

func _Fulcrum_CRUD_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FulcrumServer).CRUD(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Fulcrum/CRUD",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FulcrumServer).CRUD(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fulcrum_GetRebelds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRebeldsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FulcrumServer).GetRebelds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Fulcrum/GetRebelds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FulcrumServer).GetRebelds(ctx, req.(*GetRebeldsReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Fulcrum_ServiceDesc is the grpc.ServiceDesc for Fulcrum service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Fulcrum_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Fulcrum",
	HandlerType: (*FulcrumServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CRUD",
			Handler:    _Fulcrum_CRUD_Handler,
		},
		{
			MethodName: "GetRebelds",
			Handler:    _Fulcrum_GetRebelds_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/services.proto",
}