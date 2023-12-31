package timeutil

const (
	MilliSecondPerSecond = 1000
	SecondPerMinute      = 60
)

func SecToMS(sec int64) int64 {
	return sec * MilliSecondPerSecond
}

func MinToMS(min int64) int64 {
	return min * SecondPerMinute * MilliSecondPerSecond
}

func MinToSec(min int64) int64 {
	return min * SecondPerMinute
}
