package parser

import (
	"github.com/google/gopacket"
	"github.com/jschwinger233/grpcdump/grpchelper"
)

type Parser interface {
	Parse(gopacket.Packet) ([]grpchelper.Message, error)
}
