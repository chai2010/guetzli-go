# [Guetzli](https://github.com/google/guetzli) perceptual JPEG encoder for Go

- `guetzli-go`(package): [![GoDoc](https://godoc.org/github.com/chai2010/guetzli-go?status.svg)](https://godoc.org/github.com/chai2010/guetzli-go)
- `guetzli`(command): [![GoDoc](https://godoc.org/github.com/chai2010/guetzli-go/apps/guetzli?status.svg)](https://godoc.org/github.com/chai2010/guetzli-go/apps/guetzli)


Install
=======

Install `GCC` or `MinGW` ([download here](http://tdm-gcc.tdragon.net/download)) at first,
and then run these commands:

1. `go get github.com/chai2010/guetzli-go`
1. `go get github.com/chai2010/guetzli-go/apps/guetzli`
1. `go run hello.go`
1. `guetzli -h`


Command: [guetzli](apps/guetzli/main.go)
========================================

```
Guetzli JPEG compressor

Usage:

    guetzli [flags] input_filename output_filename
    guetzli [flags] input_dir output_dir [ext...]

      -quality int
            Expressed as a JPEG quality value(>=84 and <= 110). (default 84)
      -regexp string
            regexp for base filename.
      -version
            Show version and exit.

Example:

    guetzli [-quality=95] original.png output.jpg
    guetzli [-quality=95] original.jpg output.jpg

    guetzli [-quality=95] input_dir output_dir .png
    guetzli [-quality=95] input_dir output_dir .png .jpg .jpeg
    guetzli [-quality=95 -regexp="^\d+"] input_dir output_dir .png

Note: Default image ext is: .jpeg .jpg .png

Note: Supported formats: .gif, .jpeg, .jpg, .png

See https://godoc.org/github.com/chai2010/guetzli-go See
https://github.com/google/guetzli

Report bugs to <chaishushan{AT}gmail.com>.
```

Build more image format support(use "all_formats" tag):

```
go get   -tags="all_formats" github.com/chai2010/guetzli-go/apps/guetzli
go build -tags="all_formats"
```

See [main_all_formats.go](apps/guetzli/main_all_formats.go)

Benchmark
=========

```
go test -test.bench=.*
BenchmarkEncode_guetzli_quality_95-4   	       1	1043846782 ns/op
BenchmarkEncode_jpeg_quality_95-4      	    2000	   1130081 ns/op
BenchmarkEncode_png-4                  	     200	   6423543 ns/op
PASS
ok  	github.com/chai2010/guetzli-go	5.817s
```

Example
=======

This is a simple example:

```Go
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/chai2010/guetzli-go"
)

func main() {
	m0 := loadImage("./testdata/video-001.png")

	data1 := jpegEncode(m0, 95)
	data2 := guetzliEncode(m0, 95)

	fmt.Println("jpeg encoded size:", len(data1))
	fmt.Println("guetzli encoded size:", len(data2))

	if err := ioutil.WriteFile("a.out.jpeg", data1, 0666); err != nil {
		log.Println(err)
	}
	if err := ioutil.WriteFile("a.out.guetzli.jpeg", data2, 0666); err != nil {
		log.Println(err)
	}

	fmt.Println("Done")
}

func loadImage(name string) image.Image {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	m, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func jpegEncode(m image.Image, quality int) []byte {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, m, &jpeg.Options{Quality: quality})
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func guetzliEncode(m image.Image, quality float32) []byte {
	var buf bytes.Buffer
	err := guetzli.Encode(&buf, m, &guetzli.Options{Quality: quality})
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
```

BUGS
====

Report bugs to <chaishushan@gmail.com>.

Thanks!
