// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go test -test.bench=.*

package guetzli

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"testing"
)

func BenchmarkEncode_guetzli_quality_95(b *testing.B) {
	m, err := loadImage("./testdata/video-001.png")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, ok := EncodeImage(m, 95)
		if !ok {
			b.Fatal("EncodeImage failed")
		}
		_ = data
	}
}

func BenchmarkEncode_jpeg_quality_95(b *testing.B) {
	m, err := loadImage("./testdata/video-001.png")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := jpeg.Encode(&buf, m, &jpeg.Options{Quality: 95})
		if err != nil {
			b.Fatal(err)
		}
		_ = buf.Bytes()
	}
}
func BenchmarkEncode_png(b *testing.B) {
	m, err := loadImage("./testdata/video-001.png")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := png.Encode(&buf, m)
		if err != nil {
			b.Fatal(err)
		}
		_ = buf.Bytes()
	}
}
