package utils

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBytesReader(t *testing.T) {
	should := require.New(t)
	data := []byte("Hello, World!")
	reader := NewBytesReader(data)

	// Read 示例
	buf := make([]byte, 5)
	n, err := reader.Read(buf)
	should.Nil(err)
	t.Logf("Read %d bytes: %s, Remaining: %d\n", n, buf, reader.InboundBuffered())

	// Next 示例
	next, err := reader.Next(5)
	should.Nil(err)
	t.Logf("Next 5 bytes: %s, Remaining: %d\n", next, reader.InboundBuffered())

	// Peek 示例
	peek, err := reader.Peek(2)
	should.Nil(err)
	t.Logf("Peek 5 bytes: %s, Remaining: %d\n", peek, reader.InboundBuffered())

	// Discard 示例
	discarded, err := reader.Discard(1)
	should.Nil(err)
	t.Logf("Discarded %d bytes, Remaining: %d\n", discarded, reader.InboundBuffered())

	// 继续读取
	var writer bytes.Buffer
	if n64, wErr := reader.WriteTo(&writer); wErr != nil {
		should.Nil(wErr)
	} else {
		t.Logf("Write %d bytes: %s, Remaining: %d\n", n64, writer.String(), reader.InboundBuffered())
	}
}
