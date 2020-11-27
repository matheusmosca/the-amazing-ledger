// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.13.0
// source: ledger/ledger.proto

package proto

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Operation has the possible operations to be used in Entry.
type Operation int32

const (
	// Don't use. It's just the default value.
	Operation_OPERATION_UNSPECIFIED Operation = 0
	// Debit operation.
	Operation_OPERATION_DEBIT Operation = 1
	// Credit operation.
	Operation_OPERATION_CREDIT Operation = 2
)

// Enum value maps for Operation.
var (
	Operation_name = map[int32]string{
		0: "OPERATION_UNSPECIFIED",
		1: "OPERATION_DEBIT",
		2: "OPERATION_CREDIT",
	}
	Operation_value = map[string]int32{
		"OPERATION_UNSPECIFIED": 0,
		"OPERATION_DEBIT":       1,
		"OPERATION_CREDIT":      2,
	}
)

func (x Operation) Enum() *Operation {
	p := new(Operation)
	*p = x
	return p
}

func (x Operation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Operation) Descriptor() protoreflect.EnumDescriptor {
	return file_ledger_ledger_proto_enumTypes[0].Descriptor()
}

func (Operation) Type() protoreflect.EnumType {
	return &file_ledger_ledger_proto_enumTypes[0]
}

func (x Operation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Operation.Descriptor instead.
func (Operation) EnumDescriptor() ([]byte, []int) {
	return file_ledger_ledger_proto_rawDescGZIP(), []int{0}
}

// ServingStatus is the enum of the possible health check status
type HealthCheckResponse_ServingStatus int32

const (
	// Don't use. It's just the default value.
	HealthCheckResponse_SERVING_STATUS_UNKNOWN_UNSPECIFIED HealthCheckResponse_ServingStatus = 0
	// Healthy
	HealthCheckResponse_SERVING_STATUS_SERVING HealthCheckResponse_ServingStatus = 1
	// Unhealthy
	HealthCheckResponse_SERVING_STATUS_NOT_SERVING HealthCheckResponse_ServingStatus = 2
	// Used only when streaming
	HealthCheckResponse_SERVING_STATUS_SERVICE_UNKNOWN HealthCheckResponse_ServingStatus = 3
)

// Enum value maps for HealthCheckResponse_ServingStatus.
var (
	HealthCheckResponse_ServingStatus_name = map[int32]string{
		0: "SERVING_STATUS_UNKNOWN_UNSPECIFIED",
		1: "SERVING_STATUS_SERVING",
		2: "SERVING_STATUS_NOT_SERVING",
		3: "SERVING_STATUS_SERVICE_UNKNOWN",
	}
	HealthCheckResponse_ServingStatus_value = map[string]int32{
		"SERVING_STATUS_UNKNOWN_UNSPECIFIED": 0,
		"SERVING_STATUS_SERVING":             1,
		"SERVING_STATUS_NOT_SERVING":         2,
		"SERVING_STATUS_SERVICE_UNKNOWN":     3,
	}
)

func (x HealthCheckResponse_ServingStatus) Enum() *HealthCheckResponse_ServingStatus {
	p := new(HealthCheckResponse_ServingStatus)
	*p = x
	return p
}

func (x HealthCheckResponse_ServingStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HealthCheckResponse_ServingStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_ledger_ledger_proto_enumTypes[1].Descriptor()
}

func (HealthCheckResponse_ServingStatus) Type() protoreflect.EnumType {
	return &file_ledger_ledger_proto_enumTypes[1]
}

func (x HealthCheckResponse_ServingStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HealthCheckResponse_ServingStatus.Descriptor instead.
func (HealthCheckResponse_ServingStatus) EnumDescriptor() ([]byte, []int) {
	return file_ledger_ledger_proto_rawDescGZIP(), []int{4, 0}
}

// SaveTransactionRequest represents a transaction to be saved. A transaction must
// have at least two entries, with a valid balance. More info here:
// https://en.wikipedia.org/wiki/Double-entry_bookkeeping
type CreateTransactionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID (UUID) to link the entries to a transaction.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The list of entries, where len(entries) must be >= 2.
	Entries []*Entry `protobuf:"bytes,2,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *CreateTransactionRequest) Reset() {
	*x = CreateTransactionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ledger_ledger_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTransactionRequest) ProtoMessage() {}

func (x *CreateTransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_ledger_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTransactionRequest.ProtoReflect.Descriptor instead.
func (*CreateTransactionRequest) Descriptor() ([]byte, []int) {
	return file_ledger_ledger_proto_rawDescGZIP(), []int{0}
}

func (x *CreateTransactionRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CreateTransactionRequest) GetEntries() []*Entry {
	if x != nil {
		return x.Entries
	}
	return nil
}

// Entry represents a new entry on the Ledger.
type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// It's the idempotency key, and must be unique (UUID).
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Account involved in the operation.
	AccountId string `protobuf:"bytes,2,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// To deal with optimistic lock.
	ExpectedVersion uint64 `protobuf:"varint,3,opt,name=expected_version,json=expectedVersion,proto3" json:"expected_version,omitempty"`
	// Operation: debit or credit.
	Operation Operation `protobuf:"varint,4,opt,name=operation,proto3,enum=ledger.Operation" json:"operation,omitempty"`
	// Amount (in cents).
	Amount int32 `protobuf:"varint,5,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ledger_ledger_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_ledger_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_ledger_ledger_proto_rawDescGZIP(), []int{1}
}

func (x *Entry) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Entry) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *Entry) GetExpectedVersion() uint64 {
	if x != nil {
		return x.ExpectedVersion
	}
	return 0
}

func (x *Entry) GetOperation() Operation {
	if x != nil {
		return x.Operation
	}
	return Operation_OPERATION_UNSPECIFIED
}

func (x *Entry) GetAmount() int32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

// GetAccountBalance Request
type GetAccountBalanceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The account name
	AccountPath string `protobuf:"bytes,1,opt,name=account_path,json=accountPath,proto3" json:"account_path,omitempty"`
}

func (x *GetAccountBalanceRequest) Reset() {
	*x = GetAccountBalanceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ledger_ledger_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAccountBalanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccountBalanceRequest) ProtoMessage() {}

func (x *GetAccountBalanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_ledger_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccountBalanceRequest.ProtoReflect.Descriptor instead.
func (*GetAccountBalanceRequest) Descriptor() ([]byte, []int) {
	return file_ledger_ledger_proto_rawDescGZIP(), []int{2}
}

func (x *GetAccountBalanceRequest) GetAccountPath() string {
	if x != nil {
		return x.AccountPath
	}
	return ""
}

// GetAccountBalance Response
type GetAccountBalanceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The account name
	AccountPath string `protobuf:"bytes,1,opt,name=account_path,json=accountPath,proto3" json:"account_path,omitempty"`
	// The account version
	CurrentVersion uint64 `protobuf:"varint,2,opt,name=current_version,json=currentVersion,proto3" json:"current_version,omitempty"`
	// All credit accumulated
	TotalCredit int64 `protobuf:"varint,3,opt,name=total_credit,json=totalCredit,proto3" json:"total_credit,omitempty"`
	// All debit accumulated
	TotalDebit int64 `protobuf:"varint,4,opt,name=total_debit,json=totalDebit,proto3" json:"total_debit,omitempty"`
	// The Account balance
	Balance int64 `protobuf:"varint,5,opt,name=balance,proto3" json:"balance,omitempty"`
}

func (x *GetAccountBalanceResponse) Reset() {
	*x = GetAccountBalanceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ledger_ledger_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAccountBalanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccountBalanceResponse) ProtoMessage() {}

func (x *GetAccountBalanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_ledger_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccountBalanceResponse.ProtoReflect.Descriptor instead.
func (*GetAccountBalanceResponse) Descriptor() ([]byte, []int) {
	return file_ledger_ledger_proto_rawDescGZIP(), []int{3}
}

func (x *GetAccountBalanceResponse) GetAccountPath() string {
	if x != nil {
		return x.AccountPath
	}
	return ""
}

func (x *GetAccountBalanceResponse) GetCurrentVersion() uint64 {
	if x != nil {
		return x.CurrentVersion
	}
	return 0
}

func (x *GetAccountBalanceResponse) GetTotalCredit() int64 {
	if x != nil {
		return x.TotalCredit
	}
	return 0
}

func (x *GetAccountBalanceResponse) GetTotalDebit() int64 {
	if x != nil {
		return x.TotalDebit
	}
	return 0
}

func (x *GetAccountBalanceResponse) GetBalance() int64 {
	if x != nil {
		return x.Balance
	}
	return 0
}

//https://github.com/grpc/grpc/blob/master/doc/health-checking.md
// HealthCheckResponse is the health check status
type HealthCheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Server status.
	Status HealthCheckResponse_ServingStatus `protobuf:"varint,1,opt,name=status,proto3,enum=ledger.HealthCheckResponse_ServingStatus" json:"status,omitempty"`
}

func (x *HealthCheckResponse) Reset() {
	*x = HealthCheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ledger_ledger_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthCheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckResponse) ProtoMessage() {}

func (x *HealthCheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ledger_ledger_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheckResponse.ProtoReflect.Descriptor instead.
func (*HealthCheckResponse) Descriptor() ([]byte, []int) {
	return file_ledger_ledger_proto_rawDescGZIP(), []int{4}
}

func (x *HealthCheckResponse) GetStatus() HealthCheckResponse_ServingStatus {
	if x != nil {
		return x.Status
	}
	return HealthCheckResponse_SERVING_STATUS_UNKNOWN_UNSPECIFIED
}

var File_ledger_ledger_proto protoreflect.FileDescriptor

var file_ledger_ledger_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2f, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x53, 0x0a, 0x18, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x27, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0xaa, 0x01,
	0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x10, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74,
	0x65, 0x64, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0f, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x2f, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3d, 0x0a, 0x18, 0x47, 0x65,
	0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x61, 0x74, 0x68, 0x22, 0xc5, 0x01, 0x0a, 0x19, 0x47, 0x65,
	0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x61, 0x74, 0x68, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x72, 0x65,
	0x64, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f,
	0x64, 0x65, 0x62, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x44, 0x65, 0x62, 0x69, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x22, 0xf2, 0x01, 0x0a, 0x13, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29, 0x2e, 0x6c, 0x65, 0x64, 0x67,
	0x65, 0x72, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x6e, 0x67, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x97, 0x01, 0x0a,
	0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x26,
	0x0a, 0x22, 0x53, 0x45, 0x52, 0x56, 0x49, 0x4e, 0x47, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53,
	0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x52, 0x56, 0x49, 0x4e,
	0x47, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x45, 0x52, 0x56, 0x49, 0x4e, 0x47,
	0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x45, 0x52, 0x56, 0x49, 0x4e, 0x47, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x53, 0x45, 0x52, 0x56, 0x49, 0x4e, 0x47,
	0x10, 0x02, 0x12, 0x22, 0x0a, 0x1e, 0x53, 0x45, 0x52, 0x56, 0x49, 0x4e, 0x47, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x03, 0x2a, 0x51, 0x0a, 0x09, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x15, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x13,
	0x0a, 0x0f, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x42, 0x49,
	0x54, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x5f, 0x43, 0x52, 0x45, 0x44, 0x49, 0x54, 0x10, 0x02, 0x32, 0x8b, 0x02, 0x0a, 0x0d, 0x4c, 0x65,
	0x64, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6e, 0x0a, 0x11, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1f, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x19, 0x22, 0x14, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x3a, 0x01, 0x2a, 0x12, 0x89, 0x01, 0x0a, 0x11,
	0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x12, 0x20, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x29, 0x12, 0x27,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73,
	0x2f, 0x7b, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x7d, 0x2f,
	0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x32, 0x5e, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x12, 0x54, 0x0a, 0x05, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x1b, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x48, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x2d, 0x63, 0x6f, 0x2f, 0x74,
	0x68, 0x65, 0x2d, 0x61, 0x6d, 0x61, 0x7a, 0x69, 0x6e, 0x67, 0x2d, 0x6c, 0x65, 0x64, 0x67, 0x65,
	0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ledger_ledger_proto_rawDescOnce sync.Once
	file_ledger_ledger_proto_rawDescData = file_ledger_ledger_proto_rawDesc
)

func file_ledger_ledger_proto_rawDescGZIP() []byte {
	file_ledger_ledger_proto_rawDescOnce.Do(func() {
		file_ledger_ledger_proto_rawDescData = protoimpl.X.CompressGZIP(file_ledger_ledger_proto_rawDescData)
	})
	return file_ledger_ledger_proto_rawDescData
}

var file_ledger_ledger_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_ledger_ledger_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_ledger_ledger_proto_goTypes = []interface{}{
	(Operation)(0),                         // 0: ledger.Operation
	(HealthCheckResponse_ServingStatus)(0), // 1: ledger.HealthCheckResponse.ServingStatus
	(*CreateTransactionRequest)(nil),       // 2: ledger.CreateTransactionRequest
	(*Entry)(nil),                          // 3: ledger.Entry
	(*GetAccountBalanceRequest)(nil),       // 4: ledger.GetAccountBalanceRequest
	(*GetAccountBalanceResponse)(nil),      // 5: ledger.GetAccountBalanceResponse
	(*HealthCheckResponse)(nil),            // 6: ledger.HealthCheckResponse
	(*empty.Empty)(nil),                    // 7: google.protobuf.Empty
}
var file_ledger_ledger_proto_depIdxs = []int32{
	3, // 0: ledger.CreateTransactionRequest.entries:type_name -> ledger.Entry
	0, // 1: ledger.Entry.operation:type_name -> ledger.Operation
	1, // 2: ledger.HealthCheckResponse.status:type_name -> ledger.HealthCheckResponse.ServingStatus
	2, // 3: ledger.LedgerService.CreateTransaction:input_type -> ledger.CreateTransactionRequest
	4, // 4: ledger.LedgerService.GetAccountBalance:input_type -> ledger.GetAccountBalanceRequest
	7, // 5: ledger.Health.Check:input_type -> google.protobuf.Empty
	7, // 6: ledger.LedgerService.CreateTransaction:output_type -> google.protobuf.Empty
	5, // 7: ledger.LedgerService.GetAccountBalance:output_type -> ledger.GetAccountBalanceResponse
	6, // 8: ledger.Health.Check:output_type -> ledger.HealthCheckResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_ledger_ledger_proto_init() }
func file_ledger_ledger_proto_init() {
	if File_ledger_ledger_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ledger_ledger_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTransactionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ledger_ledger_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ledger_ledger_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAccountBalanceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ledger_ledger_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAccountBalanceResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ledger_ledger_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthCheckResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ledger_ledger_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_ledger_ledger_proto_goTypes,
		DependencyIndexes: file_ledger_ledger_proto_depIdxs,
		EnumInfos:         file_ledger_ledger_proto_enumTypes,
		MessageInfos:      file_ledger_ledger_proto_msgTypes,
	}.Build()
	File_ledger_ledger_proto = out.File
	file_ledger_ledger_proto_rawDesc = nil
	file_ledger_ledger_proto_goTypes = nil
	file_ledger_ledger_proto_depIdxs = nil
}
