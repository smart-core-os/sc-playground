// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.1
// source: pkg/playpb/playground.proto

package playpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AddDeviceTraitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	TraitName string `protobuf:"bytes,2,opt,name=trait_name,json=traitName,proto3" json:"trait_name,omitempty"`
}

func (x *AddDeviceTraitRequest) Reset() {
	*x = AddDeviceTraitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_playpb_playground_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddDeviceTraitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddDeviceTraitRequest) ProtoMessage() {}

func (x *AddDeviceTraitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_playpb_playground_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddDeviceTraitRequest.ProtoReflect.Descriptor instead.
func (*AddDeviceTraitRequest) Descriptor() ([]byte, []int) {
	return file_pkg_playpb_playground_proto_rawDescGZIP(), []int{0}
}

func (x *AddDeviceTraitRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AddDeviceTraitRequest) GetTraitName() string {
	if x != nil {
		return x.TraitName
	}
	return ""
}

type AddDeviceTraitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddDeviceTraitResponse) Reset() {
	*x = AddDeviceTraitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_playpb_playground_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddDeviceTraitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddDeviceTraitResponse) ProtoMessage() {}

func (x *AddDeviceTraitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_playpb_playground_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddDeviceTraitResponse.ProtoReflect.Descriptor instead.
func (*AddDeviceTraitResponse) Descriptor() ([]byte, []int) {
	return file_pkg_playpb_playground_proto_rawDescGZIP(), []int{1}
}

type ListSupportedTraitsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListSupportedTraitsRequest) Reset() {
	*x = ListSupportedTraitsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_playpb_playground_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListSupportedTraitsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSupportedTraitsRequest) ProtoMessage() {}

func (x *ListSupportedTraitsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_playpb_playground_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSupportedTraitsRequest.ProtoReflect.Descriptor instead.
func (*ListSupportedTraitsRequest) Descriptor() ([]byte, []int) {
	return file_pkg_playpb_playground_proto_rawDescGZIP(), []int{2}
}

type ListSupportedTraitsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TraitName []string `protobuf:"bytes,1,rep,name=trait_name,json=traitName,proto3" json:"trait_name,omitempty"`
}

func (x *ListSupportedTraitsResponse) Reset() {
	*x = ListSupportedTraitsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_playpb_playground_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListSupportedTraitsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSupportedTraitsResponse) ProtoMessage() {}

