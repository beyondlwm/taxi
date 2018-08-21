// Code generated by protoc-gen-go. DO NOT EDIT.
// source: descriptor.proto

package descriptor

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// available type names:
//  nil,bool,int8,uint8,int16,uint16,int32,uint32,int64,uint64,float32,float64,string,bytes,datetime
type TypeEnum int32

const (
	TypeEnum_Unknown  TypeEnum = 0
	TypeEnum_Nil      TypeEnum = 1
	TypeEnum_Bool     TypeEnum = 2
	TypeEnum_Int8     TypeEnum = 3
	TypeEnum_Uint8    TypeEnum = 4
	TypeEnum_Int16    TypeEnum = 5
	TypeEnum_Uint16   TypeEnum = 6
	TypeEnum_Int      TypeEnum = 7
	TypeEnum_Int32    TypeEnum = 8
	TypeEnum_Uint32   TypeEnum = 9
	TypeEnum_Int64    TypeEnum = 10
	TypeEnum_Uint64   TypeEnum = 11
	TypeEnum_Float    TypeEnum = 12
	TypeEnum_Float32  TypeEnum = 13
	TypeEnum_Float64  TypeEnum = 14
	TypeEnum_String   TypeEnum = 15
	TypeEnum_Enum     TypeEnum = 16
	TypeEnum_Bytes    TypeEnum = 17
	TypeEnum_DateTime TypeEnum = 18
	TypeEnum_Json     TypeEnum = 19
	TypeEnum_Array    TypeEnum = 20
	TypeEnum_Any      TypeEnum = 21
)

var TypeEnum_name = map[int32]string{
	0:  "Unknown",
	1:  "Nil",
	2:  "Bool",
	3:  "Int8",
	4:  "Uint8",
	5:  "Int16",
	6:  "Uint16",
	7:  "Int",
	8:  "Int32",
	9:  "Uint32",
	10: "Int64",
	11: "Uint64",
	12: "Float",
	13: "Float32",
	14: "Float64",
	15: "String",
	16: "Enum",
	17: "Bytes",
	18: "DateTime",
	19: "Json",
	20: "Array",
	21: "Any",
}
var TypeEnum_value = map[string]int32{
	"Unknown":  0,
	"Nil":      1,
	"Bool":     2,
	"Int8":     3,
	"Uint8":    4,
	"Int16":    5,
	"Uint16":   6,
	"Int":      7,
	"Int32":    8,
	"Uint32":   9,
	"Int64":    10,
	"Uint64":   11,
	"Float":    12,
	"Float32":  13,
	"Float64":  14,
	"String":   15,
	"Enum":     16,
	"Bytes":    17,
	"DateTime": 18,
	"Json":     19,
	"Array":    20,
	"Any":      21,
}

func (x TypeEnum) String() string {
	return proto.EnumName(TypeEnum_name, int32(x))
}
func (TypeEnum) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_18a642691b671acd, []int{0}
}

