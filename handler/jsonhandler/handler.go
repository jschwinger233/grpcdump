package jsonhandler

import (
	"encoding/json"
	"fmt"
	"time"

	grpc "github.com/jschwinger23/grpcdump/grpchelper"
	"github.com/jschwinger23/grpcdump/handler"
)

type JsonHandler struct{}

func New() handler.GrpcHandler {
	return &JsonHandler{}
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
	case grpc.RequestType:
		o.Type = "Data"
		if bytes, err = msg.Response.MarshalJSON(); err != nil {
			return
		}
		if err = json.Unmarshal(bytes, &o.Payload); err != nil {
			return
		}
	case grpc.ResponseType:
		o.Type = "Data"
		if bytes, err = msg.Response.MarshalJSON(); err != nil {
			return
		}
		if err = json.Unmarshal(bytes, &o.Payload); err != nil {
			return
		}
	case grpc.HeaderType:
		o.Type = "Header"
		if bytes, err = json.Marshal(msg.Header); err != nil {
			return
		}
		if err = json.Unmarshal(bytes, &o.Payload); err != nil {
			return
		}
	case grpc.UnknownType:
		o.Type = "Unknown"
		o.Payload = "unknown data frame"
	}
	if bytes, err = json.Marshal(o); err != nil {
		return
	}
	fmt.Printf("%s\n", string(bytes))
	return
}
