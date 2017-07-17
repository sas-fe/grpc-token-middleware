// Code generated by protoc-gen-go.
// source: tokenstore.proto
// DO NOT EDIT!

/*
Package tokenstore is a generated protocol buffer package.

It is generated from these files:
	tokenstore.proto

It has these top-level messages:
	Token
	TokenMsg
	TokenStatusMsg
	RpcStatus
*/
package tokenstore

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Token is a message containing the ID.
type Token struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Token) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// TokenMsg contains token ID and limit.
type TokenMsg struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Limit int32  `protobuf:"varint,2,opt,name=limit" json:"limit,omitempty"`
}

func (m *TokenMsg) Reset()                    { *m = TokenMsg{} }
func (m *TokenMsg) String() string            { return proto.CompactTextString(m) }
func (*TokenMsg) ProtoMessage()               {}
func (*TokenMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TokenMsg) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *TokenMsg) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

// TokenStatusMsg contains token ID, limit, and current usage.
type TokenStatusMsg struct {
	Id      string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Limit   int32  `protobuf:"varint,2,opt,name=limit" json:"limit,omitempty"`
	Usage   int32  `protobuf:"varint,3,opt,name=usage" json:"usage,omitempty"`
	Allowed bool   `protobuf:"varint,4,opt,name=allowed" json:"allowed,omitempty"`
}

func (m *TokenStatusMsg) Reset()                    { *m = TokenStatusMsg{} }
func (m *TokenStatusMsg) String() string            { return proto.CompactTextString(m) }
func (*TokenStatusMsg) ProtoMessage()               {}
func (*TokenStatusMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *TokenStatusMsg) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *TokenStatusMsg) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *TokenStatusMsg) GetUsage() int32 {
	if m != nil {
		return m.Usage
	}
	return 0
}

func (m *TokenStatusMsg) GetAllowed() bool {
	if m != nil {
		return m.Allowed
	}
	return false
}

