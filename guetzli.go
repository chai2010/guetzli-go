// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Guetzli perceptual JPEG encoder for Go.
package guetzli // import "github.com/chai2010/guetzli-go"

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"
	"os"
	"reflect"
)

const DefaultQuality = 75

var errEncodeFailed = errors.New("guetzli: encode failed!")

type Options struct {
	Quality float32
}

// MemP Image Spec (Native Endian), see https://github.com/chai2010/image.
type MemP interface {
	MemPMagic() string
	Bounds() image.Rectangle
	Channels() int
	DataType() reflect.Kind
	Pix() []byte // PixSlice type

	// Stride is the Pix stride (in bytes, must align with SizeofKind(p.DataType))
	// between vertically adjacent pixels.
	Stride() int
}

func EncodeImage(m image.Image, quality float32) (data []byte, ok bool) {
	return encodeImage(m, quality)
}

func EncodeGray(m *image.Gray, quality float32) (data []byte, ok bool) {
	return encodeGray(m.Pix, m.Bounds().Dx(), m.Bounds().Dy(), m.Stride, quality)
}
func EncodeRGBA(m *image.RGBA, quality float32) (data []byte, ok bool) {
	return encodeRGBA(m.Pix, m.Bounds().Dx(), m.Bounds().Dy(), m.Stride, quality)
}
func EncodeRGB(pix []byte, w, h, stride int, quality float32) (data []byte, ok bool) {
	return encodeRGB(pix, w, h, stride, quality)
}

func Encode(w io.Writer, m image.Image, o *Options) error {
	var quality = float32(DefaultQuality)
	if o != nil {
		quality = o.Quality
	}

	data, ok := encodeImage(m, quality)
	if !ok {
		return errEncodeFailed
	}
	_, err := io.Copy(w, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return nil
}

func Save(name string, m image.Image, o *Options) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return Encode(f, m, o)
}

func encodeImage(m image.Image, quality float32) (data []byte, ok bool) {
	b := m.Bounds()
	if memp, ok := m.(MemP); ok {
		switch {
		case memp.Channels() == 1 && memp.DataType() == reflect.Uint8:
			return encodeGray(memp.Pix(), b.Dx(), b.Dy(), memp.Stride(), quality)
		case memp.Channels() == 3 && memp.DataType() == reflect.Uint8:
			return encodeRGB(memp.Pix(), b.Dx(), b.Dy(), memp.Stride(), quality)
		case memp.Channels() == 4 && memp.DataType() == reflect.Uint8:
			return encodeRGBA(memp.Pix(), b.Dx(), b.Dy(), memp.Stride(), quality)
		}
	}
	switch m := m.(type) {
	case *image.Gray:
		return encodeGray(m.Pix, b.Dx(), b.Dy(), m.Stride, quality)
	case *image.RGBA:
		return encodeRGBA(m.Pix, b.Dx(), b.Dy(), m.Stride, quality)
	default:
		rgba := toRGBAImage(m)
		return encodeRGBA(rgba.Pix, b.Dx(), b.Dy(), rgba.Stride, quality)
	}
}

func toRGBAImage(m image.Image) *image.RGBA {
	if m, ok := m.(*image.RGBA); ok {
		return m
	}
	b := m.Bounds()
	rgba := image.NewRGBA(b)
	dstColorRGBA64 := &color.RGBA64{}
	dstColor := color.Color(dstColorRGBA64)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pr, pg, pb, pa := m.At(x, y).RGBA()
			dstColorRGBA64.R = uint16(pr)
			dstColorRGBA64.G = uint16(pg)
			dstColorRGBA64.B = uint16(pb)
			dstColorRGBA64.A = uint16(pa)
			rgba.Set(x, y, dstColor)
		}
	}
	return rgba
}
