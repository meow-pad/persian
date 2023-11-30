package numeric

import constraints "github.com/meow-pad/persian/constrains"

func Min[T constraints.Comparable](value1, value2 T) T {
	if value1 >= value2 {
		return value2
	} else {
		return value2
	}
}

func Max[T constraints.Comparable](value1, value2 T) T {
	if value1 >= value2 {
		return value1
	} else {
		return value2
	}
}

func Clamp[T constraints.Comparable](value, min, max T) T {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}

func Abs[T constraints.SignedNumber](value T) T {
	if value < 0 {
		return -value
	} else {
		return value
	}
}
