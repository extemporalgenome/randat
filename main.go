// Copyright 2012 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt" // flag uses fmt anyway
	"io"
	"os"
)

func main() {
	const maxbuf = 4096
	var (
		n   int64
		r   int
		b64 bool
		raw bool
	)

	flag.Int64Var(&n, "n", 16, "number of (unencoded) bytes to output")
	flag.IntVar(&r, "r", 1, "number of repititions (output lines)")
	flag.BoolVar(&raw, "raw", false, "raw output")
	flag.BoolVar(&b64, "64", false, "base64 output")
	flag.Parse()

	var w io.WriteCloser = NopWriteCloser(os.Stdout)

	if !raw {
		if b64 {
			w = base64.NewEncoder(base64.StdEncoding, w)
		} else {
			w = NopWriteCloser(NewHexWriter(w))
		}
	}

	for r > 0 {
		_, err := io.CopyN(w, rand.Reader, n)
		w.Close() // ignoring error
		fmt.Println()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		r--
	}
}

type writecloser struct {
	io.Writer
}

func (w writecloser) Close() error {
	return nil
}

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return writecloser{w}
}

type hexwriter struct {
	w io.Writer
	n int
	b []byte
}

var ErrOddSizeWrite = errors.New("Odd sized write")

func NewHexWriter(w io.Writer) io.Writer {
	return &hexwriter{w: w}
}

func (h *hexwriter) Write(data []byte) (n int, err error) {
	if len(data) > h.n {
		h.n = len(data)
		h.b = make([]byte, hex.EncodedLen(h.n))
	}
	hex.Encode(h.b, data)
	n, err = h.w.Write(h.b)
	if n != len(h.b) {
		if err == nil {
			if n%2 > 0 {
				err = ErrOddSizeWrite
			} else {
				err = io.ErrShortWrite
			}
		}
	}
	n = hex.DecodedLen(n)
	return
}
