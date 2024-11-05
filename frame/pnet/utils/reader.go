package utils

import (
	"errors"
	"io"
)

func NewBytesReader(data []byte) *BytesReader {
	return &BytesReader{data: data}
}

type BytesReader struct {
	data []byte
	pos  int
}

func (r *BytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (r *BytesReader) Next(n int) (buf []byte, err error) {
	if r.pos+n > len(r.data) {
		n = len(r.data) - r.pos
		if n <= 0 {
			r.pos = len(r.data)
			return nil, nil
		}
	}
	buf = r.data[r.pos : r.pos+n]
	r.pos += n
	return buf, nil
}

func (r *BytesReader) Peek(n int) (buf []byte, err error) {
	if r.pos+n > len(r.data) {
		return nil, errors.New("out of range")
	}
	return r.data[r.pos : r.pos+n], nil
}

func (r *BytesReader) Discard(n int) (discarded int, err error) {
	if r.pos+n > len(r.data) {
		return 0, errors.New("out of range")
	}
	r.pos += n
	return n, nil
}

func (r *BytesReader) InboundBuffered() int {
	return len(r.data) - r.pos
}

func (r *BytesReader) WriteTo(w io.Writer) (n int64, err error) {
	for {
		if r.pos >= len(r.data) {
			break
		}

		bytesToWrite := r.data[r.pos:]
		m, wErr := w.Write(bytesToWrite)
		if wErr != nil {
			return n, wErr
		}

		n += int64(m)
		r.pos += m
	}
	return n, nil
}
