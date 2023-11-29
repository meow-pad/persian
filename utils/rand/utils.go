package rand

import (
	"github.com/valyala/fastrand"
	"math"
)

func Uint32() uint32 {
	return fastrand.Uint32()
}

// Uint32n
//
//	@Description: 返回一个非负的32位整数，随机数范围为[0, maxN)
//	@param maxN
//	@return uint32
func Uint32n(maxN uint32) uint32 {
	return fastrand.Uint32n(maxN)
}

// Int32
//
//	@Description: 返回一个非负的32位整数，随机数范围为[0, math.MaxInt32]
//	@return int32
func Int32() int32 {
	return int32(fastrand.Uint32() << 1 >> 1)
}

// Int32n
//
//	@Description: 返回一个非负的32位整数，随机数范围为[0, maxN)
//	@param maxN
//	@return int32
func Int32n(maxN int32) int32 {
	if maxN <= 0 {
		panic("maxN should be >= 0")
	}
	return int32(fastrand.Uint32n(uint32(maxN)))
}

func Uint64() uint64 {
	return uint64(fastrand.Uint32()) | (uint64(fastrand.Uint32()) << 32)
}

func Uint64n(maxN uint64) uint64 {
	return uint64(Float64() * float64(maxN))
}

func Int64() int64 {
	return int64(Uint64())
}

func Float32() float32 {
	return float32(fastrand.Uint32()) / math.MaxUint32
}

func Float64() float64 {
	return float64(Uint64()) / math.MaxUint64
}
