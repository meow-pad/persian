package rand

import (
	"github.com/meow-pad/persian/utils/numeric"
	"math"
	"testing"
)

func variance(targetRate float32, counts []float32, times float32) (variance float32, maxDiff float32) {
	for i := 0; i < len(counts); i++ {
		rate := counts[i] / times
		diff := numeric.Abs[float32](targetRate - rate)
		maxDiff = numeric.Max[float32](maxDiff, diff)
		variance += (diff / targetRate) * (diff / targetRate) // 相对方差
	}
	variance /= float32(len(counts))
	return
}

func TestUint64(t *testing.T) {
	partition := uint64(100)
	count := make([]float32, partition)
	size := math.MaxUint64 / partition
	times := 1_000_000
	for i := 0; i < times; i++ {
		num := Uint64()
		index := num / size
		count[index] += 1
	}
	expectRate := 1 / float32(partition)
	rateVariance, maxDiff := variance(expectRate, count, float32(times))
	t.Logf("variance %f ,max diff %f", rateVariance, maxDiff)
}

func TestUint64n(t *testing.T) {
	maxN := 100
	times := 1_000_000
	count := make([]float32, maxN)
	for i := 0; i < times; i++ {
		num := Uint64n(uint64(maxN))
		count[num] += 1
	}
	expectRate := 1 / float32(maxN)
	rateVariance, maxDiff := variance(expectRate, count, float32(times))
	t.Logf("variance %f ,max diff %f", rateVariance, maxDiff)
}
