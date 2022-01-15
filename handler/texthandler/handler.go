package texthandler

import (
	"fmt"
	"strings"
	"time"

	grpc "github.com/jschwinger23/grpcdump/grpchelper"
	"github.com/jschwinger23/grpcdump/grpchelper/grpcurl"
	"github.com/jschwinger23/grpcdump/handler"
)

type TextHandler struct {
	grpcurlManager *grpcurl.Manager
}

func New(grpcurlManager *grpcurl.Manager) handler.GrpcHandler {
	return &TextHandler{
		grpcurlManager: grpcurlManager,
	}
}

func (h *TextHandler) Handle(msg grpc.Message) (err error) {
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
			"%s\t%s:%d->%s:%d\tpacketno:%d\tstreamid:%d\tdata:%s%s\n",
			msg.CaptureInfo.Timestamp.Format(time.StampMicro),
			msg.Src,
			msg.Sport,
			msg.Dst,
			msg.Dport,
			msg.PacketNumber,
			msg.HTTP2Header.StreamID,
			indicator,
			data,
		)

		if h.grpcurlManager != nil && msg.Ext[grpc.DataDirection] == grpc.C2S {
			cmd, err := h.grpcurlManager.Render(grpcurl.RenderContext{
				Payload: msg.Data,
				Dst:     msg.Dst,
				Dport:   msg.Dport,
				Path:    strings.TrimPrefix(msg.Ext[grpc.DataPath], "/"),
			})
			if err != nil {
				return err
			}
			fmt.Println(cmd)
		}

	case grpc.HeaderType:
		partialIndicator := ""
		if _, ok := msg.Ext[grpc.HeaderPartiallyParsed]; ok {
			partialIndicator = "(partial)"
		}
		fmt.Printf(
			"%s\t%s:%d->%s:%d\tpacketno:%d\tstreamid:%d\theader:%s%+v\n",
			msg.CaptureInfo.Timestamp.Format(time.StampMicro),
			msg.Src,
			msg.Sport,
			msg.Dst,
			msg.Dport,
			msg.PacketNumber,
			msg.HTTP2Header.StreamID,
			partialIndicator,
			msg.Header,
		)
	}
	return
}
