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
