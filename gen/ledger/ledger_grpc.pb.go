// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package ledger

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LedgerServiceClient is the client API for LedgerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LedgerServiceClient interface {
	CreateTransaction(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetAccountBalance(ctx context.Context, in *GetAccountBalanceRequest, opts ...grpc.CallOption) (*GetAccountBalanceResponse, error)
	ListAccountEntries(ctx context.Context, in *ListAccountEntriesRequest, opts ...grpc.CallOption) (*ListAccountEntriesResponse, error)
	GetSyntheticReport(ctx context.Context, in *GetSyntheticReportRequest, opts ...grpc.CallOption) (*GetSyntheticReportResponse, error)
}

type ledgerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLedgerServiceClient(cc grpc.ClientConnInterface) LedgerServiceClient {
	return &ledgerServiceClient{cc}
}

func (c *ledgerServiceClient) CreateTransaction(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/ledger.LedgerService/CreateTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) GetAccountBalance(ctx context.Context, in *GetAccountBalanceRequest, opts ...grpc.CallOption) (*GetAccountBalanceResponse, error) {
	out := new(GetAccountBalanceResponse)
	err := c.cc.Invoke(ctx, "/ledger.LedgerService/GetAccountBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) ListAccountEntries(ctx context.Context, in *ListAccountEntriesRequest, opts ...grpc.CallOption) (*ListAccountEntriesResponse, error) {
	out := new(ListAccountEntriesResponse)
	err := c.cc.Invoke(ctx, "/ledger.LedgerService/ListAccountEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) GetSyntheticReport(ctx context.Context, in *GetSyntheticReportRequest, opts ...grpc.CallOption) (*GetSyntheticReportResponse, error) {
	out := new(GetSyntheticReportResponse)
	err := c.cc.Invoke(ctx, "/ledger.LedgerService/GetSyntheticReport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LedgerServiceServer is the server API for LedgerService service.
// All implementations should embed UnimplementedLedgerServiceServer
// for forward compatibility
type LedgerServiceServer interface {
	CreateTransaction(context.Context, *CreateTransactionRequest) (*emptypb.Empty, error)
	GetAccountBalance(context.Context, *GetAccountBalanceRequest) (*GetAccountBalanceResponse, error)
	ListAccountEntries(context.Context, *ListAccountEntriesRequest) (*ListAccountEntriesResponse, error)
	GetSyntheticReport(context.Context, *GetSyntheticReportRequest) (*GetSyntheticReportResponse, error)
}

// UnimplementedLedgerServiceServer should be embedded to have forward compatible implementations.
type UnimplementedLedgerServiceServer struct {
}

func (UnimplementedLedgerServiceServer) CreateTransaction(context.Context, *CreateTransactionRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTransaction not implemented")
}
func (UnimplementedLedgerServiceServer) GetAccountBalance(context.Context, *GetAccountBalanceRequest) (*GetAccountBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccountBalance not implemented")
}
func (UnimplementedLedgerServiceServer) ListAccountEntries(context.Context, *ListAccountEntriesRequest) (*ListAccountEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccountEntries not implemented")
}
func (UnimplementedLedgerServiceServer) GetSyntheticReport(context.Context, *GetSyntheticReportRequest) (*GetSyntheticReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSyntheticReport not implemented")
}

// UnsafeLedgerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LedgerServiceServer will
// result in compilation errors.
type UnsafeLedgerServiceServer interface {
	mustEmbedUnimplementedLedgerServiceServer()
}

func RegisterLedgerServiceServer(s grpc.ServiceRegistrar, srv LedgerServiceServer) {
	s.RegisterService(&LedgerService_ServiceDesc, srv)
}

func _LedgerService_CreateTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).CreateTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.LedgerService/CreateTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).CreateTransaction(ctx, req.(*CreateTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_GetAccountBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).GetAccountBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.LedgerService/GetAccountBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).GetAccountBalance(ctx, req.(*GetAccountBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_ListAccountEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAccountEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).ListAccountEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.LedgerService/ListAccountEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).ListAccountEntries(ctx, req.(*ListAccountEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_GetSyntheticReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSyntheticReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).GetSyntheticReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.LedgerService/GetSyntheticReport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).GetSyntheticReport(ctx, req.(*GetSyntheticReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LedgerService_ServiceDesc is the grpc.ServiceDesc for LedgerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LedgerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ledger.LedgerService",
	HandlerType: (*LedgerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTransaction",
			Handler:    _LedgerService_CreateTransaction_Handler,
		},
		{
			MethodName: "GetAccountBalance",
			Handler:    _LedgerService_GetAccountBalance_Handler,
		},
		{
			MethodName: "ListAccountEntries",
			Handler:    _LedgerService_ListAccountEntries_Handler,
		},
		{
			MethodName: "GetSyntheticReport",
			Handler:    _LedgerService_GetSyntheticReport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ledger/ledger.proto",
}

// HealthClient is the client API for Health service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HealthClient interface {
	// Check - checks the system health.
	Check(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type healthClient struct {
	cc grpc.ClientConnInterface
}

func NewHealthClient(cc grpc.ClientConnInterface) HealthClient {
	return &healthClient{cc}
}

func (c *healthClient) Check(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, "/ledger.Health/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HealthServer is the server API for Health service.
// All implementations should embed UnimplementedHealthServer
// for forward compatibility
type HealthServer interface {
	// Check - checks the system health.
	Check(context.Context, *emptypb.Empty) (*HealthCheckResponse, error)
}

// UnimplementedHealthServer should be embedded to have forward compatible implementations.
type UnimplementedHealthServer struct {
}

func (UnimplementedHealthServer) Check(context.Context, *emptypb.Empty) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}

// UnsafeHealthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HealthServer will
// result in compilation errors.
type UnsafeHealthServer interface {
	mustEmbedUnimplementedHealthServer()
}

func RegisterHealthServer(s grpc.ServiceRegistrar, srv HealthServer) {
	s.RegisterService(&Health_ServiceDesc, srv)
}

func _Health_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.Health/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthServer).Check(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Health_ServiceDesc is the grpc.ServiceDesc for Health service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Health_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ledger.Health",
	HandlerType: (*HealthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _Health_Check_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ledger/ledger.proto",
}
