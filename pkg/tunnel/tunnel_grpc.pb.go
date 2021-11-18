// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package tunnel

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

// AgentTunnelServiceClient is the client API for AgentTunnelService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentTunnelServiceClient interface {
	EventTunnel(ctx context.Context, opts ...grpc.CallOption) (AgentTunnelService_EventTunnelClient, error)
}

type agentTunnelServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentTunnelServiceClient(cc grpc.ClientConnInterface) AgentTunnelServiceClient {
	return &agentTunnelServiceClient{cc}
}

func (c *agentTunnelServiceClient) EventTunnel(ctx context.Context, opts ...grpc.CallOption) (AgentTunnelService_EventTunnelClient, error) {
	stream, err := c.cc.NewStream(ctx, &AgentTunnelService_ServiceDesc.Streams[0], "/tunnel.AgentTunnelService/EventTunnel", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentTunnelServiceEventTunnelClient{stream}
	return x, nil
}

type AgentTunnelService_EventTunnelClient interface {
	Send(*MessageWrapper) error
	Recv() (*MessageWrapper, error)
	grpc.ClientStream
}

type agentTunnelServiceEventTunnelClient struct {
	grpc.ClientStream
}

func (x *agentTunnelServiceEventTunnelClient) Send(m *MessageWrapper) error {
	return x.ClientStream.SendMsg(m)
}

func (x *agentTunnelServiceEventTunnelClient) Recv() (*MessageWrapper, error) {
	m := new(MessageWrapper)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AgentTunnelServiceServer is the server API for AgentTunnelService service.
// All implementations must embed UnimplementedAgentTunnelServiceServer
// for forward compatibility
type AgentTunnelServiceServer interface {
	EventTunnel(AgentTunnelService_EventTunnelServer) error
	mustEmbedUnimplementedAgentTunnelServiceServer()
}

// UnimplementedAgentTunnelServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAgentTunnelServiceServer struct {
}

func (UnimplementedAgentTunnelServiceServer) EventTunnel(AgentTunnelService_EventTunnelServer) error {
	return status.Errorf(codes.Unimplemented, "method EventTunnel not implemented")
}
func (UnimplementedAgentTunnelServiceServer) mustEmbedUnimplementedAgentTunnelServiceServer() {}

// UnsafeAgentTunnelServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentTunnelServiceServer will
// result in compilation errors.
type UnsafeAgentTunnelServiceServer interface {
	mustEmbedUnimplementedAgentTunnelServiceServer()
}

func RegisterAgentTunnelServiceServer(s grpc.ServiceRegistrar, srv AgentTunnelServiceServer) {
	s.RegisterService(&AgentTunnelService_ServiceDesc, srv)
}

func _AgentTunnelService_EventTunnel_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentTunnelServiceServer).EventTunnel(&agentTunnelServiceEventTunnelServer{stream})
}

type AgentTunnelService_EventTunnelServer interface {
	Send(*MessageWrapper) error
	Recv() (*MessageWrapper, error)
	grpc.ServerStream
}

type agentTunnelServiceEventTunnelServer struct {
	grpc.ServerStream
}

func (x *agentTunnelServiceEventTunnelServer) Send(m *MessageWrapper) error {
	return x.ServerStream.SendMsg(m)
}

func (x *agentTunnelServiceEventTunnelServer) Recv() (*MessageWrapper, error) {
	m := new(MessageWrapper)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AgentTunnelService_ServiceDesc is the grpc.ServiceDesc for AgentTunnelService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AgentTunnelService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tunnel.AgentTunnelService",
	HandlerType: (*AgentTunnelServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "EventTunnel",
			Handler:       _AgentTunnelService_EventTunnel_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/tunnel/tunnel.proto",
}

// CmdToolTunnelServiceClient is the client API for CmdToolTunnelService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CmdToolTunnelServiceClient interface {
	EventTunnel(ctx context.Context, opts ...grpc.CallOption) (CmdToolTunnelService_EventTunnelClient, error)
}

type cmdToolTunnelServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCmdToolTunnelServiceClient(cc grpc.ClientConnInterface) CmdToolTunnelServiceClient {
	return &cmdToolTunnelServiceClient{cc}
}

func (c *cmdToolTunnelServiceClient) EventTunnel(ctx context.Context, opts ...grpc.CallOption) (CmdToolTunnelService_EventTunnelClient, error) {
	stream, err := c.cc.NewStream(ctx, &CmdToolTunnelService_ServiceDesc.Streams[0], "/tunnel.CmdToolTunnelService/EventTunnel", opts...)
	if err != nil {
		return nil, err
	}
	x := &cmdToolTunnelServiceEventTunnelClient{stream}
	return x, nil
}

type CmdToolTunnelService_EventTunnelClient interface {
	Send(*CmdToolToControllerWrapper) error
	Recv() (*ControllerToCmdToolWrapper, error)
	grpc.ClientStream
}

type cmdToolTunnelServiceEventTunnelClient struct {
	grpc.ClientStream
}

func (x *cmdToolTunnelServiceEventTunnelClient) Send(m *CmdToolToControllerWrapper) error {
	return x.ClientStream.SendMsg(m)
}

func (x *cmdToolTunnelServiceEventTunnelClient) Recv() (*ControllerToCmdToolWrapper, error) {
	m := new(ControllerToCmdToolWrapper)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CmdToolTunnelServiceServer is the server API for CmdToolTunnelService service.
// All implementations must embed UnimplementedCmdToolTunnelServiceServer
// for forward compatibility
type CmdToolTunnelServiceServer interface {
	EventTunnel(CmdToolTunnelService_EventTunnelServer) error
	mustEmbedUnimplementedCmdToolTunnelServiceServer()
}

// UnimplementedCmdToolTunnelServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCmdToolTunnelServiceServer struct {
}

func (UnimplementedCmdToolTunnelServiceServer) EventTunnel(CmdToolTunnelService_EventTunnelServer) error {
	return status.Errorf(codes.Unimplemented, "method EventTunnel not implemented")
}
func (UnimplementedCmdToolTunnelServiceServer) mustEmbedUnimplementedCmdToolTunnelServiceServer() {}

// UnsafeCmdToolTunnelServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CmdToolTunnelServiceServer will
// result in compilation errors.
type UnsafeCmdToolTunnelServiceServer interface {
	mustEmbedUnimplementedCmdToolTunnelServiceServer()
}

func RegisterCmdToolTunnelServiceServer(s grpc.ServiceRegistrar, srv CmdToolTunnelServiceServer) {
	s.RegisterService(&CmdToolTunnelService_ServiceDesc, srv)
}

func _CmdToolTunnelService_EventTunnel_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CmdToolTunnelServiceServer).EventTunnel(&cmdToolTunnelServiceEventTunnelServer{stream})
}

type CmdToolTunnelService_EventTunnelServer interface {
	Send(*ControllerToCmdToolWrapper) error
	Recv() (*CmdToolToControllerWrapper, error)
	grpc.ServerStream
}

type cmdToolTunnelServiceEventTunnelServer struct {
	grpc.ServerStream
}

func (x *cmdToolTunnelServiceEventTunnelServer) Send(m *ControllerToCmdToolWrapper) error {
	return x.ServerStream.SendMsg(m)
}

func (x *cmdToolTunnelServiceEventTunnelServer) Recv() (*CmdToolToControllerWrapper, error) {
	m := new(CmdToolToControllerWrapper)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CmdToolTunnelService_ServiceDesc is the grpc.ServiceDesc for CmdToolTunnelService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CmdToolTunnelService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tunnel.CmdToolTunnelService",
	HandlerType: (*CmdToolTunnelServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "EventTunnel",
			Handler:       _CmdToolTunnelService_EventTunnel_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/tunnel/tunnel.proto",
}
