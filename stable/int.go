package stable

import "strconv"

// IntMaxIndex returns index of max number (either 0 or 1)
func IntMaxIndex(i, j int) int {
	return ((i - j) >> (strconv.IntSize - 1)) & 1
}

// IntZeroIfEqual returns 0 if equal else returns 1
func IntZeroIfEqual(i, j int) int {
	ij := i ^ j
	return IntZeroIfZero(ij)
}

// IntZeroIfZero returns 0 if 0 else 1
func IntZeroIfZero(i int) int {
	return ((i | -i) >> (strconv.IntSize - 1)) & 1
}
