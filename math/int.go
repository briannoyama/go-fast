package stable

import "strconv"

func intMaxIndex(i, j int) int {
	return ((i - j) >> (strconv.IntSize - 1)) & 1
}

// May not be worth it. Use built in max instead?
func intMax(i, j int) int {
	index := intMaxIndex(i, j)
	return (i & (index - 1)) | (j & (0 - index))
}
