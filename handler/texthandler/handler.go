package texthandler

import (
	"fmt"
	"time"

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

	case grpc.DataType:
		var indicator, data string
		if _, ok := msg.Ext[grpc.DataGuessed]; ok {
			indicator = "(guess)"
		}
		if msg.Data == nil {
			indicator = "(unknown)"
		}

		data = ""
		if msg.Data != nil {
			data = msg.Data.String()
		}
		fmt.Printf(
			"%s\t%s\tpacketno:%d\tstreamid:%d\tdata:%s%s\n",
			msg.CaptureInfo.Timestamp.Format(time.StampMicro),
			msg.ConnID(),
			msg.PacketNumber,
			msg.HTTP2Header.StreamID,
			indicator,
			data,
		)

	case grpc.HeaderType:
		partialIndicator := ""
		if _, ok := msg.Ext[grpc.HeaderPartiallyParsed]; ok {
			partialIndicator = "(partial)"
		}
		fmt.Printf(
			"%s\t%s\tpacketno:%d\tstreamid:%d\theader:%s%+v\n",
			msg.CaptureInfo.Timestamp.Format(time.StampMicro),
			msg.ConnID(),
			msg.PacketNumber,
			msg.HTTP2Header.StreamID,
			partialIndicator,
			msg.Header,
		)
	}
	return
}
