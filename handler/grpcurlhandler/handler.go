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

	case grpc.RequestType:
		// grpcurl -plaintext -proto rpc/gen/core.proto -d '{"appname":"zc","entrypoint":"zc"}' localhost:5001 pb.CoreRPC/WorkloadStatusStream
		bytes, err := msg.Request.MarshalJSON()
		if err != nil {
			return err
		}
		fmt.Printf("grpcurl -plaintext -proto %s -d '%s' %s:%d %s\n", h.protoFilename, string(bytes), msg.Dst, msg.Dport, strings.TrimPrefix(msg.Ext[grpc.DataPath], "/"))

	case grpc.ResponseType:
		guessIndicator := ""
		if _, ok := msg.Ext[grpc.DataGuessed]; ok {
			guessIndicator = "(guess)"
		}
		fmt.Printf("%s\t%s\tstreamid:%d\tdata:%s%s\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID, guessIndicator, msg.Response.String())

	case grpc.HeaderType:
		partialIndicator := ""
		if _, ok := msg.Ext[grpc.HeaderPartiallyParsed]; ok {
			partialIndicator = "(partial)"
		}
		fmt.Printf("%s\t%s\tstreamid:%d\theader:%s%+v\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID, partialIndicator, msg.Header)

	case grpc.UnknownType:
		fmt.Printf("%s\t%s\tstreamid:%d\tunknown data frame\n", msg.CaptureInfo.Timestamp, msg.ConnID(), msg.HTTP2Header.StreamID, msg.Header)
	}
	return
}
