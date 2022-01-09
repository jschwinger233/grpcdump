package grpchelper

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

type (
	ServiceName = string
	MethodName  = string
)

type ProtoParser interface {
	MarshalRequest(service, method string, message []byte) (*dynamic.Message, error)
	MarshalResponse(service, method string, message []byte) (*dynamic.Message, error)
}

type parser struct {
	Requests  map[ServiceName]map[MethodName]*desc.MessageDescriptor
	Responses map[ServiceName]map[MethodName]*desc.MessageDescriptor
}

func NewProtoParser(filename string) (_ ProtoParser, err error) {
	requests := make(map[ServiceName]map[MethodName]*desc.MessageDescriptor)
	response := make(map[ServiceName]map[MethodName]*desc.MessageDescriptor)

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
			serviceName := "pb." + service.GetName()
			requests[serviceName] = make(map[MethodName]*desc.MessageDescriptor)
			response[serviceName] = make(map[MethodName]*desc.MessageDescriptor)
			for _, method := range service.GetMethods() {
				requests[serviceName][method.GetName()] = method.GetInputType()
				response[serviceName][method.GetName()] = method.GetOutputType()
			}
		}
	}
	return &parser{
		Requests:  requests,
		Responses: response,
	}, nil
}

func (p *parser) MarshalRequest(service, method string, message []byte) (*dynamic.Message, error) {
	methods, ok := p.Requests[service]
	if !ok {
		return nil, fmt.Errorf("service not found: %s", service)
	}
	descriptor, ok := methods[method]
	if !ok {
		return nil, fmt.Errorf("method not found: %s", method)
	}
	msg := dynamic.NewMessage(descriptor)
	return msg, msg.Unmarshal(message)
}

func (p *parser) MarshalResponse(service, method string, message []byte) (*dynamic.Message, error) {
	methods, ok := p.Responses[service]
	if !ok {
		return nil, fmt.Errorf("service not found: %s", service)
	}
	descriptor, ok := methods[method]
	if !ok {
		return nil, fmt.Errorf("method not found: %s", method)
	}
	msg := dynamic.NewMessage(descriptor)
	return msg, msg.Unmarshal(message)
}
