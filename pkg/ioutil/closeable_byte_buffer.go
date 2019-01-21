package ioutil

import (
	"bytes"
)

// CloseableByteBuffer it's a bytes.Buffer that implements io.Closer
type CloseableByteBuffer struct {
	*bytes.Buffer
}

// Close is a no-op for a string io.Buffer
func (bb *CloseableByteBuffer) Close() error {
	return nil
}

// CreateCloseableBufferString creates a CloseableByteBuffer from a string
func CreateCloseableBufferString(data string) *CloseableByteBuffer {
	return &CloseableByteBuffer{bytes.NewBufferString(data)}
}
