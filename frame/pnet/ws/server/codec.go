package server

import (
	"bytes"
	"errors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/meow-pad/persian/errdef"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"github.com/meow-pad/persian/frame/pnet/message"
	"github.com/panjf2000/gnet/v2"
	"io"
)

type wsReadWrite struct {
	io.Reader
	io.Writer
}

func newWsCodec(msgCodec message.Codec) (*wsCodec, error) {
	if msgCodec == nil {
		return nil, errdef.ErrInvalidParams
	}
	return &wsCodec{msgCodec: msgCodec}, nil
}

type wsCodec struct {
	msgCodec message.Codec
}

func (codec *wsCodec) Encode(msg any) (out []byte, err error) {
	out, err = codec.msgCodec.Encode(msg)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(make([]byte, 0, len(out)+ws.MaxHeaderSize))
	err = wsutil.WriteServerMessage(buf, ws.OpBinary, out)
	if err != nil {
		out = nil
		return
	}
	out = buf.Bytes()
	return
}

func (codec *wsCodec) Decode(conn *Conn) ([]any, int, gnet.Action) {
	if !conn.upgraded {
		ok, action := codec.upgrade(conn)
		if action == gnet.Close {
			return nil, 0, gnet.Close
		}
		if !ok {
			return nil, 0, gnet.None
		}
	}
	if conn.InboundBuffered() <= 0 {
		return nil, 0, gnet.None
	}
	wsMsgArr, err := codec.readWsMessages(conn)
	if err != nil {
		plog.Error("read ws messages error:", pfield.Uint64("conn", conn.Hash()), pfield.Error(err))
	}
	arrLen := len(wsMsgArr)
	if arrLen <= 0 {
		return nil, 0, gnet.None
	}
	msgArr := make([]any, 0, arrLen)
	totalLen := 0
	for _, wsMsg := range wsMsgArr {
		switch wsMsg.OpCode {
		case ws.OpText, ws.OpBinary:
			totalLen += len(wsMsg.Payload)
			msg, dErr := codec.msgCodec.Decode(wsMsg.Payload)
			if dErr != nil {
				plog.Error("decode message error:", pfield.Uint64("conn", conn.Hash()), pfield.Error(dErr))
			} else {
				msgArr = append(msgArr, msg)
			}
		case ws.OpClose:
			err = wsutil.WriteServerMessage(conn, ws.OpClose, nil)
			if err != nil {
				return nil, 0, gnet.Close
			}
		case ws.OpPing:
			err = wsutil.WriteServerMessage(conn, ws.OpPong, nil)
			if err != nil {
				return nil, 0, gnet.Close
			}
		case ws.OpPong:
		default:
			plog.Warn("unknown ws opCode", pfield.Uint8("opCoded", uint8(wsMsg.OpCode)))
		}
	}
	return msgArr, totalLen, gnet.None
}

func (codec *wsCodec) upgrade(conn *Conn) (ok bool, action gnet.Action) {
	size := conn.InboundBuffered()
	buf, err := conn.Peek(size)
	if err != nil {
		plog.Error("peek bytes error:", pfield.Uint64("conn", conn.Hash()), pfield.Error(err))
		action = gnet.Close
		return
	}
	read := len(buf)
	if read < size {
		plog.Error("peek bytes len error",
			pfield.Uint64("conn", conn.Hash()),
			pfield.Int("need", size), pfield.Int("read", read))
		action = gnet.Close
		return
	}
	tmpReader := bytes.NewReader(buf)
	oldLen := tmpReader.Len()
	hs, err := ws.Upgrade(wsReadWrite{tmpReader, conn})
	skipN := oldLen - tmpReader.Len()
	if err != nil {
		if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
			//数据不完整
			return
		}
		plog.Error("upgraded conn error", pfield.Uint64("conn", conn.Hash()), pfield.Error(err))
		action = gnet.Close
		return
	}
	if _, err = conn.Discard(skipN); err != nil {
		plog.Error("discard bytes error", pfield.Uint64("conn", conn.Hash()), pfield.Error(err))
		action = gnet.Close
		return
	}
	ok = true
	conn.upgraded = true
	plog.Debug("upgraded websocket protocol",
		pfield.Uint64("conn", conn.Hash()), pfield.Any("handshake", hs))
	return
}

func (codec *wsCodec) readWsMessages(conn *Conn) (messages []wsutil.Message, err error) {
	msgBuf := &conn.wsMsgBuf
	for {
		if msgBuf.curHeader == nil {
			if conn.InboundBuffered() < ws.MinHeaderSize {
				//头长度至少是2
				return
			}
			var head ws.Header
			if conn.InboundBuffered() >= ws.MaxHeaderSize {
				head, err = ws.ReadHeader(conn)
				if err != nil {
					return
				}
			} else {
				//有可能不完整，构建新的 reader 读取 head 读取成功才实际对 in 进行读操作
				var pBuf []byte
				pBuf, err = conn.Peek(conn.InboundBuffered())
				if err != nil {
					return
				}
				tmpReader := bytes.NewReader(pBuf)
				oldLen := tmpReader.Len()
				head, err = ws.ReadHeader(tmpReader)
				skipN := oldLen - tmpReader.Len()
				if err != nil {
					if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
						//数据不完整
						err = nil
						return
					}
					if _, dErr := conn.Discard(skipN); dErr != nil {
						plog.Error("discard bytes error:",
							pfield.Uint64("conn", conn.Hash()), pfield.Error(err))
					}
					return
				}
				if _, err = conn.Discard(skipN); err != nil {
					return
				}
			}

			msgBuf.curHeader = &head
			err = ws.WriteHeader(&msgBuf.inboundCached, head)
			if err != nil {
				return
			}
		}
		dataLen := msgBuf.curHeader.Length
		if dataLen > 0 {
			if (int64)(conn.InboundBuffered()) >= dataLen {
				_, err = io.CopyN(&msgBuf.inboundCached, conn, dataLen)
				if err != nil {
					return
				}
			} else {
				//数据不完整
				return
			}
		}
		if msgBuf.curHeader.Fin {
			//当前 header 已经是一个完整消息
			messages, err = wsutil.ReadClientMessage(&msgBuf.inboundCached, messages)
			if err != nil {
				return nil, err
			}
			msgBuf.inboundCached.Reset()
		} else {
			plog.Debug("data is split into multiple frames", pfield.Uint64("conn", conn.Hash()))
		}
		msgBuf.curHeader = nil
	}
}
