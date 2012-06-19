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
		input    string
		n        int64
		rep      int
		raw, b64 bool
		code     bool
		cols     uint
		str, esc bool
	)

	flag.StringVar(&input, "i", "", "input file or stdin if '-' (default is urandom)")
	flag.Int64Var(&n, "n", 16, "number of (unencoded) bytes to output")
	flag.IntVar(&rep, "r", 1, "number of repititions (output lines)")
	flag.BoolVar(&raw, "raw", false, "raw output")
	flag.BoolVar(&b64, "64", false, "base64 output")
	flag.BoolVar(&code, "code", false, "comma-separated hex literals (typical array syntax)")
	flag.UintVar(&cols, "cols", 12, "number of columns to use with -code")
	flag.BoolVar(&str, "str", false, "double-quoted string output")
	flag.BoolVar(&esc, "esc", false, "hex-escape all characters (C compat) when using -str")
	flag.Parse()

	var (
		w, sink io.WriteCloser
		r       io.Reader
	)

	switch input {
	case "":
		r = rand.Reader
	case "-":
		r = os.Stdin
	default:
		if f, err := os.Open(input); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else {
			if n < 0 {
				if info, err := f.Stat(); err == nil {
					n = info.Size()
				}
			}
			r = f
		}
	}

	sink = NopWriteCloser(os.Stdout)

	for rep > 0 {
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

		_, err := io.CopyN(w, r, n)
		w.Close() // ignoring error
		if !raw {
			fmt.Println()
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		rep--
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
