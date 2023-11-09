package hash

import (
	"github.com/spaolacci/murmur3"
)

func Murmur3Hash32(data []byte) (uint32, error) {
	h := murmur3.New32()
	_, err := h.Write(data)
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}

func Murmur3Hash64(data []byte) (uint64, error) {
	h := murmur3.New64()
	_, err := h.Write(data)
	if err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}
