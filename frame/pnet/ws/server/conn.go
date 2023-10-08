package server

import (
	"bytes"
	"github.com/gobwas/ws"
	"github.com/meow-pad/persian/frame/pnet/tcp/session"
	"github.com/panjf2000/gnet/v2"
)

type wsMessageBuf struct {
	firstHeader   *ws.Header
	curHeader     *ws.Header
	inboundCached bytes.Buffer
}

func newConn(gConn gnet.Conn) (*Conn, error) {
	conn := &Conn{}
	err := conn.Init(gConn)
	if err != nil {
		return nil, err
	}
	return conn, err
}

type Conn struct {
	gnet.Conn
	session.BaseConn

	upgraded bool
	wsMsgBuf wsMessageBuf
}

func (conn *Conn) Init(gConn gnet.Conn) error {
	conn.Conn = gConn
	return conn.BaseConn.Init(gConn, false)
}

func (conn *Conn) AsyncWrite(buf []byte, callback func(c session.Conn, err error) error) (err error) {
	err = conn.Conn.AsyncWrite(buf, func(c gnet.Conn, err error) error {
		return callback(conn, err)
	})
	return
}

func (conn *Conn) AsyncWritev(bs [][]byte, callback func(c session.Conn, err error) error) (err error) {
	err = conn.Conn.AsyncWritev(bs, func(c gnet.Conn, err error) error {
		return callback(conn, err)
	})
	return
}
