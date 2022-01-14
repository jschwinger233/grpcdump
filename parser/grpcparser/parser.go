package grpcparser

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/google/gopacket"
	grpc "github.com/jschwinger23/grpcdump/grpchelper"
	"github.com/jschwinger23/grpcdump/parser"
	"golang.org/x/net/http2"
)

type (
	StreamID      = uint32
	ConnID        = string
	DataDirection int
)

const (
	ToService DataDirection = iota
	ToClient
)

type Parser struct {
	protoFilename string
	servicePort   int
	guessPaths    []string
	protoParser   grpc.ProtoParser
	streams       map[ConnID]map[StreamID][]grpc.Message
	packetCount   int
}

func New(protoFilename string, servicePort int, guessPaths []string) (_ parser.Parser, err error) {
	protoParser, err := grpc.NewProtoParser(protoFilename)
	if err != nil {
		return
	}

	return &Parser{
		protoFilename: protoFilename,
		servicePort:   servicePort,
		guessPaths:    guessPaths,
		protoParser:   protoParser,
		streams:       map[ConnID]map[StreamID][]grpc.Message{},
	}, nil
}

func (p *Parser) Parse(packet gopacket.Packet) (messages []grpc.Message, err error) {
	p.packetCount++
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
		var message grpc.Message = grpc.Message{
			Meta: grpc.Meta{
				PacketNumber: p.packetCount,
				CaptureInfo:  packet.Metadata().CaptureInfo,
				HTTP2Header:  frame.Header(),
				Src:          packet.NetworkLayer().NetworkFlow().Src().String(),
				Dst:          packet.NetworkLayer().NetworkFlow().Src().String(),
				Sport:        sport,
				Dport:        dport,
			},
			Ext: make(map[grpc.ExtKey]string),
		}
		if _, ok := p.streams[message.ConnID()]; !ok {
			p.streams[message.ConnID()] = make(map[StreamID][]grpc.Message)
		}
		switch frame := frame.(type) {
		case *http2.HeadersFrame:
			payload := map[string]string{}
			headerFields, err := hpackDecode(message.ConnID(), frame)
			if err == nil {
				for _, field := range headerFields {
					payload[field.Name] = field.Value
				}
			} else {
				buf := frame.HeaderBlockFragment()
				for _, field := range hpackDecodePartial(buf) {
					payload[field.Name] = field.Value
				}
				message.Ext[grpc.HeaderPartiallyParsed] = ""
			}

			message.Type = grpc.HeaderType
			message.Header = payload

		case *http2.DataFrame:
			var (
				dataDirection DataDirection
				possiblePaths []string
			)

			dataDirection = ToClient
			if p.servicePort == message.Dport {
				dataDirection = ToService
			}

			searchStream := p.streams[message.ConnID()][streamID]
			if dataDirection == ToClient {
				searchStream = p.streams[message.RevConnID()][streamID]
			}
			for _, msg := range searchStream {
				if msg.Type == grpc.HeaderType {
					for key := range msg.Header {
						if key == ":path" {
							possiblePaths = []string{msg.Header[key]}
						}
					}
				}
			}

			if len(possiblePaths) == 0 {
				message.Ext[grpc.DataGuessed] = ""
				possiblePaths = p.guessPaths
				if len(possiblePaths) == 1 && possiblePaths[0] == "AUTO" {
					possiblePaths = p.protoParser.GetAllPaths()
				}
			}

			msgs := []grpc.Message{}
			for _, possiblePath := range possiblePaths {
				msg, err := p.unmarshalDataFrame(dataDirection, possiblePath, frame)
				if err == nil {
					msgs = append(msgs, msg)
				}
			}

			var msg grpc.Message
			curMax := -1
			for _, m := range msgs {
				var n int
				switch m.Type {
				case grpc.DataType:
					n = len(m.Data.String())
				}
				if n > curMax {
					curMax = n
					msg = m
				}
			}
			message.Type = grpc.DataType
			message.Data = msg.Data
			for k, v := range msg.Ext {
				message.Ext[k] = v
			}

		default:
			continue
		}
		p.streams[message.ConnID()][streamID] = append(p.streams[message.ConnID()][streamID], message)
		messages = append(messages, message)
	}
	return messages, nil
}

func (p *Parser) unmarshalDataFrame(dataDirection DataDirection, path string, frame *http2.DataFrame) (message grpc.Message, err error) {
	message.Ext = map[grpc.ExtKey]string{grpc.DataPath: path}
	message.Data, err = p.protoParser.MarshalResponse(path, frame.Data()[5:])
	if dataDirection == ToService {
		message.Data, err = p.protoParser.MarshalRequest(path, frame.Data()[5:])
	}
	return
}
