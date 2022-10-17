// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: heartbeat.proto

package heartbeat

import (
	bytes "bytes"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// HeartbeatV2 represents the heartbeat message that is sent between peers from the same shard containing
// current node status
type HeartbeatV2 struct {
	Payload         []byte `protobuf:"bytes,1,opt,name=Payload,proto3" json:"Payload,omitempty"`
	VersionNumber   string `protobuf:"bytes,2,opt,name=VersionNumber,proto3" json:"VersionNumber,omitempty"`
	NodeDisplayName string `protobuf:"bytes,3,opt,name=NodeDisplayName,proto3" json:"NodeDisplayName,omitempty"`
	Identity        string `protobuf:"bytes,4,opt,name=Identity,proto3" json:"Identity,omitempty"`
	Nonce           uint64 `protobuf:"varint,5,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	PeerSubType     uint32 `protobuf:"varint,6,opt,name=PeerSubType,proto3" json:"PeerSubType,omitempty"`
	Pubkey          []byte `protobuf:"bytes,7,opt,name=Pubkey,proto3" json:"Pubkey,omitempty"`
}

func (m *HeartbeatV2) Reset()      { *m = HeartbeatV2{} }
func (*HeartbeatV2) ProtoMessage() {}
func (*HeartbeatV2) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c667767fb9826a9, []int{0}
}
func (m *HeartbeatV2) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HeartbeatV2) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HeartbeatV2.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HeartbeatV2) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatV2.Merge(m, src)
}
func (m *HeartbeatV2) XXX_Size() int {
	return m.Size()
}
func (m *HeartbeatV2) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatV2.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatV2 proto.InternalMessageInfo

func (m *HeartbeatV2) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *HeartbeatV2) GetVersionNumber() string {
	if m != nil {
		return m.VersionNumber
	}
	return ""
}

func (m *HeartbeatV2) GetNodeDisplayName() string {
	if m != nil {
		return m.NodeDisplayName
	}
	return ""
}

func (m *HeartbeatV2) GetIdentity() string {
	if m != nil {
		return m.Identity
	}
	return ""
}

func (m *HeartbeatV2) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *HeartbeatV2) GetPeerSubType() uint32 {
	if m != nil {
		return m.PeerSubType
	}
	return 0
}

func (m *HeartbeatV2) GetPubkey() []byte {
	if m != nil {
		return m.Pubkey
	}
	return nil
}

// PeerAuthentication represents the DTO used to pass peer authentication information such as public key, peer id,
// signature, payload and the signature. This message is used to link the peerID with the associated public key
type PeerAuthentication struct {
	Pubkey           []byte `protobuf:"bytes,1,opt,name=Pubkey,proto3" json:"Pubkey,omitempty"`
	Signature        []byte `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
	Pid              []byte `protobuf:"bytes,3,opt,name=Pid,proto3" json:"Pid,omitempty"`
	Payload          []byte `protobuf:"bytes,4,opt,name=Payload,proto3" json:"Payload,omitempty"`
	PayloadSignature []byte `protobuf:"bytes,5,opt,name=PayloadSignature,proto3" json:"PayloadSignature,omitempty"`
}

func (m *PeerAuthentication) Reset()      { *m = PeerAuthentication{} }
func (*PeerAuthentication) ProtoMessage() {}
func (*PeerAuthentication) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c667767fb9826a9, []int{1}
}
func (m *PeerAuthentication) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PeerAuthentication) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PeerAuthentication.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PeerAuthentication) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerAuthentication.Merge(m, src)
}
func (m *PeerAuthentication) XXX_Size() int {
	return m.Size()
}
func (m *PeerAuthentication) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerAuthentication.DiscardUnknown(m)
}

var xxx_messageInfo_PeerAuthentication proto.InternalMessageInfo

func (m *PeerAuthentication) GetPubkey() []byte {
	if m != nil {
		return m.Pubkey
	}
	return nil
}

func (m *PeerAuthentication) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *PeerAuthentication) GetPid() []byte {
	if m != nil {
		return m.Pid
	}
	return nil
}

func (m *PeerAuthentication) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *PeerAuthentication) GetPayloadSignature() []byte {
	if m != nil {
		return m.PayloadSignature
	}
	return nil
}

