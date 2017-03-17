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
)

const DefaultQuality = 75

type Options struct {
	Quality float32
}

var errEncodeFailed = errors.New("guetzli: encode failed!")

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

	data, err := encode(m, quality)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return nil
}

func encode(m image.Image, quality float32) (data []byte, err error) {
	b := m.Bounds()
	switch m := m.(type) {
	case *image.Gray:
		data, ok := encodeGray(m.Pix, b.Dx(), b.Dy(), m.Stride, quality)
		if !ok {
			return nil, errEncodeFailed
		}
		return data, nil
	case *image.RGBA:
		data, ok := encodeRGBA(m.Pix, b.Dx(), b.Dy(), m.Stride, quality)
		if !ok {
			return nil, errEncodeFailed
		}
		return data, nil
	default:
		rgba := toRGBAImage(m)
		data, ok := encodeRGBA(rgba.Pix, b.Dx(), b.Dy(), rgba.Stride, quality)
		if !ok {
			return nil, errEncodeFailed
		}
		return data, nil
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
