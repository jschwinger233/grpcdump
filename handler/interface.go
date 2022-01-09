package handler

import "github.com/jschwinger23/grpcdump/grpchelper"

type GrpcHandler interface {
	Handle(grpchelper.Message) error
}