// Payload represents the DTO used as payload for both HeartbeatV2 and PeerAuthentication messages
type Payload struct {
	Timestamp          int64  `protobuf:"varint,1,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	HardforkMessage    string `protobuf:"bytes,2,opt,name=HardforkMessage,proto3" json:"HardforkMessage,omitempty"`
	NumTrieNodesSynced uint64 `protobuf:"varint,3,opt,name=NumTrieNodesSynced,proto3" json:"NumTrieNodesSynced,omitempty"`
}

func (m *Payload) Reset()      { *m = Payload{} }
func (*Payload) ProtoMessage() {}
func (*Payload) Descriptor() ([]byte, []int) {
	return fileDescriptor_3c667767fb9826a9, []int{2}
}
func (m *Payload) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Payload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Payload.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Payload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Payload.Merge(m, src)
}
func (m *Payload) XXX_Size() int {
	return m.Size()
}
func (m *Payload) XXX_DiscardUnknown() {
	xxx_messageInfo_Payload.DiscardUnknown(m)
}

var xxx_messageInfo_Payload proto.InternalMessageInfo

func (m *Payload) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Payload) GetHardforkMessage() string {
	if m != nil {
		return m.HardforkMessage
	}
	return ""
}

func (m *Payload) GetNumTrieNodesSynced() uint64 {
	if m != nil {
		return m.NumTrieNodesSynced
	}
	return 0
}

func init() {
	proto.RegisterType((*HeartbeatV2)(nil), "proto.HeartbeatV2")
	proto.RegisterType((*PeerAuthentication)(nil), "proto.PeerAuthentication")
	proto.RegisterType((*Payload)(nil), "proto.Payload")
}

func init() { proto.RegisterFile("heartbeat.proto", fileDescriptor_3c667767fb9826a9) }

var fileDescriptor_3c667767fb9826a9 = []byte{
	// 398 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xb1, 0x6e, 0xd4, 0x40,
	0x10, 0xf5, 0x70, 0xf6, 0x85, 0xdb, 0x5c, 0x94, 0x68, 0x85, 0xd0, 0x0a, 0xa1, 0x95, 0x75, 0xa2,
	0xb0, 0x28, 0x52, 0xc0, 0x07, 0x20, 0x10, 0x45, 0x28, 0xb0, 0xac, 0xbd, 0x53, 0x0a, 0xba, 0xf5,
	0x79, 0x48, 0x56, 0x89, 0xbd, 0xd6, 0x7a, 0x5d, 0xb8, 0x83, 0x3f, 0xe0, 0x1b, 0xa8, 0xf8, 0x14,
	0xca, 0x2b, 0x53, 0x72, 0x76, 0x43, 0x99, 0x4f, 0x40, 0xde, 0x5c, 0xee, 0x9c, 0x23, 0x95, 0xe7,
	0xbd, 0x79, 0x1a, 0xbf, 0x79, 0x3b, 0xe4, 0xf8, 0x12, 0xa5, 0xb1, 0x29, 0x4a, 0x7b, 0x5a, 0x1a,
	0x6d, 0x35, 0x0d, 0xdc, 0x67, 0xd6, 0x01, 0x39, 0x3c, 0xbb, 0x6f, 0x9d, 0xbf, 0xa1, 0x8c, 0x1c,
	0x24, 0xb2, 0xb9, 0xd6, 0x32, 0x63, 0x10, 0x42, 0x34, 0x15, 0xf7, 0x90, 0xbe, 0x22, 0x47, 0xe7,
	0x68, 0x2a, 0xa5, 0x8b, 0xb8, 0xce, 0x53, 0x34, 0xec, 0x49, 0x08, 0xd1, 0x44, 0x3c, 0x24, 0x69,
	0x44, 0x8e, 0x63, 0x9d, 0xe1, 0x47, 0x55, 0x95, 0xd7, 0xb2, 0x89, 0x65, 0x8e, 0x6c, 0xe4, 0x74,
	0xfb, 0x34, 0x7d, 0x41, 0x9e, 0x7e, 0xca, 0xb0, 0xb0, 0xca, 0x36, 0xcc, 0x77, 0x92, 0x2d, 0xa6,
	0xcf, 0x48, 0x10, 0xeb, 0x62, 0x89, 0x2c, 0x08, 0x21, 0xf2, 0xc5, 0x1d, 0xa0, 0x21, 0x39, 0x4c,
	0x10, 0xcd, 0xbc, 0x4e, 0x17, 0x4d, 0x89, 0x6c, 0x1c, 0x42, 0x74, 0x24, 0x86, 0x14, 0x7d, 0x4e,
	0xc6, 0x49, 0x9d, 0x5e, 0x61, 0xc3, 0x0e, 0x9c, 0xf9, 0x0d, 0x9a, 0xfd, 0x04, 0x42, 0x7b, 0xdd,
	0xfb, 0xda, 0x5e, 0xf6, 0xbf, 0x58, 0x4a, 0xab, 0x74, 0x31, 0x90, 0xc3, 0x50, 0x4e, 0x5f, 0x92,
	0xc9, 0x5c, 0x5d, 0x14, 0xd2, 0xd6, 0x06, 0xdd, 0x9a, 0x53, 0xb1, 0x23, 0xe8, 0x09, 0x19, 0x25,
	0x2a, 0x73, 0x6b, 0x4d, 0x45, 0x5f, 0x0e, 0x43, 0xf3, 0x1f, 0x86, 0xf6, 0x9a, 0x9c, 0x6c, 0xca,
	0xdd, 0xc0, 0xc0, 0x49, 0xfe, 0xe3, 0x67, 0xdf, 0x61, 0x3b, 0xa6, 0x77, 0xb0, 0x50, 0x39, 0x56,
	0x56, 0xe6, 0xa5, 0x33, 0x37, 0x12, 0x3b, 0xa2, 0x0f, 0xf9, 0x4c, 0x9a, 0xec, 0xab, 0x36, 0x57,
	0x9f, 0xb1, 0xaa, 0xe4, 0x05, 0x6e, 0x1e, 0x63, 0x9f, 0xa6, 0xa7, 0x84, 0xc6, 0x75, 0xbe, 0x30,
	0x0a, 0xfb, 0xf8, 0xab, 0x79, 0x53, 0x2c, 0xf1, 0xce, 0xba, 0x2f, 0x1e, 0xe9, 0x7c, 0x78, 0xb7,
	0x5a, 0x73, 0xef, 0x66, 0xcd, 0xbd, 0xdb, 0x35, 0x87, 0x6f, 0x2d, 0x87, 0x5f, 0x2d, 0x87, 0xdf,
	0x2d, 0x87, 0x55, 0xcb, 0xe1, 0x4f, 0xcb, 0xe1, 0x6f, 0xcb, 0xbd, 0xdb, 0x96, 0xc3, 0x8f, 0x8e,
	0x7b, 0xab, 0x8e, 0x7b, 0x37, 0x1d, 0xf7, 0xbe, 0x4c, 0xb6, 0xc7, 0x95, 0x8e, 0xdd, 0x59, 0xbd,
	0xfd, 0x17, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x96, 0x8d, 0x30, 0x70, 0x02, 0x00, 0x00,
}

func (this *HeartbeatV2) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HeartbeatV2)
	if !ok {
		that2, ok := that.(HeartbeatV2)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.Payload, that1.Payload) {
		return false
	}
	if this.VersionNumber != that1.VersionNumber {
		return false
	}
	if this.NodeDisplayName != that1.NodeDisplayName {
		return false
	}
	if this.Identity != that1.Identity {
		return false
	}
	if this.Nonce != that1.Nonce {
		return false
	}
	if this.PeerSubType != that1.PeerSubType {
		return false
	}
	if !bytes.Equal(this.Pubkey, that1.Pubkey) {
		return false
	}
	return true
}
func (this *PeerAuthentication) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PeerAuthentication)
	if !ok {
		that2, ok := that.(PeerAuthentication)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.Pubkey, that1.Pubkey) {
		return false
	}
	if !bytes.Equal(this.Signature, that1.Signature) {
		return false
	}
	if !bytes.Equal(this.Pid, that1.Pid) {
		return false
	}
	if !bytes.Equal(this.Payload, that1.Payload) {
		return false
	}
	if !bytes.Equal(this.PayloadSignature, that1.PayloadSignature) {
		return false
	}
	return true
}
func (this *Payload) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Payload)
	if !ok {
		that2, ok := that.(Payload)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Timestamp != that1.Timestamp {
		return false
	}
	if this.HardforkMessage != that1.HardforkMessage {
		return false
	}
	if this.NumTrieNodesSynced != that1.NumTrieNodesSynced {
		return false
	}
	return true
}
func (this *HeartbeatV2) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 11)
	s = append(s, "&heartbeat.HeartbeatV2{")
	s = append(s, "Payload: "+fmt.Sprintf("%#v", this.Payload)+",\n")
	s = append(s, "VersionNumber: "+fmt.Sprintf("%#v", this.VersionNumber)+",\n")
	s = append(s, "NodeDisplayName: "+fmt.Sprintf("%#v", this.NodeDisplayName)+",\n")
	s = append(s, "Identity: "+fmt.Sprintf("%#v", this.Identity)+",\n")
	s = append(s, "Nonce: "+fmt.Sprintf("%#v", this.Nonce)+",\n")
	s = append(s, "PeerSubType: "+fmt.Sprintf("%#v", this.PeerSubType)+",\n")
	s = append(s, "Pubkey: "+fmt.Sprintf("%#v", this.Pubkey)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *PeerAuthentication) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&heartbeat.PeerAuthentication{")
	s = append(s, "Pubkey: "+fmt.Sprintf("%#v", this.Pubkey)+",\n")
	s = append(s, "Signature: "+fmt.Sprintf("%#v", this.Signature)+",\n")
	s = append(s, "Pid: "+fmt.Sprintf("%#v", this.Pid)+",\n")
	s = append(s, "Payload: "+fmt.Sprintf("%#v", this.Payload)+",\n")
	s = append(s, "PayloadSignature: "+fmt.Sprintf("%#v", this.PayloadSignature)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Payload) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&heartbeat.Payload{")
	s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	s = append(s, "HardforkMessage: "+fmt.Sprintf("%#v", this.HardforkMessage)+",\n")
	s = append(s, "NumTrieNodesSynced: "+fmt.Sprintf("%#v", this.NumTrieNodesSynced)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringHeartbeat(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *HeartbeatV2) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HeartbeatV2) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HeartbeatV2) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Pubkey) > 0 {
		i -= len(m.Pubkey)
		copy(dAtA[i:], m.Pubkey)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.Pubkey)))
		i--
		dAtA[i] = 0x3a
	}
	if m.PeerSubType != 0 {
		i = encodeVarintHeartbeat(dAtA, i, uint64(m.PeerSubType))
		i--
		dAtA[i] = 0x30
	}
	if m.Nonce != 0 {
		i = encodeVarintHeartbeat(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Identity) > 0 {
		i -= len(m.Identity)
		copy(dAtA[i:], m.Identity)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.Identity)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.NodeDisplayName) > 0 {
		i -= len(m.NodeDisplayName)
		copy(dAtA[i:], m.NodeDisplayName)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.NodeDisplayName)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.VersionNumber) > 0 {
		i -= len(m.VersionNumber)
		copy(dAtA[i:], m.VersionNumber)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.VersionNumber)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PeerAuthentication) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PeerAuthentication) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PeerAuthentication) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PayloadSignature) > 0 {
		i -= len(m.PayloadSignature)
		copy(dAtA[i:], m.PayloadSignature)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.PayloadSignature)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Pid) > 0 {
		i -= len(m.Pid)
		copy(dAtA[i:], m.Pid)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.Pid)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Pubkey) > 0 {
		i -= len(m.Pubkey)
		copy(dAtA[i:], m.Pubkey)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.Pubkey)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Payload) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Payload) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Payload) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.NumTrieNodesSynced != 0 {
		i = encodeVarintHeartbeat(dAtA, i, uint64(m.NumTrieNodesSynced))
		i--
		dAtA[i] = 0x18
	}
	if len(m.HardforkMessage) > 0 {
		i -= len(m.HardforkMessage)
		copy(dAtA[i:], m.HardforkMessage)
		i = encodeVarintHeartbeat(dAtA, i, uint64(len(m.HardforkMessage)))
		i--
		dAtA[i] = 0x12
	}
	if m.Timestamp != 0 {
		i = encodeVarintHeartbeat(dAtA, i, uint64(m.Timestamp))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintHeartbeat(dAtA []byte, offset int, v uint64) int {
	offset -= sovHeartbeat(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *HeartbeatV2) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	l = len(m.VersionNumber)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	l = len(m.NodeDisplayName)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	l = len(m.Identity)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	if m.Nonce != 0 {
		n += 1 + sovHeartbeat(uint64(m.Nonce))
	}
	if m.PeerSubType != 0 {
		n += 1 + sovHeartbeat(uint64(m.PeerSubType))
	}
	l = len(m.Pubkey)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	return n
}

func (m *PeerAuthentication) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Pubkey)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	l = len(m.Pid)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	l = len(m.PayloadSignature)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	return n
}

func (m *Payload) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Timestamp != 0 {
		n += 1 + sovHeartbeat(uint64(m.Timestamp))
	}
	l = len(m.HardforkMessage)
	if l > 0 {
		n += 1 + l + sovHeartbeat(uint64(l))
	}
	if m.NumTrieNodesSynced != 0 {
		n += 1 + sovHeartbeat(uint64(m.NumTrieNodesSynced))
	}
	return n
}

func sovHeartbeat(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozHeartbeat(x uint64) (n int) {
	return sovHeartbeat(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *HeartbeatV2) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&HeartbeatV2{`,
		`Payload:` + fmt.Sprintf("%v", this.Payload) + `,`,
		`VersionNumber:` + fmt.Sprintf("%v", this.VersionNumber) + `,`,
		`NodeDisplayName:` + fmt.Sprintf("%v", this.NodeDisplayName) + `,`,
		`Identity:` + fmt.Sprintf("%v", this.Identity) + `,`,
		`Nonce:` + fmt.Sprintf("%v", this.Nonce) + `,`,
		`PeerSubType:` + fmt.Sprintf("%v", this.PeerSubType) + `,`,
		`Pubkey:` + fmt.Sprintf("%v", this.Pubkey) + `,`,
		`}`,
	}, "")
	return s
}
func (this *PeerAuthentication) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PeerAuthentication{`,
		`Pubkey:` + fmt.Sprintf("%v", this.Pubkey) + `,`,
		`Signature:` + fmt.Sprintf("%v", this.Signature) + `,`,
		`Pid:` + fmt.Sprintf("%v", this.Pid) + `,`,
		`Payload:` + fmt.Sprintf("%v", this.Payload) + `,`,
		`PayloadSignature:` + fmt.Sprintf("%v", this.PayloadSignature) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Payload) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Payload{`,
		`Timestamp:` + fmt.Sprintf("%v", this.Timestamp) + `,`,
		`HardforkMessage:` + fmt.Sprintf("%v", this.HardforkMessage) + `,`,
		`NumTrieNodesSynced:` + fmt.Sprintf("%v", this.NumTrieNodesSynced) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringHeartbeat(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *HeartbeatV2) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHeartbeat
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HeartbeatV2: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HeartbeatV2: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VersionNumber", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.VersionNumber = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeDisplayName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NodeDisplayName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Identity", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Identity = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PeerSubType", wireType)
			}
			m.PeerSubType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PeerSubType |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pubkey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pubkey = append(m.Pubkey[:0], dAtA[iNdEx:postIndex]...)
			if m.Pubkey == nil {
				m.Pubkey = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHeartbeat(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PeerAuthentication) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHeartbeat
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PeerAuthentication: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PeerAuthentication: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pubkey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pubkey = append(m.Pubkey[:0], dAtA[iNdEx:postIndex]...)
			if m.Pubkey == nil {
				m.Pubkey = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pid", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pid = append(m.Pid[:0], dAtA[iNdEx:postIndex]...)
			if m.Pid == nil {
				m.Pid = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PayloadSignature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PayloadSignature = append(m.PayloadSignature[:0], dAtA[iNdEx:postIndex]...)
			if m.PayloadSignature == nil {
				m.PayloadSignature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHeartbeat(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Payload) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHeartbeat
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Payload: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Payload: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HardforkMessage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthHeartbeat
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HardforkMessage = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumTrieNodesSynced", wireType)
			}
			m.NumTrieNodesSynced = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumTrieNodesSynced |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipHeartbeat(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthHeartbeat
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipHeartbeat(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHeartbeat
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowHeartbeat
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthHeartbeat
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupHeartbeat
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthHeartbeat
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthHeartbeat        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHeartbeat          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupHeartbeat = fmt.Errorf("proto: unexpected end of group")
)
