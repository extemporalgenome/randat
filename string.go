package main

import (
	"encoding/hex"
	"io"
	"unicode"
)

func NewQuotedWriter(w io.Writer, esc bool) io.WriteCloser {
	return &quotedwriter{w: w, esc: esc}
}

type quotedwriter struct {
	w        io.Writer
	b        []byte
	esc, add bool
}

func (q *quotedwriter) Write(data []byte) (n int, err error) {
	b := q.b
	if !q.add {
		b = append(b, '"')
		q.add = true
	}
	for _, c := range data {
		switch {
		case q.esc, c > unicode.MaxASCII:
			fallthrough
		default:
			x := []byte{'\\', 'x', 0, c}
			hex.Encode(x[2:], x[3:])
			b = append(b, x...)
		case c == '"':
			b = append(b, '\\', '"')
		case c == '\\':
			b = append(b, '\\', '\\')
		case unicode.IsPrint(rune(c)):
			b = append(b, c)
		}
	}
	n, err = q.w.Write(b)
	q.b = b[:0]
	if err == nil {
		if n != len(b) {
			err = io.ErrShortWrite
		} else {
			n = len(data)
		}
	}
	return
}

func (q *quotedwriter) Close() error {
	q.add = false
	_, err := q.w.Write([]byte{'"'})
	return err
}