type StructDescriptor struct {
	Name                 string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	CamelCaseName        string             `protobuf:"bytes,2,opt,name=camel_case_name,json=camelCaseName,proto3" json:"camel_case_name,omitempty"`
	Comment              string             `protobuf:"bytes,3,opt,name=comment,proto3" json:"comment,omitempty"`
	Fields               []*FieldDescriptor `protobuf:"bytes,4,rep,name=fields,proto3" json:"fields,omitempty"`
	Options              map[string]string  `protobuf:"bytes,5,rep,name=options,proto3" json:"options,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *StructDescriptor) Reset()         { *m = StructDescriptor{} }
func (m *StructDescriptor) String() string { return proto.CompactTextString(m) }
func (*StructDescriptor) ProtoMessage()    {}
func (*StructDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_18a642691b671acd, []int{0}
}
func (m *StructDescriptor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StructDescriptor.Unmarshal(m, b)
}
func (m *StructDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StructDescriptor.Marshal(b, m, deterministic)
}
func (dst *StructDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StructDescriptor.Merge(dst, src)
}
func (m *StructDescriptor) XXX_Size() int {
	return xxx_messageInfo_StructDescriptor.Size(m)
}
func (m *StructDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_StructDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_StructDescriptor proto.InternalMessageInfo

func (m *StructDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *StructDescriptor) GetCamelCaseName() string {
	if m != nil {
		return m.CamelCaseName
	}
	return ""
}

func (m *StructDescriptor) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *StructDescriptor) GetFields() []*FieldDescriptor {
	if m != nil {
		return m.Fields
	}
	return nil
}

func (m *StructDescriptor) GetOptions() map[string]string {
	if m != nil {
		return m.Options
	}
	return nil
}

type FieldDescriptor struct {
	Name                 string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	CamelCaseName        string            `protobuf:"bytes,2,opt,name=camel_case_name,json=camelCaseName,proto3" json:"camel_case_name,omitempty"`
	Type                 TypeEnum          `protobuf:"varint,3,opt,name=type,proto3,enum=descriptor.TypeEnum" json:"type,omitempty"`
	TypeName             string            `protobuf:"bytes,4,opt,name=type_name,json=typeName,proto3" json:"type_name,omitempty"`
	OriginalTypeName     string            `protobuf:"bytes,5,opt,name=original_type_name,json=originalTypeName,proto3" json:"original_type_name,omitempty"`
	Comment              string            `protobuf:"bytes,6,opt,name=comment,proto3" json:"comment,omitempty"`
	ColumnIndex          uint32            `protobuf:"varint,7,opt,name=column_index,json=columnIndex,proto3" json:"column_index,omitempty"`
	IsVector             bool              `protobuf:"varint,8,opt,name=is_vector,json=isVector,proto3" json:"is_vector,omitempty"`
	Options              map[string]string `protobuf:"bytes,10,rep,name=options,proto3" json:"options,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *FieldDescriptor) Reset()         { *m = FieldDescriptor{} }
func (m *FieldDescriptor) String() string { return proto.CompactTextString(m) }
func (*FieldDescriptor) ProtoMessage()    {}
func (*FieldDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_18a642691b671acd, []int{1}
}
func (m *FieldDescriptor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FieldDescriptor.Unmarshal(m, b)
}
func (m *FieldDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FieldDescriptor.Marshal(b, m, deterministic)
}
func (dst *FieldDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FieldDescriptor.Merge(dst, src)
}
func (m *FieldDescriptor) XXX_Size() int {
	return xxx_messageInfo_FieldDescriptor.Size(m)
}
func (m *FieldDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_FieldDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_FieldDescriptor proto.InternalMessageInfo

func (m *FieldDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *FieldDescriptor) GetCamelCaseName() string {
	if m != nil {
		return m.CamelCaseName
	}
	return ""
}

func (m *FieldDescriptor) GetType() TypeEnum {
	if m != nil {
		return m.Type
	}
	return TypeEnum_Unknown
}

func (m *FieldDescriptor) GetTypeName() string {
	if m != nil {
		return m.TypeName
	}
	return ""
}

func (m *FieldDescriptor) GetOriginalTypeName() string {
	if m != nil {
		return m.OriginalTypeName
	}
	return ""
}

func (m *FieldDescriptor) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *FieldDescriptor) GetColumnIndex() uint32 {
	if m != nil {
		return m.ColumnIndex
	}
	return 0
}

func (m *FieldDescriptor) GetIsVector() bool {
	if m != nil {
		return m.IsVector
	}
	return false
}

func (m *FieldDescriptor) GetOptions() map[string]string {
	if m != nil {
		return m.Options
	}
	return nil
}

type ImportResult struct {
	Version              string              `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Comment              string              `protobuf:"bytes,2,opt,name=comment,proto3" json:"comment,omitempty"`
	Timestamp            string              `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Options              map[string]string   `protobuf:"bytes,4,rep,name=options,proto3" json:"options,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Descriptors          []*StructDescriptor `protobuf:"bytes,5,rep,name=descriptors,proto3" json:"descriptors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *ImportResult) Reset()         { *m = ImportResult{} }
func (m *ImportResult) String() string { return proto.CompactTextString(m) }
func (*ImportResult) ProtoMessage()    {}
func (*ImportResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_18a642691b671acd, []int{2}
}
func (m *ImportResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImportResult.Unmarshal(m, b)
}
func (m *ImportResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImportResult.Marshal(b, m, deterministic)
}
func (dst *ImportResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImportResult.Merge(dst, src)
}
func (m *ImportResult) XXX_Size() int {
	return xxx_messageInfo_ImportResult.Size(m)
}
func (m *ImportResult) XXX_DiscardUnknown() {
	xxx_messageInfo_ImportResult.DiscardUnknown(m)
}

var xxx_messageInfo_ImportResult proto.InternalMessageInfo

func (m *ImportResult) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ImportResult) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *ImportResult) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func (m *ImportResult) GetOptions() map[string]string {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *ImportResult) GetDescriptors() []*StructDescriptor {
	if m != nil {
		return m.Descriptors
	}
	return nil
}

