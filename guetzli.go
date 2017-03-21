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
	"io/ioutil"
	"reflect"
)

// Guetzli algorithm is beneficial only for this quality or higher.
// If you work around it and lower the quality, you'll only
// waste your time generating average JPEGs, without Guetzli improvement.
//
// If you need even smaller files and can tolerate bigger distortions,
// MozJPEG is a better choice for quality < 84.
const (
	MinQuality = 84
	MaxQuality = 110

	DefaultQuality = 84
)

var (
	errEncodeFailed   = errors.New("guetzli: encode failed!")
	errInvalidQuality = errors.New("guetzli: invalid quality (must >=84 and <= 110)")
)

func isInvalidQuality(quality int) bool {
	return quality < MinQuality || quality > MaxQuality
}

type Options struct {
	Quality int // 84 <= quality <= 110
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

func EncodeImage(m image.Image, quality int) (data []byte, ok bool) {
	if isInvalidQuality(quality) {
		return nil, false
	}
	return encodeImage(m, quality)
}

func EncodeGray(m *image.Gray, quality int) (data []byte, ok bool) {
	if isInvalidQuality(quality) {
		return nil, false
	}
	return encodeGray(m.Pix, m.Bounds().Dx(), m.Bounds().Dy(), m.Stride, quality)
}
func EncodeRGBA(m *image.RGBA, quality int) (data []byte, ok bool) {
	if isInvalidQuality(quality) {
		return nil, false
	}
	return encodeRGBA(m.Pix, m.Bounds().Dx(), m.Bounds().Dy(), m.Stride, quality)
}
func EncodeRGB(pix []byte, w, h, stride int, quality int) (data []byte, ok bool) {
	if isInvalidQuality(quality) {
		return nil, false
	}
	return encodeRGB(pix, w, h, stride, quality)
}

func Encode(w io.Writer, m image.Image, o *Options) error {
	var quality = DefaultQuality
	if o != nil {
		quality = o.Quality
	}
	if isInvalidQuality(quality) {
		return errInvalidQuality
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
	var quality = DefaultQuality
	if o != nil {
		quality = o.Quality
	}
	if isInvalidQuality(quality) {
		return errInvalidQuality
	}

	data, ok := encodeImage(m, quality)
	if !ok {
		return errEncodeFailed
	}

	return ioutil.WriteFile(name, data, 0666)
}

func encodeImage(m image.Image, quality int) (data []byte, ok bool) {
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
