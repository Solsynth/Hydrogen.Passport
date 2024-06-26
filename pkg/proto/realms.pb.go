// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: realms.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RealmLookupWithUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *RealmLookupWithUserRequest) Reset() {
	*x = RealmLookupWithUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_realms_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RealmLookupWithUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RealmLookupWithUserRequest) ProtoMessage() {}

func (x *RealmLookupWithUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_realms_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RealmLookupWithUserRequest.ProtoReflect.Descriptor instead.
func (*RealmLookupWithUserRequest) Descriptor() ([]byte, []int) {
	return file_realms_proto_rawDescGZIP(), []int{0}
}

func (x *RealmLookupWithUserRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type RealmLookupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          *uint64 `protobuf:"varint,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	Alias       *string `protobuf:"bytes,2,opt,name=alias,proto3,oneof" json:"alias,omitempty"`
	IsPublic    *bool   `protobuf:"varint,3,opt,name=is_public,json=isPublic,proto3,oneof" json:"is_public,omitempty"`
	IsCommunity *bool   `protobuf:"varint,4,opt,name=is_community,json=isCommunity,proto3,oneof" json:"is_community,omitempty"`
}

func (x *RealmLookupRequest) Reset() {
	*x = RealmLookupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_realms_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RealmLookupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RealmLookupRequest) ProtoMessage() {}

func (x *RealmLookupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_realms_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RealmLookupRequest.ProtoReflect.Descriptor instead.
func (*RealmLookupRequest) Descriptor() ([]byte, []int) {
	return file_realms_proto_rawDescGZIP(), []int{1}
}

func (x *RealmLookupRequest) GetId() uint64 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return 0
}

func (x *RealmLookupRequest) GetAlias() string {
	if x != nil && x.Alias != nil {
		return *x.Alias
	}
	return ""
}

func (x *RealmLookupRequest) GetIsPublic() bool {
	if x != nil && x.IsPublic != nil {
		return *x.IsPublic
	}
	return false
}

func (x *RealmLookupRequest) GetIsCommunity() bool {
	if x != nil && x.IsCommunity != nil {
		return *x.IsCommunity
	}
	return false
}

type RealmResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Alias       string `protobuf:"bytes,2,opt,name=alias,proto3" json:"alias,omitempty"`
	Name        string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	IsPublic    bool   `protobuf:"varint,5,opt,name=is_public,json=isPublic,proto3" json:"is_public,omitempty"`
	IsCommunity bool   `protobuf:"varint,6,opt,name=is_community,json=isCommunity,proto3" json:"is_community,omitempty"`
}

func (x *RealmResponse) Reset() {
	*x = RealmResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_realms_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RealmResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RealmResponse) ProtoMessage() {}

func (x *RealmResponse) ProtoReflect() protoreflect.Message {
	mi := &file_realms_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RealmResponse.ProtoReflect.Descriptor instead.
func (*RealmResponse) Descriptor() ([]byte, []int) {
	return file_realms_proto_rawDescGZIP(), []int{2}
}

func (x *RealmResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RealmResponse) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

func (x *RealmResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RealmResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *RealmResponse) GetIsPublic() bool {
	if x != nil {
		return x.IsPublic
	}
	return false
}

func (x *RealmResponse) GetIsCommunity() bool {
	if x != nil {
		return x.IsCommunity
	}
	return false
}

type ListRealmResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*RealmResponse `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *ListRealmResponse) Reset() {
	*x = ListRealmResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_realms_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRealmResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRealmResponse) ProtoMessage() {}

func (x *ListRealmResponse) ProtoReflect() protoreflect.Message {
	mi := &file_realms_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRealmResponse.ProtoReflect.Descriptor instead.
func (*ListRealmResponse) Descriptor() ([]byte, []int) {
	return file_realms_proto_rawDescGZIP(), []int{3}
}

func (x *ListRealmResponse) GetData() []*RealmResponse {
	if x != nil {
		return x.Data
	}
	return nil
}

type RealmMemberLookupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RealmId uint64  `protobuf:"varint,1,opt,name=realm_id,json=realmId,proto3" json:"realm_id,omitempty"`
	UserId  *uint64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
}

func (x *RealmMemberLookupRequest) Reset() {
	*x = RealmMemberLookupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_realms_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RealmMemberLookupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RealmMemberLookupRequest) ProtoMessage() {}

func (x *RealmMemberLookupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_realms_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RealmMemberLookupRequest.ProtoReflect.Descriptor instead.
func (*RealmMemberLookupRequest) Descriptor() ([]byte, []int) {
	return file_realms_proto_rawDescGZIP(), []int{4}
}

func (x *RealmMemberLookupRequest) GetRealmId() uint64 {
	if x != nil {
		return x.RealmId
	}
	return 0
}

func (x *RealmMemberLookupRequest) GetUserId() uint64 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

type RealmMemberResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RealmId    uint64 `protobuf:"varint,1,opt,name=realm_id,json=realmId,proto3" json:"realm_id,omitempty"`
	UserId     uint64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	PowerLevel int32  `protobuf:"varint,3,opt,name=power_level,json=powerLevel,proto3" json:"power_level,omitempty"`
}

func (x *RealmMemberResponse) Reset() {
	*x = RealmMemberResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_realms_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RealmMemberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RealmMemberResponse) ProtoMessage() {}

func (x *RealmMemberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_realms_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RealmMemberResponse.ProtoReflect.Descriptor instead.
func (*RealmMemberResponse) Descriptor() ([]byte, []int) {
	return file_realms_proto_rawDescGZIP(), []int{5}
}

func (x *RealmMemberResponse) GetRealmId() uint64 {
	if x != nil {
		return x.RealmId
	}
	return 0
}

func (x *RealmMemberResponse) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *RealmMemberResponse) GetPowerLevel() int32 {
	if x != nil {
		return x.PowerLevel
	}
	return 0
}

type ListRealmMemberResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*RealmMemberResponse `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *ListRealmMemberResponse) Reset() {
	*x = ListRealmMemberResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_realms_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRealmMemberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRealmMemberResponse) ProtoMessage() {}

func (x *ListRealmMemberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_realms_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRealmMemberResponse.ProtoReflect.Descriptor instead.
func (*ListRealmMemberResponse) Descriptor() ([]byte, []int) {
	return file_realms_proto_rawDescGZIP(), []int{6}
}

func (x *ListRealmMemberResponse) GetData() []*RealmMemberResponse {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_realms_proto protoreflect.FileDescriptor

var file_realms_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x35, 0x0a, 0x1a, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4c, 0x6f, 0x6f, 0x6b, 0x75,
	0x70, 0x57, 0x69, 0x74, 0x68, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0xbe, 0x01, 0x0a, 0x12, 0x52, 0x65,
	0x61, 0x6c, 0x6d, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x02,
	0x69, 0x64, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x88, 0x01, 0x01,
	0x12, 0x20, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x48, 0x02, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x88,
	0x01, 0x01, 0x12, 0x26, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69,
	0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x48, 0x03, 0x52, 0x0b, 0x69, 0x73, 0x43, 0x6f,
	0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79, 0x88, 0x01, 0x01, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x69,
	0x64, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x42, 0x0c, 0x0a, 0x0a, 0x5f,
	0x69, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x69, 0x73,
	0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79, 0x22, 0xab, 0x01, 0x0a, 0x0d, 0x52,
	0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x6c, 0x69,
	0x61, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x6d,
	0x75, 0x6e, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x43,
	0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79, 0x22, 0x3d, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x5f, 0x0a, 0x18, 0x52, 0x65, 0x61, 0x6c, 0x6d,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x49, 0x64, 0x12, 0x1c,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x48,
	0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x22, 0x6a, 0x0a, 0x13, 0x52, 0x65, 0x61, 0x6c,
	0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x07, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x6c, 0x65, 0x76,
	0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x4c,
	0x65, 0x76, 0x65, 0x6c, 0x22, 0x49, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x6c,
	0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32,
	0xde, 0x03, 0x0a, 0x06, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x73, 0x12, 0x48, 0x0a, 0x12, 0x4c, 0x69,
	0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x74, 0x79, 0x52, 0x65, 0x61, 0x6c, 0x6d,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x53, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x21, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x57, 0x69,
	0x74, 0x68, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4f, 0x0a, 0x0e, 0x4c, 0x69, 0x73,
	0x74, 0x4f, 0x77, 0x6e, 0x65, 0x64, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x21, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x57,
	0x69, 0x74, 0x68, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x08, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52,
	0x65, 0x61, 0x6c, 0x6d, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x0f, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x4f, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_realms_proto_rawDescOnce sync.Once
	file_realms_proto_rawDescData = file_realms_proto_rawDesc
)

func file_realms_proto_rawDescGZIP() []byte {
	file_realms_proto_rawDescOnce.Do(func() {
		file_realms_proto_rawDescData = protoimpl.X.CompressGZIP(file_realms_proto_rawDescData)
	})
	return file_realms_proto_rawDescData
}

var file_realms_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_realms_proto_goTypes = []interface{}{
	(*RealmLookupWithUserRequest)(nil), // 0: proto.RealmLookupWithUserRequest
	(*RealmLookupRequest)(nil),         // 1: proto.RealmLookupRequest
	(*RealmResponse)(nil),              // 2: proto.RealmResponse
	(*ListRealmResponse)(nil),          // 3: proto.ListRealmResponse
	(*RealmMemberLookupRequest)(nil),   // 4: proto.RealmMemberLookupRequest
	(*RealmMemberResponse)(nil),        // 5: proto.RealmMemberResponse
	(*ListRealmMemberResponse)(nil),    // 6: proto.ListRealmMemberResponse
	(*emptypb.Empty)(nil),              // 7: google.protobuf.Empty
}
var file_realms_proto_depIdxs = []int32{
	2, // 0: proto.ListRealmResponse.data:type_name -> proto.RealmResponse
	5, // 1: proto.ListRealmMemberResponse.data:type_name -> proto.RealmMemberResponse
	7, // 2: proto.Realms.ListCommunityRealm:input_type -> google.protobuf.Empty
	0, // 3: proto.Realms.ListAvailableRealm:input_type -> proto.RealmLookupWithUserRequest
	0, // 4: proto.Realms.ListOwnedRealm:input_type -> proto.RealmLookupWithUserRequest
	1, // 5: proto.Realms.GetRealm:input_type -> proto.RealmLookupRequest
	4, // 6: proto.Realms.ListRealmMember:input_type -> proto.RealmMemberLookupRequest
	4, // 7: proto.Realms.GetRealmMember:input_type -> proto.RealmMemberLookupRequest
	3, // 8: proto.Realms.ListCommunityRealm:output_type -> proto.ListRealmResponse
	3, // 9: proto.Realms.ListAvailableRealm:output_type -> proto.ListRealmResponse
	3, // 10: proto.Realms.ListOwnedRealm:output_type -> proto.ListRealmResponse
	2, // 11: proto.Realms.GetRealm:output_type -> proto.RealmResponse
	6, // 12: proto.Realms.ListRealmMember:output_type -> proto.ListRealmMemberResponse
	5, // 13: proto.Realms.GetRealmMember:output_type -> proto.RealmMemberResponse
	8, // [8:14] is the sub-list for method output_type
	2, // [2:8] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_realms_proto_init() }
func file_realms_proto_init() {
	if File_realms_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_realms_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RealmLookupWithUserRequest); i {
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
		file_realms_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RealmLookupRequest); i {
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
		file_realms_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RealmResponse); i {
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
		file_realms_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRealmResponse); i {
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
		file_realms_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RealmMemberLookupRequest); i {
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
		file_realms_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RealmMemberResponse); i {
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
		file_realms_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRealmMemberResponse); i {
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
	file_realms_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_realms_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_realms_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_realms_proto_goTypes,
		DependencyIndexes: file_realms_proto_depIdxs,
		MessageInfos:      file_realms_proto_msgTypes,
	}.Build()
	File_realms_proto = out.File
	file_realms_proto_rawDesc = nil
	file_realms_proto_goTypes = nil
	file_realms_proto_depIdxs = nil
}
