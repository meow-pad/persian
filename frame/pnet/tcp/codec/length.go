package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/v2"
	"io"
	"persian/frame/plog"
	"persian/frame/plog/cfield"
	"persian/frame/pnet"
	"persian/utils/numeric"
)

var (
	DefaultLengthByteOrder = binary.LittleEndian
	DefaultLengthSize      = 2
	MaxLengthSize          = 4
)

type LengthOptions struct {
	OptionsBase
	// 字节序
	ByteOrder binary.ByteOrder
	// 消息长度所占字节数
	LengthSize int
}

func (opts *LengthOptions) Complete() error {
	if err := opts.OptionsBase.Complete(); err != nil {
		return err
	}
	if opts.ByteOrder == nil {
		opts.ByteOrder = DefaultLengthByteOrder
	}
	if opts.LengthSize <= 0 {
		opts.LengthSize = DefaultLengthSize
	}
	opts.LengthSize = numeric.Min[int](opts.LengthSize, MaxLengthSize)
	maxMsgLen := 1<<(opts.LengthSize*8-1) - 1
	opts.maxDecodedLength = numeric.Min(maxMsgLen, opts.maxDecodedLength)
	opts.maxEncodedLength = numeric.Min(maxMsgLen, opts.maxEncodedLength)
	opts.warningEncodedLength = numeric.Min(maxMsgLen, opts.warningEncodedLength)
	return nil
}

func WithByteOrder(value binary.ByteOrder) Option[*LengthOptions] {
	return func(options *LengthOptions) {
		options.ByteOrder = value
	}
}

func WithLengthSize(value int) Option[*LengthOptions] {
	return func(options *LengthOptions) {
		options.LengthSize = value
	}
}

func NewLengthFieldCodec(opts ...Option[*LengthOptions]) (Codec, error) {
	options := &LengthOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if err := options.Complete(); err != nil {
		return nil, err
	}
	return &lengthFieldCodec{
		LengthOptions: options,
	}, nil
}

// lengthFieldCodec
//
//	@Description: 带长度编码
//
// * 0       magicSize               lengthSize
// * +-----------+-----------------------+
// * |   magic   |       body len        |
// * +-----------+-----------+-----------+
// * |                                   |
// * +                                   +
// * |           body bytes              |
// * +                                   +
// * |            ... ...                |
// * +-----------------------------------+
type lengthFieldCodec struct {
	*LengthOptions
}

func (codec *lengthFieldCodec) Encode(msg any) (out []byte, err error) {
	if msg == nil {
		err = pnet.ErrNilMessage
		return
	}
	var bodyBuf []byte
	bodyBuf, err = codec.messageCodec.Encode(msg)
	if err != nil {
		return
	}
	bodyLen := len(bodyBuf)
	if bodyLen > 0 {
		if bodyLen > codec.maxEncodedLength {
			err = pnet.ErrMessageTooLarge
			return
		}
		if bodyLen > codec.warningEncodedLength {
			plog.Warn("encoded message is too long", cfield.Int("bodyLen", bodyLen))
		}
		// 写入魔数、消息长度和消息内容
		bodyOffset := codec.magicSize + codec.LengthSize
		msgLen := bodyOffset + bodyLen
		out = make([]byte, msgLen)
		copy(out, codec.magicBytes)
		codec.ByteOrder.PutUint32(out[codec.magicSize:bodyOffset], uint32(bodyLen))
		copy(out[bodyOffset:msgLen], bodyBuf)
		return
	} else {
		// ??? 为什么会出现长度为0的情况
		err = pnet.ErrEmptyEncodeBuffer
		return
	}
}

func (codec *lengthFieldCodec) Decode(reader gnet.Reader) (result []any, totalLen int, err error) {
	for {
		var msg any
		var msgLen int
		msg, msgLen, err = codec.decodeOne(reader)
		if err != nil {
			// 数据不足，稍后读取
			if errors.Is(err, io.ErrShortBuffer) {
				err = nil
			}
			return
		}
		if msg != nil {
			if result == nil {
				result = make([]any, 0)
			}
			result = append(result, msg)
			totalLen += msgLen
		}
		if reader.InboundBuffered() > 0 {
			// 还有可解析数据，尝试继续解析
		} else {
			break
		}
	}
	return
}

// decodeOne
//
//	@Description: 从数据中解析一条结果
//	@receiver codec 解码器
//	@param reader 读数据
//	@return msg 解析结果
//	@return msgLen 消息长度
//	@return err 解析异常
func (codec *lengthFieldCodec) decodeOne(reader gnet.Reader) (msg any, msgLen int, err error) {
	var (
		bodyOffset int
		headerBuf  []byte
	)
	// 判定消息头
	bodyOffset = codec.magicSize + codec.LengthSize
	headerBuf, err = reader.Peek(bodyOffset)
	if err != nil || headerBuf == nil {
		return
	}
	if !bytes.Equal(codec.magicBytes, headerBuf[:codec.magicSize]) {
		err = pnet.ErrInvalidMagic
		return
	}
	bodyLen := int(codec.ByteOrder.Uint32(headerBuf[codec.magicSize:bodyOffset]))
	if bodyLen > codec.maxDecodedLength {
		err = pnet.ErrMessageTooLarge
		return
	}
	// 读消息体
	msgLen = bodyOffset + bodyLen
	var msgBuf []byte
	msgBuf, err = reader.Peek(msgLen)
	if err != nil || msgBuf == nil {
		return
	}
	_, err = reader.Discard(msgLen)
	if err != nil {
		return
	}
	msg, err = codec.messageCodec.Decode(msgBuf[bodyOffset:msgLen])
	return
}
