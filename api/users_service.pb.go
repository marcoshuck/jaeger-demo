// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.20.0
// source: api/users_service.proto

package api

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

var File_api_users_service_proto protoreflect.FileDescriptor

var file_api_users_service_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x17,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x2f, 0x0a, 0x0c, 0x55, 0x73, 0x65,
	0x72, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x03, 0x47, 0x65, 0x74,
	0x12, 0x0b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x09, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x22, 0x00, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x72, 0x63, 0x6f, 0x73, 0x68,
	0x75, 0x63, 0x6b, 0x2f, 0x6a, 0x61, 0x65, 0x67, 0x65, 0x72, 0x2d, 0x64, 0x65, 0x6d, 0x6f, 0x2f,
	0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_api_users_service_proto_goTypes = []interface{}{
	(*UserID)(nil), // 0: api.UserID
	(*User)(nil),   // 1: api.User
}
var file_api_users_service_proto_depIdxs = []int32{
	0, // 0: api.UsersService.Get:input_type -> api.UserID
	1, // 1: api.UsersService.Get:output_type -> api.User
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_users_service_proto_init() }
func file_api_users_service_proto_init() {
	if File_api_users_service_proto != nil {
		return
	}
	file_api_value_objects_proto_init()
	file_api_domain_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_users_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_users_service_proto_goTypes,
		DependencyIndexes: file_api_users_service_proto_depIdxs,
	}.Build()
	File_api_users_service_proto = out.File
	file_api_users_service_proto_rawDesc = nil
	file_api_users_service_proto_goTypes = nil
	file_api_users_service_proto_depIdxs = nil
}
