// Code generated by protoc-gen-go. DO NOT EDIT.
// source: request_reply.proto

package model

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type TypeRequests int32

const (
	TypeRequests_PUSHPULL_REQUEST TypeRequests = 0
)

var TypeRequests_name = map[int32]string{
	0: "PUSHPULL_REQUEST",
}

var TypeRequests_value = map[string]int32{
	"PUSHPULL_REQUEST": 0,
}

func (x TypeRequests) String() string {
	return proto.EnumName(TypeRequests_name, int32(x))
}

func (TypeRequests) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_7bf67568605725bb, []int{0}
}

type TypeReplies int32

const (
	TypeReplies_PUSHPULL_REPLY TypeReplies = 0
)

var TypeReplies_name = map[int32]string{
	0: "PUSHPULL_REPLY",
}

var TypeReplies_value = map[string]int32{
	"PUSHPULL_REPLY": 0,
}

func (x TypeReplies) String() string {
	return proto.EnumName(TypeReplies_name, int32(x))
}

func (TypeReplies) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_7bf67568605725bb, []int{1}
}

type RequestHeader struct {
	Type                 TypeRequests `protobuf:"varint,1,opt,name=type,proto3,enum=model.TypeRequests" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *RequestHeader) Reset()         { *m = RequestHeader{} }
func (m *RequestHeader) String() string { return proto.CompactTextString(m) }
func (*RequestHeader) ProtoMessage()    {}
func (*RequestHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bf67568605725bb, []int{0}
}

func (m *RequestHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestHeader.Unmarshal(m, b)
}
func (m *RequestHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestHeader.Marshal(b, m, deterministic)
}
func (m *RequestHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestHeader.Merge(m, src)
}
func (m *RequestHeader) XXX_Size() int {
	return xxx_messageInfo_RequestHeader.Size(m)
}
func (m *RequestHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestHeader.DiscardUnknown(m)
}

var xxx_messageInfo_RequestHeader proto.InternalMessageInfo

func (m *RequestHeader) GetType() TypeRequests {
	if m != nil {
		return m.Type
	}
	return TypeRequests_PUSHPULL_REQUEST
}

type ReplyHeader struct {
	Type                 TypeReplies `protobuf:"varint,1,opt,name=type,proto3,enum=model.TypeReplies" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ReplyHeader) Reset()         { *m = ReplyHeader{} }
func (m *ReplyHeader) String() string { return proto.CompactTextString(m) }
func (*ReplyHeader) ProtoMessage()    {}
func (*ReplyHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bf67568605725bb, []int{1}
}

func (m *ReplyHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReplyHeader.Unmarshal(m, b)
}
func (m *ReplyHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReplyHeader.Marshal(b, m, deterministic)
}
func (m *ReplyHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReplyHeader.Merge(m, src)
}
func (m *ReplyHeader) XXX_Size() int {
	return xxx_messageInfo_ReplyHeader.Size(m)
}
func (m *ReplyHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_ReplyHeader.DiscardUnknown(m)
}

var xxx_messageInfo_ReplyHeader proto.InternalMessageInfo

func (m *ReplyHeader) GetType() TypeReplies {
	if m != nil {
		return m.Type
	}
	return TypeReplies_PUSHPULL_REPLY
}

