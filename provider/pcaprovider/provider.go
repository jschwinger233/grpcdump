package pcaprovider

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/jschwinger233/grpcdump/provider"
)

type PcapProvider struct {
	pcapFilename string
}

func New(source string) provider.Provider {
	return &PcapProvider{
		pcapFilename: source,
	}
}

func (p *PcapProvider) PacketStream() (<-chan gopacket.Packet, error) {
	handle, err := pcap.OpenOffline(p.pcapFilename)
	if err != nil {
		return nil, err
	}

	ch := make(chan gopacket.Packet)
	go func() {
		defer close(ch)
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			ch <- packet
		}
	}()

	return ch, nil
}
