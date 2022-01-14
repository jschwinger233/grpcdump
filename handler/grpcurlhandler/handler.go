package grpcurlhandler

import (
	"fmt"
	"strings"

	grpc "github.com/jschwinger23/grpcdump/grpchelper"
	"github.com/jschwinger23/grpcdump/handler"
)

type GrpcurlHandler struct {
	protoFilename string
}

func New(protoFilename string) handler.GrpcHandler {
	return &GrpcurlHandler{protoFilename}
}

func (h *GrpcurlHandler) Handle(msg grpc.Message) (err error) {
	// time, conn, streamid, data
	switch msg.Type {

	case grpc.DataType:
		// grpcurl -plaintext -proto rpc/gen/core.proto -d '{"appname":"zc","entrypoint":"zc"}' localhost:5001 pb.CoreRPC/WorkloadStatusStream
		bytes, err := msg.Data.MarshalJSON()
		if err != nil {
			return err
		}
		fmt.Printf("grpcurl -plaintext -proto %s -d '%s' %s:%d %s\n", h.protoFilename, string(bytes), msg.Dst, msg.Dport, strings.TrimPrefix(msg.Ext[grpc.DataPath], "/"))
	}
	return
}
