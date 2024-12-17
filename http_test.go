// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

package jsoncall

import (
	"bytes"
	"net/http"
	"testing"
)

type testResponseWriter struct {
	header     http.Header
	buf        bytes.Buffer
	statusCode int
}

func (w *testResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *testResponseWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *testResponseWriter) Header() http.Header {
	return w.header
}

func (w *testResponseWriter) CheckStatus(t *testing.T, expected int) {
	if expected != w.statusCode {
		t.Fatalf("status: %d, expected: %d", w.statusCode, expected)
	}
}

func newResponseWriter() *testResponseWriter {
	var t testResponseWriter
	t.header = make(http.Header)
	return &t
}

func TestResponse(t *testing.T) {
	var w http.ResponseWriter
	var r testResponseWriter
	w = &r
	_ = w
}
