package client

import (
	"github.com/gorilla/websocket"
	"github.com/meow-pad/persian/frame/pnet/tcp/session"
	"io"
	"time"
)

func NewConn(client *Client, pConn *websocket.Conn) (*Conn, error) {
	conn := &Conn{}
	err := conn.Init(client, pConn)
	if err != nil {
		return nil, err
	}
	return conn, nil

}

type Conn struct {
	*websocket.Conn
	session.BaseConn

	client *Client
}

func (conn *Conn) Init(client *Client, pConn *websocket.Conn) error {
	conn.Conn = pConn
	conn.client = client
	if conn.client.MaxMessageLength > 0 {
		conn.SetReadLimit(conn.client.MaxMessageLength)
	}
	//conn.SetCloseHandler(nil)
	//conn.SetPingHandler(nil)
	return conn.BaseConn.Init(conn, true)
}

func (conn *Conn) WriteTo(w io.Writer) (n int64, err error) {
	panic("not implemented")
}

func (conn *Conn) Next(n int) (buf []byte, err error) {
	panic("not implemented")
}

func (conn *Conn) Peek(n int) (buf []byte, err error) {
	panic("not implemented")
}

func (conn *Conn) Discard(n int) (discarded int, err error) {
	panic("not implemented")
}

func (conn *Conn) InboundBuffered() (n int) {
	panic("not implemented")
}

func (conn *Conn) ReadFrom(r io.Reader) (n int64, err error) {
	panic("not implemented")
}

func (conn *Conn) Writev(bs [][]byte) (n int, err error) {
	panic("not implemented")
}

func (conn *Conn) Flush() (err error) {
	panic("not implemented")
}

func (conn *Conn) OutboundBuffered() (n int) {
	panic("not implemented")
}

func (conn *Conn) AsyncWrite(buf []byte, callback func(c session.Conn, err error) error) (err error) {
	panic("not implemented")
}

func (conn *Conn) AsyncWritev(bs [][]byte, callback func(c session.Conn, err error) error) (err error) {
	panic("not implemented")
}

func (conn *Conn) Context() (ctx any) {
	panic("not implemented")
}

func (conn *Conn) SetContext(ctx any) {
	panic("not implemented")
}

func (conn *Conn) Read(_ []byte) (n int, err error) {
	panic("not implemented")
}

func (conn *Conn) Write(_ []byte) (n int, err error) {
	panic("not implemented")
}

func (conn *Conn) SetDeadline(t time.Time) error {
	panic("not implemented")
}
