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
	Time     time.Time `json:"time"`
	Src      string    `json:"src"`
	Dst      string    `json:"dst"`
	Sport    int       `json:"sport"`
	Dport    int       `json:"dport"`
	StreamID uint32    `json:"stream_id"`
	Type     string    `json:"type"`
	Payload  string    `json:"payload"`
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
		o.Payload = msg.Request.String()
		bytes, err = json.Marshal(o)
	case grpchelper.ResponseType:
		o.Type = "Data"
		o.Payload = msg.Response.String()
		bytes, err = json.Marshal(0)
	case grpchelper.HeaderType:
		o.Type = "Header"
		o.Payload = fmt.Sprintf("%q", msg.Header)
		bytes, err = json.Marshal(o)
	case grpchelper.UnknownType:
		o.Type = "Unknown"
		o.Payload = "unknown data frame"
		bytes, err = json.Marshal(o)
	}
	if err != nil {
		return
	}
	fmt.Printf("%s\n", string(bytes))
	return
}