type ExportRequest struct {
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Format               string   `protobuf:"bytes,3,opt,name=format,proto3" json:"format,omitempty"`
	Filepath             string   `protobuf:"bytes,4,opt,name=filepath,proto3" json:"filepath,omitempty"`
	Datafile             string   `protobuf:"bytes,5,opt,name=datafile,proto3" json:"datafile,omitempty"`
	Args                 string   `protobuf:"bytes,6,opt,name=args,proto3" json:"args,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExportRequest) Reset()         { *m = ExportRequest{} }
func (m *ExportRequest) String() string { return proto.CompactTextString(m) }
func (*ExportRequest) ProtoMessage()    {}
func (*ExportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_18a642691b671acd, []int{3}
}
func (m *ExportRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExportRequest.Unmarshal(m, b)
}
func (m *ExportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExportRequest.Marshal(b, m, deterministic)
}
func (dst *ExportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExportRequest.Merge(dst, src)
}
func (m *ExportRequest) XXX_Size() int {
	return xxx_messageInfo_ExportRequest.Size(m)
}
func (m *ExportRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExportRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExportRequest proto.InternalMessageInfo

func (m *ExportRequest) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ExportRequest) GetFormat() string {
	if m != nil {
		return m.Format
	}
	return ""
}

func (m *ExportRequest) GetFilepath() string {
	if m != nil {
		return m.Filepath
	}
	return ""
}

func (m *ExportRequest) GetDatafile() string {
	if m != nil {
		return m.Datafile
	}
	return ""
}

func (m *ExportRequest) GetArgs() string {
	if m != nil {
		return m.Args
	}
	return ""
}

func init() {
	proto.RegisterType((*StructDescriptor)(nil), "descriptor.StructDescriptor")
	proto.RegisterMapType((map[string]string)(nil), "descriptor.StructDescriptor.OptionsEntry")
	proto.RegisterType((*FieldDescriptor)(nil), "descriptor.FieldDescriptor")
	proto.RegisterMapType((map[string]string)(nil), "descriptor.FieldDescriptor.OptionsEntry")
	proto.RegisterType((*ImportResult)(nil), "descriptor.ImportResult")
	proto.RegisterMapType((map[string]string)(nil), "descriptor.ImportResult.OptionsEntry")
	proto.RegisterType((*ExportRequest)(nil), "descriptor.ExportRequest")
	proto.RegisterEnum("descriptor.TypeEnum", TypeEnum_name, TypeEnum_value)
}

func init() { proto.RegisterFile("descriptor.proto", fileDescriptor_descriptor_18a642691b671acd) }

var fileDescriptor_descriptor_18a642691b671acd = []byte{
	// 638 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x5d, 0x4f, 0x13, 0x41,
	0x14, 0x75, 0xfb, 0xb9, 0xbd, 0x6d, 0xe1, 0x3a, 0xa2, 0x69, 0x80, 0x87, 0x4a, 0xa2, 0xa9, 0xc6,
	0x90, 0xd0, 0x92, 0x86, 0xf0, 0xa0, 0xe1, 0x33, 0xa9, 0x0f, 0x98, 0x2c, 0xe0, 0x6b, 0x33, 0x6e,
	0x07, 0x9c, 0xb0, 0x3b, 0xb3, 0xce, 0x4c, 0x91, 0xfd, 0x11, 0xfe, 0x03, 0xf5, 0xaf, 0x6a, 0x66,
	0xba, 0xcb, 0x2e, 0xc4, 0xf0, 0x42, 0x7c, 0xe2, 0xde, 0x73, 0xce, 0x1d, 0xce, 0xfd, 0xd8, 0x02,
	0xce, 0x98, 0x0e, 0x15, 0x4f, 0x8c, 0x54, 0x9b, 0x89, 0x92, 0x46, 0x12, 0x28, 0x90, 0x8d, 0xdf,
	0x15, 0xc0, 0x53, 0xa3, 0xe6, 0xa1, 0x39, 0xbc, 0x05, 0x09, 0x81, 0x9a, 0xa0, 0x31, 0xeb, 0x79,
	0x7d, 0x6f, 0xd0, 0x0a, 0x5c, 0x4c, 0x5e, 0xc3, 0x72, 0x48, 0x63, 0x16, 0x4d, 0x43, 0xaa, 0xd9,
	0xd4, 0xd1, 0x15, 0x47, 0x77, 0x1d, 0x7c, 0x40, 0x35, 0x3b, 0xb1, 0xba, 0x1e, 0x34, 0x43, 0x19,
	0xc7, 0x4c, 0x98, 0x5e, 0xd5, 0xf1, 0x79, 0x4a, 0x46, 0xd0, 0xb8, 0xe0, 0x2c, 0x9a, 0xe9, 0x5e,
	0xad, 0x5f, 0x1d, 0xb4, 0x87, 0x6b, 0x9b, 0x25, 0x67, 0xc7, 0x96, 0x29, 0x2c, 0x04, 0x99, 0x94,
	0x1c, 0x40, 0x53, 0x26, 0x86, 0x4b, 0xa1, 0x7b, 0x75, 0x57, 0xf5, 0xa6, 0x5c, 0x75, 0xdf, 0xf9,
	0xe6, 0xa7, 0x85, 0xf6, 0x48, 0x18, 0x95, 0x06, 0x79, 0xe5, 0xea, 0x2e, 0x74, 0xca, 0x04, 0x41,
	0xa8, 0x5e, 0xb1, 0x34, 0x6b, 0xcf, 0x86, 0x64, 0x05, 0xea, 0xd7, 0x34, 0x9a, 0xe7, 0x3d, 0x2d,
	0x92, 0xdd, 0xca, 0x8e, 0xb7, 0xf1, 0xb3, 0x0a, 0xcb, 0xf7, 0xcc, 0x3d, 0x6a, 0x3e, 0x03, 0xa8,
	0x99, 0x34, 0x61, 0x6e, 0x38, 0x4b, 0xc3, 0x95, 0x72, 0x37, 0x67, 0x69, 0xc2, 0x8e, 0xc4, 0x3c,
	0x0e, 0x9c, 0x82, 0xac, 0x41, 0xcb, 0xfe, 0x5d, 0xbc, 0x55, 0x73, 0x6f, 0xf9, 0x16, 0x70, 0xcf,
	0xbc, 0x03, 0x22, 0x15, 0xbf, 0xe4, 0x82, 0x46, 0xd3, 0x42, 0x55, 0x77, 0x2a, 0xcc, 0x99, 0xb3,
	0x5c, 0x5d, 0x5a, 0x4a, 0xe3, 0xee, 0x52, 0x5e, 0x42, 0x27, 0x94, 0xd1, 0x3c, 0x16, 0x53, 0x2e,
	0x66, 0xec, 0xa6, 0xd7, 0xec, 0x7b, 0x83, 0x6e, 0xd0, 0x5e, 0x60, 0x13, 0x0b, 0x59, 0x1f, 0x5c,
	0x4f, 0xaf, 0x59, 0x68, 0xa4, 0xea, 0xf9, 0x7d, 0x6f, 0xe0, 0x07, 0x3e, 0xd7, 0x9f, 0x5d, 0x4e,
	0xf6, 0x8b, 0xfd, 0x80, 0xdb, 0xcf, 0xe0, 0x81, 0xad, 0xfe, 0x87, 0xf5, 0xfc, 0xaa, 0x40, 0x67,
	0x12, 0x27, 0x52, 0x99, 0x80, 0xe9, 0x79, 0x64, 0x6c, 0xab, 0xd7, 0x4c, 0x69, 0x2e, 0x45, 0xf6,
	0x40, 0x9e, 0x96, 0x87, 0x50, 0xb9, 0x3b, 0x84, 0x75, 0x68, 0x19, 0x1e, 0x33, 0x6d, 0x68, 0x9c,
	0x64, 0x57, 0x5b, 0x00, 0xe4, 0x43, 0xd1, 0xe2, 0xe2, 0x70, 0x5f, 0x95, 0x5b, 0x2c, 0xff, 0xf3,
	0x7f, 0xf7, 0x47, 0xde, 0x43, 0xbb, 0x28, 0xc8, 0xef, 0x78, 0xfd, 0xa1, 0x3b, 0x0e, 0xca, 0x05,
	0x8f, 0x9a, 0xcf, 0x0f, 0x0f, 0xba, 0x47, 0x37, 0x0b, 0x8b, 0xdf, 0xe6, 0x4c, 0x3f, 0x34, 0xa0,
	0x17, 0xd0, 0xb8, 0x90, 0x2a, 0xa6, 0xf9, 0x97, 0x9b, 0x65, 0x64, 0x15, 0xfc, 0x0b, 0x1e, 0xb1,
	0x84, 0x9a, 0xaf, 0xf9, 0x1d, 0xe6, 0xb9, 0xe5, 0x66, 0xd4, 0x50, 0x9b, 0x67, 0xd7, 0x77, 0x9b,
	0xdb, 0xcf, 0x84, 0xaa, 0x4b, 0x9d, 0x9d, 0x9c, 0x8b, 0xdf, 0xfe, 0xf1, 0xc0, 0xcf, 0xef, 0x9c,
	0xb4, 0xa1, 0x79, 0x2e, 0xae, 0x84, 0xfc, 0x2e, 0xf0, 0x09, 0x69, 0x42, 0xf5, 0x84, 0x47, 0xe8,
	0x11, 0x1f, 0x6a, 0xfb, 0x52, 0x46, 0x58, 0xb1, 0xd1, 0x44, 0x98, 0x1d, 0xac, 0x92, 0x16, 0xd4,
	0xcf, 0xb9, 0x0d, 0x6b, 0x36, 0x9c, 0x08, 0xb3, 0x35, 0xc6, 0x3a, 0x01, 0x68, 0x58, 0x74, 0x6b,
	0x8c, 0x0d, 0x5b, 0x3e, 0x11, 0x06, 0x9b, 0x19, 0x3f, 0x1a, 0xa2, 0x9f, 0xf3, 0xa3, 0x21, 0xb6,
	0x32, 0x78, 0xbc, 0x8d, 0x90, 0xc3, 0xe3, 0x6d, 0x6c, 0x5b, 0xf8, 0x38, 0x92, 0xd4, 0x60, 0xc7,
	0xba, 0x71, 0xe1, 0x68, 0x88, 0xdd, 0xdb, 0x64, 0xbc, 0x8d, 0x4b, 0xb6, 0xe0, 0xd4, 0x28, 0x2e,
	0x2e, 0x71, 0xd9, 0x7a, 0xb2, 0xde, 0x11, 0x6d, 0xe9, 0x7e, 0x6a, 0x98, 0xc6, 0xa7, 0xa4, 0x03,
	0xfe, 0x21, 0x35, 0xec, 0x8c, 0xc7, 0x0c, 0x89, 0x95, 0x7c, 0xd4, 0x52, 0xe0, 0x33, 0x2b, 0xd9,
	0x53, 0x8a, 0xa6, 0xb8, 0x62, 0xfd, 0xed, 0x89, 0x14, 0x9f, 0x7f, 0x69, 0xb8, 0x1f, 0xe1, 0xd1,
	0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xbb, 0xa4, 0x21, 0xad, 0x98, 0x05, 0x00, 0x00,
}
