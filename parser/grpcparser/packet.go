package grpcparser

import (
	"fmt"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type TCPSegment struct {
	gopacket.Packet
}

func (s TCPSegment) HasApplicationLayer() bool {
	return s.ApplicationLayer() != nil
}

func (s TCPSegment) ConnID() string {
	return fmt.Sprintf(
		"%s:%d->%s:%d",
		s.Src(),
		s.Sport(),
		s.Dst(),
		s.Dport(),
	)
}

func (s TCPSegment) RevConnID() string {
	return fmt.Sprintf(
		"%s:%d->%s:%d",
		s.Dst(),
		s.Dport(),
		s.Src(),
		s.Sport(),
	)
}

func (s TCPSegment) Src() string {
	return s.NetworkLayer().NetworkFlow().Src().String()
}

func (s TCPSegment) Dst() string {
	return s.NetworkLayer().NetworkFlow().Dst().String()
}

func (s TCPSegment) Sport() int {
	sport, _ := strconv.Atoi(s.TransportLayer().TransportFlow().Src().String())
	return sport
}

func (s TCPSegment) Dport() int {
	dport, _ := strconv.Atoi(s.TransportLayer().TransportFlow().Dst().String())
	return dport
}

func (s TCPSegment) FIN() bool {
	tcp, _ := s.Layer(layers.LayerTypeTCP).(*layers.TCP)
	return tcp.FIN
}

func (s TCPSegment) Payload() (payload []byte) {
	appLayer := s.ApplicationLayer()
	if appLayer == nil {
		return
	}
	return appLayer.Payload()
}
