// Copyright 2018-2019 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: cs3/app/registry/v1beta1/registry_api.proto

package registryv1beta1

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

const (
	RegistryAPI_GetAppProviders_FullMethodName                  = "/cs3.app.registry.v1beta1.RegistryAPI/GetAppProviders"
	RegistryAPI_AddAppProvider_FullMethodName                   = "/cs3.app.registry.v1beta1.RegistryAPI/AddAppProvider"
	RegistryAPI_ListAppProviders_FullMethodName                 = "/cs3.app.registry.v1beta1.RegistryAPI/ListAppProviders"
	RegistryAPI_ListSupportedMimeTypes_FullMethodName           = "/cs3.app.registry.v1beta1.RegistryAPI/ListSupportedMimeTypes"
	RegistryAPI_GetDefaultAppProviderForMimeType_FullMethodName = "/cs3.app.registry.v1beta1.RegistryAPI/GetDefaultAppProviderForMimeType"
	RegistryAPI_SetDefaultAppProviderForMimeType_FullMethodName = "/cs3.app.registry.v1beta1.RegistryAPI/SetDefaultAppProviderForMimeType"
)

// RegistryAPIClient is the client API for RegistryAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegistryAPIClient interface {
	// Returns the app providers that are capable of handling this resource info.
	// MUST return CODE_NOT_FOUND if no providers are available.
	GetAppProviders(ctx context.Context, in *GetAppProvidersRequest, opts ...grpc.CallOption) (*GetAppProvidersResponse, error)
	// Registers a new app provider to the registry.
	AddAppProvider(ctx context.Context, in *AddAppProviderRequest, opts ...grpc.CallOption) (*AddAppProviderResponse, error)
	// Returns a list of the available app providers known by this registry.
	ListAppProviders(ctx context.Context, in *ListAppProvidersRequest, opts ...grpc.CallOption) (*ListAppProvidersResponse, error)
	// Returns a list of the supported mime types along with the apps which they can be opened with.
	ListSupportedMimeTypes(ctx context.Context, in *ListSupportedMimeTypesRequest, opts ...grpc.CallOption) (*ListSupportedMimeTypesResponse, error)
	// Returns the default app provider which serves a specified mime type.
	GetDefaultAppProviderForMimeType(ctx context.Context, in *GetDefaultAppProviderForMimeTypeRequest, opts ...grpc.CallOption) (*GetDefaultAppProviderForMimeTypeResponse, error)
	// Sets the default app provider for a specified mime type.
	SetDefaultAppProviderForMimeType(ctx context.Context, in *SetDefaultAppProviderForMimeTypeRequest, opts ...grpc.CallOption) (*SetDefaultAppProviderForMimeTypeResponse, error)
}

type registryAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewRegistryAPIClient(cc grpc.ClientConnInterface) RegistryAPIClient {
	return &registryAPIClient{cc}
}

