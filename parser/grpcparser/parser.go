package grpcparser

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/google/gopacket"
	"github.com/jschwinger23/grpcdump/grpchelper"
	"github.com/jschwinger23/grpcdump/parser"
	"golang.org/x/net/http2"
)

type (
	StreamID = uint32
	ConnID   = string
)

type Parser struct {
	protoFilename             string
	guessService, guessMethod string
	protoParser               grpchelper.ProtoParser
	streams                   map[ConnID]map[StreamID][]grpchelper.Message
}

type GrpcStreamMessage struct {
	ConnID
}

func New(protoFilename, guessPath string) (_ parser.Parser, err error) {
	protoParser, err := grpchelper.NewProtoParser(protoFilename)
	if err != nil {
		return
	}

	var guessService, guessMethod string
	if guessPath != "" {
		parts := strings.Split(guessPath, "/")
		guessService, guessMethod = parts[len(parts)-2], parts[len(parts)-1]
	}
	return &Parser{
		protoFilename: protoFilename,
		guessService:  guessService,
		guessMethod:   guessMethod,
		protoParser:   protoParser,
		streams:       map[ConnID]map[StreamID][]grpchelper.Message{},
	}, nil
}

func (p *Parser) Parse(packet gopacket.Packet) (messages []grpchelper.Message, err error) {
	appLayer := packet.ApplicationLayer()
	if appLayer == nil {
		return
	}
	payload := appLayer.Payload()
	framer := http2.NewFramer(ioutil.Discard, bytes.NewReader(payload))
	for {
		frame, err := framer.ReadFrame()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		streamID := frame.Header().StreamID
		connID := getConnID(packet)
		if _, ok := p.streams[connID]; !ok {
			p.streams[connID] = make(map[StreamID][]grpchelper.Message)
		}

		var message grpchelper.Message
		switch frame := frame.(type) {
		case *http2.HeadersFrame:
			payload := map[string]string{}
			buf := frame.HeaderBlockFragment()
			for _, field := range hpackDecodePartial(buf) {
				payload[field.Name] = field.Value
			}

			message = grpchelper.Message{
				Type: grpchelper.HeaderType,
				Header: grpchelper.Header{
					HTTP2Header: frame.Header(),
					Payload:     payload,
				},
			}

		case *http2.DataFrame:
			var (
				request         bool
				service, method string
			)
			if p.guessMethod != "" {
				service, method = p.guessService, p.guessMethod
			}

			for _, msg := range p.streams[connID][streamID] {
				if msg.Type == grpchelper.HeaderType {
					for key := range msg.Header.Payload {
						if key == ":path" {
							request = true
							parts := strings.Split(msg.Header.Payload[key], "/")
							service, method = parts[1], parts[2]
						}
						if key == ":status" {
							request = false
						}
					}
				}
			}

			// search opposite flow
			if !request {
				for _, msg := range p.streams[revConnID(connID)][streamID] {
					if msg.Type == grpchelper.HeaderType {
						for key := range msg.Header.Payload {
							if key == ":path" {
								parts := strings.Split(msg.Header.Payload[key], "/")
								service, method = parts[1], parts[2]
							}
						}
					}
				}
			}

			if service == "" || method == "" {
				continue
			}

			if request {
				message = grpchelper.Message{
					Type: grpchelper.RequestType,
					Request: grpchelper.Request{
						HTTP2Header: frame.Header(),
					},
				}
				message.Request.Payload, err = p.protoParser.MarshalRequest(service, method, frame.Data()[5:])
			} else {
				message = grpchelper.Message{
					Type: grpchelper.ResponseType,
					Response: grpchelper.Response{
						HTTP2Header: frame.Header(),
					},
				}
				message.Response.Payload, err = p.protoParser.MarshalResponse(service, method, frame.Data()[5:])
			}
			if err != nil {
				continue
			}
		}
		p.streams[connID][streamID] = append(p.streams[connID][streamID], message)
		messages = append(messages, message)
	}
	return messages, nil
}
