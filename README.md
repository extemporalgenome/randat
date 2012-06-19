# randat (random data)

randat outputs random bytes (from /dev/urandom) in raw or encoded form
(hex encoded is the default).

Basic usage:

	$ randat
	011648b9bf1354bf788cb0c57285e939


## Encoding options

	-raw  raw binary
	-64   base64
	-hex  hex encoded (default)
	-str  C style strings (pass -esc for C compatibility).
	      Should work syntactically with Go, Python, Javascript,
	      and languages with similar conventions (that support
	      \x00 style escapes).
	-code comma-separated hex integer literals, suitable for
	      inclusion in array/list literals. Use -cols to specify
		  the number of elements per output line (set to zero to
		  disable multiline output).

## Output control options

	-r    number of repititions (input chunks/output lines,
	      generally). -raw does not separate repititions with
		  newlines. Normally defaults to 1. If negative, a nearly
		  endless (~2^63) number of lines will be output.
	-n    size of each input chunk, in bytes. Normally defaults
	      to 16.
	-o    A filename, or optionally, a printf-style filename
	      pattern. If a pattern, it should have one Go format
		  specifier suitable for an int. Ex: `-o output.%03d.txt`.
		  If a pattern is given, one file will be created for each
		  repitition. If unspecified, stdout is used.

## Input options

	-i    input source. If left blank, it defaults to /dev/random.
	      '-' is stdin. Anything else is a file.

When `-i` is given a file with a detectable size (e.g. a regular file),
`-r` and `-n` take on new behaviors: -r represents the maximum number of
chunks (or output files if `-o` is passed a pattern) that will be
produced. If negative, `-r` will produce as many chunks/files as needed
to output the entire input file with the given chunk size of `-n`.

If `-n` is negative or not specified, the chunk size will be auto-sized
produce no more than `-r` chunks/files. If both are negative or
unspecified (default), the entire file will be output in one chunk/file.

## Recipes

Given that arbitrary (non)random input can be specified, `randat` can be
used to encode normal files (like image files) in any of the above
formats, or can act as a cat/dd replacement. `-raw` output is as fast as
dd/cat anyway. Equivalences:

	cat file
	randat -i file -raw

	dd if=infile of=outfile bs=4096 count=1024 -raw
	randat -i infile -o outfile -n 4096 -r 1024

Note that randat may not read/write or buffer in the specified `-n` and
`-r` sizes -- unlike dd, randat's sizes are semantic, not literal.

### Quoted strings

10 (visually aligned) random 32-byte strings:

	randat -r 10 -n 32 -str -esc

	"\xd9\xcd\x2a\x7b\x75\x6a\x6c\x2e\x33\xa9\x8d\x96\x0a\x9a\x74\x98"
	"\x59\x83\x88\x52\xab\xeb\xf2\x68\x8f\xb1\x15\x44\x73\xc6\x49\x05"
	"\x74\x5b\xa3\xe4\xb7\x82\x77\x21\x8d\x6f\x5b\x49\x61\x5d\xe1\x4d"
	"\x05\xce\x1f\x15\x9b\x78\x44\xf8\x43\x87\x7d\x78\x71\x8c\x1f\xdf"
	"\x2a\xde\xc0\xe4\x72\x3b\x8a\x54\x22\x4d\x48\xff\x2f\xbb\x1f\x07"
	"\x61\x42\x86\x59\xff\x63\x1c\x4f\xb0\x5e\x83\xbc\x21\x5c\xcc\xf6"
	"\xc9\xec\xd3\xc7\xa7\x3a\x4b\xa8\x70\x2f\x3d\xe8\xf9\xc7\xe0\x16"
	"\xb1\x3d\xe2\x01\xc5\x97\xd8\x7d\x35\x6d\x2a\xd5\x3c\x44\xc6\x47"
	"\xb5\x4e\x2f\xae\xc4\xbf\x73\xec\x24\xdc\x47\x81\xcb\x25\x4f\x33"
	"\xd0\x10\xf9\xa1\x94\xc5\x3b\x9e\x49\xe3\x52\xd7\x1c\xf2\x46\x47"

