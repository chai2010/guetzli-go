// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "./capi.h"

#include <string>
#include <vector>
#include <memory>

#include <guetzli/processor.h>
#include <guetzli/quality.h>
#include <guetzli/stats.h>

struct guetzli_string_t {
	std::string str_;
};

guetzli_string_t* guetzli_string_new(int size) {
	auto p = new guetzli_string_t();
	p->str_.resize(size);
	return p;
}
void guetzli_data_delete(guetzli_string_t* p) {
	delete p;
}

void guetzli_string_resize(guetzli_string_t* p, int size) {
	p->str_.resize(size);
}
int guetzli_string_size(guetzli_string_t* p) {
	return int(p->str_.size());
}
char* guetzli_string_data(guetzli_string_t* p) {
	return (char*)(p->str_.data());
}

static void grayToRGBVector(std::vector<uint8_t>* rgb, const uint8_t* pix, int w, int h, int stride) {
	rgb->resize(w*h*3);
	if(stride == 0) {
		stride = w;
	}
	int off = 0;
	for(int i = 0; i < h; i++) {
		auto p = pix + i*stride;
		for(int j = 0; j < w; j++) {
			(*rgb)[off++] = p[j]; // R
			(*rgb)[off++] = p[j]; // G
			(*rgb)[off++] = p[j]; // B
		}
	}
	return;
}

static void rgbToRGBVector(std::vector<uint8_t>* rgb, const uint8_t* pix, int w, int h, int stride) {
	rgb->resize(w*h*3);
	if(stride == 0) {
		memcpy(rgb->data(), pix, rgb->size());
		return;
	}
	for(int i = 0; i < h; i++) {
		memcpy(&(*rgb)[i*w*3], pix+i*stride, w*3);
	}
	return;
}

static void rgbaToRGBVector(std::vector<uint8_t>* rgb, const uint8_t* pix, int w, int h, int stride) {
	rgb->resize(w*h*3);
	if(stride == 0) {
		stride = w*4;
	}
	int off = 0;
	for(int i = 0; i < h; i++) {
		auto p = pix + i*stride;
		for(int j = 0; j < w; j++) {
			(*rgb)[off++] = p[j*4+0]; // R
			(*rgb)[off++] = p[j*4+1]; // G
			(*rgb)[off++] = p[j*4+2]; // B
		}
	}
	return;
}

static guetzli_string_t* encodeRGB(const std::vector<uint8_t>& rgb, int w, int h, float quality) {
	guetzli::Params params;
	params.butteraugli_target = guetzli::ButteraugliScoreForQuality(quality);
	guetzli::ProcessStats stats;

	auto p = guetzli_string_new(0);
	if(!guetzli::Process(params, &stats, rgb, w, h, &p->str_)) {
		guetzli_data_delete(p);
		return NULL;
	}
	return p;
}

guetzli_string_t* guetzliEncodeGray(const uint8_t* pix, int w, int h, int stride, float quality) {
	std::vector<uint8_t> rgb;
	grayToRGBVector(&rgb, pix, w, h, stride);
	return encodeRGB(rgb, w, h, quality);
}

guetzli_string_t* guetzliEncodeRGB(const uint8_t* pix, int w, int h, int stride, float quality) {
	std::vector<uint8_t> rgb;
	rgbToRGBVector(&rgb, pix, w, h, stride);
	return encodeRGB(rgb, w, h, quality);
}

guetzli_string_t* guetzliEncodeRGBA(const uint8_t* pix, int w, int h, int stride, float quality) {
	std::vector<uint8_t> rgb;
	rgbaToRGBVector(&rgb, pix, w, h, stride);
	return encodeRGB(rgb, w, h, quality);
}
