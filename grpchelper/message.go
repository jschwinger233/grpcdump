package grpchelper

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/jhump/protoreflect/dynamic"
	"golang.org/x/net/http2"
)

type Type int

const (
	HeaderType Type = iota
	DataType
)

type Meta struct {
	PacketNumber int
	gopacket.CaptureInfo
	Src, Dst     string
	Sport, Dport int
	HTTP2Header  http2.FrameHeader
}

type ExtKey string

const (
	HeaderPartiallyParsed ExtKey = "partial_header"
	DataGuessed           ExtKey = "data_guessed"
	DataPath              ExtKey = "data_path"
)

type Message struct {
	Meta
	Type
	Header map[string]string
	Data   *dynamic.Message

	Ext map[ExtKey]string
}

func (m Message) ConnID() string {
	return fmt.Sprintf("%s:%d->%s:%d", m.Src, m.Sport, m.Dst, m.Dport)
}

func (m Message) RevConnID() string {
	return fmt.Sprintf("%s:%d->%s:%d", m.Dst, m.Dport, m.Src, m.Sport)
}
