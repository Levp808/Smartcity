// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: petition.proto

package petition_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Petition_CreatePetition_FullMethodName = "/petition.Petition/CreatePetition"
	Petition_GetPetition_FullMethodName    = "/petition.Petition/GetPetition"
	Petition_UpdatePetition_FullMethodName = "/petition.Petition/UpdatePetition"
	Petition_DeletePetition_FullMethodName = "/petition.Petition/DeletePetition"
)

// PetitionClient is the client API for Petition service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PetitionClient interface {
	CreatePetition(ctx context.Context, in *CreatePetitionRequest, opts ...grpc.CallOption) (*CreatePetitionResponse, error)
	GetPetition(ctx context.Context, in *GetPetitionRequest, opts ...grpc.CallOption) (*GetPetitionResponse, error)
	UpdatePetition(ctx context.Context, in *UpdatePetitionRequest, opts ...grpc.CallOption) (*UpdatePetitionResponse, error)
	DeletePetition(ctx context.Context, in *DeletePetitionRequest, opts ...grpc.CallOption) (*DeletePetitionResponse, error)
}

type petitionClient struct {
	cc grpc.ClientConnInterface
}

func NewPetitionClient(cc grpc.ClientConnInterface) PetitionClient {
	return &petitionClient{cc}
}

func (c *petitionClient) CreatePetition(ctx context.Context, in *CreatePetitionRequest, opts ...grpc.CallOption) (*CreatePetitionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreatePetitionResponse)
	err := c.cc.Invoke(ctx, Petition_CreatePetition_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petitionClient) GetPetition(ctx context.Context, in *GetPetitionRequest, opts ...grpc.CallOption) (*GetPetitionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPetitionResponse)
	err := c.cc.Invoke(ctx, Petition_GetPetition_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petitionClient) UpdatePetition(ctx context.Context, in *UpdatePetitionRequest, opts ...grpc.CallOption) (*UpdatePetitionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdatePetitionResponse)
	err := c.cc.Invoke(ctx, Petition_UpdatePetition_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petitionClient) DeletePetition(ctx context.Context, in *DeletePetitionRequest, opts ...grpc.CallOption) (*DeletePetitionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeletePetitionResponse)
	err := c.cc.Invoke(ctx, Petition_DeletePetition_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PetitionServer is the server API for Petition service.
// All implementations must embed UnimplementedPetitionServer
// for forward compatibility.
type PetitionServer interface {
	CreatePetition(context.Context, *CreatePetitionRequest) (*CreatePetitionResponse, error)
	GetPetition(context.Context, *GetPetitionRequest) (*GetPetitionResponse, error)
	UpdatePetition(context.Context, *UpdatePetitionRequest) (*UpdatePetitionResponse, error)
	DeletePetition(context.Context, *DeletePetitionRequest) (*DeletePetitionResponse, error)
	mustEmbedUnimplementedPetitionServer()
}

// UnimplementedPetitionServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPetitionServer struct{}

func (UnimplementedPetitionServer) CreatePetition(context.Context, *CreatePetitionRequest) (*CreatePetitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePetition not implemented")
}
func (UnimplementedPetitionServer) GetPetition(context.Context, *GetPetitionRequest) (*GetPetitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPetition not implemented")
}
func (UnimplementedPetitionServer) UpdatePetition(context.Context, *UpdatePetitionRequest) (*UpdatePetitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePetition not implemented")
}
func (UnimplementedPetitionServer) DeletePetition(context.Context, *DeletePetitionRequest) (*DeletePetitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePetition not implemented")
}
func (UnimplementedPetitionServer) mustEmbedUnimplementedPetitionServer() {}
func (UnimplementedPetitionServer) testEmbeddedByValue()                  {}

// UnsafePetitionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PetitionServer will
// result in compilation errors.
type UnsafePetitionServer interface {
	mustEmbedUnimplementedPetitionServer()
}

func RegisterPetitionServer(s grpc.ServiceRegistrar, srv PetitionServer) {
	// If the following call pancis, it indicates UnimplementedPetitionServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Petition_ServiceDesc, srv)
}

func _Petition_CreatePetition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePetitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetitionServer).CreatePetition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Petition_CreatePetition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetitionServer).CreatePetition(ctx, req.(*CreatePetitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Petition_GetPetition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPetitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetitionServer).GetPetition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Petition_GetPetition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetitionServer).GetPetition(ctx, req.(*GetPetitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Petition_UpdatePetition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePetitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetitionServer).UpdatePetition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Petition_UpdatePetition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetitionServer).UpdatePetition(ctx, req.(*UpdatePetitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Petition_DeletePetition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePetitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetitionServer).DeletePetition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Petition_DeletePetition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetitionServer).DeletePetition(ctx, req.(*DeletePetitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Petition_ServiceDesc is the grpc.ServiceDesc for Petition service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Petition_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "petition.Petition",
	HandlerType: (*PetitionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePetition",
			Handler:    _Petition_CreatePetition_Handler,
		},
		{
			MethodName: "GetPetition",
			Handler:    _Petition_GetPetition_Handler,
		},
		{
			MethodName: "UpdatePetition",
			Handler:    _Petition_UpdatePetition_Handler,
		},
		{
			MethodName: "DeletePetition",
			Handler:    _Petition_DeletePetition_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "petition.proto",
}
