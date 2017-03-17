// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <stdint.h>
#include <stdio.h>

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

	FILE* fin = fopen("./testdata/lena.jpg", "rb");
	if (!fin) {
		fprintf(stderr, "Can't open input file\n");
		return 1;
	}

	std::string in_data = ReadFileOrDie(fin);
	std::string out_data;

	if (!guetzli::Process(params, &stats, in_data, &out_data)) {
		fprintf(stderr, "Guetzli processing failed\n");
		return 1;
	}

	FILE* fout = fopen("a.out.jpg", "wb");
	if (!fout) {
		fprintf(stderr, "Can't open output file for writing\n");
		return 1;
	}

	WriteFileOrDie(fout, out_data);
	return 0;
}
