// Copyright 2024 Matt Schultz <schultz@sent.com>. All rights reserved.
// Use of this source code is governed by an ISC license that can be found in the LICENSE file.

/*
Package limitedio provides wrappers for [io.Reader] and [io.Writer] that limit input and output
by configurable parameters.

A LimitedReader type is excluded from this package as one is already defined in the standard
library [io] package. The implementations in this package are inspired heavily by
[io.LimitedReader]. Call-limited wrappers target test cases where underlying readers and writers
are inaccessible or not introspectable. None of the implementations in this package are inherently
safe for concurrent use.
*/
package limitedio

import "io"

// A LimitedWriter writes to W but limits the amount of data returned to just N bytes. Each call to
// Write updates N to reflect the new amount remaining. Write returns EOF when N <= 0 or when the
// underlying W returns EOF.
type LimitedWriter struct {
	W io.Writer // W is the underlying io.Writer.
	N int64     // N is the maximum number of bytes remaining.
}

// LimitWriter returns a Writer that writes to w but stops with EOF after n bytes. The underlying
// implementation is a *LimitedWriter.
func LimitWriter(w io.Writer, n int64) io.Writer { return &LimitedWriter{w, n} }

// Write sends the provided bytes to the underlying [io.Writer], limiting the output to the
// remaining bytes. When there are no bytes remaining in the limit, an EOF is returned. Any error
// from the underlying writer will be returned.
func (l *LimitedWriter) Write(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.W.Write(p)
	l.N -= int64(n)
	return
}

// A CallLimitedReader reads from R but limits the number of calls to Read to N. Each call to Read
// updates N to reflect the new call count remaining. Read returns EOF when N <= 0 or when the
// underlying R returns EOF.
type CallLimitedReader struct {
	R io.Reader // R is the underlying io.Reader.
	N int64     // N is the maximum number of calls remaining.
}

// CallLimitReader returns a Reader that reads from r but stops with EOF after n calls to Read. The
// underlying implementation is a *CallLimitedReader.
func CallLimitReader(r io.Reader, n int64) io.Reader { return &CallLimitedReader{r, n} }

// Read receives bytes from the underlying [io.Reader], limiting the output based on the number of
// remaining calls allowed to Read. When there are no more Read calls remaining in the limit, an
// EOF is returned. Any error from the underlying reader will be returned.
func (l *CallLimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, io.EOF
	}
	n, err = l.R.Read(p)
	l.N -= 1
	return
}

// A CallLimitedWriter writes to W but limits the number of calls to Write to N. Each call to Write
// updates N to reflect the new call count remaining. Write returns an EOF when N <= 0 or when the
// underlying W returns EOF.
type CallLimitedWriter struct {
	W io.Writer // W is the underlying io.Writer.
	N int64     // N is the maximum number of calls remaining.
}

// CallLimitWriter returns a Writer that writes to w but stops with EOF after n calls to Write. The
// underlying implementation is a *CallLimitedWriter.
func CallLimitWriter(w io.Writer, n int64) io.Writer { return &CallLimitedWriter{w, n} }

// Write sends the provided bytes to the underlying [io.Writer], limiting the output based on the
// number of remaining calls allowed to Write. When there are no more Write calls remaining in the
// limit, an EOF is returned. Any error from the underlying writer will be returned.
func (l *CallLimitedWriter) Write(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, io.EOF
	}
	n, err = l.W.Write(p)
	l.N -= 1
	return
}
