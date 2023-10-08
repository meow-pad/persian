package codec

import (
	"github.com/panjf2000/gnet/v2"
	"math"
	"persian/frame/pnet/message"
)

const (
	DefaultMaxMessageLength     int = 8 * 1024
	DefaultWarningEncodedMsgLen int = 4 * 1024
)

// Codec
//
//	@Description: 编解码器
type Codec interface {

	// Encode
	//	@Description: 消息编码
	//	@param msg 输入消息
	//	@return []bytes 编码结果输出
	//	@return error
	//
	Encode(msg any) ([]byte, error)

	// Decode
	//	@Description: 消息解码
	//	@param []byte 输入
	//	@return []Msg 解码消息输出
	//	@return totalLen 消息所占字节大小
	//	@return error
	//
	Decode(reader gnet.Reader) ([]any, int, error)
}

type Options interface {
	GetMagicBytes() []byte
	SetMagicBytes([]byte)
	GetMagicSize() int

	GetMaxDecodedLength() int
	SetMaxDecodedLength(int)

	GetMaxEncodedLength() int
	SetMaxEncodedLength(int)

	GetWarningEncodedLength() int
	SetWarningEncodedLength(int)

	GetMessageCodec() message.Codec
	SetMessageCodec(message.Codec)
}

type Option[T Options] func(o T)

type OptionsBase struct {
	// 魔数
	magicBytes []byte
	// 魔术长度
	magicSize int
	// 最大解码消息长度
	maxDecodedLength int
	// 最大编码消息长度
	maxEncodedLength int
	// 编码消息长度告警
	warningEncodedLength int
	// 消息解析器
	messageCodec message.Codec
}

func (opts *OptionsBase) Complete() error {
	opts.magicSize = len(opts.magicBytes)
	if opts.maxDecodedLength == 0 {
		opts.maxDecodedLength = DefaultMaxMessageLength
	} else if opts.maxDecodedLength < 0 {
		opts.maxDecodedLength = math.MaxInt
	}

	if opts.maxEncodedLength == 0 {
		opts.maxEncodedLength = DefaultMaxMessageLength
	} else if opts.maxEncodedLength < 0 {
		opts.maxEncodedLength = math.MaxInt
	}

	if opts.warningEncodedLength == 0 {
		opts.warningEncodedLength = DefaultWarningEncodedMsgLen
	} else if opts.warningEncodedLength < 0 {
		opts.warningEncodedLength = math.MaxInt
	}

	if opts.messageCodec == nil {
		opts.messageCodec = &message.TextCodec{}
	}
	return nil
}

func (opts *OptionsBase) GetMagicBytes() []byte {
	return opts.magicBytes
}

func (opts *OptionsBase) SetMagicBytes(value []byte) {
	opts.magicBytes = value
}

func (opts *OptionsBase) GetMagicSize() int {
	return opts.magicSize
}

func (opts *OptionsBase) GetMaxDecodedLength() int {
	return opts.maxDecodedLength
}

func (opts *OptionsBase) SetMaxDecodedLength(value int) {
	opts.maxEncodedLength = value
}

func (opts *OptionsBase) GetMaxEncodedLength() int {
	return opts.maxEncodedLength
}

func (opts *OptionsBase) SetMaxEncodedLength(value int) {
	opts.maxEncodedLength = value
}

func (opts *OptionsBase) GetWarningEncodedLength() int {
	return opts.warningEncodedLength
}

func (opts *OptionsBase) SetWarningEncodedLength(value int) {
	opts.warningEncodedLength = value
}

func (opts *OptionsBase) GetMessageCodec() message.Codec {
	return opts.messageCodec
}

func (opts *OptionsBase) SetMessageCodec(value message.Codec) {
	opts.messageCodec = value
}

func WithMagic[T Options](magicBytes []byte) Option[T] {
	return func(options T) {
		options.SetMagicBytes(magicBytes)
	}
}

func WithMaxDecodedLength[T Options](maxDecodedLength int) Option[T] {
	return func(options T) {
		options.SetMaxDecodedLength(maxDecodedLength)
	}
}

func WithMaxEncodedLength[T Options](maxEncodedLength int) Option[T] {
	return func(options T) {
		options.SetMaxEncodedLength(maxEncodedLength)
	}
}

func WithWarningEncodedLength[T Options](warningEncodedLength int) Option[T] {
	return func(options T) {
		options.SetWarningEncodedLength(warningEncodedLength)
	}
}

func WithMessageCodec[T Options](messageCodec message.Codec) Option[T] {
	return func(options T) {
		options.SetMessageCodec(messageCodec)
	}
}
