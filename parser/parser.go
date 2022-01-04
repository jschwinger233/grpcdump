package parser

import (
	"github.com/google/gopacket"
	"golang.org/x/net/http2"
)

type Parser struct{}

func New(protoFilename, guessMethod string) *Parser {
	return nil
}

func (p *Parser) Parse(packet gopacket.Packet) (frame http2.DataFrame) {
	return
}
