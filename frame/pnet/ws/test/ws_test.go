package test

import (
	"context"
	"github.com/stretchr/testify/require"
	"persian/frame/pnet/message"
	"persian/frame/pnet/tcp/session"
	"persian/frame/pnet/ws/client"
	"persian/frame/pnet/ws/server"
	"testing"
	"time"
)

func newCodec() message.Codec {
	return &message.TextCodec{}
}

func newServer(protoAddr string, listener session.Listener) (*server.Server, error) {
	return server.NewServer("test-server", protoAddr, newCodec(), listener)
}

func newClient(listener session.Listener) (*client.Client, error) {
	return client.NewClient(newCodec(), listener, client.WithName("test-client"))
}

type svrListener struct {
	session.EmptyListener
	t *testing.T
}

func (listener *svrListener) OnOpened(session session.Session) {
	listener.t.Logf("ws-svr open conn:%v", session.Connection().Hash())
}

func (listener *svrListener) OnClosed(session session.Session) {
	listener.t.Logf("ws-svr close conn:%v", session.Connection().Hash())
}

func (listener *svrListener) OnReceive(session session.Session, msg any, msgLen int) (err error) {
	listener.t.Logf("ws-svr received msg:%v", msg)
	session.SendMessage(msg)
	return nil
}

func (listener *svrListener) OnReceiveMulti(session session.Session, msgArr []any, totalLen int) error {
	for _, msg := range msgArr {
		listener.t.Logf("ws-svr received _msg:%v", msg)
		session.SendMessage(msg)
	}
	return nil
}

type cliListener struct {
	session.EmptyListener
	t *testing.T
}

func (listener *cliListener) OnOpened(session session.Session) {
	listener.t.Logf("ws-cli open conn:%v", session.Connection().Hash())
}

func (listener *cliListener) OnClosed(session session.Session) {
	listener.t.Logf("ws-cli close conn:%v", session.Connection().Hash())
}

func (listener *cliListener) OnReceive(session session.Session, msg any, msgLen int) (err error) {
	listener.t.Logf("ws-cli received msg:%v", msg)
	return nil
}

func (listener *cliListener) OnReceiveMulti(session session.Session, msgArr []any, totalLen int) error {
	for _, msg := range msgArr {
		listener.t.Logf("ws-cli received _msg:%v", msg)
	}
	return nil
}

func TestWS_Echo(t *testing.T) {
	should := require.New(t)
	addr := "127.0.0.1:9080"
	echoSvr, err := newServer(addr, &svrListener{t: t})
	should.Nil(err)
	err = echoSvr.Start(context.Background())
	should.Nil(err)
	echoCLi, err := newClient(&cliListener{t: t})
	should.Nil(err)
	err = echoCLi.Dial(context.Background(), addr)
	should.Nil(err)
	echoCLi.SendMessage("123")
	echoCLi.SendMessage("456")
	echoCLi.SendMessage("789")
	time.Sleep(2 * time.Second)
	err = echoCLi.Close()
	should.Nil(err)
	err = echoSvr.Stop(context.Background())
	should.Nil(err)
}
