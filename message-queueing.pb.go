// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: message-queueing.proto

package queueing

import (
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

// Queue message used as input
//
// Data type used when submitting new messages to the system.
type RawQueueMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// raw message data as a byte array, the bytes will be returned as they are received
	Data []byte `protobuf:"bytes,1,req,name=data" json:"data,omitempty"`
	// attributes associated with the message, can be used to represent metadata
	Attributes map[string]string `protobuf:"bytes,2,rep,name=attributes" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (x *RawQueueMessage) Reset() {
	*x = RawQueueMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_queueing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RawQueueMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawQueueMessage) ProtoMessage() {}

func (x *RawQueueMessage) ProtoReflect() protoreflect.Message {
	mi := &file_message_queueing_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawQueueMessage.ProtoReflect.Descriptor instead.
func (*RawQueueMessage) Descriptor() ([]byte, []int) {
	return file_message_queueing_proto_rawDescGZIP(), []int{0}
}

func (x *RawQueueMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *RawQueueMessage) GetAttributes() map[string]string {
	if x != nil {
		return x.Attributes
	}
	return nil
}

// Queue message representation in the system
//
// QueueMessage represents the a RawQueueMessage after it has been submitted
// into the system. The RawQueueMessage is enhanced with metadata and identifiers
type QueueMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// raw message data as a byte array
	Data []byte `protobuf:"bytes,1,req,name=data" json:"data,omitempty"`
	// the message unique id (UUID) as a byte array
	MessageID []byte `protobuf:"bytes,2,req,name=messageID" json:"messageID,omitempty"`
	// timestamp when the message entered the system
	Timestamp *int64 `protobuf:"varint,3,req,name=timestamp" json:"timestamp,omitempty"`
	// md5 hash of the raw message data
	DataHash []byte `protobuf:"bytes,4,req,name=dataHash" json:"dataHash,omitempty"`
	// attributes associated with the message
	Attributes map[string]string `protobuf:"bytes,5,rep,name=attributes" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// lastRetrieved
	LastRetrieved *int64 `protobuf:"varint,6,opt,name=lastRetrieved" json:"lastRetrieved,omitempty"`
}

func (x *QueueMessage) Reset() {
	*x = QueueMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_queueing_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueMessage) ProtoMessage() {}

