package main

import (
	"encoding/hex"
	"io"
)

func NewHexWriter(w io.Writer) io.Writer {
	return &hexwriter{w: w}
}

type hexwriter struct {
	w io.Writer
	b []byte
}

func (h *hexwriter) Write(data []byte) (n int, err error) {
	length := hex.EncodedLen(len(data))
	if length > cap(h.b) {
		h.b = make([]byte, length)
	}
	b := h.b[:length]
	hex.Encode(b, data)
	n, err = h.w.Write(b)
	if n != length {
		if err == nil {
			if n%2 > 0 {
				err = ErrPartialChunkWrite
			} else {
				err = io.ErrShortWrite
			}
		}
	}
	n /= 2
	return
}
