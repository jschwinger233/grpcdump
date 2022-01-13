package texthandler

import (
	"fmt"

	grpc "github.com/jschwinger23/grpcdump/grpchelper"
	"github.com/jschwinger23/grpcdump/handler"
)

type TextHandler struct{}

func New() handler.GrpcHandler {
	return &TextHandler{}
}

func (h *TextHandler) Handle(msg grpc.Message) (err error) {
	// time, conn, streamid, data
	switch msg.Type {

	case grpc.RequestType:
		guessIndicator := ""
		if _, ok := msg.Ext[grpc.DataGuessed]; ok {
			guessIndicator = "(guess)"
		}
		fmt.Printf("%s\t%s\tstreamid:%d\tdata:%s%s\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID, guessIndicator, msg.Request.String())

	case grpc.ResponseType:
		guessIndicator := ""
		if _, ok := msg.Ext[grpc.DataGuessed]; ok {
			guessIndicator = "(guess)"
		}
		fmt.Printf("%s\t%s\tstreamid:%d\tdata:%s%s\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID, guessIndicator, msg.Response.String())

	case grpc.HeaderType:
		partialIndicator := ""
		if _, ok := msg.Ext[grpc.HeaderPartiallyParsed]; ok {
			partialIndicator = "(partial)"
		}
		fmt.Printf("%s\t%s\tstreamid:%d\theader:%s%+v\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID, partialIndicator, msg.Header)

	case grpc.UnknownType:
		fmt.Printf("%s\t%s\tstreamid:%d\tunknown data frame\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID)
	}
	return
}
