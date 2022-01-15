package grpcurl

import (
	"fmt"
	"strings"

	"github.com/jhump/protoreflect/dynamic"
)

type Manager struct {
	protoFilename string
}

func New(protoFilename string) *Manager {
	return &Manager{protoFilename: protoFilename}
}

type RenderContext struct {
	Payload *dynamic.Message
	Dst     string
	Dport   int
	Path    string
}

func (m *Manager) Render(ctx RenderContext) (cmd string, err error) {
	payload, err := ctx.Payload.MarshalJSON()
	return fmt.Sprintf(
		"grpcurl -plaintext -proto %s -d '%s' %s:%d %s",
		m.protoFilename,
		payload,
		ctx.Dst,
		ctx.Dport,
		strings.TrimPrefix(ctx.Path, "/"),
	), err
}
