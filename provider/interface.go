package provider

import "github.com/google/gopacket"

type Provider interface {
	PacketStream() <-chan gopacket.Packet
}
