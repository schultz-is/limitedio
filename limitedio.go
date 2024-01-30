package limitedio

import "io"

// A LimitedWriter writes to W but limits the amount of data returned to just N bytes. Each call to
// Write updates N to reflect the new amount remaining. Write returns EOF when N <= 0 or when the
// underlying W returns EOF.
type LimitedWriter struct {
	W io.Writer // W is the underlying io.Writer.
	N int64     // N is the maximum number of bytes remaining.
}

// Write sends the provided bytes to the underlying [io.Writer], limiting the output to the
// remaining bytes. When there are no bytes remaining in the limit, an EOF is returned. If there
// are insufficient bytes remaining to write the provided bytes in their entirety, only N bytes
// will be written and an EOF will be returned.
func (l *LimitedWriter) Write(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.W.Write(p)
	l.N -= int64(n)
	return
}
