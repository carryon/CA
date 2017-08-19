// Copyright (C) 2017, Beijing Bochen Technology Co.,Ltd.  All rights reserved.
//
// This file is part of msg-net 
// 
// The msg-net is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// The msg-net is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// 
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Code generated by protoc-gen-go.
// source: message.proto
// DO NOT EDIT!

/*
Package protos is a generated protocol buffer package.

It is generated from these files:
	message.proto

It has these top-level messages:
	Message
	Router
	Routers
	Peer
	Peers
	ChainMessage
*/
package protos

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

type Message_Type int32

const (
	Message_UNDEFINED        Message_Type = 0
	Message_ROUTER_HELLO     Message_Type = 1
	Message_ROUTER_HELLO_ACK Message_Type = 2
	Message_ROUTER_CLOSE     Message_Type = 3
	Message_ROUTER_GET       Message_Type = 4
	Message_ROUTER_GET_ACK   Message_Type = 5
	Message_ROUTER_SYNC      Message_Type = 6
	Message_PEER_HELLO       Message_Type = 11
	Message_PEER_HELLO_ACK   Message_Type = 12
	Message_PEER_CLOSE       Message_Type = 13
	Message_PEER_SYNC        Message_Type = 14
	Message_CHAIN_MESSAGE    Message_Type = 21
	Message_KEEPALIVE        Message_Type = 31
	Message_KEEPALIVE_ACK    Message_Type = 32
)

var Message_Type_name = map[int32]string{
	0:  "UNDEFINED",
	1:  "ROUTER_HELLO",
	2:  "ROUTER_HELLO_ACK",
	3:  "ROUTER_CLOSE",
	4:  "ROUTER_GET",
	5:  "ROUTER_GET_ACK",
	6:  "ROUTER_SYNC",
	11: "PEER_HELLO",
	12: "PEER_HELLO_ACK",
	13: "PEER_CLOSE",
	14: "PEER_SYNC",
	21: "CHAIN_MESSAGE",
	31: "KEEPALIVE",
	32: "KEEPALIVE_ACK",
}
var Message_Type_value = map[string]int32{
	"UNDEFINED":        0,
	"ROUTER_HELLO":     1,
	"ROUTER_HELLO_ACK": 2,
	"ROUTER_CLOSE":     3,
	"ROUTER_GET":       4,
	"ROUTER_GET_ACK":   5,
	"ROUTER_SYNC":      6,
	"PEER_HELLO":       11,
	"PEER_HELLO_ACK":   12,
	"PEER_CLOSE":       13,
	"PEER_SYNC":        14,
	"CHAIN_MESSAGE":    21,
	"KEEPALIVE":        31,
	"KEEPALIVE_ACK":    32,
}

func (x Message_Type) String() string {
	return proto.EnumName(Message_Type_name, int32(x))
}
func (Message_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Message struct {
	Type     Message_Type `protobuf:"varint,1,opt,name=type,enum=protos.Message_Type" json:"type,omitempty"`
	Payload  []byte       `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	Metadata []byte       `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetType() Message_Type {
	if m != nil {
		return m.Type
	}
	return Message_UNDEFINED
}

func (m *Message) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Message) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Router struct {
	Id      string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address" json:"address,omitempty"`
}

func (m *Router) Reset()                    { *m = Router{} }
func (m *Router) String() string            { return proto.CompactTextString(m) }
func (*Router) ProtoMessage()               {}
func (*Router) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Router) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Router) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type Routers struct {
	Id      string    `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Routers []*Router `protobuf:"bytes,2,rep,name=routers" json:"routers,omitempty"`
}

func (m *Routers) Reset()                    { *m = Routers{} }
func (m *Routers) String() string            { return proto.CompactTextString(m) }
func (*Routers) ProtoMessage()               {}
func (*Routers) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Routers) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Routers) GetRouters() []*Router {
	if m != nil {
		return m.Routers
	}
	return nil
}

type Peer struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *Peer) Reset()                    { *m = Peer{} }
func (m *Peer) String() string            { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()               {}
func (*Peer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Peer) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Peers struct {
	Id    string  `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Peers []*Peer `protobuf:"bytes,2,rep,name=peers" json:"peers,omitempty"`
}

func (m *Peers) Reset()                    { *m = Peers{} }
func (m *Peers) String() string            { return proto.CompactTextString(m) }
func (*Peers) ProtoMessage()               {}
func (*Peers) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Peers) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Peers) GetPeers() []*Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}

