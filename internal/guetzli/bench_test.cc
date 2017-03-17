// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "test.h"
#include "test_util.h"

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include <string>
#include <vector>
#include <memory>

#include <guetzli/processor.h>
#include <guetzli/quality.h>
#include <guetzli/stats.h>

BENCH(guetzli, encodeRGB) {
	for(int i = 0; i < BenchN(); ++i) {
		//
	}
}
