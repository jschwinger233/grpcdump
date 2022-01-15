package handler

import "github.com/jschwinger233/grpcdump/grpchelper"

type GrpcHandler interface {
	Handle(grpchelper.Message) error
}
