package grpcurl

import (
	"fmt"
	"strings"

	"github.com/jhump/protoreflect/dynamic"
)

type Manager struct {
	protoFilenames []string
}

func New(protoFilenames []string) *Manager {
	return &Manager{protoFilenames: protoFilenames}
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
		m.protoFilenames[0], // TODO@gray
		payload,
		ctx.Dst,
		ctx.Dport,
		strings.TrimPrefix(ctx.Path, "/"),
	), err
}
