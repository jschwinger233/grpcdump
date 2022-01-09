package texthandler

import (
	"github.com/jschwinger23/grpcdump/grpchelper"
	"github.com/jschwinger23/grpcdump/handler"
)

type TextHandler struct{}

func New(verbose bool) handler.GrpcHandler {
	return &TextHandler{}
}

func (h *TextHandler) Handle(msg grpchelper.Message) (err error) {
	if msg.Type == grpchelper.RequestType {
		println(msg.Request.Payload.String())
	} else if msg.Type == grpchelper.ResponseType {
		println(msg.Response.Payload.String())
	}
	return
}