type PushPullRequest struct {
	Header               *RequestHeader  `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   int32           `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	PushPullPacks        []*PushPullPack `protobuf:"bytes,3,rep,name=pushPullPacks,proto3" json:"pushPullPacks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *PushPullRequest) Reset()         { *m = PushPullRequest{} }
func (m *PushPullRequest) String() string { return proto.CompactTextString(m) }
func (*PushPullRequest) ProtoMessage()    {}
func (*PushPullRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bf67568605725bb, []int{2}
}

func (m *PushPullRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushPullRequest.Unmarshal(m, b)
}
func (m *PushPullRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushPullRequest.Marshal(b, m, deterministic)
}
func (m *PushPullRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushPullRequest.Merge(m, src)
}
func (m *PushPullRequest) XXX_Size() int {
	return xxx_messageInfo_PushPullRequest.Size(m)
}
func (m *PushPullRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PushPullRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PushPullRequest proto.InternalMessageInfo

func (m *PushPullRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *PushPullRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PushPullRequest) GetPushPullPacks() []*PushPullPack {
	if m != nil {
		return m.PushPullPacks
	}
	return nil
}

type PushPullReply struct {
	Header               *ReplyHeader    `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   int32           `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	PushPullPacks        []*PushPullPack `protobuf:"bytes,3,rep,name=pushPullPacks,proto3" json:"pushPullPacks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *PushPullReply) Reset()         { *m = PushPullReply{} }
func (m *PushPullReply) String() string { return proto.CompactTextString(m) }
func (*PushPullReply) ProtoMessage()    {}
func (*PushPullReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bf67568605725bb, []int{3}
}

func (m *PushPullReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushPullReply.Unmarshal(m, b)
}
func (m *PushPullReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushPullReply.Marshal(b, m, deterministic)
}
func (m *PushPullReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushPullReply.Merge(m, src)
}
func (m *PushPullReply) XXX_Size() int {
	return xxx_messageInfo_PushPullReply.Size(m)
}
func (m *PushPullReply) XXX_DiscardUnknown() {
	xxx_messageInfo_PushPullReply.DiscardUnknown(m)
}

var xxx_messageInfo_PushPullReply proto.InternalMessageInfo

func (m *PushPullReply) GetHeader() *ReplyHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *PushPullReply) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PushPullReply) GetPushPullPacks() []*PushPullPack {
	if m != nil {
		return m.PushPullPacks
	}
	return nil
}

func init() {
	proto.RegisterEnum("model.TypeRequests", TypeRequests_name, TypeRequests_value)
	proto.RegisterEnum("model.TypeReplies", TypeReplies_name, TypeReplies_value)
	proto.RegisterType((*RequestHeader)(nil), "model.RequestHeader")
	proto.RegisterType((*ReplyHeader)(nil), "model.ReplyHeader")
	proto.RegisterType((*PushPullRequest)(nil), "model.PushPullRequest")
	proto.RegisterType((*PushPullReply)(nil), "model.PushPullReply")
}

func init() { proto.RegisterFile("request_reply.proto", fileDescriptor_7bf67568605725bb) }

