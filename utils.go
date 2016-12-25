package main

// Abs function
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Sign check
func Sign(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}
