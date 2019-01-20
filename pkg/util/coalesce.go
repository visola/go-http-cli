package util

// FirstOrZero returns the first non-zero or zero if none is found
func FirstOrZero(values ...int) int {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}
