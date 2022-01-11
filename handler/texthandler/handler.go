package texthandler

import (
	"fmt"

	"github.com/jschwinger23/grpcdump/grpchelper"
	"github.com/jschwinger23/grpcdump/handler"
)

type TextHandler struct{}

func New() handler.GrpcHandler {
	return &TextHandler{}
}

func (h *TextHandler) Handle(msg grpchelper.Message) (err error) {
	// time, conn, streamid, data
	switch msg.Type {
	case grpchelper.RequestType:
		fmt.Printf("%s\t%s\tstreamid:%d\tdata:%s\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID, msg.Request.String())
	case grpchelper.ResponseType:
		fmt.Printf("%s\t%s\tstreamid:%d\tdata:%s\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID, msg.Response.String())
	case grpchelper.HeaderType:
		fmt.Printf("%s\t%s\tstreamid:%d\theader:%+v\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID, msg.Header)
	case grpchelper.UnknownType:
		fmt.Printf("%s\t%s\tstreamid:%d\tunknown data frame\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID, msg.Header)
	}
	return
}
