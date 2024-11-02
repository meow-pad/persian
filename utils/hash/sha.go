package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1(s string) string {
	sum := sha1.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}
