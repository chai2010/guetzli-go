// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Guetzli JPEG compressor
//
// Usage:
//
//	guetzli [flags] input_filename output_filename
//
//	  -quality float
//	        Expressed as a JPEG quality value. (default 95)
//	  -version
//	        Show version and exit.
//
// Example:
//
//	guetzli [-quality=95] original.png output.jpg
//	guetzli [-quality=95] original.jpg output.jpg
//
// See https://godoc.org/github.com/chai2010/guetzli-go
// See https://github.com/google/guetzli
//
// Report bugs to <chaishushan{AT}gmail.com>.
//
package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/chai2010/guetzli-go"
)

const Version = "1.0"

var (
	flagQuality = flag.Float64("quality", 95, "Expressed as a JPEG quality value.")
	flagVersion = flag.Bool("version", false, "Show version and exit.")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `Guetzli JPEG compressor.

Usage: guetzli [flags] input_filename output_filename
`)
		flag.PrintDefaults()
		fmt.Println(`
Example:

  guetzli [-quality=95] original.png output.jpg
  guetzli [-quality=95] original.jpg output.jpg

See https://godoc.org/github.com/chai2010/guetzli-go
See https://github.com/google/guetzli

Report bugs to <chaishushan{AT}gmail.com>.
`)
	}
}
func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf("guetzli-%s\n", Version)
		os.Exit(0)
	}

	if flag.NArg() != 2 {
		fmt.Printf("guetzli-%s\n", Version)
		os.Exit(0)
	}

	var (
		inputFilename  = flag.Arg(0)
		outputFilename = flag.Arg(1)
	)

	fin, err := os.Open(inputFilename)
	if err != nil {
		log.Fatalf("open %q failed, err = %v", inputFilename, err)
	}
	defer fin.Close()

	m, _, err := image.Decode(fin)
	if err != nil {
		log.Fatalf("decode %q failed, err = %v", inputFilename, err)
	}

	err = guetzli.Save(outputFilename, m, &guetzli.Options{Quality: float32(*flagQuality)})
	if err != nil {
		log.Fatalf("save %q failed, err = %v", outputFilename, err)
	}
}
