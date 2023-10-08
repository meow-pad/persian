package utils

import (
	"fmt"
	"strings"
)

const (
	ProtoTCP       = "tcp"
	ProtoUDP       = "udp"
	ProtoWebsocket = "ws"
)

// CompleteAddress
//
//	@Description: 补完地址
//	@param address 原始地址，格式如 `192.168.0.10:9851` 或 `tcp://192.168.0.10:9851`
//	@param proto  协议，如 `tcp`、`udp`
//	@return string 地址格式如 `tcp://192.168.0.10:9851`
//	@return error
func CompleteAddress(address, proto string) (string, error) {
	address = strings.ToLower(address)
	index := strings.Index(address, "://")
	if index >= 0 {
		addrProto := address[:index]
		if addrProto != proto {
			return "", fmt.Errorf("protocol(%s) not match protocol(%s)", addrProto, proto)
		}
	} else {
		address = fmt.Sprintf("%s://%s", proto, address)
	}
	return address, nil
}
