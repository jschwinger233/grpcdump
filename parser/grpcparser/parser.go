package grpcparser

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/google/gopacket"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jschwinger23/grpcdump/parser"
	"golang.org/x/net/http2"
)

type (
	ServiceName = string
	MethodName  = string
)

type Parser struct {
	protoFilename string
	guessMethod   string

	grpcInputTypes  map[ServiceName]map[MethodName]*desc.MessageDescriptor
	grpcOutputTypes map[ServiceName]map[MethodName]*desc.MessageDescriptor
}

func New(protoFilename, guessMethod string) (_ parser.Parser, err error) {
	p := &Parser{
		protoFilename: protoFilename,
		guessMethod:   guessMethod,
	}
	p.grpcInputTypes, p.grpcOutputTypes, err = loadProto(protoFilename)
	return nil, err
}

func loadProto(filename string) (input, output map[ServiceName]map[MethodName]*desc.MessageDescriptor, err error) {
	input = make(map[ServiceName]map[MethodName]*desc.MessageDescriptor)
	output = make(map[ServiceName]map[MethodName]*desc.MessageDescriptor)

	if filename, err = filepath.Abs(filename); err != nil {
		return
	}
	dir, base := filepath.Dir(filename), filepath.Base(filename)
	fileNames, err := protoparse.ResolveFilenames([]string{dir}, base)
	if err != nil {
		return
	}
	p := protoparse.Parser{
		ImportPaths:           []string{dir},
		IncludeSourceCodeInfo: true,
	}
	parsedFiles, err := p.ParseFiles(fileNames...)
	if err != nil {
		return
	}

	if len(parsedFiles) < 1 {
		err = errors.New("proto file not found")
		return
	}

	for _, parsedFile := range parsedFiles {
		for _, service := range parsedFile.GetServices() {
			input[service.GetName()] = make(map[MethodName]*desc.MessageDescriptor)
			output[service.GetName()] = make(map[MethodName]*desc.MessageDescriptor)
			for _, method := range service.GetMethods() {
				input[service.GetName()][method.GetName()] = method.GetInputType()
				output[service.GetName()][method.GetName()] = method.GetOutputType()
			}
		}
	}

}

func (p *Parser) Parse(packet gopacket.Packet) (msg dynamic.Message, err error) {
	payload := packet.ApplicationLayer().Payload()
	framer := http2.NewFramer(ioutil.Discard, bytes.NewReader(payload))
	frame, err := framer.ReadFrame()
	if err != nil {
		return
	}

	_, ok := frame.(*http2.DataFrame)
	if !ok {
		err = errors.New("failed to cast type from http.Frame to http2.DataFrame")
		return
	}

	return
}
