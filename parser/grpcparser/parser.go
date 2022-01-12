package grpcparser

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"

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
	protoParser   grpchelper.ProtoParser
	streams       map[ConnID]map[StreamID][]grpchelper.Message
}

func New(protoFilename string, guessPaths []string) (_ parser.Parser, err error) {
	protoParser, err := grpchelper.NewProtoParser(protoFilename)
	if err != nil {
		return
	}

	return &Parser{
		protoFilename: protoFilename,
		guessPaths:    guessPaths,
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
		sport, _ := strconv.Atoi(packet.TransportLayer().TransportFlow().Src().String())
		dport, _ := strconv.Atoi(packet.TransportLayer().TransportFlow().Dst().String())
		var message grpchelper.Message = grpchelper.Message{
			Meta: grpchelper.Meta{
				CaptureInfo: packet.Metadata().CaptureInfo,
				HTTP2Header: frame.Header(),
				Src:         packet.NetworkLayer().NetworkFlow().Src().String(),
				Dst:         packet.NetworkLayer().NetworkFlow().Src().String(),
				Sport:       sport,
				Dport:       dport,
			},
		}
		if _, ok := p.streams[message.ConnID()]; !ok {
			p.streams[message.ConnID()] = make(map[StreamID][]grpchelper.Message)
		}
		switch frame := frame.(type) {
		case *http2.HeadersFrame:
			payload := map[string]string{}
			buf := frame.HeaderBlockFragment()
			for _, field := range hpackDecodePartial(buf) {
				payload[field.Name] = field.Value
			}

			message.Type = grpchelper.HeaderType
			message.Header = payload

		case *http2.DataFrame:
			var (
				possibleTypes []grpchelper.Type
				possiblePaths []string
			)

			for _, msg := range p.streams[message.ConnID()][streamID] {
				if msg.Type == grpchelper.HeaderType {
					for key := range msg.Header {
						if key == ":path" {
							possibleTypes = []grpchelper.Type{grpchelper.RequestType}
							possiblePaths = []string{msg.Header[key]}
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
				for _, msg := range p.streams[message.RevConnID()][streamID] {
					if msg.Type == grpchelper.HeaderType {
						for key := range msg.Header {
							if key == ":path" {
								possibleTypes = []grpchelper.Type{grpchelper.ResponseType}
								possiblePaths = []string{msg.Header[key]}
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
				if len(possiblePaths) == 1 && possiblePaths[0] == "AUTO" {
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
				msg := msgs[0]
				message.Type = msg.Type
				message.Request = msg.Request
				message.Response = msg.Response
				message.Ext = msg.Ext
			case 0:
				fallthrough
			default:
				curMax := -1
				for _, msg := range msgs {
					var n int
					switch msg.Type {
					case grpchelper.RequestType:
						n = len(msg.Request.String())
					case grpchelper.ResponseType:
						n = len(msg.Response.String())
					}
					if n > curMax {
						curMax = n
						message.Type = msg.Type
						message.Request = msg.Request
						message.Response = msg.Response
						message.Ext = msg.Ext
					}
				}
			}

		default:
			continue
		}
		p.streams[message.ConnID()][streamID] = append(p.streams[message.ConnID()][streamID], message)
		messages = append(messages, message)
	}
	return messages, nil
}

func (p *Parser) unmarshalDataFrame(dataType grpchelper.Type, path string, frame *http2.DataFrame) (message grpchelper.Message, err error) {
	message.Type = dataType
	message.Ext = map[string]string{":path": path}
	message.Response, err = p.protoParser.MarshalResponse(path, frame.Data()[5:])
	if dataType == grpchelper.RequestType {
		message.Request, err = p.protoParser.MarshalRequest(path, frame.Data()[5:])
	}
	return
}
