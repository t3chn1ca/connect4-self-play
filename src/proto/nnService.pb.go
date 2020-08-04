// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0-devel
// 	protoc        v3.7.1
// source: src/proto/nnService.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type BoardRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Board []int32 `protobuf:"varint,1,rep,packed,name=board,proto3" json:"board,omitempty"`
}

func (x *BoardRequest) Reset() {
	*x = BoardRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_proto_nnService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoardRequest) ProtoMessage() {}

func (x *BoardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_src_proto_nnService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoardRequest.ProtoReflect.Descriptor instead.
func (*BoardRequest) Descriptor() ([]byte, []int) {
	return file_src_proto_nnService_proto_rawDescGZIP(), []int{0}
}

func (x *BoardRequest) GetBoard() []int32 {
	if x != nil {
		return x.Board
	}
	return nil
}

type TrainFromIndex struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UidFrom int32 `protobuf:"varint,1,opt,name=uidFrom,proto3" json:"uidFrom,omitempty"`
	UidTo   int32 `protobuf:"varint,2,opt,name=uidTo,proto3" json:"uidTo,omitempty"`
}

func (x *TrainFromIndex) Reset() {
	*x = TrainFromIndex{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_proto_nnService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrainFromIndex) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrainFromIndex) ProtoMessage() {}

func (x *TrainFromIndex) ProtoReflect() protoreflect.Message {
	mi := &file_src_proto_nnService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrainFromIndex.ProtoReflect.Descriptor instead.
func (*TrainFromIndex) Descriptor() ([]byte, []int) {
	return file_src_proto_nnService_proto_rawDescGZIP(), []int{1}
}

func (x *TrainFromIndex) GetUidFrom() int32 {
	if x != nil {
		return x.UidFrom
	}
	return 0
}

func (x *TrainFromIndex) GetUidTo() int32 {
	if x != nil {
		return x.UidTo
	}
	return 0
}

type NNResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result []float64 `protobuf:"fixed64,1,rep,packed,name=result,proto3" json:"result,omitempty"`
}

func (x *NNResponse) Reset() {
	*x = NNResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_proto_nnService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NNResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NNResponse) ProtoMessage() {}

