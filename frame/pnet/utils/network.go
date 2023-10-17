package utils

import (
	"fmt"
	"github.com/meow-pad/persian/errdef"
	"strings"
)

const (
	ProtoTCP       = "tcp"
	ProtoUDP       = "udp"
	ProtoWebsocket = "ws"
)

const separator = "://"

// CompleteAddress
//
//	@Description: 补完地址
//	@param address 原始地址，格式如 `192.168.0.10:9851` 或 `tcp://192.168.0.10:9851`
//	@param proto  协议，如 `tcp`、`udp`
//	@return string 地址格式如 `tcp://192.168.0.10:9851`
//	@return error
func CompleteAddress(address, proto string) (string, error) {
	address = strings.ToLower(address)
	index := strings.Index(address, separator)
	if index >= 0 {
		addrProto := address[:index]
		if addrProto != proto {
			return "", fmt.Errorf("protocol(%s) not match protocol(%s)", addrProto, proto)
		}
	} else {
		address = fmt.Sprintf("%s%s%s", proto, separator, address)
	}
	return address, nil
}

// GetAddress
//
//	@Description: 分拆协议和地址
//	@param protoAddr 格式如：`tcp://192.168.0.10:9851`
//	@return proto
//	@return address
//	@return err
func GetAddress(protoAddr string) (proto string, address string, err error) {
	protoAddr = strings.ToLower(protoAddr)
	index := strings.Index(protoAddr, separator)
	if index >= 0 {
		proto = protoAddr[:index]
		address = protoAddr[index+len(separator):]
		return
	} else {
		err = errdef.ErrInvalidParams
		return
	}
}