type ChainMessage struct {
	SrcId     string `protobuf:"bytes,1,opt,name=srcId" json:"srcId,omitempty"`
	DstId     string `protobuf:"bytes,2,opt,name=dstId" json:"dstId,omitempty"`
	Payload   []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	Signature []byte `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *ChainMessage) Reset()                    { *m = ChainMessage{} }
func (m *ChainMessage) String() string            { return proto.CompactTextString(m) }
func (*ChainMessage) ProtoMessage()               {}
func (*ChainMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ChainMessage) GetSrcId() string {
	if m != nil {
		return m.SrcId
	}
	return ""
}

func (m *ChainMessage) GetDstId() string {
	if m != nil {
		return m.DstId
	}
	return ""
}

func (m *ChainMessage) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *ChainMessage) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "protos.Message")
	proto.RegisterType((*Router)(nil), "protos.Router")
	proto.RegisterType((*Routers)(nil), "protos.Routers")
	proto.RegisterType((*Peer)(nil), "protos.Peer")
	proto.RegisterType((*Peers)(nil), "protos.Peers")
	proto.RegisterType((*ChainMessage)(nil), "protos.ChainMessage")
	proto.RegisterEnum("protos.Message_Type", Message_Type_name, Message_Type_value)
}

func init() { proto.RegisterFile("message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 406 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x52, 0x4f, 0x8f, 0x93, 0x40,
	0x14, 0x17, 0x4a, 0x5b, 0xfb, 0x4a, 0x71, 0x9c, 0x54, 0x43, 0x8c, 0x89, 0xcd, 0x9c, 0x38, 0xf5,
	0x50, 0x8f, 0x9e, 0x1a, 0x76, 0xdc, 0x25, 0xcb, 0xd2, 0x66, 0xe8, 0x9a, 0x78, 0x6a, 0x46, 0x67,
	0xb2, 0x92, 0xb8, 0x0b, 0x61, 0xe8, 0xa1, 0xdf, 0xd8, 0x0f, 0xe1, 0xc1, 0xcc, 0x0c, 0x7f, 0x6a,
	0x7a, 0x22, 0xbf, 0xbf, 0x0f, 0x1e, 0x0f, 0x16, 0xcf, 0x52, 0x29, 0xfe, 0x24, 0xd7, 0x55, 0x5d,
	0x36, 0x25, 0x9e, 0x98, 0x87, 0x22, 0x7f, 0x5c, 0x98, 0x3e, 0x58, 0x05, 0x47, 0xe0, 0x35, 0xe7,
	0x4a, 0x86, 0xce, 0xca, 0x89, 0x82, 0xcd, 0xd2, 0x3a, 0xd5, 0xba, 0x95, 0xd7, 0x87, 0x73, 0x25,
	0x99, 0x71, 0xe0, 0x10, 0xa6, 0x15, 0x3f, 0xff, 0x2e, 0xb9, 0x08, 0xdd, 0x95, 0x13, 0xf9, 0xac,
	0x83, 0xf8, 0x03, 0xbc, 0x7e, 0x96, 0x0d, 0x17, 0xbc, 0xe1, 0xe1, 0xc8, 0x48, 0x3d, 0x26, 0x7f,
	0x1d, 0xf0, 0x74, 0x09, 0x5e, 0xc0, 0xec, 0x31, 0xbb, 0xa1, 0x5f, 0x93, 0x8c, 0xde, 0xa0, 0x57,
	0x18, 0x81, 0xcf, 0x76, 0x8f, 0x07, 0xca, 0x8e, 0x77, 0x34, 0x4d, 0x77, 0xc8, 0xc1, 0x4b, 0x40,
	0x97, 0xcc, 0x71, 0x1b, 0xdf, 0x23, 0xf7, 0xc2, 0x17, 0xa7, 0xbb, 0x9c, 0xa2, 0x11, 0x0e, 0x00,
	0x5a, 0xe6, 0x96, 0x1e, 0x90, 0x87, 0x31, 0x04, 0x03, 0x36, 0xa9, 0x31, 0x7e, 0x03, 0xf3, 0x96,
	0xcb, 0xbf, 0x67, 0x31, 0x9a, 0xe8, 0xd0, 0x9e, 0xf6, 0xc3, 0xe6, 0x3a, 0x34, 0x60, 0x13, 0xf2,
	0x7b, 0x8f, 0x1d, 0xb4, 0xd0, 0x6f, 0x6c, 0xb0, 0xa9, 0x08, 0xf0, 0x5b, 0x58, 0xc4, 0x77, 0xdb,
	0x24, 0x3b, 0x3e, 0xd0, 0x3c, 0xdf, 0xde, 0x52, 0xf4, 0x4e, 0x3b, 0xee, 0x29, 0xdd, 0x6f, 0xd3,
	0xe4, 0x1b, 0x45, 0x9f, 0xb4, 0xa3, 0x87, 0xa6, 0x73, 0x45, 0x36, 0x30, 0x61, 0xe5, 0xa9, 0x91,
	0x35, 0x0e, 0xc0, 0x2d, 0x84, 0x59, 0xf3, 0x8c, 0xb9, 0x85, 0xd0, 0xeb, 0xe4, 0x42, 0xd4, 0x52,
	0x29, 0xb3, 0xce, 0x19, 0xeb, 0x20, 0x89, 0x61, 0x6a, 0x33, 0xea, 0x2a, 0x14, 0xc1, 0xb4, 0xb6,
	0x52, 0xe8, 0xae, 0x46, 0xd1, 0x7c, 0x13, 0x74, 0x3f, 0xcc, 0x26, 0x58, 0x27, 0x93, 0xf7, 0xe0,
	0xed, 0xe5, 0xf5, 0x58, 0xf2, 0x05, 0xc6, 0x9a, 0xbf, 0xae, 0x26, 0x30, 0xae, 0xe4, 0x50, 0xec,
	0x77, 0xc5, 0xda, 0xcd, 0xac, 0x44, 0x6a, 0xf0, 0xe3, 0x5f, 0xbc, 0x78, 0xe9, 0x8e, 0x67, 0x09,
	0x63, 0x55, 0xff, 0x4c, 0xba, 0x1a, 0x0b, 0x34, 0x2b, 0x54, 0x93, 0x88, 0xf6, 0xbb, 0x2c, 0xb8,
	0x3c, 0x9f, 0xd1, 0xff, 0xe7, 0xf3, 0x11, 0x66, 0xaa, 0x78, 0x7a, 0xe1, 0xcd, 0xa9, 0x96, 0xa1,
	0x67, 0xb4, 0x81, 0xf8, 0x61, 0x8f, 0xf6, 0xf3, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4a, 0xb0,
	0x6e, 0x8f, 0xcc, 0x02, 0x00, 0x00,
}
