package hash

import (
	"github.com/spaolacci/murmur3"
	"hash/fnv"
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

// Fnv32 更快，但唯一性差
func Fnv32(data []byte) (uint32, error) {
	h32 := fnv.New32()
	_, err := h32.Write(data)
	if err != nil {
		return 0, err
	}
	return h32.Sum32(), nil
}

// Fnv64 更快，但唯一性差
func Fnv64(data []byte) (uint64, error) {
	h64 := fnv.New64()
	_, err := h64.Write(data)
	if err != nil {
		return 0, err
	}
	return h64.Sum64(), nil
}