func (x *QueueMessage) ProtoReflect() protoreflect.Message {
	mi := &file_message_queueing_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueMessage.ProtoReflect.Descriptor instead.
func (*QueueMessage) Descriptor() ([]byte, []int) {
	return file_message_queueing_proto_rawDescGZIP(), []int{1}
}

func (x *QueueMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *QueueMessage) GetMessageID() []byte {
	if x != nil {
		return x.MessageID
	}
	return nil
}

func (x *QueueMessage) GetTimestamp() int64 {
	if x != nil && x.Timestamp != nil {
		return *x.Timestamp
	}
	return 0
}

func (x *QueueMessage) GetDataHash() []byte {
	if x != nil {
		return x.DataHash
	}
	return nil
}

func (x *QueueMessage) GetAttributes() map[string]string {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *QueueMessage) GetLastRetrieved() int64 {
	if x != nil && x.LastRetrieved != nil {
		return *x.LastRetrieved
	}
	return 0
}

type SubmitReceipt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok        *bool   `protobuf:"varint,1,req,name=ok" json:"ok,omitempty"`
	MessageID *string `protobuf:"bytes,2,req,name=messageID" json:"messageID,omitempty"`
	Timestamp *int64  `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
	DataHash  []byte  `protobuf:"bytes,4,opt,name=dataHash" json:"dataHash,omitempty"`
}

func (x *SubmitReceipt) Reset() {
	*x = SubmitReceipt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_queueing_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitReceipt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitReceipt) ProtoMessage() {}

func (x *SubmitReceipt) ProtoReflect() protoreflect.Message {
	mi := &file_message_queueing_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitReceipt.ProtoReflect.Descriptor instead.
func (*SubmitReceipt) Descriptor() ([]byte, []int) {
	return file_message_queueing_proto_rawDescGZIP(), []int{2}
}

func (x *SubmitReceipt) GetOk() bool {
	if x != nil && x.Ok != nil {
		return *x.Ok
	}
	return false
}

func (x *SubmitReceipt) GetMessageID() string {
	if x != nil && x.MessageID != nil {
		return *x.MessageID
	}
	return ""
}

func (x *SubmitReceipt) GetTimestamp() int64 {
	if x != nil && x.Timestamp != nil {
		return *x.Timestamp
	}
	return 0
}

func (x *SubmitReceipt) GetDataHash() []byte {
	if x != nil {
		return x.DataHash
	}
	return nil
}

type SubmitMessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueueID  *string            `protobuf:"bytes,1,req,name=queueID" json:"queueID,omitempty"`
	Messages []*RawQueueMessage `protobuf:"bytes,2,rep,name=messages" json:"messages,omitempty"`
}

func (x *SubmitMessagesRequest) Reset() {
	*x = SubmitMessagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_queueing_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitMessagesRequest) ProtoMessage() {}

func (x *SubmitMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_queueing_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitMessagesRequest.ProtoReflect.Descriptor instead.
func (*SubmitMessagesRequest) Descriptor() ([]byte, []int) {
	return file_message_queueing_proto_rawDescGZIP(), []int{3}
}

func (x *SubmitMessagesRequest) GetQueueID() string {
	if x != nil && x.QueueID != nil {
		return *x.QueueID
	}
	return ""
}

func (x *SubmitMessagesRequest) GetMessages() []*RawQueueMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

type SubmitMessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Receipts []*SubmitReceipt `protobuf:"bytes,1,rep,name=receipts" json:"receipts,omitempty"`
}

func (x *SubmitMessagesResponse) Reset() {
	*x = SubmitMessagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_queueing_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitMessagesResponse) ProtoMessage() {}

func (x *SubmitMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_message_queueing_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitMessagesResponse.ProtoReflect.Descriptor instead.
func (*SubmitMessagesResponse) Descriptor() ([]byte, []int) {
	return file_message_queueing_proto_rawDescGZIP(), []int{4}
}

func (x *SubmitMessagesResponse) GetReceipts() []*SubmitReceipt {
	if x != nil {
		return x.Receipts
	}
	return nil
}

type RetrieveMessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count *int32 `protobuf:"varint,1,req,name=count" json:"count,omitempty"`
}

func (x *RetrieveMessagesRequest) Reset() {
	*x = RetrieveMessagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_queueing_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetrieveMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetrieveMessagesRequest) ProtoMessage() {}

func (x *RetrieveMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_queueing_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetrieveMessagesRequest.ProtoReflect.Descriptor instead.
func (*RetrieveMessagesRequest) Descriptor() ([]byte, []int) {
	return file_message_queueing_proto_rawDescGZIP(), []int{5}
}

func (x *RetrieveMessagesRequest) GetCount() int32 {
	if x != nil && x.Count != nil {
		return *x.Count
	}
	return 0
}

type RetrieveMessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count    *int32          `protobuf:"varint,1,req,name=count" json:"count,omitempty"`
	Messages []*QueueMessage `protobuf:"bytes,2,rep,name=messages" json:"messages,omitempty"`
}

func (x *RetrieveMessagesResponse) Reset() {
	*x = RetrieveMessagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_queueing_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetrieveMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetrieveMessagesResponse) ProtoMessage() {}

func (x *RetrieveMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_message_queueing_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetrieveMessagesResponse.ProtoReflect.Descriptor instead.
func (*RetrieveMessagesResponse) Descriptor() ([]byte, []int) {
	return file_message_queueing_proto_rawDescGZIP(), []int{6}
}

func (x *RetrieveMessagesResponse) GetCount() int32 {
	if x != nil && x.Count != nil {
		return *x.Count
	}
	return 0
}

func (x *RetrieveMessagesResponse) GetMessages() []*QueueMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

type AcknowledgeMessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AcknowledgeMessagesRequest) Reset() {
	*x = AcknowledgeMessagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_queueing_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcknowledgeMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcknowledgeMessagesRequest) ProtoMessage() {}

func (x *AcknowledgeMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_queueing_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcknowledgeMessagesRequest.ProtoReflect.Descriptor instead.
func (*AcknowledgeMessagesRequest) Descriptor() ([]byte, []int) {
	return file_message_queueing_proto_rawDescGZIP(), []int{7}
}

type AcknowledgeMessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AcknowledgeMessagesResponse) Reset() {
	*x = AcknowledgeMessagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_queueing_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcknowledgeMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcknowledgeMessagesResponse) ProtoMessage() {}

func (x *AcknowledgeMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_message_queueing_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcknowledgeMessagesResponse.ProtoReflect.Descriptor instead.
func (*AcknowledgeMessagesResponse) Descriptor() ([]byte, []int) {
	return file_message_queueing_proto_rawDescGZIP(), []int{8}
}

var File_message_queueing_proto protoreflect.FileDescriptor

var file_message_queueing_proto_rawDesc = []byte{
	0x0a, 0x16, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2d, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69,
	0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67, 0x22, 0xb7, 0x01, 0x0a, 0x0f, 0x52,
	0x61, 0x77, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x02, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x51, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67, 0x2e, 0x52, 0x61, 0x77, 0x51, 0x75, 0x65,
	0x75, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x1a, 0x3d, 0x0a, 0x0f, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0xaf, 0x02, 0x0a, 0x0c, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x02, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x02, 0x28, 0x0c, 0x52, 0x09, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x02, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73,
	0x68, 0x18, 0x04, 0x20, 0x02, 0x28, 0x0c, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x4e, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18,
	0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x73, 0x12, 0x24, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76,
	0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x52, 0x65,
	0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x64, 0x1a, 0x3d, 0x0a, 0x0f, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x77, 0x0a, 0x0d, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74,
	0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20,
	0x02, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x02, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73, 0x68, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x48, 0x61, 0x73, 0x68, 0x22,
	0x70, 0x0a, 0x15, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x71, 0x75, 0x65, 0x75,
	0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x07, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x49, 0x44, 0x12, 0x3d, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67, 0x2e, 0x52, 0x61, 0x77, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x22, 0x55, 0x0a, 0x16, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x08, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67,
	0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x08,
	0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x22, 0x2f, 0x0a, 0x17, 0x52, 0x65, 0x74, 0x72,
	0x69, 0x65, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x02,
	0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x6c, 0x0a, 0x18, 0x52, 0x65, 0x74,
	0x72, 0x69, 0x65, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x02, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x3a, 0x0a, 0x08, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67,
	0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22, 0x1c, 0x0a, 0x1a, 0x41, 0x63, 0x6b, 0x6e, 0x6f,
	0x77, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x1d, 0x0a, 0x1b, 0x41, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c,
	0x65, 0x64, 0x67, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0xd2, 0x02, 0x0a, 0x0c, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x63, 0x0a, 0x0e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x27, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69,
	0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x28, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x69, 0x6e, 0x67, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x69, 0x0a, 0x10, 0x52, 0x65,
	0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x29,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x6e,
	0x67, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67, 0x2e, 0x52, 0x65, 0x74,
	0x72, 0x69, 0x65, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x72, 0x0a, 0x13, 0x41, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c,
	0x65, 0x64, 0x67, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x2c, 0x2e, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67, 0x2e,
	0x41, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67, 0x2e, 0x41, 0x63,
	0x6b, 0x6e, 0x6f, 0x77, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x69, 0x6e, 0x67,
}

var (
	file_message_queueing_proto_rawDescOnce sync.Once
	file_message_queueing_proto_rawDescData = file_message_queueing_proto_rawDesc
)

func file_message_queueing_proto_rawDescGZIP() []byte {
	file_message_queueing_proto_rawDescOnce.Do(func() {
		file_message_queueing_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_queueing_proto_rawDescData)
	})
	return file_message_queueing_proto_rawDescData
}

var file_message_queueing_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_message_queueing_proto_goTypes = []interface{}{
	(*RawQueueMessage)(nil),             // 0: message_queueing.RawQueueMessage
	(*QueueMessage)(nil),                // 1: message_queueing.QueueMessage
	(*SubmitReceipt)(nil),               // 2: message_queueing.SubmitReceipt
	(*SubmitMessagesRequest)(nil),       // 3: message_queueing.SubmitMessagesRequest
	(*SubmitMessagesResponse)(nil),      // 4: message_queueing.SubmitMessagesResponse
	(*RetrieveMessagesRequest)(nil),     // 5: message_queueing.RetrieveMessagesRequest
	(*RetrieveMessagesResponse)(nil),    // 6: message_queueing.RetrieveMessagesResponse
	(*AcknowledgeMessagesRequest)(nil),  // 7: message_queueing.AcknowledgeMessagesRequest
	(*AcknowledgeMessagesResponse)(nil), // 8: message_queueing.AcknowledgeMessagesResponse
	nil,                                 // 9: message_queueing.RawQueueMessage.AttributesEntry
	nil,                                 // 10: message_queueing.QueueMessage.AttributesEntry
}
var file_message_queueing_proto_depIdxs = []int32{
	9,  // 0: message_queueing.RawQueueMessage.attributes:type_name -> message_queueing.RawQueueMessage.AttributesEntry
	10, // 1: message_queueing.QueueMessage.attributes:type_name -> message_queueing.QueueMessage.AttributesEntry
	0,  // 2: message_queueing.SubmitMessagesRequest.messages:type_name -> message_queueing.RawQueueMessage
	2,  // 3: message_queueing.SubmitMessagesResponse.receipts:type_name -> message_queueing.SubmitReceipt
	1,  // 4: message_queueing.RetrieveMessagesResponse.messages:type_name -> message_queueing.QueueMessage
	3,  // 5: message_queueing.QueueService.SubmitMessages:input_type -> message_queueing.SubmitMessagesRequest
	5,  // 6: message_queueing.QueueService.RetrieveMessages:input_type -> message_queueing.RetrieveMessagesRequest
	7,  // 7: message_queueing.QueueService.AcknowledgeMessages:input_type -> message_queueing.AcknowledgeMessagesRequest
	4,  // 8: message_queueing.QueueService.SubmitMessages:output_type -> message_queueing.SubmitMessagesResponse
	6,  // 9: message_queueing.QueueService.RetrieveMessages:output_type -> message_queueing.RetrieveMessagesResponse
	8,  // 10: message_queueing.QueueService.AcknowledgeMessages:output_type -> message_queueing.AcknowledgeMessagesResponse
	8,  // [8:11] is the sub-list for method output_type
	5,  // [5:8] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_message_queueing_proto_init() }
func file_message_queueing_proto_init() {
	if File_message_queueing_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_queueing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RawQueueMessage); i {
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
		file_message_queueing_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueMessage); i {
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
		file_message_queueing_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitReceipt); i {
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
		file_message_queueing_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitMessagesRequest); i {
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
		file_message_queueing_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitMessagesResponse); i {
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
		file_message_queueing_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetrieveMessagesRequest); i {
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
		file_message_queueing_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetrieveMessagesResponse); i {
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
		file_message_queueing_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcknowledgeMessagesRequest); i {
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
		file_message_queueing_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcknowledgeMessagesResponse); i {
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
			RawDescriptor: file_message_queueing_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_message_queueing_proto_goTypes,
		DependencyIndexes: file_message_queueing_proto_depIdxs,
		MessageInfos:      file_message_queueing_proto_msgTypes,
	}.Build()
	File_message_queueing_proto = out.File
	file_message_queueing_proto_rawDesc = nil
	file_message_queueing_proto_goTypes = nil
	file_message_queueing_proto_depIdxs = nil
}
