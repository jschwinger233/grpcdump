package grpchelper

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

type ProtoParser interface {
	MarshalRequest(path string, message []byte) (*dynamic.Message, error)
	MarshalResponse(path string, message []byte) (*dynamic.Message, error)
	GetAllPaths() []string
}

type parser struct {
	Requests  map[string]*desc.MessageDescriptor
	Responses map[string]*desc.MessageDescriptor
}

func NewProtoParser(filenames []string) (_ ProtoParser, err error) {
	requests := make(map[string]*desc.MessageDescriptor)
	response := make(map[string]*desc.MessageDescriptor)

	for _, filename := range filenames {
		if filename, err = filepath.Abs(filename); err != nil {
			return
		}
		dir, base := filepath.Dir(filename), filepath.Base(filename)
		fileNames, err := protoparse.ResolveFilenames([]string{dir}, base)
		if err != nil {
			return nil, err
		}
		p := protoparse.Parser{
			ImportPaths:           []string{dir},
			IncludeSourceCodeInfo: true,
		}
		parsedFiles, err := p.ParseFiles(fileNames...)
		if err != nil {
			return nil, err
		}

		if len(parsedFiles) < 1 {
			err = errors.New("proto file not found")
			return nil, err
		}

		for _, parsedFile := range parsedFiles {
			for _, service := range parsedFile.GetServices() {
				serviceName := fmt.Sprintf("%s.%s", parsedFile.GetPackage(), service.GetName())
				for _, method := range service.GetMethods() {
					path := fmt.Sprintf("/%s/%s", serviceName, method.GetName())
					requests[path] = method.GetInputType()
					response[path] = method.GetOutputType()
				}
			}
		}
	}
	return &parser{
		Requests:  requests,
		Responses: response,
	}, nil
}

func (p *parser) MarshalRequest(path string, message []byte) (*dynamic.Message, error) {
	descriptor, ok := p.Requests[path]
	if !ok {
		return nil, fmt.Errorf("path not found: %s", path)
	}
	msg := dynamic.NewMessage(descriptor)
	return msg, msg.Unmarshal(message)
}

func (p *parser) MarshalResponse(path string, message []byte) (*dynamic.Message, error) {
	descriptor, ok := p.Responses[path]
	if !ok {
		return nil, fmt.Errorf("path not found: %s", path)
	}
	msg := dynamic.NewMessage(descriptor)
	return msg, msg.Unmarshal(message)
}

func (p *parser) GetAllPaths() (paths []string) {
	for path := range p.Requests {
		paths = append(paths, path)
	}
	return
}
