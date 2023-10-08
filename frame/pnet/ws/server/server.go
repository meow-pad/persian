package server

import (
	"context"
	"errors"
	"github.com/panjf2000/gnet/v2"
	"persian/frame/plog"
	"persian/frame/plog/cfield"
	"persian/frame/pnet/message"
	"persian/frame/pnet/tcp/session"
	"persian/frame/pnet/utils"
)

func NewServer(name string, protoAddr string,
	codec message.Codec, listener session.Listener, opts ...Option) (server *Server, err error) {
	var options *Options
	options, err = NewOptions(opts...)
	if err != nil {
		return
	}
	protoAddr, err = utils.CompleteAddress(protoAddr, utils.ProtoTCP)
	if err != nil {
		return
	}
	var swCodec *wsCodec
	if swCodec, err = newWsCodec(codec); err != nil {
		return
	}
	if listener == nil {
		err = errors.New("less listener")
		return
	}
	var manager *session.Manager
	if manager, err = session.NewManager(name, options.UnregisterSessionLife); err != nil {
		return
	}
	return &Server{
		Manager:   manager,
		options:   options,
		name:      name,
		protoAddr: protoAddr,
		codec:     swCodec,
		listener:  listener,
	}, nil
}

type Server struct {
	*session.Manager
	options *Options

	// 服务名称
	name string
	// 带协议地址
	protoAddr string
	// gnet引擎
	engine gnet.Engine
	// ws编解码器
	codec *wsCodec
	// 会话监听器
	listener session.Listener
	// 启动通知
	startChan chan gnet.Engine
}

func (server *Server) Start(ctx context.Context) error {
	server.startChan = make(chan gnet.Engine)
	go func() {
		err := gnet.Run(
			&eventHandler{server: server},
			server.protoAddr,
			gnet.WithOptions(server.options.GNetOptions))
		if err != nil {
			plog.Error("run server error:", cfield.String("server", server.name), cfield.Error(err))
		}
	}()
	select {
	case eng := <-server.startChan:
		if err := eng.Validate(); err != nil {
			return err
		}
		server.engine = eng
	case <-ctx.Done():
		return errors.New("start server timeout")
	}
	return nil
}

func (server *Server) Stop(ctx context.Context) error {
	return server.engine.Stop(ctx)
}

func (server *Server) CName() string {
	return server.name
}