var fileDescriptor_7bf67568605725bb = []byte{
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x92, 0x4d, 0x6b, 0xc2, 0x40,
	0x10, 0x86, 0x8d, 0x56, 0x0f, 0x13, 0xbf, 0x18, 0xa5, 0x88, 0x27, 0x2b, 0xa5, 0x15, 0x29, 0x1e,
	0x2c, 0x85, 0xf6, 0xd4, 0x93, 0xe0, 0x41, 0x70, 0xbb, 0x9a, 0x43, 0x4f, 0x62, 0x93, 0x01, 0x43,
	0xb7, 0xec, 0x76, 0x37, 0x29, 0xe4, 0x0f, 0xf4, 0xd0, 0x5f, 0x5d, 0x9a, 0xac, 0x69, 0x62, 0xcf,
	0x3d, 0x05, 0x26, 0xcf, 0xfb, 0xce, 0x93, 0x21, 0xd0, 0xd3, 0xf4, 0x1e, 0x93, 0x89, 0x76, 0x9a,
	0x94, 0x48, 0x66, 0x4a, 0xcb, 0x48, 0x62, 0xfd, 0x4d, 0x06, 0x24, 0x86, 0x6e, 0xfa, 0xc8, 0x66,
	0xe3, 0x7b, 0x68, 0xf1, 0x0c, 0x5d, 0xd2, 0x3e, 0x20, 0x8d, 0xd7, 0x70, 0x16, 0x25, 0x8a, 0x06,
	0xce, 0xc8, 0x99, 0xb4, 0xe7, 0xbd, 0x59, 0x06, 0x6f, 0x13, 0x45, 0x96, 0x33, 0x3c, 0x05, 0xc6,
	0x77, 0xe0, 0xf2, 0x9f, 0x72, 0x9b, 0xbb, 0x2a, 0xe5, 0xb0, 0x94, 0x53, 0x22, 0xa4, 0x63, 0xec,
	0xcb, 0x81, 0x0e, 0x8b, 0xcd, 0x81, 0xc5, 0x42, 0xd8, 0x46, 0xbc, 0x81, 0xc6, 0x21, 0x6d, 0x49,
	0xd3, 0xee, 0xbc, 0x6f, 0xd3, 0x25, 0x33, 0x6e, 0x19, 0x6c, 0x43, 0x35, 0x0c, 0x06, 0xd5, 0x91,
	0x33, 0xa9, 0xf3, 0x6a, 0x18, 0xe0, 0x03, 0xb4, 0x94, 0x2d, 0x64, 0x7b, 0xff, 0xd5, 0x0c, 0x6a,
	0xa3, 0xda, 0xc4, 0xcd, 0xd5, 0x59, 0xe1, 0x1d, 0x2f, 0x93, 0xe3, 0x4f, 0x07, 0x5a, 0xbf, 0x32,
	0x4a, 0x24, 0x38, 0x3d, 0x51, 0xc1, 0x5c, 0x25, 0xff, 0xd4, 0x7f, 0x10, 0x99, 0x5e, 0x42, 0xb3,
	0x78, 0x62, 0xec, 0x43, 0x97, 0x79, 0x9b, 0x25, 0xf3, 0x56, 0xab, 0x1d, 0x5f, 0x3c, 0x79, 0x8b,
	0xcd, 0xb6, 0x5b, 0x99, 0x5e, 0x80, 0x5b, 0x38, 0x28, 0x22, 0xb4, 0x0b, 0x10, 0x5b, 0x3d, 0x77,
	0x2b, 0xf3, 0x35, 0x34, 0xd7, 0x3a, 0x92, 0x72, 0x43, 0xfa, 0x23, 0xf4, 0x09, 0x1f, 0xa1, 0xc3,
	0xb4, 0xf4, 0xc9, 0x98, 0xe3, 0x7a, 0x3c, 0x3f, 0xf1, 0xb1, 0x4b, 0x87, 0xfd, 0x3f, 0x73, 0x25,
	0x92, 0x97, 0x46, 0xfa, 0x9f, 0xdc, 0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0x85, 0x4b, 0x17, 0x67,
	0x52, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OrtooServiceClient is the client API for OrtooService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrtooServiceClient interface {
	ProcessPushPull(ctx context.Context, in *PushPullRequest, opts ...grpc.CallOption) (*PushPullReply, error)
}

type ortooServiceClient struct {
	cc *grpc.ClientConn
}

func NewOrtooServiceClient(cc *grpc.ClientConn) OrtooServiceClient {
	return &ortooServiceClient{cc}
}

func (c *ortooServiceClient) ProcessPushPull(ctx context.Context, in *PushPullRequest, opts ...grpc.CallOption) (*PushPullReply, error) {
	out := new(PushPullReply)
	err := c.cc.Invoke(ctx, "/model.OrtooService/ProcessPushPull", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrtooServiceServer is the server API for OrtooService service.
type OrtooServiceServer interface {
	ProcessPushPull(context.Context, *PushPullRequest) (*PushPullReply, error)
}

func RegisterOrtooServiceServer(s *grpc.Server, srv OrtooServiceServer) {
	s.RegisterService(&_OrtooService_serviceDesc, srv)
}

func _OrtooService_ProcessPushPull_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushPullRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrtooServiceServer).ProcessPushPull(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/model.OrtooService/ProcessPushPull",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrtooServiceServer).ProcessPushPull(ctx, req.(*PushPullRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrtooService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "model.OrtooService",
	HandlerType: (*OrtooServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessPushPull",
			Handler:    _OrtooService_ProcessPushPull_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "request_reply.proto",
}