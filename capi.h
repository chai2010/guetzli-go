// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef guetzli_h_
#define guetzli_h_

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct guetzli_string_t guetzli_string_t;

guetzli_string_t* guetzli_string_new(int size);
void guetzli_string_delete(guetzli_string_t* p);

void guetzli_string_resize(guetzli_string_t* p, int size);
int guetzli_string_size(guetzli_string_t* p);
char* guetzli_string_data(guetzli_string_t* p);

guetzli_string_t* guetzliEncodeGray(const uint8_t* pix, int w, int h, int stride, float quality);
guetzli_string_t* guetzliEncodeRGB(const uint8_t* pix, int w, int h, int stride, float quality);
guetzli_string_t* guetzliEncodeRGBA(const uint8_t* pix, int w, int h, int stride, float quality);

#ifdef __cplusplus
}
#endif
#endif // guetzli_h_
