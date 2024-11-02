package hash

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func UpperMD5(s string) string {
	return strings.ToUpper(MD5(s))
}
