package handler

import "github.com/jschwinger23/grpcdump/grpc"

type GrpcHandler interface {
	HandleGrpcHeader(grpc.GrpcHeader) error
	HandleGrpcRequest(grpc.GrpcRequest) error
	HandleGrpcResponse(grpc.GrpcRequest) error
}
