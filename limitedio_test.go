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
	lw := LimitWriter(&buf, -1)
	n, err := lw.Write([]byte("howdy"))
	if n != 0 {
		t.Fatalf("Write on LimitedWriter with N<0 wrote %d bytes; expected 0", n)
	}
	if err != io.EOF {
		t.Fatalf("Write on LimitedWriter with N<0 did not return EOF: %s", err)
	}

	// Calling Write with a zero value for N should result in zero bytes written and an EOF.
	buf.Reset()
	lw = LimitWriter(&buf, 0)
	n, err = lw.Write([]byte("howdy"))
	if n != 0 {
		t.Fatalf("Write on LimitedWriter with N=0 wrote %d bytes; expected 0", n)
	}
	if err != io.EOF {
		t.Fatalf("Write on LimitedWriter with N=0 did not return EOF: %s", err)
	}

	// When N is smaller than the length of the input, Write should output N bytes and return an EOF.
	buf.Reset()
	lw = LimitWriter(&buf, 3)
	n, err = lw.Write([]byte("howdy"))
	if n != 3 {
		t.Fatalf("Write on LimitedWriter with N<len(input) wrote %d bytes; expected 3", n)
	}
	if err != nil {
		t.Fatalf("Write on LimitedWriter returned an unexpected error: %s", err)
	}

	// Errors from the underlying writer should propagate up.
	buf.Reset()
	lw = LimitWriter(LimitWriter(&buf, 0), 50)
	_, err = lw.Write([]byte("howdy"))
	if err == nil {
		t.Fatal("Write on LimitedWriter did not propagate underlying error")
	}
}