func (c *registryAPIClient) GetAppProviders(ctx context.Context, in *GetAppProvidersRequest, opts ...grpc.CallOption) (*GetAppProvidersResponse, error) {
	out := new(GetAppProvidersResponse)
	err := c.cc.Invoke(ctx, RegistryAPI_GetAppProviders_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryAPIClient) AddAppProvider(ctx context.Context, in *AddAppProviderRequest, opts ...grpc.CallOption) (*AddAppProviderResponse, error) {
	out := new(AddAppProviderResponse)
	err := c.cc.Invoke(ctx, RegistryAPI_AddAppProvider_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryAPIClient) ListAppProviders(ctx context.Context, in *ListAppProvidersRequest, opts ...grpc.CallOption) (*ListAppProvidersResponse, error) {
	out := new(ListAppProvidersResponse)
	err := c.cc.Invoke(ctx, RegistryAPI_ListAppProviders_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryAPIClient) ListSupportedMimeTypes(ctx context.Context, in *ListSupportedMimeTypesRequest, opts ...grpc.CallOption) (*ListSupportedMimeTypesResponse, error) {
	out := new(ListSupportedMimeTypesResponse)
	err := c.cc.Invoke(ctx, RegistryAPI_ListSupportedMimeTypes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryAPIClient) GetDefaultAppProviderForMimeType(ctx context.Context, in *GetDefaultAppProviderForMimeTypeRequest, opts ...grpc.CallOption) (*GetDefaultAppProviderForMimeTypeResponse, error) {
	out := new(GetDefaultAppProviderForMimeTypeResponse)
	err := c.cc.Invoke(ctx, RegistryAPI_GetDefaultAppProviderForMimeType_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryAPIClient) SetDefaultAppProviderForMimeType(ctx context.Context, in *SetDefaultAppProviderForMimeTypeRequest, opts ...grpc.CallOption) (*SetDefaultAppProviderForMimeTypeResponse, error) {
	out := new(SetDefaultAppProviderForMimeTypeResponse)
	err := c.cc.Invoke(ctx, RegistryAPI_SetDefaultAppProviderForMimeType_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegistryAPIServer is the server API for RegistryAPI service.
// All implementations should embed UnimplementedRegistryAPIServer
// for forward compatibility
type RegistryAPIServer interface {
	// Returns the app providers that are capable of handling this resource info.
	// MUST return CODE_NOT_FOUND if no providers are available.
	GetAppProviders(context.Context, *GetAppProvidersRequest) (*GetAppProvidersResponse, error)
	// Registers a new app provider to the registry.
	AddAppProvider(context.Context, *AddAppProviderRequest) (*AddAppProviderResponse, error)
	// Returns a list of the available app providers known by this registry.
	ListAppProviders(context.Context, *ListAppProvidersRequest) (*ListAppProvidersResponse, error)
	// Returns a list of the supported mime types along with the apps which they can be opened with.
	ListSupportedMimeTypes(context.Context, *ListSupportedMimeTypesRequest) (*ListSupportedMimeTypesResponse, error)
	// Returns the default app provider which serves a specified mime type.
	GetDefaultAppProviderForMimeType(context.Context, *GetDefaultAppProviderForMimeTypeRequest) (*GetDefaultAppProviderForMimeTypeResponse, error)
	// Sets the default app provider for a specified mime type.
	SetDefaultAppProviderForMimeType(context.Context, *SetDefaultAppProviderForMimeTypeRequest) (*SetDefaultAppProviderForMimeTypeResponse, error)
}

// UnimplementedRegistryAPIServer should be embedded to have forward compatible implementations.
type UnimplementedRegistryAPIServer struct {
}

func (UnimplementedRegistryAPIServer) GetAppProviders(context.Context, *GetAppProvidersRequest) (*GetAppProvidersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAppProviders not implemented")
}
func (UnimplementedRegistryAPIServer) AddAppProvider(context.Context, *AddAppProviderRequest) (*AddAppProviderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAppProvider not implemented")
}
func (UnimplementedRegistryAPIServer) ListAppProviders(context.Context, *ListAppProvidersRequest) (*ListAppProvidersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAppProviders not implemented")
}
func (UnimplementedRegistryAPIServer) ListSupportedMimeTypes(context.Context, *ListSupportedMimeTypesRequest) (*ListSupportedMimeTypesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSupportedMimeTypes not implemented")
}
func (UnimplementedRegistryAPIServer) GetDefaultAppProviderForMimeType(context.Context, *GetDefaultAppProviderForMimeTypeRequest) (*GetDefaultAppProviderForMimeTypeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefaultAppProviderForMimeType not implemented")
}
func (UnimplementedRegistryAPIServer) SetDefaultAppProviderForMimeType(context.Context, *SetDefaultAppProviderForMimeTypeRequest) (*SetDefaultAppProviderForMimeTypeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDefaultAppProviderForMimeType not implemented")
}

// UnsafeRegistryAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegistryAPIServer will
// result in compilation errors.
type UnsafeRegistryAPIServer interface {
	mustEmbedUnimplementedRegistryAPIServer()
}

func RegisterRegistryAPIServer(s grpc.ServiceRegistrar, srv RegistryAPIServer) {
	s.RegisterService(&RegistryAPI_ServiceDesc, srv)
}

func _RegistryAPI_GetAppProviders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAppProvidersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryAPIServer).GetAppProviders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryAPI_GetAppProviders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryAPIServer).GetAppProviders(ctx, req.(*GetAppProvidersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryAPI_AddAppProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddAppProviderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryAPIServer).AddAppProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryAPI_AddAppProvider_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryAPIServer).AddAppProvider(ctx, req.(*AddAppProviderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryAPI_ListAppProviders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAppProvidersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryAPIServer).ListAppProviders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryAPI_ListAppProviders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryAPIServer).ListAppProviders(ctx, req.(*ListAppProvidersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryAPI_ListSupportedMimeTypes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSupportedMimeTypesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryAPIServer).ListSupportedMimeTypes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryAPI_ListSupportedMimeTypes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryAPIServer).ListSupportedMimeTypes(ctx, req.(*ListSupportedMimeTypesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryAPI_GetDefaultAppProviderForMimeType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDefaultAppProviderForMimeTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryAPIServer).GetDefaultAppProviderForMimeType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryAPI_GetDefaultAppProviderForMimeType_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryAPIServer).GetDefaultAppProviderForMimeType(ctx, req.(*GetDefaultAppProviderForMimeTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryAPI_SetDefaultAppProviderForMimeType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDefaultAppProviderForMimeTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryAPIServer).SetDefaultAppProviderForMimeType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryAPI_SetDefaultAppProviderForMimeType_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryAPIServer).SetDefaultAppProviderForMimeType(ctx, req.(*SetDefaultAppProviderForMimeTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegistryAPI_ServiceDesc is the grpc.ServiceDesc for RegistryAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegistryAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cs3.app.registry.v1beta1.RegistryAPI",
	HandlerType: (*RegistryAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAppProviders",
			Handler:    _RegistryAPI_GetAppProviders_Handler,
		},
		{
			MethodName: "AddAppProvider",
			Handler:    _RegistryAPI_AddAppProvider_Handler,
		},
		{
			MethodName: "ListAppProviders",
			Handler:    _RegistryAPI_ListAppProviders_Handler,
		},
		{
			MethodName: "ListSupportedMimeTypes",
			Handler:    _RegistryAPI_ListSupportedMimeTypes_Handler,
		},
		{
			MethodName: "GetDefaultAppProviderForMimeType",
			Handler:    _RegistryAPI_GetDefaultAppProviderForMimeType_Handler,
		},
		{
			MethodName: "SetDefaultAppProviderForMimeType",
			Handler:    _RegistryAPI_SetDefaultAppProviderForMimeType_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cs3/app/registry/v1beta1/registry_api.proto",
}