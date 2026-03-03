package utils

import (
	"bytes"
	"io"
)

var _ io.Writer = &WriteBuffer{}

type WriteBuffer struct {
	data []byte
}

// Write implements io.Writer.
func (w *WriteBuffer) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func NewWriteBuffer() *WriteBuffer {
	return &WriteBuffer{
		data: []byte{},
	}
}

// Flush writes the buffer to the target writer and resets it.
// If replacements are provided, all occurrences of the first element in each pair are replaced with the second element.
// In case of an uneven number of replacements, the last replacement is ignored.
func (w *WriteBuffer) Flush(target io.Writer, replacements ...string) error {
	for i := 0; i < len(replacements) && i+1 < len(replacements); i += 2 {
		w.data = bytes.ReplaceAll(w.data, []byte(replacements[i]), []byte(replacements[i+1]))
	}
	_, err := target.Write(w.data)
	w.data = []byte{}
	return err
}

func (w *WriteBuffer) FlushToString() string {
	return string(w.data)
}

func (w *WriteBuffer) Data() []byte {
	return w.data
}
