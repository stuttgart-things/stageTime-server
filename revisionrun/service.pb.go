// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: revisionrun/service.proto

package stageTime_server

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_revisionrun_service_proto protoreflect.FileDescriptor

var file_revisionrun_service_proto_rawDesc = []byte{
	0x0a, 0x19, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x75, 0x6e, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x72, 0x65, 0x76,
	0x69, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x75, 0x6e, 0x1a, 0x1a, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x72, 0x75, 0x6e, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x75,
	0x6e, 0x2f, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x75, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x32, 0x72, 0x0a, 0x1b, 0x53, 0x74, 0x61, 0x67, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x53, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x76, 0x69,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x75, 0x6e, 0x12, 0x25, 0x2e, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x72, 0x75, 0x6e, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x76, 0x69,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x75, 0x6e, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x75, 0x74, 0x74, 0x67, 0x61, 0x72, 0x74, 0x2d,
	0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x67, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_revisionrun_service_proto_goTypes = []interface{}{
	(*CreateRevisionRunRequest)(nil), // 0: revisionrun.CreateRevisionRunRequest
	(*Response)(nil),                 // 1: revisionrun.Response
}
var file_revisionrun_service_proto_depIdxs = []int32{
	0, // 0: revisionrun.StageTimeApplicationService.CreateRevisionRun:input_type -> revisionrun.CreateRevisionRunRequest
	1, // 1: revisionrun.StageTimeApplicationService.CreateRevisionRun:output_type -> revisionrun.Response
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_revisionrun_service_proto_init() }
func file_revisionrun_service_proto_init() {
	if File_revisionrun_service_proto != nil {
		return
	}
	file_revisionrun_response_proto_init()
	file_revisionrun_revisionrun_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_revisionrun_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_revisionrun_service_proto_goTypes,
		DependencyIndexes: file_revisionrun_service_proto_depIdxs,
	}.Build()
	File_revisionrun_service_proto = out.File
	file_revisionrun_service_proto_rawDesc = nil
	file_revisionrun_service_proto_goTypes = nil
	file_revisionrun_service_proto_depIdxs = nil
}
