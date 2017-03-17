// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "./lodepng.h"

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

#include <string>
#include <vector>
#include <memory>

#include <guetzli/processor.h>
#include <guetzli/quality.h>
#include <guetzli/stats.h>

std::string ReadFileOrDie(FILE* f) {
	if (fseek(f, 0, SEEK_END) != 0) {
		perror("fseek");
		exit(1);
	}
	off_t size = ftell(f);
	if (size < 0) {
		perror("ftell");
		exit(1);
	}
	if (fseek(f, 0, SEEK_SET) != 0) {
		perror("fseek");
		exit(1);
	}
	std::unique_ptr<char[]> buf(new char[size]);
	if (fread(buf.get(), 1, size, f) != (size_t)size) {
		perror("fread");
		exit(1);
	}
	std::string result(buf.get(), size);
	return result;
}

void WriteFileOrDie(FILE* f, const std::string& contents) {
	if (fwrite(contents.data(), 1, contents.size(), f) != contents.size()) {
		perror("fwrite");
		exit(1);
	}
	if (fclose(f) < 0) {
		perror("fclose");
		exit(1);
	}
}

int main() {
	guetzli::Params params;
	params.butteraugli_target = guetzli::ButteraugliScoreForQuality(95);
	guetzli::ProcessStats stats;

	FILE* fin = fopen("./testdata/video-001.png", "rb");
	if (!fin) {
		fprintf(stderr, "Can't open input file\n");
		return 1;
	}

	std::string in_data = ReadFileOrDie(fin);
	std::string out_data;

	unsigned char* img;
	unsigned w, h;

	auto err = lodepng_decode24(&img, &w, &h, (const unsigned char*)in_data.data(), in_data.size());
	if(err != 0) {
		fprintf(stderr, "lodepng_decode24 failed\n");
		return 1;
	}

	std::vector<uint8_t> rgb;
	rgb.assign((uint8_t*)img, (uint8_t*)(img+w*h*3));
	free(img);

	clock_t startTime = clock();
	if(!guetzli::Process(params, &stats, rgb, w, h, &out_data)) {
		fprintf(stderr, "Guetzli processing failed\n");
		return 1;
	}

	double duration = (clock() - startTime) / (double) CLOCKS_PER_SEC;
	printf("duration: %.2fs\n", duration); // 5s

	FILE* fout = fopen("a.out.jpg", "wb");
	if(!fout) {
		fprintf(stderr, "Can't open output file for writing\n");
		return 1;
	}

	WriteFileOrDie(fout, out_data);
	return 0;
}
