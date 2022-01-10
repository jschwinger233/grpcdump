package grpchelper

import (
	"github.com/jhump/protoreflect/dynamic"
	"golang.org/x/net/http2"
)

type Header struct {
	HTTP2Header http2.FrameHeader
	Payload     map[string]string
}

type Request struct {
	HTTP2Header http2.FrameHeader
	Payload     *dynamic.Message
}

type Response struct {
	HTTP2Header http2.FrameHeader
	Payload     *dynamic.Message
}

type Type int

const (
	UnknownType Type = iota
	HeaderType
	RequestType
	ResponseType
)

type Message struct {
	Type
	Header
	Request
	Response
}
