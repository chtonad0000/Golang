//go:build !solution

package gzep

import (
	"compress/gzip"
	"io"
	"sync"
)

var writerPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(io.Discard)
	},
}

func Encode(data []byte, w io.Writer) error {
	writer := writerPool.Get().(*gzip.Writer)
	defer writerPool.Put(writer)
	writer.Reset(w)
	defer func() {
		if err := writer.Close(); err != nil {
			return
		}
	}()

	_, err := writer.Write(data)
	if err != nil {
		return err
	}
	return writer.Flush()
}