func (x *ListSupportedTraitsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_playpb_playground_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSupportedTraitsResponse.ProtoReflect.Descriptor instead.
func (*ListSupportedTraitsResponse) Descriptor() ([]byte, []int) {
	return file_pkg_playpb_playground_proto_rawDescGZIP(), []int{3}
}

func (x *ListSupportedTraitsResponse) GetTraitName() []string {
	if x != nil {
		return x.TraitName
	}
	return nil
}

type AddRemoteDeviceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Endpoint  string     `protobuf:"bytes,2,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	TraitName []string   `protobuf:"bytes,3,rep,name=trait_name,json=traitName,proto3" json:"trait_name,omitempty"`
	Tls       *RemoteTLS `protobuf:"bytes,4,opt,name=tls,proto3" json:"tls,omitempty"`
	Insecure  bool       `protobuf:"varint,5,opt,name=insecure,proto3" json:"insecure,omitempty"`
}

func (x *AddRemoteDeviceRequest) Reset() {
	*x = AddRemoteDeviceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_playpb_playground_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRemoteDeviceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRemoteDeviceRequest) ProtoMessage() {}

func (x *AddRemoteDeviceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_playpb_playground_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRemoteDeviceRequest.ProtoReflect.Descriptor instead.
func (*AddRemoteDeviceRequest) Descriptor() ([]byte, []int) {
	return file_pkg_playpb_playground_proto_rawDescGZIP(), []int{4}
}

func (x *AddRemoteDeviceRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AddRemoteDeviceRequest) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *AddRemoteDeviceRequest) GetTraitName() []string {
	if x != nil {
		return x.TraitName
	}
	return nil
}

func (x *AddRemoteDeviceRequest) GetTls() *RemoteTLS {
	if x != nil {
		return x.Tls
	}
	return nil
}

func (x *AddRemoteDeviceRequest) GetInsecure() bool {
	if x != nil {
		return x.Insecure
	}
	return false
}

type RemoteTLS struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerCaCert string `protobuf:"bytes,1,opt,name=server_ca_cert,json=serverCaCert,proto3" json:"server_ca_cert,omitempty"`
	SkipVerify   bool   `protobuf:"varint,2,opt,name=skip_verify,json=skipVerify,proto3" json:"skip_verify,omitempty"`
}

func (x *RemoteTLS) Reset() {
	*x = RemoteTLS{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_playpb_playground_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoteTLS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoteTLS) ProtoMessage() {}

func (x *RemoteTLS) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_playpb_playground_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoteTLS.ProtoReflect.Descriptor instead.
func (*RemoteTLS) Descriptor() ([]byte, []int) {
	return file_pkg_playpb_playground_proto_rawDescGZIP(), []int{5}
}

func (x *RemoteTLS) GetServerCaCert() string {
	if x != nil {
		return x.ServerCaCert
	}
	return ""
}

func (x *RemoteTLS) GetSkipVerify() bool {
	if x != nil {
		return x.SkipVerify
	}
	return false
}

type AddRemoteDeviceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddRemoteDeviceResponse) Reset() {
	*x = AddRemoteDeviceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_playpb_playground_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRemoteDeviceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRemoteDeviceResponse) ProtoMessage() {}

func (x *AddRemoteDeviceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_playpb_playground_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRemoteDeviceResponse.ProtoReflect.Descriptor instead.
func (*AddRemoteDeviceResponse) Descriptor() ([]byte, []int) {
	return file_pkg_playpb_playground_proto_rawDescGZIP(), []int{6}
}

type Performance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Frame   *durationpb.Duration `protobuf:"bytes,1,opt,name=frame,proto3" json:"frame,omitempty"`
	Capture *durationpb.Duration `protobuf:"bytes,2,opt,name=capture,proto3" json:"capture,omitempty"`
	Scrub   *durationpb.Duration `protobuf:"bytes,3,opt,name=scrub,proto3" json:"scrub,omitempty"`
	Respond *durationpb.Duration `protobuf:"bytes,4,opt,name=respond,proto3" json:"respond,omitempty"`
	Idle    *durationpb.Duration `protobuf:"bytes,5,opt,name=idle,proto3" json:"idle,omitempty"`
}

func (x *Performance) Reset() {
	*x = Performance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_playpb_playground_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Performance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Performance) ProtoMessage() {}

func (x *Performance) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_playpb_playground_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Performance.ProtoReflect.Descriptor instead.
func (*Performance) Descriptor() ([]byte, []int) {
	return file_pkg_playpb_playground_proto_rawDescGZIP(), []int{7}
}

func (x *Performance) GetFrame() *durationpb.Duration {
	if x != nil {
		return x.Frame
	}
	return nil
}

func (x *Performance) GetCapture() *durationpb.Duration {
	if x != nil {
		return x.Capture
	}
	return nil
}

func (x *Performance) GetScrub() *durationpb.Duration {
	if x != nil {
		return x.Scrub
	}
	return nil
}

func (x *Performance) GetRespond() *durationpb.Duration {
	if x != nil {
		return x.Respond
	}
	return nil
}

func (x *Performance) GetIdle() *durationpb.Duration {
	if x != nil {
		return x.Idle
	}
	return nil
}

type PullPerformanceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PullPerformanceRequest) Reset() {
	*x = PullPerformanceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_playpb_playground_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullPerformanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullPerformanceRequest) ProtoMessage() {}

func (x *PullPerformanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_playpb_playground_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullPerformanceRequest.ProtoReflect.Descriptor instead.
func (*PullPerformanceRequest) Descriptor() ([]byte, []int) {
	return file_pkg_playpb_playground_proto_rawDescGZIP(), []int{8}
}

type PullPerformanceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Performance *Performance           `protobuf:"bytes,1,opt,name=performance,proto3" json:"performance,omitempty"`
	ChangeTime  *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=change_time,json=changeTime,proto3" json:"change_time,omitempty"`
}

func (x *PullPerformanceResponse) Reset() {
	*x = PullPerformanceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_playpb_playground_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullPerformanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullPerformanceResponse) ProtoMessage() {}

func (x *PullPerformanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_playpb_playground_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullPerformanceResponse.ProtoReflect.Descriptor instead.
func (*PullPerformanceResponse) Descriptor() ([]byte, []int) {
	return file_pkg_playpb_playground_proto_rawDescGZIP(), []int{9}
}

func (x *PullPerformanceResponse) GetPerformance() *Performance {
	if x != nil {
		return x.Performance
	}
	return nil
}

func (x *PullPerformanceResponse) GetChangeTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ChangeTime
	}
	return nil
}

var File_pkg_playpb_playground_proto protoreflect.FileDescriptor

