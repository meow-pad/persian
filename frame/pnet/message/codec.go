package message

import "persian/errdef"

// Codec
//
//	@Description: 消息内容解析器
type Codec interface {
	// Encode
	//	@Description: 编码指定消息
	//	@param msg 消息对象
	//	@return []byte 消息编码对象
	//	@return error
	//
	Encode(msg any) ([]byte, error)

	// Decode
	//	@Description: 解码消息
	//	@param in 输入字节数据
	//	@return Msg 消息对象
	//	@return error
	//
	Decode(in []byte) (any, error)
}

type TextCodec struct {
}

func (codec *TextCodec) Encode(msg any) ([]byte, error) {
	text, ok := msg.(string)
	if !ok {
		return nil, errdef.ErrInvalidParams
	}
	return []byte(text), nil
}

func (codec *TextCodec) Decode(in []byte) (any, error) {
	return string(in), nil
}
