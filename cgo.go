// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package guetzli

/*
#cgo CXXFLAGS: -std=c++11
#cgo CPPFLAGS: -I./internal/guetzli-1.0
#cgo CPPFLAGS: -I./internal/guetzli-1.0/third_party/butteraugli

#cgo windows LDFLAGS: -Wl,--allow-multiple-definition

#cgo !windows LDFLAGS: -lm

#include "./capi.h"
*/
import "C"
import "unsafe"

func encodeGray(pix []byte, w, h, stride int, quality float32) (data []byte, ok bool) {
	p := C.guetzliEncodeGray(
		(*C.uint8_t)(unsafe.Pointer(&pix[0])),
		C.int(w), C.int(h), C.int(stride),
		C.float(quality),
	)
	if p == nil {
		return nil, false
	}
	defer C.guetzli_string_delete(p)

	cptr := C.guetzli_string_data(p)
	data = make([]byte, int(C.guetzli_string_size(p)))
	copy(data, ((*[1 << 30]byte)(unsafe.Pointer(cptr)))[0:len(data):len(data)])
	return data, true
}

func encodeRGB(pix []byte, w, h, stride int, quality float32) (data []byte, ok bool) {
	p := C.guetzliEncodeRGB(
		(*C.uint8_t)(unsafe.Pointer(&pix[0])),
		C.int(w), C.int(h), C.int(stride),
		C.float(quality),
	)
	if p == nil {
		return nil, false
	}
	defer C.guetzli_string_delete(p)

	cptr := C.guetzli_string_data(p)
	data = make([]byte, int(C.guetzli_string_size(p)))
	copy(data, ((*[1 << 30]byte)(unsafe.Pointer(cptr)))[0:len(data):len(data)])
	return data, true
}

func encodeRGBA(pix []byte, w, h, stride int, quality float32) (data []byte, ok bool) {
	p := C.guetzliEncodeRGBA(
		(*C.uint8_t)(unsafe.Pointer(&pix[0])),
		C.int(w), C.int(h), C.int(stride),
		C.float(quality),
	)
	if p == nil {
		return nil, false
	}
	defer C.guetzli_string_delete(p)

	cptr := C.guetzli_string_data(p)
	data = make([]byte, int(C.guetzli_string_size(p)))
	copy(data, ((*[1 << 30]byte)(unsafe.Pointer(cptr)))[0:len(data):len(data)])
	return data, true
}
