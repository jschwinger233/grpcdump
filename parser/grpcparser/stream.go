package grpcparser

import grpc "github.com/jschwinger23/grpcdump/grpchelper"

type (
	StreamID = uint32
	ConnID   = string
)

type HTTP2Stream struct {
	streams map[ConnID]map[StreamID][]grpc.Message
}

func (s *HTTPStream) DeleteStream(segment TCPSegment, id StreamID) {
	delete(s.streams[segment.ConnID()], id)
}
