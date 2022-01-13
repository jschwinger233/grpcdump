package grpcparser

import (
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

func hpackDecodePartial(p []byte) []hpack.HeaderField {
	var hf []hpack.HeaderField
	decoder := hpack.NewDecoder(65536, func(f hpack.HeaderField) { hf = append(hf, f) })
	decoder.Write(p)
	decoder.Close()
	return hf
}

var hdecs map[string]*hpack.Decoder = make(map[string]*hpack.Decoder)

func hpackDecode(connID string, hf *http2.HeadersFrame) ([]hpack.HeaderField, error) {
	if _, ok := hdecs[connID]; !ok {
		hdecs[connID] = hpack.NewDecoder(65536, nil)
	}
	hdec := hdecs[connID]
	return hdec.DecodeFull(hf.HeaderBlockFragment())
}
