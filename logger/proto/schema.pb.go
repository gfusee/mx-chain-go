// Code generated by protoc-gen-go. DO NOT EDIT.
// source: schema.proto

package protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type LogLineMessage struct {
	Message              string   `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	LogLevel             int32    `protobuf:"varint,2,opt,name=LogLevel,proto3" json:"LogLevel,omitempty"`
	Args                 []string `protobuf:"bytes,3,rep,name=Args,proto3" json:"Args,omitempty"`
	Timestamp            int64    `protobuf:"varint,4,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogLineMessage) Reset()         { *m = LogLineMessage{} }
func (m *LogLineMessage) String() string { return proto.CompactTextString(m) }
func (*LogLineMessage) ProtoMessage()    {}
func (*LogLineMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c5fb4d8cc22d66a, []int{0}
}

func (m *LogLineMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogLineMessage.Unmarshal(m, b)
}
func (m *LogLineMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogLineMessage.Marshal(b, m, deterministic)
}
func (m *LogLineMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogLineMessage.Merge(m, src)
}
func (m *LogLineMessage) XXX_Size() int {
	return xxx_messageInfo_LogLineMessage.Size(m)
}
func (m *LogLineMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_LogLineMessage.DiscardUnknown(m)
}

var xxx_messageInfo_LogLineMessage proto.InternalMessageInfo

func (m *LogLineMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *LogLineMessage) GetLogLevel() int32 {
	if m != nil {
		return m.LogLevel
	}
	return 0
}

func (m *LogLineMessage) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *LogLineMessage) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func init() {
	proto.RegisterType((*LogLineMessage)(nil), "protobuf.LogLineMessage")
}

func init() { proto.RegisterFile("schema.proto", fileDescriptor_1c5fb4d8cc22d66a) }

var fileDescriptor_1c5fb4d8cc22d66a = []byte{
	// 137 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4e, 0xce, 0x48,
	0xcd, 0x4d, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53, 0x49, 0xa5, 0x69, 0x4a,
	0x15, 0x5c, 0x7c, 0x3e, 0xf9, 0xe9, 0x3e, 0x99, 0x79, 0xa9, 0xbe, 0xa9, 0xc5, 0xc5, 0x89, 0xe9,
	0xa9, 0x42, 0x12, 0x5c, 0xec, 0x50, 0xa6, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x8c, 0x2b,
	0x24, 0xc5, 0xc5, 0x01, 0x52, 0x9b, 0x5a, 0x96, 0x9a, 0x23, 0xc1, 0x04, 0x94, 0x62, 0x0d, 0x82,
	0xf3, 0x85, 0x84, 0xb8, 0x58, 0x1c, 0x8b, 0xd2, 0x8b, 0x25, 0x98, 0x15, 0x98, 0x81, 0x5a, 0xc0,
	0x6c, 0x21, 0x19, 0x2e, 0xce, 0x90, 0xcc, 0xdc, 0xd4, 0xe2, 0x92, 0xc4, 0xdc, 0x02, 0x09, 0x16,
	0xa0, 0x06, 0xe6, 0x20, 0x84, 0x40, 0x12, 0x1b, 0xd8, 0x0d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x7f, 0x2f, 0xbc, 0xc7, 0x9a, 0x00, 0x00, 0x00,
}
