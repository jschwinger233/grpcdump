package grpcparser

import (
	"fmt"
	"strings"

	"github.com/google/gopacket"
	"golang.org/x/net/http2/hpack"
)

func getConnID(packet gopacket.Packet) (id string) {
	netFlow := packet.NetworkLayer().NetworkFlow()
	transFlow := packet.TransportLayer().TransportFlow()
	src := fmt.Sprintf("%s:%s", netFlow.Src().String(), transFlow.Src().String())
	dst := fmt.Sprintf("%s:%s", netFlow.Dst().String(), transFlow.Dst().String())
	return fmt.Sprintf("%s-%s", src, dst)
}

func revConnID(id string) string {
	parts := strings.Split(id, "-")
	return fmt.Sprintf("%s-%s", parts[1], parts[0])
}

func hpackDecodePartial(p []byte) []hpack.HeaderField {
	var hf []hpack.HeaderField
	decoder := hpack.NewDecoder(4096, func(f hpack.HeaderField) { hf = append(hf, f) })
	decoder.Write(p)
	decoder.Close()
	return hf
}
