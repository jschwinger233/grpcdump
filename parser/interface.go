package parser

import (
	"github.com/google/gopacket"
	"github.com/jschwinger23/grpcdump/handler"
)

type Parser interface {
	Parse(gopacket.Packet) (func(handler.GrpcHandler) error, error)
}
