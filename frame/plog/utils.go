package plog

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"persian/utils/coding"
	"strings"
	"time"
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

// 添加颜色描述
func addColorDesc(color int, s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(color), s)
}

// 带颜色通用格式的时间编码
func colorTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encodeTimeLayout(t, "2006-01-02 15:04:05.000Z0700", enc, true)
}

// 通用格式时间编码
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encodeTimeLayout(t, "2006-01-02 15:04:05.000Z0700", enc, false)
}

// 时间编码
func encodeTimeLayout(t time.Time, layout string, enc zapcore.PrimitiveArrayEncoder, withColor bool) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if tEnc, ok := enc.(appendTimeEncoder); ok {
		tEnc.AppendTimeLayout(t, layout)
		return
	}

	timeStr := t.Format(layout)
	if withColor {
		timeStr = addColorDesc(colorDarkGray, timeStr)
	}
	enc.AppendString(timeStr)
}

// stackTracer
//
//	@Description: 堆栈格式输出
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// errorFormatState
//
//	@Description: 错误格式状态
type errorFormatState struct {
	buf []byte
}

// Write implement fmt.State interface.
func (s *errorFormatState) Write(buf []byte) (n int, err error) {
	s.buf = append(s.buf, buf...)
	return len(buf), nil
}

// Width implement fmt.State interface.
func (s *errorFormatState) Width() (wid int, ok bool) {
	return 0, false
}

// Precision implement fmt.State interface.
func (s *errorFormatState) Precision() (precision int, ok bool) {
	return 0, false
}

// Flag implement fmt.State interface.
func (s *errorFormatState) Flag(c int) bool {
	if c == '+' {
		return true
	}
	return false
}

// Clean 清理缓存信息
func (s *errorFormatState) Clean() {
	s.buf = s.buf[:0]
	return
}

// 错误堆栈格式化字符串
func errorFrameField(f errors.Frame, s *errorFormatState, c rune) string {
	f.Format(s, c)
	return string(s.buf)
}

// 构造堆栈格式化字符串
func marshalStack(err error) any {
	var stErr stackTracer
	// 找到第一个内部包装的可用错误
	for {
		if st, ok := err.(stackTracer); ok {
			stErr = st
			break
		} else {
			err = errors.Unwrap(err)
			if err == nil {
				return nil
			}
		} // end of else
	} // end of for
	st := stErr.StackTrace()
	str := strings.Builder{}
	state := &errorFormatState{}
	for _, frame := range st {
		str.WriteString("\n")
		str.WriteString(errorFrameField(frame, state, 'v'))
		state.Clean()
	}
	str.WriteRune('\n')
	return addColorDesc(colorRed, str.String())
}

// 错误字段格式整理
func encodeErrorFieldLayout(msg string, fields ...zap.Field) (string, []zap.Field) {
	for i := 0; i < len(fields); {
		field := fields[i]
		if field.Type == zapcore.ErrorType {
			if err := coding.Cast[error](field.Interface); err != nil {
				//if err, ok := field.Interface.(error); ok {
				fields = append(fields[:i], fields[i+1:]...)
				stack := marshalStack(err)
				if stack == nil {
					msg += fmt.Sprintf("\n%s:\n%s", field.Key, err.Error())
				} else {
					msg += fmt.Sprintf("\n%s:\n%s%v", field.Key, err.Error(), stack)
				}
			} // end of if
		} else {
			i++
		}
	} // end of for
	return msg, fields
}

// 合并字段数据到消息
func toMessageWithFields(msg string, fields ...zap.Field) string {
	if len(fields) <= 0 {
		return msg
	}
	builder := strings.Builder{}
	builder.WriteString(msg)
	builder.WriteRune(':')
	for _, field := range fields {
		builder.WriteString(fmt.Sprintf("(\"%s\",%v)", field.Key, field))
	}
	return builder.String()
}
