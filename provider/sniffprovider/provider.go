package sniffprovider

import (
	"log"
	"strconv"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/jschwinger23/grpcdump/provider"
	"golang.org/x/net/bpf"
)

type SniffProvider struct {
	Device string
	Port   int
}

func New(source string) provider.Provider {
	parts := strings.Split(source, ":")
	if len(parts) != 2 {
		log.Fatalf("invalid sniff source: %s", source)
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatalf("invalid sniff port: %s", parts[1])
	}
	return &SniffProvider{
		Device: parts[0],
		Port:   port,
	}
}

func (p *SniffProvider) PacketStream() (_ <-chan gopacket.Packet, err error) {
	handler, err := pcapgo.NewEthernetHandle(p.Device)
	if err != nil {
		return nil, err
	}
	err = handler.SetBPF(makeBPFFilter(p.Port))
	packetSource := gopacket.NewPacketSource(handler, layers.LayerTypeEthernet)
	ch := make(chan gopacket.Packet)
	go func() {
		defer close(ch)
		var lastId *uint16
		for packet := range packetSource.Packets() {
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer == nil {
				continue
			}
			ip, _ := ipLayer.(*layers.IPv4)
			if lastId != nil && *lastId == ip.Id {
				continue
			}
			lastId = &ip.Id
			ch <- packet
		}
	}()
	return ch, nil
}

func makeBPFFilter(port int) (filter []bpf.RawInstruction) {
	p := uint32(port)
	return []bpf.RawInstruction{
		{Op: 40, Jt: 0, Jf: 0, K: 12},
		{Op: 21, Jt: 0, Jf: 6, K: 34525},
		{Op: 48, Jt: 0, Jf: 0, K: 20},
		{Op: 21, Jt: 0, Jf: 15, K: 6},
		{Op: 40, Jt: 0, Jf: 0, K: 54},
		{Op: 21, Jt: 12, Jf: 0, K: p},
		{Op: 40, Jt: 0, Jf: 0, K: 56},
		{Op: 21, Jt: 10, Jf: 11, K: p},
		{Op: 21, Jt: 0, Jf: 10, K: 2048},
		{Op: 48, Jt: 0, Jf: 0, K: 23},
		{Op: 21, Jt: 0, Jf: 8, K: 6},
		{Op: 40, Jt: 0, Jf: 0, K: 20},
		{Op: 69, Jt: 6, Jf: 0, K: 8191},
		{Op: 177, Jt: 0, Jf: 0, K: 14},
		{Op: 72, Jt: 0, Jf: 0, K: 14},
		{Op: 21, Jt: 2, Jf: 0, K: p},
		{Op: 72, Jt: 0, Jf: 0, K: 16},
		{Op: 21, Jt: 0, Jf: 1, K: p},
		{Op: 6, Jt: 0, Jf: 0, K: 1000},
		{Op: 6, Jt: 0, Jf: 0, K: 0},
	}

}
