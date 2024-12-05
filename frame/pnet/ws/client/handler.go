package client

import (
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
)

type action uint32

const (
	actionNone action = iota
	actionClose
)

func newEventHandler() *eventHandler {
	return &eventHandler{}
}

type eventHandler struct {
}

func (handler *eventHandler) OnOpen(conn *Conn) (action action) {
	conn.client.listener.OnOpened(conn.client)
	plog.Debug("open connecting:",
		pfield.String("client", conn.client.Name),
		pfield.Uint64("conn", conn.Hash()))
	return
}

func (handler *eventHandler) OnClose(conn *Conn, err error) (action action) {
	conn.ToClosed(err)
	conn.client.listener.OnClosed(conn.client)
	plog.Debug("close ws-client connecting:",
		pfield.String("client", conn.client.Name),
		pfield.Uint64("conn", conn.Hash()))
	return
}

func (handler *eventHandler) OnTraffic(conn *Conn, data []byte) (action action) {
	msg, err := conn.client.codec.Decode(data)
	if err != nil {
		plog.Error("decode error", pfield.Error(err))
		action = actionClose
		return
	}
	_ = conn.client.listener.OnReceive(conn.client, msg, len(data))
	return
}
