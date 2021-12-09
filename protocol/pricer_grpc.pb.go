// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protocol

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

// PricerClient is the client API for Pricer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PricerClient interface {
	Send(ctx context.Context, in *Id, opts ...grpc.CallOption) (Pricer_SendClient, error)
}

type pricerClient struct {
	cc grpc.ClientConnInterface
}

func NewPricerClient(cc grpc.ClientConnInterface) PricerClient {
	return &pricerClient{cc}
}

func (c *pricerClient) Send(ctx context.Context, in *Id, opts ...grpc.CallOption) (Pricer_SendClient, error) {
	stream, err := c.cc.NewStream(ctx, &Pricer_ServiceDesc.Streams[0], "/pgrpc.Pricer/Send", opts...)
	if err != nil {
		return nil, err
	}
	x := &pricerSendClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Pricer_SendClient interface {
	Recv() (*Stock, error)
	grpc.ClientStream
}

type pricerSendClient struct {
	grpc.ClientStream
}

func (x *pricerSendClient) Recv() (*Stock, error) {
	m := new(Stock)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PricerServer is the server API for Pricer service.
// All implementations must embed UnimplementedPricerServer
// for forward compatibility
type PricerServer interface {
	Send(*Id, Pricer_SendServer) error
	mustEmbedUnimplementedPricerServer()
}

// UnimplementedPricerServer must be embedded to have forward compatible implementations.
type UnimplementedPricerServer struct {
}

func (UnimplementedPricerServer) Send(*Id, Pricer_SendServer) error {
	return status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedPricerServer) mustEmbedUnimplementedPricerServer() {}

// UnsafePricerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PricerServer will
// result in compilation errors.
type UnsafePricerServer interface {
	mustEmbedUnimplementedPricerServer()
}

func RegisterPricerServer(s grpc.ServiceRegistrar, srv PricerServer) {
	s.RegisterService(&Pricer_ServiceDesc, srv)
}

func _Pricer_Send_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Id)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PricerServer).Send(m, &pricerSendServer{stream})
}

type Pricer_SendServer interface {
	Send(*Stock) error
	grpc.ServerStream
}

type pricerSendServer struct {
	grpc.ServerStream
}

func (x *pricerSendServer) Send(m *Stock) error {
	return x.ServerStream.SendMsg(m)
}

// Pricer_ServiceDesc is the grpc.ServiceDesc for Pricer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pricer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pgrpc.Pricer",
	HandlerType: (*PricerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Send",
			Handler:       _Pricer_Send_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pricer.proto",
}