//go:build !solution

package otp

import (
	"io"
)

type reader struct {
	r    io.Reader
	prng io.Reader
}

type writer struct {
	w    io.Writer
	prng io.Reader
}

func (r *reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	prngBytes := make([]byte, n)
	_, _ = r.prng.Read(prngBytes)
	for i := 0; i < len(prngBytes); i++ {
		p[i] = prngBytes[i] ^ p[i]
	}

	return n, err
}

func (w *writer) Write(p []byte) (n int, err error) {
	prngBytes := make([]byte, len(p))
	_, _ = w.prng.Read(prngBytes)
	tmpWriter := make([]byte, len(p))
	copy(tmpWriter, p)
	for i := 0; i < len(prngBytes); i++ {
		tmpWriter[i] = prngBytes[i] ^ tmpWriter[i]
	}

	return w.w.Write(tmpWriter)
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &reader{r, prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &writer{w, prng}
}
