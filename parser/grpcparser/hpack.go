package grpcparser

import (
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

type HpackDecoder struct {
	hdecs map[string]*hpack.Decoder
}

func newHpackDecoder() *HpackDecoder {
	return &HpackDecoder{
		hdecs: make(map[string]*hpack.Decoder),
	}
}

func (h *HpackDecoder) DecodePartial(p []byte) (hf []hpack.HeaderField) {
	decoder := hpack.NewDecoder(65536, func(f hpack.HeaderField) { hf = append(hf, f) })
	decoder.Write(p)
	decoder.Close()
	return hf
}

func (h *HpackDecoder) Decode(connID string, hf *http2.HeadersFrame) ([]hpack.HeaderField, error) {
	if _, ok := h.hdecs[connID]; !ok {
		h.hdecs[connID] = hpack.NewDecoder(65536, nil)
	}
	hdec := h.hdecs[connID]
	return hdec.DecodeFull(hf.HeaderBlockFragment())
}

func (h *HpackDecoder) Clear(connID string) {
	delete(h.hdecs, connID)
}
