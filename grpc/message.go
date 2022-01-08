package grpc

import (
	"github.com/jhump/protoreflect/dynamic"
	"golang.org/x/net/http2"
)

type GrpcHeader struct {
	HTTP2Header http2.FrameHeader
	Payload     map[string]string
}

type GrpcRequest struct {
	HTTP2Heaer http2.FrameHeader
	request    dynamic.Message
}

type GrpcResponse struct {
	HTTP2Heaer http2.FrameHeader
	response   dynamic.Message
}
