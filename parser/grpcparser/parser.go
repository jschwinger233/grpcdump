package grpcparser

import (
	"bytes"
	"io"
	"io/ioutil"

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
	protoFilename string
	guessPaths    []string
	autoGuess     bool
	protoParser   grpchelper.ProtoParser
	streams       map[ConnID]map[StreamID][]grpchelper.Message
}

type GrpcStreamMessage struct {
	ConnID
}

func New(protoFilename string, guessPaths []string, autoGuess bool) (_ parser.Parser, err error) {
	protoParser, err := grpchelper.NewProtoParser(protoFilename)
	if err != nil {
		return
	}

	return &Parser{
		protoFilename: protoFilename,
		guessPaths:    guessPaths,
		autoGuess:     autoGuess,
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
				possibleTypes []grpchelper.Type
				possiblePaths []string
			)

			for _, msg := range p.streams[connID][streamID] {
				if msg.Type == grpchelper.HeaderType {
					for key := range msg.Header.Payload {
						if key == ":path" {
							possibleTypes = []grpchelper.Type{grpchelper.RequestType}
							possiblePaths = []string{msg.Header.Payload[key]}
						}
						if key == ":status" {
							possibleTypes = []grpchelper.Type{grpchelper.ResponseType}
						}
					}
				}
			}

			// search opposite flow
			responseStream := len(possibleTypes) == 1 && possibleTypes[0] == grpchelper.ResponseType
			unknownStream := len(possibleTypes) == 0
			if responseStream || unknownStream {
				for _, msg := range p.streams[revConnID(connID)][streamID] {
					if msg.Type == grpchelper.HeaderType {
						for key := range msg.Header.Payload {
							if key == ":path" {
								possibleTypes = []grpchelper.Type{grpchelper.ResponseType}
								possiblePaths = []string{msg.Header.Payload[key]}
							}
							if key == ":status" {
								possibleTypes = []grpchelper.Type{grpchelper.RequestType}
							}
						}
					}
				}
			}

			if len(possiblePaths) == 0 {
				possibleTypes = []grpchelper.Type{grpchelper.RequestType, grpchelper.ResponseType}
				possiblePaths = p.guessPaths
				if p.autoGuess {
					possiblePaths = p.protoParser.GetAllPaths()
				}
			}

			msgs := []grpchelper.Message{}
			for _, possiblePath := range possiblePaths {
				for _, possibleType := range possibleTypes {
					msg, err := p.unmarshalDataFrame(possibleType, possiblePath, frame)
					if err == nil {
						msgs = append(msgs, msg)
					}
				}
			}

			switch len(msgs) {
			case 1:
				message = msgs[0]
			case 0:
				println("unknown data frame")
				continue
			default:
				curMax := 0
				for _, msg := range msgs {
					var n int
					switch msg.Type {
					case grpchelper.RequestType:
						n = len(msg.Request.Payload.String())
					case grpchelper.ResponseType:
						n = len(msg.Response.Payload.String())
					}
					if n >= curMax {
						curMax = n
						message = msg
					}
				}
			}

		}
		p.streams[connID][streamID] = append(p.streams[connID][streamID], message)
		messages = append(messages, message)
	}
	return messages, nil
}

func (p *Parser) unmarshalDataFrame(dataType grpchelper.Type, path string, frame *http2.DataFrame) (message grpchelper.Message, err error) {

	if dataType == grpchelper.RequestType {
		message = grpchelper.Message{
			Type: grpchelper.RequestType,
			Request: grpchelper.Request{
				HTTP2Header: frame.Header(),
			},
		}
		message.Request.Payload, err = p.protoParser.MarshalRequest(path, frame.Data()[5:])
		return
	}

	message = grpchelper.Message{
		Type: grpchelper.ResponseType,
		Response: grpchelper.Response{
			HTTP2Header: frame.Header(),
		},
	}
	message.Response.Payload, err = p.protoParser.MarshalResponse(path, frame.Data()[5:])
	return
}
