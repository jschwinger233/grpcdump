package jsonhandler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jschwinger23/grpcdump/grpchelper"
	"github.com/jschwinger23/grpcdump/handler"
)

type JsonHandler struct{}

func New() handler.GrpcHandler {
	return &JsonHandler{}
}

type Output struct {
	Time     time.Time   `json:"time"`
	Src      string      `json:"src"`
	Dst      string      `json:"dst"`
	Sport    int         `json:"sport"`
	Dport    int         `json:"dport"`
	StreamID uint32      `json:"stream_id"`
	Type     string      `json:"type"`
	Payload  interface{} `json:"payload"`
}

func (h *JsonHandler) Handle(msg grpchelper.Message) (err error) {
	var bytes []byte
	o := Output{
		Time:     msg.CaptureInfo.Timestamp,
		Src:      msg.Src,
		Dst:      msg.Dst,
		Sport:    msg.Sport,
		Dport:    msg.Dport,
		StreamID: msg.HTTP2Header.StreamID,
	}
	switch msg.Type {
	case grpchelper.RequestType:
		o.Type = "Data"
		if bytes, err = msg.Response.MarshalJSON(); err != nil {
			return
		}
		if err = json.Unmarshal(bytes, &o.Payload); err != nil {
			return
		}
	case grpchelper.ResponseType:
		o.Type = "Data"
		if bytes, err = msg.Response.MarshalJSON(); err != nil {
			return
		}
		if err = json.Unmarshal(bytes, &o.Payload); err != nil {
			return
		}
	case grpchelper.HeaderType:
		o.Type = "Header"
		if bytes, err = json.Marshal(msg.Header); err != nil {
			return
		}
		if err = json.Unmarshal(bytes, &o.Payload); err != nil {
			return
		}
	case grpchelper.UnknownType:
		o.Type = "Unknown"
		o.Payload = "unknown data frame"
	}
	if bytes, err = json.Marshal(o); err != nil {
		return
	}
	fmt.Printf("%s\n", string(bytes))
	return
}
