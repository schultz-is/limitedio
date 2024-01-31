// Copyright 2024 Matt Schultz <schultz@sent.com>. All rights reserved.
// Use of this source code is governed by an ISC license that can be found in the LICENSE file.

package limitedio_test

import (
	"bytes"
	"fmt"

	"github.com/schultz-is/limitedio"
)

func ExampleLimitedWriter_Write() {
	var buf bytes.Buffer

	w := limitedio.LimitWriter(&buf, 3)

	n, err := w.Write([]byte("howdy"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", n)
	fmt.Printf("%q\n", buf.String())

	n, err = w.Write([]byte("howdy"))
	fmt.Printf("%d\n", n)
	fmt.Printf("%s\n", err)

	// Output:
	// 3
	// "how"
	// 0
	// EOF
}

func ExampleCallLimitedReader_Read() {
	rdr := bytes.NewReader([]byte("howdy"))
	b := make([]byte, 5)

	r := limitedio.CallLimitReader(rdr, 1)

	n, err := r.Read(b[:3])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", n)
	fmt.Printf("%q\n", b)

	n, err = r.Read(b[3:])
	fmt.Printf("%d\n", n)
	fmt.Printf("%s\n", err)

	// Output:
	// 3
	// "how\x00\x00"
	// 0
	// EOF
}
