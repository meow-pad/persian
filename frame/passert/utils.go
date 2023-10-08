package passert

import (
	"go.uber.org/zap"
	"persian/frame/plog"
)

func IsTrue(value bool, msg string, args ...zap.Field) {
	if !value {
		if len(msg) <= 0 {
			msg = "value should be true"
		}
		plog.Panic(msg, args...)
	}
}

func IsFalse(value bool, msg string, args ...zap.Field) {
	if value {
		if len(msg) <= 0 {
			msg = "value should be false"
		}
		plog.Panic(msg, args...)
	}
}

func NotNil(value any, msg string, args ...zap.Field) {
	if value == nil {
		if len(msg) <= 0 {
			msg = "value should not be nil"
		}
		plog.Panic(msg, args...)
	}
}

func IsNil(value any, msg string, args ...zap.Field) {
	if value != nil {
		if len(msg) <= 0 {
			msg = "value should be nil"
		}
		plog.Panic(msg, args...)
	}
}
