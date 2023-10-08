package client

import (
	"bytes"
	"github.com/panjf2000/gnet/v2/pkg/buffer/elastic"
	"io"
	"net"
	"persian/frame/plog"
	"persian/frame/pnet"
	"persian/frame/pnet/tcp/session"
)

func NewConn(client *Client, pConn *net.TCPConn) (*Conn, error) {
	conn := &Conn{}
	err := conn.Init(client, pConn)
	if err != nil {
		return nil, err
	}
	return conn, nil

}

type Conn struct {
	*net.TCPConn
	session.BaseConn

	client   *Client
	inbound  elastic.RingBuffer
	outbound bytes.Buffer
	context  any
}

func (conn *Conn) Init(client *Client, pConn *net.TCPConn) error {
	conn.TCPConn = pConn
	conn.client = client
	if client.TCPKeepAlive > 0 {
		if err := conn.TCPConn.SetKeepAlive(true); err != nil {
			return err
		}
		if err := conn.TCPConn.SetKeepAlivePeriod(client.TCPKeepAlive); err != nil {
			return err
		}
	}
	if err := conn.TCPConn.SetNoDelay(client.TCPNoDelay); err != nil {
		return err
	}
	if err := conn.TCPConn.SetReadBuffer(client.SocketRecvBuffer); err != nil {
		return err
	}
	if err := conn.TCPConn.SetWriteBuffer(client.SocketSendBuffer); err != nil {
		return err
	}
	return conn.BaseConn.Init(pConn, true)
}

func (conn *Conn) Context() (ctx any) {
	return conn.client
}

func (conn *Conn) SetContext(_ interface{}) {
	plog.Error("cant set client connection context")
}

func (conn *Conn) Read(b []byte) (n int, err error) {
	return conn.inbound.Read(b)
}

func (conn *Conn) WriteTo(w io.Writer) (n int64, err error) {
	return conn.inbound.WriteTo(w)
}

func (conn *Conn) Next(n int) (buf []byte, err error) {
	inBufferLen := conn.inbound.Buffered()
	if totalLen := inBufferLen; n > totalLen {
		err = io.ErrShortBuffer
		return
	} else if n <= 0 {
		n = totalLen
	}
	head, tail := conn.inbound.Peek(n)
	_, _ = conn.inbound.Discard(n)
	conn.client.loop.cache.Reset()
	conn.client.loop.cache.Write(head)
	conn.client.loop.cache.Write(tail)
	buf = conn.client.loop.cache.Bytes()
	return
}

func (conn *Conn) Peek(n int) (buf []byte, err error) {
	inBufferLen := conn.inbound.Buffered()
	if totalLen := inBufferLen; n > totalLen {
		err = io.ErrShortBuffer
		return
	} else if n <= 0 {
		n = totalLen
	}
	head, tail := conn.inbound.Peek(n)
	conn.client.loop.cache.Reset()
	conn.client.loop.cache.Write(head)
	conn.client.loop.cache.Write(tail)
	buf = conn.client.loop.cache.Bytes()
	return
}

func (conn *Conn) Discard(n int) (discarded int, err error) {
	inBufferLen := conn.inbound.Buffered()
	if totalLen := inBufferLen; n > totalLen {
		err = io.ErrShortBuffer
		return
	} else if n <= 0 {
		n = totalLen
	}
	discarded, err = conn.inbound.Discard(n)
	return
}

func (conn *Conn) InboundBuffered() (n int) {
	return conn.inbound.Buffered()
}

func (conn *Conn) Write(b []byte) (n int, err error) {
	if len(b) <= 0 {
		return
	}
	if conn.outbound.Len() > conn.client.WriteBufferCap {
		err = pnet.ErrOutOfWriteCap
		return
	}
	return conn.outbound.Write(b)
}

func (conn *Conn) ReadFrom(r io.Reader) (n int64, err error) {
	return conn.outbound.ReadFrom(r)
}

func (conn *Conn) Writev(bs [][]byte) (n int, err error) {
	if conn.outbound.Len() > conn.client.WriteBufferCap {
		err = pnet.ErrOutOfWriteCap
		return
	}
	for _, b := range bs {
		var bn int
		bn, err = conn.outbound.Write(b)
		if err != nil {
			return
		}
		n += bn
	}
	return
}

func (conn *Conn) Flush() (err error) {
	return conn.client.loop.flush()
}

func (conn *Conn) OutboundBuffered() (n int) {
	return conn.outbound.Len()
}

func (conn *Conn) AsyncWrite(buf []byte, callback func(c session.Conn, err error) error) (err error) {
	err = conn.client.loop.asyncWrite(buf, callback)
	return
}

func (conn *Conn) AsyncWritev(bs [][]byte, callback func(c session.Conn, err error) error) (err error) {
	err = conn.client.loop.asyncWritev(bs, callback)
	return
}

func (conn *Conn) Close() error {
	return conn.TCPConn.Close()
}

func (conn *Conn) ToClosed(reason error) bool {
	if conn.BaseConn.ToClosed(reason) {
		conn.release()
		return true
	}
	return false
}

func (conn *Conn) release() {
	conn.inbound.Done()
	conn.outbound.Reset()
}
