package numeric

import constraints "github.com/meow-pad/persian/constrains"

func Increase[T constraints.SignedNumber](value T, increment T, max T) T {
	if increment == 0 {
		return value
	}
	if increment < 0 {
		panic("invalid increment")
	}
	if value > 0 {
		remaining := max - value
		if increment > remaining {
			return max
		} else {
			return value + increment
		}
	} else {
		result := value + increment
		if result > max {
			return max
		} else {
			return result
		}
	}
}

func Decrease[T constraints.SignedNumber](value T, decrement T, min T) T {
	if decrement == 0 {
		return value
	}
	if decrement < 0 {
		panic("invalid decrement")
	}
	if value < 0 {
		remaining := value - min
		if decrement > remaining {
			return min
		} else {
			return value - decrement
		}
	} else {
		result := value - decrement
		if result < min {
			return min
		} else {
			return result
		}
	}
}
