package app

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func ClampInt(v, min, max int) int {
	if v > max {
		return max
	}
	if v < min {
		return min
	}
	return v
}
