package parser

import (
	"github.com/google/gopacket"
	"github.com/jhump/protoreflect/dynamic"
)

type Parser struct{}

func New(protoFilename, guessMethod string) *Parser {
	return nil
}

func (p *Parser) Parse(packet gopacket.Packet) (msg dynamic.Message, err error) {
	return
}