var file_pkg_playpb_playground_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x6c, 0x61, 0x79, 0x70, 0x62, 0x2f, 0x70, 0x6c, 0x61,
	0x79, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x73,
	0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x67, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x15, 0x41, 0x64, 0x64, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x72, 0x61, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x72, 0x61, 0x69, 0x74, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x72, 0x61, 0x69, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x22, 0x18, 0x0a, 0x16, 0x41, 0x64, 0x64, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x54, 0x72, 0x61, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1c,
	0x0a, 0x1a, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x54,
	0x72, 0x61, 0x69, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3c, 0x0a, 0x1b,
	0x4c, 0x69, 0x73, 0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x54, 0x72, 0x61,
	0x69, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74,
	0x72, 0x61, 0x69, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x09, 0x74, 0x72, 0x61, 0x69, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xba, 0x01, 0x0a, 0x16, 0x41,
	0x64, 0x64, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x72, 0x61, 0x69, 0x74, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x74, 0x72, 0x61, 0x69, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x03, 0x74, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x23, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x6c,
	0x61, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x6d,
	0x6f, 0x74, 0x65, 0x54, 0x4c, 0x53, 0x52, 0x03, 0x74, 0x6c, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x69,
	0x6e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69,
	0x6e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x65, 0x22, 0x52, 0x0a, 0x09, 0x52, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x54, 0x4c, 0x53, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x63,
	0x61, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x43, 0x61, 0x43, 0x65, 0x72, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x6b,
	0x69, 0x70, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0a, 0x73, 0x6b, 0x69, 0x70, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x22, 0x19, 0x0a, 0x17, 0x41,
	0x64, 0x64, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x88, 0x02, 0x0a, 0x0b, 0x50, 0x65, 0x72, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x2f, 0x0a, 0x05, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x05, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x33, 0x0a, 0x07, 0x63, 0x61, 0x70, 0x74, 0x75,
	0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x07, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x12, 0x2f, 0x0a, 0x05,
	0x73, 0x63, 0x72, 0x75, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x73, 0x63, 0x72, 0x75, 0x62, 0x12, 0x33, 0x0a,
	0x07, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x64, 0x12, 0x2d, 0x0a, 0x04, 0x69, 0x64, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x69, 0x64, 0x6c,
	0x65, 0x22, 0x18, 0x0a, 0x16, 0x50, 0x75, 0x6c, 0x6c, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x9f, 0x01, 0x0a, 0x17,
	0x50, 0x75, 0x6c, 0x6c, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x73,
	0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x67, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x6e, 0x63, 0x65, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65,
	0x12, 0x3b, 0x0a, 0x0b, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x32, 0xfb, 0x03,
	0x0a, 0x0d, 0x50, 0x6c, 0x61, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x41, 0x70, 0x69, 0x12,
	0x73, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x72, 0x61, 0x69,
	0x74, 0x12, 0x2f, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x6c,
	0x61, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x64, 0x64,
	0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x72, 0x61, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x30, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70,
	0x6c, 0x61, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x64,
	0x64, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x72, 0x61, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x82, 0x01, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x75, 0x70,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x54, 0x72, 0x61, 0x69, 0x74, 0x73, 0x12, 0x34, 0x2e, 0x73,
	0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x67, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x75, 0x70, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x64, 0x54, 0x72, 0x61, 0x69, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x35, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70,
	0x6c, 0x61, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x54, 0x72, 0x61, 0x69, 0x74,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x76, 0x0a, 0x0f, 0x41, 0x64, 0x64,
	0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x30, 0x2e, 0x73,
	0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x67, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31,
	0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x67,
	0x72, 0x6f, 0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x6d,
	0x6f, 0x74, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x78, 0x0a, 0x0f, 0x50, 0x75, 0x6c, 0x6c, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d,
	0x61, 0x6e, 0x63, 0x65, 0x12, 0x30, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x70, 0x6c, 0x61, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x50, 0x75, 0x6c, 0x6c, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x50, 0x75, 0x6c, 0x6c, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x33, 0x5a, 0x31, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6d, 0x61, 0x72, 0x74, 0x2d,
	0x63, 0x6f, 0x72, 0x65, 0x2d, 0x6f, 0x73, 0x2f, 0x73, 0x63, 0x2d, 0x70, 0x6c, 0x61, 0x79, 0x67,
	0x72, 0x6f, 0x75, 0x6e, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x6c, 0x61, 0x79, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_playpb_playground_proto_rawDescOnce sync.Once
	file_pkg_playpb_playground_proto_rawDescData = file_pkg_playpb_playground_proto_rawDesc
)

func file_pkg_playpb_playground_proto_rawDescGZIP() []byte {
	file_pkg_playpb_playground_proto_rawDescOnce.Do(func() {
		file_pkg_playpb_playground_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_playpb_playground_proto_rawDescData)
	})
	return file_pkg_playpb_playground_proto_rawDescData
}

