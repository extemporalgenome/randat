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
		input, output string
		n, rep, size  int64
		raw, b64      bool
		code          bool
		cols          uint
		str, esc      bool
	)

	flag.StringVar(&input, "i", "", "input file or stdin if '-' (default is urandom)")
	flag.StringVar(&output, "o", "", "output filename printf-style pattern (default is stdout)")
	flag.Int64Var(&n, "n", 0, "number of (unencoded) bytes to output. Default 16")
	flag.Int64Var(&rep, "r", 0, "number of repititions (output lines). Default 1")
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
			if n <= 0 || rep <= 0 {
				if info, err := f.Stat(); err == nil {
					size = info.Size()
					switch {
					case n <= 0 && rep <= 0:
						n = size
						rep = 1
					case n <= 0:
						n = size / rep
						if size%rep > 0 {
							n++
						}
					case rep <= 0:
						rep = size / n
						if size%rep > 0 {
							rep++
						}
					}
				}
			}
			r = f
		}
	}

	if n == 0 {
		n = 16
	}
	if rep == 0 {
		rep = 1
	}

	var multioutput bool

	if output == "" {
		sink = NopWriteCloser(os.Stdout)
	} else {
		i := 0
		for _, r := range output {
			if r == '%' {
				i++
				break
			}
		}
		switch i {
		case 0:
			if f, err := os.Create(output); err != nil {
				fmt.Fprintln(os.Stderr, "Output file error:", err)
				os.Exit(1)
			} else {
				sink = NopWriteCloser(f)
			}
		case 1:
			multioutput = true
		case 2:
			// ignored for now
			fmt.Fprintf(os.Stderr, "Invalid filename pattern: %q\n", output)
			os.Exit(2)
		}
	}

	for i := int64(0); i != rep; i++ {
		if multioutput {
			if f, err := os.Create(fmt.Sprintf(output, i)); err != nil {
				fmt.Fprintln(os.Stderr, "Output file error:", err)
				os.Exit(1)
			} else {
				sink = f
			}
		}

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
		if !raw && !multioutput {
			fmt.Fprintln(sink)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		if size > 0 {
			size -= n
			if size == 0 {
				break
			} else if size < n {
				n = size
				// one more iteration
				i = rep - 2
			}
		}
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
