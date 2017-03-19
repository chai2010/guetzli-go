// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build all_formats

// go get   -tags="all_formats" github.com/chai2010/guetzli-go/apps/guetzli
// go build -tags="all_formats"

package main

import (
	"sort"

	_ "github.com/chai2010/bpg"
	_ "github.com/chai2010/webp"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

func init() {
	supportFormatExtList = append(supportFormatExtList,
		".bmp",
		".bpg",
		".tif",
		".tiff",
		".webp",
	)
	sort.Strings(supportFormatExtList)
}
