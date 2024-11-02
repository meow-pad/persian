package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"github.com/meow-pad/persian/utils/loggers"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(value any) ([]byte, error) {
	return json.Marshal(value)
}

func Unmarshal(data []byte, value any) error {
	if len(data) <= 0 {
		return nil
	}
	return json.Unmarshal(data, value)
}

func ToString(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		loggers.Error("marshal error:", pfield.Error(err))
	}
	return string(data)
}