func (x *NNResponse) ProtoReflect() protoreflect.Message {
	mi := &file_src_proto_nnService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NNResponse.ProtoReflect.Descriptor instead.
func (*NNResponse) Descriptor() ([]byte, []int) {
	return file_src_proto_nnService_proto_rawDescGZIP(), []int{2}
}

func (x *NNResponse) GetResult() []float64 {
	if x != nil {
		return x.Result
	}
	return nil
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status int32 `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_proto_nnService_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_src_proto_nnService_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_src_proto_nnService_proto_rawDescGZIP(), []int{3}
}

func (x *Status) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type NoArg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoArg) Reset() {
	*x = NoArg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_src_proto_nnService_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoArg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoArg) ProtoMessage() {}

func (x *NoArg) ProtoReflect() protoreflect.Message {
	mi := &file_src_proto_nnService_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoArg.ProtoReflect.Descriptor instead.
func (*NoArg) Descriptor() ([]byte, []int) {
	return file_src_proto_nnService_proto_rawDescGZIP(), []int{4}
}

var File_src_proto_nnService_proto protoreflect.FileDescriptor

var file_src_proto_nnService_proto_rawDesc = []byte{
	0x0a, 0x19, 0x73, 0x72, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x24, 0x0a, 0x0c, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x05, 0x52, 0x05, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x22, 0x40, 0x0a, 0x0e, 0x54, 0x72, 0x61, 0x69,
	0x6e, 0x46, 0x72, 0x6f, 0x6d, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x75, 0x69,
	0x64, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x75, 0x69, 0x64,
	0x46, 0x72, 0x6f, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x69, 0x64, 0x54, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x75, 0x69, 0x64, 0x54, 0x6f, 0x22, 0x24, 0x0a, 0x0a, 0x4e, 0x4e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x01, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x22, 0x20, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x07, 0x0a, 0x05, 0x4e, 0x6f, 0x41, 0x72, 0x67, 0x32, 0xc0, 0x02, 0x0a, 0x0a,
	0x41, 0x64, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x0b, 0x46, 0x6f,
	0x72, 0x77, 0x61, 0x72, 0x64, 0x50, 0x61, 0x73, 0x73, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x4e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2d, 0x0a, 0x05, 0x54, 0x72, 0x61, 0x69, 0x6e, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x69, 0x6e, 0x46, 0x72, 0x6f, 0x6d, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x31, 0x0a, 0x12, 0x4c, 0x6f, 0x61, 0x64, 0x42, 0x65, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x65,
	0x6c, 0x54, 0x6f, 0x43, 0x70, 0x75, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e,
	0x6f, 0x41, 0x72, 0x67, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x31, 0x0a, 0x12, 0x4c, 0x6f, 0x61, 0x64, 0x42, 0x65, 0x73, 0x74, 0x4d,
	0x6f, 0x64, 0x65, 0x6c, 0x54, 0x6f, 0x47, 0x70, 0x75, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4e, 0x6f, 0x41, 0x72, 0x67, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x10, 0x53, 0x74, 0x6f, 0x70, 0x42, 0x65,
	0x73, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x43, 0x70, 0x75, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x41, 0x72, 0x67, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x35, 0x0a, 0x16, 0x53, 0x61, 0x76, 0x65, 0x43,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x54, 0x6f, 0x42, 0x65, 0x73,
	0x74, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x41, 0x72, 0x67, 0x1a,
	0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_src_proto_nnService_proto_rawDescOnce sync.Once
	file_src_proto_nnService_proto_rawDescData = file_src_proto_nnService_proto_rawDesc
)

func file_src_proto_nnService_proto_rawDescGZIP() []byte {
	file_src_proto_nnService_proto_rawDescOnce.Do(func() {
		file_src_proto_nnService_proto_rawDescData = protoimpl.X.CompressGZIP(file_src_proto_nnService_proto_rawDescData)
	})
	return file_src_proto_nnService_proto_rawDescData
}

var file_src_proto_nnService_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_src_proto_nnService_proto_goTypes = []interface{}{
	(*BoardRequest)(nil),   // 0: proto.BoardRequest
	(*TrainFromIndex)(nil), // 1: proto.TrainFromIndex
	(*NNResponse)(nil),     // 2: proto.NNResponse
	(*Status)(nil),         // 3: proto.Status
	(*NoArg)(nil),          // 4: proto.NoArg
}
var file_src_proto_nnService_proto_depIdxs = []int32{
	0, // 0: proto.AddService.ForwardPass:input_type -> proto.BoardRequest
	1, // 1: proto.AddService.Train:input_type -> proto.TrainFromIndex
	4, // 2: proto.AddService.LoadBestModelToCpu:input_type -> proto.NoArg
	4, // 3: proto.AddService.LoadBestModelToGpu:input_type -> proto.NoArg
	4, // 4: proto.AddService.StopBestModelCpu:input_type -> proto.NoArg
	4, // 5: proto.AddService.SaveCurrentModelToBest:input_type -> proto.NoArg
	2, // 6: proto.AddService.ForwardPass:output_type -> proto.NNResponse
	3, // 7: proto.AddService.Train:output_type -> proto.Status
	3, // 8: proto.AddService.LoadBestModelToCpu:output_type -> proto.Status
	3, // 9: proto.AddService.LoadBestModelToGpu:output_type -> proto.Status
	3, // 10: proto.AddService.StopBestModelCpu:output_type -> proto.Status
	3, // 11: proto.AddService.SaveCurrentModelToBest:output_type -> proto.Status
	6, // [6:12] is the sub-list for method output_type
	0, // [0:6] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_src_proto_nnService_proto_init() }
func file_src_proto_nnService_proto_init() {
	if File_src_proto_nnService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_src_proto_nnService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoardRequest); i {
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
		file_src_proto_nnService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrainFromIndex); i {
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
		file_src_proto_nnService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NNResponse); i {
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
		file_src_proto_nnService_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
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
		file_src_proto_nnService_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoArg); i {
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
			RawDescriptor: file_src_proto_nnService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_src_proto_nnService_proto_goTypes,
		DependencyIndexes: file_src_proto_nnService_proto_depIdxs,
		MessageInfos:      file_src_proto_nnService_proto_msgTypes,
	}.Build()
	File_src_proto_nnService_proto = out.File
	file_src_proto_nnService_proto_rawDesc = nil
	file_src_proto_nnService_proto_goTypes = nil
	file_src_proto_nnService_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AddServiceClient is the client API for AddService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AddServiceClient interface {
	ForwardPass(ctx context.Context, in *BoardRequest, opts ...grpc.CallOption) (*NNResponse, error)
	Train(ctx context.Context, in *TrainFromIndex, opts ...grpc.CallOption) (*Status, error)
	LoadBestModelToCpu(ctx context.Context, in *NoArg, opts ...grpc.CallOption) (*Status, error)
	LoadBestModelToGpu(ctx context.Context, in *NoArg, opts ...grpc.CallOption) (*Status, error)
	StopBestModelCpu(ctx context.Context, in *NoArg, opts ...grpc.CallOption) (*Status, error)
	SaveCurrentModelToBest(ctx context.Context, in *NoArg, opts ...grpc.CallOption) (*Status, error)
}

type addServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAddServiceClient(cc grpc.ClientConnInterface) AddServiceClient {
	return &addServiceClient{cc}
}

func (c *addServiceClient) ForwardPass(ctx context.Context, in *BoardRequest, opts ...grpc.CallOption) (*NNResponse, error) {
	out := new(NNResponse)
	err := c.cc.Invoke(ctx, "/proto.AddService/ForwardPass", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addServiceClient) Train(ctx context.Context, in *TrainFromIndex, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/proto.AddService/Train", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addServiceClient) LoadBestModelToCpu(ctx context.Context, in *NoArg, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/proto.AddService/LoadBestModelToCpu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addServiceClient) LoadBestModelToGpu(ctx context.Context, in *NoArg, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/proto.AddService/LoadBestModelToGpu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addServiceClient) StopBestModelCpu(ctx context.Context, in *NoArg, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/proto.AddService/StopBestModelCpu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addServiceClient) SaveCurrentModelToBest(ctx context.Context, in *NoArg, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/proto.AddService/SaveCurrentModelToBest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AddServiceServer is the server API for AddService service.
type AddServiceServer interface {
	ForwardPass(context.Context, *BoardRequest) (*NNResponse, error)
	Train(context.Context, *TrainFromIndex) (*Status, error)
	LoadBestModelToCpu(context.Context, *NoArg) (*Status, error)
	LoadBestModelToGpu(context.Context, *NoArg) (*Status, error)
	StopBestModelCpu(context.Context, *NoArg) (*Status, error)
	SaveCurrentModelToBest(context.Context, *NoArg) (*Status, error)
}

// UnimplementedAddServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAddServiceServer struct {
}

func (*UnimplementedAddServiceServer) ForwardPass(context.Context, *BoardRequest) (*NNResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForwardPass not implemented")
}
func (*UnimplementedAddServiceServer) Train(context.Context, *TrainFromIndex) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Train not implemented")
}
func (*UnimplementedAddServiceServer) LoadBestModelToCpu(context.Context, *NoArg) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadBestModelToCpu not implemented")
}
func (*UnimplementedAddServiceServer) LoadBestModelToGpu(context.Context, *NoArg) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadBestModelToGpu not implemented")
}
func (*UnimplementedAddServiceServer) StopBestModelCpu(context.Context, *NoArg) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopBestModelCpu not implemented")
}
func (*UnimplementedAddServiceServer) SaveCurrentModelToBest(context.Context, *NoArg) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveCurrentModelToBest not implemented")
}

func RegisterAddServiceServer(s *grpc.Server, srv AddServiceServer) {
	s.RegisterService(&_AddService_serviceDesc, srv)
}

func _AddService_ForwardPass_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BoardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServiceServer).ForwardPass(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AddService/ForwardPass",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServiceServer).ForwardPass(ctx, req.(*BoardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddService_Train_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrainFromIndex)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServiceServer).Train(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AddService/Train",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServiceServer).Train(ctx, req.(*TrainFromIndex))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddService_LoadBestModelToCpu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServiceServer).LoadBestModelToCpu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AddService/LoadBestModelToCpu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServiceServer).LoadBestModelToCpu(ctx, req.(*NoArg))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddService_LoadBestModelToGpu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServiceServer).LoadBestModelToGpu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AddService/LoadBestModelToGpu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServiceServer).LoadBestModelToGpu(ctx, req.(*NoArg))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddService_StopBestModelCpu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServiceServer).StopBestModelCpu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AddService/StopBestModelCpu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServiceServer).StopBestModelCpu(ctx, req.(*NoArg))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddService_SaveCurrentModelToBest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoArg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServiceServer).SaveCurrentModelToBest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AddService/SaveCurrentModelToBest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServiceServer).SaveCurrentModelToBest(ctx, req.(*NoArg))
	}
	return interceptor(ctx, in, info, handler)
}

var _AddService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AddService",
	HandlerType: (*AddServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ForwardPass",
			Handler:    _AddService_ForwardPass_Handler,
		},
		{
			MethodName: "Train",
			Handler:    _AddService_Train_Handler,
		},
		{
			MethodName: "LoadBestModelToCpu",
			Handler:    _AddService_LoadBestModelToCpu_Handler,
		},
		{
			MethodName: "LoadBestModelToGpu",
			Handler:    _AddService_LoadBestModelToGpu_Handler,
		},
		{
			MethodName: "StopBestModelCpu",
			Handler:    _AddService_StopBestModelCpu_Handler,
		},
		{
			MethodName: "SaveCurrentModelToBest",
			Handler:    _AddService_SaveCurrentModelToBest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/proto/nnService.proto",
}
