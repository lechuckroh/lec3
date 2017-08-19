package lecimg

import "math"

// Max returns max value
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Min returns min value
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Minf32 returns min value
func Minf32(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

// Maxf32 returns min value
func Maxf32(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

// Sincosf32 returns sin, cos values
func Sincosf32(a float32) (float32, float32) {
	sin, cos := math.Sincos(math.Pi * float64(a) / 180)
	return float32(sin), float32(cos)
}

// Floorf32 returns floor value.
func Floorf32(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

// InRangef32 checks if value is between rangeFrom to rangeTo.
func InRangef32(value, rangeFrom, rangeTo float32) bool {
	epsilon := float32(0.00001)
	return rangeFrom <= (value-epsilon) && (value+epsilon) <= rangeTo
}

// InRange checks if value is between rangeFrom to rangeTo.
func InRange(value, rangeFrom, rangeTo int) bool {
	return rangeFrom <= value && value <= rangeTo
}
