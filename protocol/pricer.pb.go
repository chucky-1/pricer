// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: protocol/pricer.proto

package protocol

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

type Stock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title string  `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Price float32 `protobuf:"fixed32,3,opt,name=price,proto3" json:"price,omitempty"`
}

func (x *Stock) Reset() {
	*x = Stock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_pricer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stock) ProtoMessage() {}

func (x *Stock) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_pricer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stock.ProtoReflect.Descriptor instead.
func (*Stock) Descriptor() ([]byte, []int) {
	return file_protocol_pricer_proto_rawDescGZIP(), []int{0}
}

func (x *Stock) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Stock) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Stock) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

type ListID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id []int32 `protobuf:"varint,1,rep,packed,name=id,proto3" json:"id,omitempty"`
}

func (x *ListID) Reset() {
	*x = ListID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_pricer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListID) ProtoMessage() {}

func (x *ListID) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_pricer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListID.ProtoReflect.Descriptor instead.
func (*ListID) Descriptor() ([]byte, []int) {
	return file_protocol_pricer_proto_rawDescGZIP(), []int{1}
}

func (x *ListID) GetId() []int32 {
	if x != nil {
		return x.Id
	}
	return nil
}

type GrpcID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GrpcID) Reset() {
	*x = GrpcID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_pricer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrpcID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcID) ProtoMessage() {}

func (x *GrpcID) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_pricer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcID.ProtoReflect.Descriptor instead.
func (*GrpcID) Descriptor() ([]byte, []int) {
	return file_protocol_pricer_proto_rawDescGZIP(), []int{2}
}

func (x *GrpcID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_pricer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_pricer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_protocol_pricer_proto_rawDescGZIP(), []int{3}
}

var File_protocol_pricer_proto protoreflect.FileDescriptor

var file_protocol_pricer_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x67, 0x72, 0x70, 0x63, 0x22, 0x43,
	0x0a, 0x05, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x70, 0x72,
	0x69, 0x63, 0x65, 0x22, 0x18, 0x0a, 0x06, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x44, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x22, 0x18, 0x0a,
	0x06, 0x47, 0x72, 0x70, 0x63, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x0a, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0x5c, 0x0a, 0x06, 0x50, 0x72, 0x69, 0x63, 0x65, 0x72, 0x12, 0x27, 0x0a,
	0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x0d, 0x2e, 0x70, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x49, 0x44, 0x1a, 0x0c, 0x2e, 0x70, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x6f,
	0x63, 0x6b, 0x22, 0x00, 0x30, 0x01, 0x12, 0x29, 0x0a, 0x05, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x12,
	0x0d, 0x2e, 0x70, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x49, 0x44, 0x1a, 0x0f,
	0x2e, 0x70, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x63, 0x68, 0x75, 0x63, 0x6b, 0x79, 0x2d, 0x31, 0x2f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x72, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protocol_pricer_proto_rawDescOnce sync.Once
	file_protocol_pricer_proto_rawDescData = file_protocol_pricer_proto_rawDesc
)

func file_protocol_pricer_proto_rawDescGZIP() []byte {
	file_protocol_pricer_proto_rawDescOnce.Do(func() {
		file_protocol_pricer_proto_rawDescData = protoimpl.X.CompressGZIP(file_protocol_pricer_proto_rawDescData)
	})
	return file_protocol_pricer_proto_rawDescData
}

var file_protocol_pricer_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protocol_pricer_proto_goTypes = []interface{}{
	(*Stock)(nil),    // 0: pgrpc.Stock
	(*ListID)(nil),   // 1: pgrpc.ListID
	(*GrpcID)(nil),   // 2: pgrpc.GrpcID
	(*Response)(nil), // 3: pgrpc.Response
}
var file_protocol_pricer_proto_depIdxs = []int32{
	1, // 0: pgrpc.Pricer.Send:input_type -> pgrpc.ListID
	2, // 1: pgrpc.Pricer.Close:input_type -> pgrpc.GrpcID
	0, // 2: pgrpc.Pricer.Send:output_type -> pgrpc.Stock
	3, // 3: pgrpc.Pricer.Close:output_type -> pgrpc.Response
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protocol_pricer_proto_init() }
func file_protocol_pricer_proto_init() {
	if File_protocol_pricer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protocol_pricer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stock); i {
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
		file_protocol_pricer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListID); i {
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
		file_protocol_pricer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrpcID); i {
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
		file_protocol_pricer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
			RawDescriptor: file_protocol_pricer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protocol_pricer_proto_goTypes,
		DependencyIndexes: file_protocol_pricer_proto_depIdxs,
		MessageInfos:      file_protocol_pricer_proto_msgTypes,
	}.Build()
	File_protocol_pricer_proto = out.File
	file_protocol_pricer_proto_rawDesc = nil
	file_protocol_pricer_proto_goTypes = nil
	file_protocol_pricer_proto_depIdxs = nil
}
