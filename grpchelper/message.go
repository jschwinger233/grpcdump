package grpchelper

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/jhump/protoreflect/dynamic"
	"golang.org/x/net/http2"
)

type Type int

const (
	UnknownType Type = iota
	HeaderType
	RequestType
	ResponseType
)

type Meta struct {
	gopacket.CaptureInfo
	Src, Dst     string
	Sport, Dport int
	HTTP2Header  http2.FrameHeader
}

type Message struct {
	Meta
	Type
	Header   map[string]string
	Request  *dynamic.Message
	Response *dynamic.Message

	Ext map[string]string
}

func (m Message) ConnID() string {
	return fmt.Sprintf("%s:%d->%s:%d", m.Src, m.Sport, m.Dst, m.Dport)
}

func (m Message) RevConnID() string {
	return fmt.Sprintf("%s:%d->%s:%d", m.Dst, m.Dport, m.Src, m.Sport)
}
