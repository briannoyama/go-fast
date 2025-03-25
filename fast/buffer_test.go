package fast

import (
	"testing"

	"go-fast/assert"
)

var buffer CPage[int]

func setupBuffer() {
	cPage = CPage[int]{}
	cPageRefs = make([]int, 3)
	cPage.Add(&cPageRefs[1], 0)
	cPage.Add(&cPageRefs[2], 1)
	cPage.Add(&cPageRefs[0], 2)
}

func TestBufferAdd(t *testing.T) {
	setupPage()

	expItems := []int{0, 1, 2}
	assert.Equals(t, cPage.items, expItems)

	expRefs := []int{2, 0, 1}
	assert.Equals(t, cPageRefs, expRefs)
}
