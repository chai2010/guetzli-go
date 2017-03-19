// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Guetzli JPEG compressor
//
// Usage:
//
//	guetzli [flags] input_filename output_filename
//	guetzli [flags] input_dir output_dir [ext...]
//
//	  -quality float
//	        Expressed as a JPEG quality value. (default 75)
//	  -regexp string
//	        regexp for base filename.
//	  -version
//	        Show version and exit.
//
// Example:
//
//	guetzli [-quality=95] original.png output.jpg
//	guetzli [-quality=95] original.jpg output.jpg
//
//	guetzli [-quality=95] input_dir output_dir .png
//	guetzli [-quality=95] input_dir output_dir .png .jpg .jpeg
//	guetzli [-quality=95 -regexp="^\d+"] input_dir output_dir .png
//
// Note: Default image ext is: .jpeg .jpg .png
//
// Note: Supported formats: .gif, .jpeg, .jpg, .png
//
// Build more image format support (use "all_formats" tag):
//
//	go get   -tags="all_formats" github.com/chai2010/guetzli-go/apps/guetzli
//	go build -tags="all_formats"
//
// See https://godoc.org/github.com/chai2010/guetzli-go
//
// See https://github.com/google/guetzli
//
// Report bugs to <chaishushan{AT}gmail.com>.
//
package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/chai2010/guetzli-go"
)

const Version = "1.0"

var (
	flagQuality = flag.Float64("quality", 75, "Expressed as a JPEG quality value.")
	flagRegexp  = flag.String("regexp", "", "regexp for base filename.")
	flagVersion = flag.Bool("version", false, "Show version and exit.")
)

var supportFormatExtList = []string{
	".png",
	".jpg",
	".jpeg",
	".gif",
}

func init() {
	sort.Strings(supportFormatExtList)
}

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `Guetzli JPEG compressor.

Usage: guetzli [flags] input_filename output_filename
       guetzli [flags] input_dir output_dir [ext...]
`)
		flag.PrintDefaults()
		fmt.Printf(`
Example:

  guetzli [-quality=95] original.png output.jpg
  guetzli [-quality=95] original.jpg output.jpg

  guetzli [-quality=95] input_dir output_dir .png
  guetzli [-quality=95] input_dir output_dir .png .jpg .jpeg
  guetzli [-quality=95 -regexp="^\d+"] input_dir output_dir .png

Note: Default image ext is: .jpeg .jpg .png

Note: Supported formats: %s

See https://godoc.org/github.com/chai2010/guetzli-go
See https://github.com/google/guetzli

Report bugs to <chaishushan{AT}gmail.com>.
`, strings.Join(supportFormatExtList, ", "))
	}
}

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf("guetzli-%s\n", Version)
		os.Exit(0)
	}

	if flag.NArg() < 2 {
		fmt.Printf("guetzli-%s\n", Version)
		os.Exit(0)
	}

	var (
		inputPath  = flag.Arg(0)
		outputPath = flag.Arg(1)
		inputDir   = flag.Arg(0)
		outputDir  = flag.Arg(1)
		extList    = flag.Args()[2:]
	)

	// default ext is only for jpg and png
	if len(extList) == 0 {
		extList = []string{".jpg", ".jpeg", ".png"}
	}

	// only for one image
	if !isDir(inputPath) {
		err := guetzliCompressImage(inputPath, outputPath, float32(*flagQuality), nil)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// parse regexp
	if *flagRegexp == "" {
		*flagRegexp = ".*"
	}
	rePath, err := regexp.Compile(*flagRegexp)
	if err != nil {
		log.Fatalf("invalid regexp: %q", *flagRegexp)
		return
	}

	// walk dir
	var seemMap = make(map[string]bool)
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !matchExtList(path, supportFormatExtList...) {
			return nil // support format
		}
		if !matchExtList(path, extList...) {
			return nil // unselected format
		}
		if !rePath.MatchString(path) {
			return nil // unselected path
		}

		inputPath = goodAbsPath(path)
		outputPath = func() string {
			relPath, err := filepath.Rel(inputDir, path)
			if err != nil {
				panic(err)
			}
			newpath := filepath.Join(outputDir, relPath)
			if !matchExtList(newpath, ".jpg", ".jpeg") {
				newpath = strings.TrimSuffix(newpath, filepath.Ext(newpath)) + ".jpg"
			}
			newpath = goodAbsPath(newpath)
			return newpath
		}()

		os.MkdirAll(filepath.Dir(outputPath), 0777)
		timeUsed, err := func() (time.Duration, error) {
			s := time.Now()

			seemMap[outputPath] = true
			err := guetzliCompressImage(path, outputPath, float32(*flagQuality), seemMap)
			seemMap[inputPath] = true

			timeUsed := time.Now().Sub(s)
			return timeUsed, err
		}()
		if err != nil {
			log.Printf("%s filed, err = %v\n", path, err)
			return nil
		}

		fmt.Println(path, "ok,", timeUsed)
		return nil
	})
}

func matchExtList(name string, extList ...string) bool {
	if len(extList) == 0 {
		return true
	}
	ext := filepath.Ext(name)
	for _, s := range extList {
		if strings.EqualFold(s, ext) {
			return true
		}
	}

	return false
}

func guetzliCompressImage(inputFilename, outputFilename string, quality float32, seen map[string]bool) error {
	if seen != nil && seen[inputFilename] {
		return nil // skip
	}
	fin, err := os.Open(inputFilename)
	if err != nil {
		return fmt.Errorf("open %q failed, err = %v", inputFilename, err)
	}
	defer fin.Close()

	m, _, err := image.Decode(fin)
	if err != nil {
		return fmt.Errorf("decode %q failed, err = %v", inputFilename, err)
	}

	err = guetzli.Save(outputFilename, m, &guetzli.Options{Quality: quality})
	if err != nil {
		return fmt.Errorf("save %q failed, err = %v", outputFilename, err)
	}
	return nil
}

func isDir(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}
	return false
}

func goodAbsPath(path string) string {
	if abspath, err := filepath.Abs(path); err == nil {
		path = abspath
	}
	return filepath.ToSlash(filepath.Clean(path))
}
