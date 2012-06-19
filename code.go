package main

import (
	"encoding/hex"
	"io"
)

func NewCodeWriter(w io.Writer, cols uint) io.Writer {
	return &codewriter{w: w, cols: cols, rem: cols}
}

type codewriter struct {
	w         io.Writer
	cols, rem uint
}

func (c *codewriter) Write(data []byte) (n int, err error) {
	var (
		buf     = []byte(`0x00, `)
		i, j, m int
		k       = len(buf)
		b       = make([]byte, len(data)*k)
	)
	for i, buf[3] = range data {
		hex.Encode(buf[2:4], buf[3:4])
		if c.cols > 0 {
			c.rem--
			if c.rem == 0 {
				buf[5] = '\n'
				c.rem = c.cols
			} else {
				buf[5] = ' '
			}
		}
		j = i * k
		copy(b[j:j+k], buf)
	}
	m, err = c.w.Write(b)
	n = m / k
	if err == nil && m%k > 0 {
		err = ErrPartialChunkWrite
	}
	return
}