### Array literal (kind of) output

Encoding a JPG in C/Go/Python style array/list literals:

	randat -code -i image.jpg

	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x16,
	0x02, 0x00, 0x00, 0x00, 0x00, 0x70, 0xd4, 0x81, 0x64, 0x00, 0x00, 0x00,
	0x04, 0x67, 0x41, 0x4d, 0x41, 0x00, 0x00, 0xb1, 0x8f, 0x0b, 0xfc, 0x61,
	0x05, 0x00, 0x00, 0x00, 0x01, 0x73, 0x52, 0x47, 0x42, 0x00, 0xae, 0xce,
	0x1c, 0xe9, 0x00, 0x00, 0x00, 0x20, 0x63, 0x48, 0x52, 0x4d, 0x00, 0x00,
	0x7a, 0x26, 0x00, 0x00, 0x80, 0x84, 0x00, 0x00, 0xfa, 0x00, 0x00, 0x00,
	0x80, 0xe8, 0x00, 0x00, 0x75, 0x30, 0x00, 0x00, 0xea, 0x60, 0x00, 0x00,
	0x3a, 0x98, 0x00, 0x00, 0x17, 0x70, 0x9c, 0xba, 0x51, 0x3c, 0x00, 0x00,
	0x00, 0x09, 0x70, 0x48, 0x59, 0x73, 0x00, 0x00, 0x0b, 0x13, 0x00, 0x00,
	0x0b, 0x13, 0x01, 0x00, 0x9a, 0x9c, 0x18, 0x00, 0x00, 0x00, 0x69, 0x49,
	0x44, 0x41, 0x54, 0x08, 0x5b, 0x63, 0xf8, 0xf5, 0xff, 0xff, 0x1f, 0x86,
	0x5f, 0xeb, 0x81, 0xc4, 0xbf, 0xa5, 0xeb, 0x81, 0xc4, 0x97, 0x30, 0x10,
	0xa1, 0xb7, 0x87, 0xe1, 0x1f, 0x83, 0xf6, 0x1c, 0x86, 0xbf, 0x21, 0x59,
	0x2b, 0x18, 0x7e, 0x6a, 0x46, 0x2d, 0x61, 0x38, 0xba, 0x6a, 0xea, 0x14,
	0x86, 0x83, 0x89, 0x3f, 0x42, 0x18, 0x9e, 0x46, 0xbd, 0xe3, 0x60, 0xf8,
	0x22, 0xbd, 0xbf, 0x83, 0xe1, 0xd9, 0xaf, 0xff, 0x6b, 0x18, 0x80, 0x06,
	0xec, 0x81, 0x10, 0xff, 0xc0, 0xc4, 0xbf, 0xfd, 0x20, 0x16, 0x90, 0xf8,
	0xbf, 0xff, 0x7f, 0x0e, 0xc3, 0xff, 0xf9, 0x7f, 0x81, 0xac, 0x55, 0xd3,
	0xfe, 0x30, 0xfc, 0xff, 0xff, 0xfa, 0x0f, 0x00, 0x30, 0xdb, 0x44, 0x1a,
	0x1c, 0x44, 0x44, 0xc8, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
	0xae, 0x42, 0x60, 0x82,

Surrounding that output block with curly braces (C array, Go array/slice) or
square brackets (Javascript array, Python list) will form a valid literal.

### Splitting archives

Splitting an input file into size-limited output files:

	randat -raw -n 0x28000000 -i somebigfile -o 640mb-file.%03d.raw

This will produce files named `640mb-file.000.raw`, `640mb-file.001.raw`, etc.

Splitting an input file into at most `-r` output files:

	randat -raw -r 10 -i somebigfile -o file%02d.raw

This will produce files `file00.raw` through `file09.raw`.
