// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: revisionrun/pipelinerun.proto

package stageTime_server

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

type Pipelinerun struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name                 string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Canfail              bool    `protobuf:"varint,2,opt,name=canfail,proto3" json:"canfail,omitempty"`
	Stage                float64 `protobuf:"fixed64,3,opt,name=stage,proto3" json:"stage,omitempty"`
	Params               string  `protobuf:"bytes,4,opt,name=params,proto3" json:"params,omitempty"`
	ResolverParams       string  `protobuf:"bytes,5,opt,name=resolverParams,proto3" json:"resolverParams,omitempty"`
	Listparams           string  `protobuf:"bytes,6,opt,name=listparams,proto3" json:"listparams,omitempty"`
	Workspaces           string  `protobuf:"bytes,7,opt,name=workspaces,proto3" json:"workspaces,omitempty"`
	VolumeClaimTemplates string  `protobuf:"bytes,8,opt,name=volumeClaimTemplates,proto3" json:"volumeClaimTemplates,omitempty"`
}

func (x *Pipelinerun) Reset() {
	*x = Pipelinerun{}
	if protoimpl.UnsafeEnabled {
		mi := &file_revisionrun_pipelinerun_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pipelinerun) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pipelinerun) ProtoMessage() {}

func (x *Pipelinerun) ProtoReflect() protoreflect.Message {
	mi := &file_revisionrun_pipelinerun_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pipelinerun.ProtoReflect.Descriptor instead.
func (*Pipelinerun) Descriptor() ([]byte, []int) {
	return file_revisionrun_pipelinerun_proto_rawDescGZIP(), []int{0}
}

func (x *Pipelinerun) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Pipelinerun) GetCanfail() bool {
	if x != nil {
		return x.Canfail
	}
	return false
}

func (x *Pipelinerun) GetStage() float64 {
	if x != nil {
		return x.Stage
	}
	return 0
}

func (x *Pipelinerun) GetParams() string {
	if x != nil {
		return x.Params
	}
	return ""
}

func (x *Pipelinerun) GetResolverParams() string {
	if x != nil {
		return x.ResolverParams
	}
	return ""
}

func (x *Pipelinerun) GetListparams() string {
	if x != nil {
		return x.Listparams
	}
	return ""
}

func (x *Pipelinerun) GetWorkspaces() string {
	if x != nil {
		return x.Workspaces
	}
	return ""
}

func (x *Pipelinerun) GetVolumeClaimTemplates() string {
	if x != nil {
		return x.VolumeClaimTemplates
	}
	return ""
}

var File_revisionrun_pipelinerun_proto protoreflect.FileDescriptor

var file_revisionrun_pipelinerun_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x75, 0x6e, 0x2f, 0x70, 0x69,
	0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x72, 0x75, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0b, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x75, 0x6e, 0x22, 0x85, 0x02, 0x0a,
	0x0b, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x72, 0x75, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x61, 0x6e, 0x66, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x63, 0x61, 0x6e, 0x66, 0x61, 0x69, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x26, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x6f,
	0x6c, 0x76, 0x65, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x69, 0x73, 0x74, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6c, 0x69, 0x73, 0x74, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x12, 0x1e, 0x0a, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73,
	0x12, 0x32, 0x0a, 0x14, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x54,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14,
	0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x54, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x73, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x75, 0x74, 0x74, 0x67, 0x61, 0x72, 0x74, 0x2d, 0x74, 0x68, 0x69,
	0x6e, 0x67, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x67, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x2d, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_revisionrun_pipelinerun_proto_rawDescOnce sync.Once
	file_revisionrun_pipelinerun_proto_rawDescData = file_revisionrun_pipelinerun_proto_rawDesc
)

func file_revisionrun_pipelinerun_proto_rawDescGZIP() []byte {
	file_revisionrun_pipelinerun_proto_rawDescOnce.Do(func() {
		file_revisionrun_pipelinerun_proto_rawDescData = protoimpl.X.CompressGZIP(file_revisionrun_pipelinerun_proto_rawDescData)
	})
	return file_revisionrun_pipelinerun_proto_rawDescData
}

var file_revisionrun_pipelinerun_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_revisionrun_pipelinerun_proto_goTypes = []interface{}{
	(*Pipelinerun)(nil), // 0: revisionrun.Pipelinerun
}
var file_revisionrun_pipelinerun_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_revisionrun_pipelinerun_proto_init() }
func file_revisionrun_pipelinerun_proto_init() {
	if File_revisionrun_pipelinerun_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_revisionrun_pipelinerun_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pipelinerun); i {
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
			RawDescriptor: file_revisionrun_pipelinerun_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_revisionrun_pipelinerun_proto_goTypes,
		DependencyIndexes: file_revisionrun_pipelinerun_proto_depIdxs,
		MessageInfos:      file_revisionrun_pipelinerun_proto_msgTypes,
	}.Build()
	File_revisionrun_pipelinerun_proto = out.File
	file_revisionrun_pipelinerun_proto_rawDesc = nil
	file_revisionrun_pipelinerun_proto_goTypes = nil
	file_revisionrun_pipelinerun_proto_depIdxs = nil
}
