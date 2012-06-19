// Copyright 2012 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt" // flag uses fmt anyway
	"io"
	"os"
)

func main() {
	var (
		n        int64
		r        int
		raw, b64 bool
		code     bool
		cols     uint
		str, esc bool
	)

	flag.Int64Var(&n, "n", 16, "number of (unencoded) bytes to output")
	flag.IntVar(&r, "r", 1, "number of repititions (output lines)")
	flag.BoolVar(&raw, "raw", false, "raw output")
	flag.BoolVar(&b64, "64", false, "base64 output")
	flag.BoolVar(&code, "code", false, "comma-separated hex literals (typical array syntax)")
	flag.UintVar(&cols, "cols", 12, "number of columns to use with -code")
	flag.BoolVar(&str, "str", false, "double-quoted string output")
	flag.BoolVar(&esc, "esc", false, "hex-escape all characters (C compat) when using -str")
	flag.Parse()

	var w, sink io.WriteCloser

	sink = NopWriteCloser(os.Stdout)

	for r > 0 {
		switch {
		case raw:
			w = sink
		case code:
			w = NopWriteCloser(NewCodeWriter(sink, cols))
		case str:
			w = NewQuotedWriter(sink, esc)
		case b64:
			w = base64.NewEncoder(base64.StdEncoding, sink)
		default:
			w = NopWriteCloser(NewHexWriter(sink))
		}

		_, err := io.CopyN(w, rand.Reader, n)
		w.Close() // ignoring error
		if !raw {
			fmt.Println()
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		r--
	}
}

var ErrPartialChunkWrite = errors.New("Partial chunk write")

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return writecloser{w}
}

type writecloser struct {
	io.Writer
}

func (w writecloser) Close() error {
	return nil
}