// generic rpc call status
type RpcStatus struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *RpcStatus) Reset()                    { *m = RpcStatus{} }
func (m *RpcStatus) String() string            { return proto.CompactTextString(m) }
func (*RpcStatus) ProtoMessage()               {}
func (*RpcStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RpcStatus) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func init() {
	proto.RegisterType((*Token)(nil), "tokenstore.Token")
	proto.RegisterType((*TokenMsg)(nil), "tokenstore.TokenMsg")
	proto.RegisterType((*TokenStatusMsg)(nil), "tokenstore.TokenStatusMsg")
	proto.RegisterType((*RpcStatus)(nil), "tokenstore.RpcStatus")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for TokenStore service

type TokenStoreClient interface {
	// add token to redis
	AddToken(ctx context.Context, in *TokenMsg, opts ...grpc.CallOption) (*RpcStatus, error)
	// update limit on token
	UpdateLimit(ctx context.Context, in *TokenMsg, opts ...grpc.CallOption) (*RpcStatus, error)
	// token status with limit and current usage
	TokenStatus(ctx context.Context, in *Token, opts ...grpc.CallOption) (*TokenStatusMsg, error)
	// increment the token current usage
	IncUsage(ctx context.Context, in *Token, opts ...grpc.CallOption) (*RpcStatus, error)
	// delete token
	DelToken(ctx context.Context, in *Token, opts ...grpc.CallOption) (*RpcStatus, error)
}

type tokenStoreClient struct {
	cc *grpc.ClientConn
}

func NewTokenStoreClient(cc *grpc.ClientConn) TokenStoreClient {
	return &tokenStoreClient{cc}
}

func (c *tokenStoreClient) AddToken(ctx context.Context, in *TokenMsg, opts ...grpc.CallOption) (*RpcStatus, error) {
	out := new(RpcStatus)
	err := grpc.Invoke(ctx, "/tokenstore.TokenStore/AddToken", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenStoreClient) UpdateLimit(ctx context.Context, in *TokenMsg, opts ...grpc.CallOption) (*RpcStatus, error) {
	out := new(RpcStatus)
	err := grpc.Invoke(ctx, "/tokenstore.TokenStore/UpdateLimit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenStoreClient) TokenStatus(ctx context.Context, in *Token, opts ...grpc.CallOption) (*TokenStatusMsg, error) {
	out := new(TokenStatusMsg)
	err := grpc.Invoke(ctx, "/tokenstore.TokenStore/TokenStatus", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenStoreClient) IncUsage(ctx context.Context, in *Token, opts ...grpc.CallOption) (*RpcStatus, error) {
	out := new(RpcStatus)
	err := grpc.Invoke(ctx, "/tokenstore.TokenStore/IncUsage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenStoreClient) DelToken(ctx context.Context, in *Token, opts ...grpc.CallOption) (*RpcStatus, error) {
	out := new(RpcStatus)
	err := grpc.Invoke(ctx, "/tokenstore.TokenStore/DelToken", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for TokenStore service

type TokenStoreServer interface {
	// add token to redis
	AddToken(context.Context, *TokenMsg) (*RpcStatus, error)
	// update limit on token
	UpdateLimit(context.Context, *TokenMsg) (*RpcStatus, error)
	// token status with limit and current usage
	TokenStatus(context.Context, *Token) (*TokenStatusMsg, error)
	// increment the token current usage
	IncUsage(context.Context, *Token) (*RpcStatus, error)
	// delete token
	DelToken(context.Context, *Token) (*RpcStatus, error)
}

func RegisterTokenStoreServer(s *grpc.Server, srv TokenStoreServer) {
	s.RegisterService(&_TokenStore_serviceDesc, srv)
}

func _TokenStore_AddToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenStoreServer).AddToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokenstore.TokenStore/AddToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenStoreServer).AddToken(ctx, req.(*TokenMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenStore_UpdateLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenStoreServer).UpdateLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokenstore.TokenStore/UpdateLimit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenStoreServer).UpdateLimit(ctx, req.(*TokenMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenStore_TokenStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenStoreServer).TokenStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokenstore.TokenStore/TokenStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenStoreServer).TokenStatus(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenStore_IncUsage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenStoreServer).IncUsage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokenstore.TokenStore/IncUsage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenStoreServer).IncUsage(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenStore_DelToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenStoreServer).DelToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tokenstore.TokenStore/DelToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenStoreServer).DelToken(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

var _TokenStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tokenstore.TokenStore",
	HandlerType: (*TokenStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddToken",
			Handler:    _TokenStore_AddToken_Handler,
		},
		{
			MethodName: "UpdateLimit",
			Handler:    _TokenStore_UpdateLimit_Handler,
		},
		{
			MethodName: "TokenStatus",
			Handler:    _TokenStore_TokenStatus_Handler,
		},
		{
			MethodName: "IncUsage",
			Handler:    _TokenStore_IncUsage_Handler,
		},
		{
			MethodName: "DelToken",
			Handler:    _TokenStore_DelToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tokenstore.proto",
}

func init() { proto.RegisterFile("tokenstore.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xd1, 0x4a, 0xc3, 0x30,
	0x14, 0x5d, 0x3b, 0xa7, 0xd9, 0x1d, 0x4c, 0x0d, 0x1b, 0x96, 0x3e, 0x95, 0x80, 0xd0, 0xa7, 0x32,
	0x14, 0x06, 0x82, 0x08, 0x0e, 0x5f, 0x44, 0x85, 0x91, 0x6d, 0x1f, 0x10, 0x93, 0x6c, 0x14, 0xbb,
	0xa5, 0x34, 0x29, 0xfe, 0x8f, 0xff, 0xe0, 0xff, 0x49, 0x13, 0xbb, 0x15, 0xa7, 0xe0, 0x1e, 0xcf,
	0xb9, 0xf7, 0x70, 0xee, 0x39, 0x5c, 0x38, 0x33, 0xea, 0x4d, 0x6e, 0xb4, 0x51, 0x85, 0x4c, 0xf2,
	0x42, 0x19, 0x85, 0x61, 0xc7, 0x90, 0x0b, 0xe8, 0xcc, 0x2b, 0x84, 0xfb, 0xe0, 0xa7, 0x22, 0xf0,
	0x22, 0x2f, 0xee, 0x52, 0x3f, 0x15, 0x64, 0x04, 0xc8, 0x0e, 0x5e, 0xf4, 0xea, 0xe7, 0x0c, 0x0f,
	0xa0, 0x93, 0xa5, 0xeb, 0xd4, 0x04, 0x7e, 0xe4, 0xc5, 0x1d, 0xea, 0x00, 0x59, 0x42, 0xdf, 0x2a,
	0x66, 0x86, 0x99, 0x52, 0xff, 0x5b, 0x57, 0xb1, 0xa5, 0x66, 0x2b, 0x19, 0xb4, 0x1d, 0x6b, 0x01,
	0x0e, 0xe0, 0x84, 0x65, 0x99, 0x7a, 0x97, 0x22, 0x38, 0x8a, 0xbc, 0x18, 0xd1, 0x1a, 0x92, 0x4b,
	0xe8, 0xd2, 0x9c, 0x3b, 0x97, 0x6a, 0x4d, 0x97, 0x9c, 0x4b, 0xad, 0xad, 0x0f, 0xa2, 0x35, 0xbc,
	0xfa, 0xf4, 0x01, 0xbe, 0xef, 0x51, 0x85, 0xc4, 0x37, 0x80, 0xee, 0x85, 0x70, 0x59, 0x07, 0x49,
	0xa3, 0x93, 0x3a, 0x65, 0x38, 0x6c, 0xb2, 0x5b, 0x07, 0xd2, 0xc2, 0xb7, 0xd0, 0x5b, 0xe4, 0x82,
	0x19, 0xf9, 0xec, 0xee, 0x3d, 0x4c, 0x7d, 0x07, 0xbd, 0x46, 0x2d, 0xf8, 0x7c, 0x4f, 0x1d, 0x86,
	0x7b, 0xd4, 0xb6, 0x42, 0xd2, 0xc2, 0x63, 0x40, 0x8f, 0x1b, 0xbe, 0xb0, 0xa5, 0xfc, 0x22, 0xfe,
	0xd3, 0x77, 0x0c, 0xe8, 0x41, 0x66, 0x2e, 0xf0, 0x01, 0xba, 0xc9, 0x08, 0x86, 0x5c, 0xad, 0x13,
	0xcd, 0x74, 0xb2, 0x94, 0x8d, 0xa5, 0xc9, 0xe9, 0xae, 0xcd, 0x69, 0xf5, 0x47, 0x53, 0xef, 0xc3,
	0x6f, 0xcf, 0x9f, 0x66, 0xaf, 0xc7, 0xf6, 0xad, 0xae, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x65,
	0x6e, 0x40, 0x77, 0x6a, 0x02, 0x00, 0x00,
}