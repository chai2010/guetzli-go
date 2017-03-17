# Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

project(GUETZLI)
cmake_minimum_required(VERSION 3.5)

set(CMAKE_CXX_STANDARD 11)

#------------------------------------------------------------------------------

include_directories(AFTER
	../guetzli-1.0
	../guetzli-1.0/third_party/butteraugli
)

set(GUETZLI_SRC
	../guetzli-1.0/guetzli/butteraugli_comparator.cc
	../guetzli-1.0/guetzli/dct_double.cc
	../guetzli-1.0/guetzli/debug_print.cc
	../guetzli-1.0/guetzli/entropy_encode.cc
	../guetzli-1.0/guetzli/fdct.cc
	../guetzli-1.0/guetzli/gamma_correct.cc
	../guetzli-1.0/guetzli/idct.cc
	../guetzli-1.0/guetzli/jpeg_data.cc
	../guetzli-1.0/guetzli/jpeg_data_decoder.cc
	../guetzli-1.0/guetzli/jpeg_data_encoder.cc
	../guetzli-1.0/guetzli/jpeg_data_reader.cc
	../guetzli-1.0/guetzli/jpeg_data_writer.cc
	../guetzli-1.0/guetzli/jpeg_huffman_decode.cc
	../guetzli-1.0/guetzli/output_image.cc
	../guetzli-1.0/guetzli/preprocess_downsample.cc
	../guetzli-1.0/guetzli/processor.cc
	../guetzli-1.0/guetzli/quality.cc
	../guetzli-1.0/guetzli/quantize.cc
	../guetzli-1.0/guetzli/score.cc

	../guetzli-1.0/third_party/butteraugli/butteraugli/butteraugli.cc
)

#------------------------------------------------------------------------------

add_library(guetzlilib STATIC
	${GUETZLI_SRC}
)
set_target_properties(guetzlilib
	PROPERTIES OUTPUT_NAME "guetzli"
)

install(TARGETS guetzlilib
	RUNTIME DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}
	LIBRARY DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}
	ARCHIVE DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}
)

#------------------------------------------------------------------------------

add_executable(guetzlihello
	./hello.cc
	./png/lodepng.cpp

	${GUETZLI_SRC}
)
install(TARGETS guetzlihello
  RUNTIME DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}
  LIBRARY DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}
  ARCHIVE DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}
)

#------------------------------------------------------------------------------

add_executable(guetzlitest
	./test.cc
	./test_util.cc
	./test_util_jpg.cc
	./test_util_png.cc
	./bench_test.cc

	${GUETZLI_SRC}
)
install(TARGETS guetzlitest
  RUNTIME DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}
  LIBRARY DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}
  ARCHIVE DESTINATION ${CMAKE_CURRENT_SOURCE_DIR}
)

#------------------------------------------------------------------------------
