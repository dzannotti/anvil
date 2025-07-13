package mathi

func Abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func Clamp(val int, low int, high int) int {
	if val > high {
		return high
	}
	if val < low {
		return low
	}
	return val
}

func Sign(val int) int {
	if val < 0 {
		return -1
	}
	return 1
}

func Sum(values ...int) int {
	total := 0
	for _, v := range values {
		total += v
	}
	return total
}
