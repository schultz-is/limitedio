// Copyright 2024 Matt Schultz <schultz@sent.com>. All rights reserved.
// Use of this source code is governed by an ISC license that can be found in the LICENSE file.

package limitedio_test

import (
	"bytes"
	"io"
	"testing"

	. "github.com/schultz-is/limitedio"
)

func TestLimitedWriter(t *testing.T) {
	// Calling Write with a negative value for N should result in zero bytes written and an EOF.
	var buf bytes.Buffer
	w := LimitWriter(&buf, -1)
	n, err := w.Write([]byte("howdy"))
	if n != 0 {
		t.Fatalf("Write on LimitedWriter with N<0 wrote %d bytes; expected 0", n)
	}
	if err != io.EOF {
		t.Fatalf("Write on LimitedWriter with N<0 did not return EOF: %s", err)
	}

	// Calling Write with a zero value for N should result in zero bytes written and an EOF.
	buf.Reset()
	w = LimitWriter(&buf, 0)
	n, err = w.Write([]byte("howdy"))
	if n != 0 {
		t.Fatalf("Write on LimitedWriter with N=0 wrote %d bytes; expected 0", n)
	}
	if err != io.EOF {
		t.Fatalf("Write on LimitedWriter with N=0 did not return EOF: %s", err)
	}

	// When N is smaller than the length of the input, Write should output N bytes and return an EOF.
	buf.Reset()
	w = LimitWriter(&buf, 3)
	n, err = w.Write([]byte("howdy"))
	if n != 3 {
		t.Fatalf("Write on LimitedWriter with N<len(input) wrote %d bytes; expected 3", n)
	}
	if err != nil {
		t.Fatalf("Write on LimitedWriter returned an unexpected error: %s", err)
	}

	// Errors from the underlying writer should propagate up.
	buf.Reset()
	w = LimitWriter(LimitWriter(&buf, 0), 50)
	_, err = w.Write([]byte("howdy"))
	if err == nil {
		t.Fatal("Write on LimitedWriter did not propagate underlying error")
	}

	// Write should properly decrement N.
	limit := 50
	buf.Reset()
	lw := LimitedWriter{&buf, int64(limit)}
	n, err = lw.Write([]byte("howdy"))
	if err != nil {
		t.Fatalf("Write on LimitedWriter unexpectedly errored: %s", err)
	}
	if lw.N != int64(limit-n) {
		t.Fatalf("Write on LimitedWriter incorrectly updated N; got %d, expected %d", limit-n, lw.N)
	}
}

func TestCallLimitedReader(t *testing.T) {
	msg := []byte("howdy")

	// Calling Read with a negative value for N should result in zero bytes read and an EOF.
	rdr := bytes.NewReader(msg)
	b := make([]byte, 5)
	r := CallLimitReader(rdr, -1)
	n, err := r.Read(b)
	if n != 0 {
		t.Fatalf("Read on CallLimitedReader with N<0 read %d bytes; expected 0", n)
	}
	if err != io.EOF {
		t.Fatalf("Read on CallLimitedReader with N<0 did not return EOF: %s", err)
	}

	// Calling Read with a zero value for N should result in zero bytes read and an EOF.
	rdr.Reset(msg)
	b = make([]byte, 5)
	r = CallLimitReader(rdr, 0)
	n, err = r.Read(b)
	if n != 0 {
		t.Fatalf("Read on CallLimitedReader with N=0 read %d bytes; expected 0", n)
	}
	if err != io.EOF {
		t.Fatalf("Read on CallLimitedReader with N=0 did not return EOF: %s", err)
	}

	// Errors from the underlying reader should propagate up.
	rdr.Reset(nil)
	b = make([]byte, 5)
	r = CallLimitReader(rdr, 50)
	_, err = r.Read(b)
	if err == nil {
		t.Fatal("Read on CallLimitedReader did not propagate underlying error")
	}

	// Read should properly decrement N.
	limit := 50
	rdr.Reset(msg)
	b = make([]byte, 5)
	lr := CallLimitedReader{rdr, int64(limit)}
	_, err = lr.Read(b)
	if err != nil {
		t.Fatalf("Read on CallLimitedReader unexpectedly errored: %s", err)
	}
	if lr.N != int64(limit-1) {
		t.Fatalf("Read on CallLimitedReader incorrectly updated N; got %d, expected %d", limit-1, lr.N)
	}
}

func TestCallLimitedWriter(t *testing.T) {
	// Calling Write with a negative value for N should result in zero bytes written and an EOF.
	var buf bytes.Buffer
	w := CallLimitWriter(&buf, -1)
	n, err := w.Write([]byte("howdy"))
	if n != 0 {
		t.Fatalf("Write on CallLimitedWriter with N<0 wrote %d bytes; expected 0", n)
	}
	if err != io.EOF {
		t.Fatalf("Write on CallLimitedWriter with N<0 did not return EOF: %s", err)
	}

	// Calling Write with a zero value for N should result in zero bytes written and an EOF.
	buf.Reset()
	w = CallLimitWriter(&buf, 0)
	n, err = w.Write([]byte("howdy"))
	if n != 0 {
		t.Fatalf("Write on CallLimitedWriter with N=0 wrote %d bytes; expected 0", n)
	}
	if err != io.EOF {
		t.Fatalf("Write on CallLimitedWriter with N=0 did not return EOF: %s", err)
	}

	// Errors from the underlying writer should propagate up.
	buf.Reset()
	w = CallLimitWriter(LimitWriter(&buf, 0), 50)
	_, err = w.Write([]byte("howdy"))
	if err == nil {
		t.Fatal("Write on CallLimitedWriter did not propagate underlying error")
	}

	// Write should properly decrement N.
	limit := 50
	buf.Reset()
	lw := CallLimitedWriter{&buf, int64(limit)}
	_, err = lw.Write([]byte("howdy"))
	if err != nil {
		t.Fatalf("Write on CallLimitedWriter unexpectedly errored: %s", err)
	}
	if lw.N != int64(limit-1) {
		t.Fatalf("Write on CallLimitedWriter incorrectly updated N; got %d, expected %d", limit-1, lw.N)
	}
}
