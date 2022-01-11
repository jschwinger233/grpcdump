package grpcparser

import (
	"golang.org/x/net/http2/hpack"
)

func hpackDecodePartial(p []byte) []hpack.HeaderField {
	var hf []hpack.HeaderField
	decoder := hpack.NewDecoder(4096, func(f hpack.HeaderField) { hf = append(hf, f) })
	decoder.Write(p)
	decoder.Close()
	return hf
}
