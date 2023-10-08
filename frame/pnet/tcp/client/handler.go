package client

import (
	"persian/frame/plog"
	"persian/frame/plog/cfield"
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
		cfield.String("client", conn.client.Name),
		cfield.Uint64("conn", conn.Hash()))
	return
}

func (handler *eventHandler) OnClose(conn *Conn, err error) (action action) {
	conn.ToClosed(err)
	conn.client.listener.OnClosed(conn.client)
	plog.Debug("close connecting:",
		cfield.String("client", conn.client.Name),
		cfield.Uint64("conn", conn.Hash()))
	return
}

func (handler *eventHandler) OnTraffic(conn *Conn) (action action) {
	msgArr, totalLen, err := conn.client.codec.Decode(conn)
	if err != nil {
		plog.Error("decode error", cfield.Error(err))
		action = actionClose
		return
	}
	msgNum := len(msgArr)
	if msgNum > 1 {
		_ = conn.client.listener.OnReceiveMulti(conn.client, msgArr, totalLen)
	} else if msgNum == 1 {
		_ = conn.client.listener.OnReceive(conn.client, msgArr[0], totalLen)
	}
	return
}