var file_pkg_playpb_playground_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_pkg_playpb_playground_proto_goTypes = []interface{}{
	(*AddDeviceTraitRequest)(nil),       // 0: smartcore.playground.api.AddDeviceTraitRequest
	(*AddDeviceTraitResponse)(nil),      // 1: smartcore.playground.api.AddDeviceTraitResponse
	(*ListSupportedTraitsRequest)(nil),  // 2: smartcore.playground.api.ListSupportedTraitsRequest
	(*ListSupportedTraitsResponse)(nil), // 3: smartcore.playground.api.ListSupportedTraitsResponse
	(*AddRemoteDeviceRequest)(nil),      // 4: smartcore.playground.api.AddRemoteDeviceRequest
	(*RemoteTLS)(nil),                   // 5: smartcore.playground.api.RemoteTLS
	(*AddRemoteDeviceResponse)(nil),     // 6: smartcore.playground.api.AddRemoteDeviceResponse
	(*Performance)(nil),                 // 7: smartcore.playground.api.Performance
	(*PullPerformanceRequest)(nil),      // 8: smartcore.playground.api.PullPerformanceRequest
	(*PullPerformanceResponse)(nil),     // 9: smartcore.playground.api.PullPerformanceResponse
	(*durationpb.Duration)(nil),         // 10: google.protobuf.Duration
	(*timestamppb.Timestamp)(nil),       // 11: google.protobuf.Timestamp
}
var file_pkg_playpb_playground_proto_depIdxs = []int32{
	5,  // 0: smartcore.playground.api.AddRemoteDeviceRequest.tls:type_name -> smartcore.playground.api.RemoteTLS
	10, // 1: smartcore.playground.api.Performance.frame:type_name -> google.protobuf.Duration
	10, // 2: smartcore.playground.api.Performance.capture:type_name -> google.protobuf.Duration
	10, // 3: smartcore.playground.api.Performance.scrub:type_name -> google.protobuf.Duration
	10, // 4: smartcore.playground.api.Performance.respond:type_name -> google.protobuf.Duration
	10, // 5: smartcore.playground.api.Performance.idle:type_name -> google.protobuf.Duration
	7,  // 6: smartcore.playground.api.PullPerformanceResponse.performance:type_name -> smartcore.playground.api.Performance
	11, // 7: smartcore.playground.api.PullPerformanceResponse.change_time:type_name -> google.protobuf.Timestamp
	0,  // 8: smartcore.playground.api.PlaygroundApi.AddDeviceTrait:input_type -> smartcore.playground.api.AddDeviceTraitRequest
	2,  // 9: smartcore.playground.api.PlaygroundApi.ListSupportedTraits:input_type -> smartcore.playground.api.ListSupportedTraitsRequest
	4,  // 10: smartcore.playground.api.PlaygroundApi.AddRemoteDevice:input_type -> smartcore.playground.api.AddRemoteDeviceRequest
	8,  // 11: smartcore.playground.api.PlaygroundApi.PullPerformance:input_type -> smartcore.playground.api.PullPerformanceRequest
	1,  // 12: smartcore.playground.api.PlaygroundApi.AddDeviceTrait:output_type -> smartcore.playground.api.AddDeviceTraitResponse
	3,  // 13: smartcore.playground.api.PlaygroundApi.ListSupportedTraits:output_type -> smartcore.playground.api.ListSupportedTraitsResponse
	6,  // 14: smartcore.playground.api.PlaygroundApi.AddRemoteDevice:output_type -> smartcore.playground.api.AddRemoteDeviceResponse
	9,  // 15: smartcore.playground.api.PlaygroundApi.PullPerformance:output_type -> smartcore.playground.api.PullPerformanceResponse
	12, // [12:16] is the sub-list for method output_type
	8,  // [8:12] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_pkg_playpb_playground_proto_init() }
func file_pkg_playpb_playground_proto_init() {
	if File_pkg_playpb_playground_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_playpb_playground_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddDeviceTraitRequest); i {
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
		file_pkg_playpb_playground_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddDeviceTraitResponse); i {
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
		file_pkg_playpb_playground_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListSupportedTraitsRequest); i {
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
		file_pkg_playpb_playground_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListSupportedTraitsResponse); i {
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
		file_pkg_playpb_playground_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRemoteDeviceRequest); i {
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
		file_pkg_playpb_playground_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoteTLS); i {
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
		file_pkg_playpb_playground_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRemoteDeviceResponse); i {
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
		file_pkg_playpb_playground_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Performance); i {
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
		file_pkg_playpb_playground_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullPerformanceRequest); i {
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
		file_pkg_playpb_playground_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullPerformanceResponse); i {
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
			RawDescriptor: file_pkg_playpb_playground_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_playpb_playground_proto_goTypes,
		DependencyIndexes: file_pkg_playpb_playground_proto_depIdxs,
		MessageInfos:      file_pkg_playpb_playground_proto_msgTypes,
	}.Build()
	File_pkg_playpb_playground_proto = out.File
	file_pkg_playpb_playground_proto_rawDesc = nil
	file_pkg_playpb_playground_proto_goTypes = nil
	file_pkg_playpb_playground_proto_depIdxs = nil
}
