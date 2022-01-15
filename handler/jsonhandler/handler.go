package jsonhandler

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	grpc "github.com/jschwinger233/grpcdump/grpchelper"
	"github.com/jschwinger233/grpcdump/grpchelper/grpcurl"
	"github.com/jschwinger233/grpcdump/handler"
)

type JsonHandler struct {
	grpcurlManager *grpcurl.Manager
}

func New(grpcurlManager *grpcurl.Manager) handler.GrpcHandler {
	return &JsonHandler{
		grpcurlManager: grpcurlManager,
	}
}

type Output struct {
	Time         time.Time              `json:"time"`
	PacketNumber int                    `json:"packet_number"`
	Src          string                 `json:"src"`
	Dst          string                 `json:"dst"`
	Sport        int                    `json:"sport"`
	Dport        int                    `json:"dport"`
	StreamID     uint32                 `json:"stream_id"`
	Type         string                 `json:"type"`
	Payload      interface{}            `json:"payload"`
	Ext          map[grpc.ExtKey]string `json:"ext"`
	Grpcurl      string                 `json:"grpcurl,omitempty"`
}

func (h *JsonHandler) Handle(msg grpc.Message) (err error) {
	var bytes []byte
	o := Output{
		Time:         msg.CaptureInfo.Timestamp,
		PacketNumber: msg.PacketNumber,
		Src:          msg.Src,
		Dst:          msg.Dst,
		Sport:        msg.Sport,
		Dport:        msg.Dport,
		StreamID:     msg.HTTP2Header.StreamID,
		Ext:          msg.Ext,
	}
	switch msg.Type {

	case grpc.DataType:
		o.Type = "Data"
		o.Payload = "unknown"
		if msg.Data != nil {
			if bytes, err = msg.Data.MarshalJSON(); err != nil {
				return
			}
			if err = json.Unmarshal(bytes, &o.Payload); err != nil {
				return
			}
		}

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
			o.Grpcurl = cmd
		}

	case grpc.HeaderType:
		o.Type = "Header"
		if bytes, err = json.Marshal(msg.Header); err != nil {
			return
		}
		if err = json.Unmarshal(bytes, &o.Payload); err != nil {
			return
		}
	}

	if bytes, err = json.Marshal(o); err != nil {
		return
	}
	fmt.Printf("%s\n", string(bytes))
	return
}
