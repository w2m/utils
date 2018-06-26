package ioStream

import (
	"io"
)

type WriterStream struct {
	writers []io.Writer
}

func (w *WriterStream) AddWriter(writer io.Writer) *WriterStream {
	w.writers = append(w.writers, writer)
	return w
}

func (w *WriterStream) Write(p []byte) (n int, err error) {
	for _, writer := range w.writers {
		n, err = writer.Write(p)
	}
	return
}
