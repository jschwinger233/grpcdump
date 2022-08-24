package grpcurl

import (
	"fmt"
	"strings"

	"github.com/jhump/protoreflect/dynamic"
)

type Manager struct {
	pathFilenames map[string]string
}

func New(pathFilenames map[string]string) *Manager {
	pathFiles := map[string]string{}
	for k, v := range pathFilenames {
		pathFiles[strings.TrimPrefix(k, "/")] = v
	}
	return &Manager{pathFilenames: pathFiles}
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
		m.pathFilenames[strings.TrimPrefix(ctx.Path, "/")],
		payload,
		ctx.Dst,
		ctx.Dport,
		strings.TrimPrefix(ctx.Path, "/"),
	), err
}
