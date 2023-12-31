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

// Uint64n
//
//	@Description: 返回一个非负的64位整数，随机数范围为[0, maxN)
//	@param maxN
//	@return uint64
func Uint64n(maxN uint64) uint64 {
	return uint64(Float64() * float64(maxN))
}

// Int64
//
//	@Description: 返回一个非负的64位整数，随机数范围为[0, math.MaxInt64]
//	@return int64
func Int64() int64 {
	return int64(Uint64() << 1 >> 1)
}

// Int64n
//
//	@Description: 返回一个非负的64位整数，随机数范围为[0, maxN)
//	@param maxN
//	@return int64
func Int64n(maxN int64) int64 {
	if maxN <= 0 {
		panic("maxN should be >= 0")
	}
	return int64(Uint64n(uint64(maxN)))
}

func Float32() float32 {
	return float32(fastrand.Uint32()) / math.MaxUint32
}

// Float32n
//
//	@Description: 返回一个非负的32浮点数，随机数范围为[0, maxN)
//	@param maxN
//	@return float32
func Float32n(maxN float32) float32 {
	return Float32() * maxN
}

func Float64() float64 {
	return float64(Uint64()) / math.MaxUint64
}

// Float64n
//
//	@Description: 返回一个非负的64浮点数，随机数范围为[0, maxN)
//	@param maxN
//	@return float64
func Float64n(maxN float64) float64 {
	return Float64() * maxN
}
