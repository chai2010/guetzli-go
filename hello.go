// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

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

func guetzliEncode(m image.Image, quality int) []byte {
	var buf bytes.Buffer
	err := guetzli.Encode(&buf, m, &guetzli.Options{Quality: quality})
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
