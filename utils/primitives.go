package utils

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 // thanks golang
}

// IndexOf returns the first index of needle in haystack
// or -1 if needle is not in haystack.
func IndexOf[T comparable](arr []T, val T) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

// Abs returns absolute value of number
func Abs[T Number](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

// Min returns the min value of two numbers
func Min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the max value of two numbers
func Max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}
