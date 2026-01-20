package fast

import (
	"testing"

	"github.com/briannoyama/go-fast/assert"
)

var (
	cPage     CPage[int]
	cPageRefs []int
)

func setupPage() {
	cPage = CPage[int]{}
	cPageRefs = make([]int, 3)
	cPage.Add(&cPageRefs[1], 0)
	cPage.Add(&cPageRefs[2], 1)
	cPage.Add(&cPageRefs[0], 2)
}

func TestCPageAdd(t *testing.T) {
	setupPage()

	expItems := []int{0, 1, 2}
	assert.Equals(t, cPage.items, expItems)

	expRefs := []int{2, 0, 1}
	assert.Equals(t, cPageRefs, expRefs)
}

func TestCPageGet(t *testing.T) {
	setupPage()

	assert.Equals(t, cPage.Get(cPageRefs[1]), 0)
	assert.Equals(t, cPage.Get(cPageRefs[2]), 1)
	assert.Equals(t, cPage.Get(cPageRefs[0]), 2)
}

func TestCPageLen(t *testing.T) {
	setupPage()

	assert.Equals(t, cPage.Len(), 3)
}

func TestCPageModify(t *testing.T) {
	setupPage()

	cPage.Modify(0, func(i *int) { *i = 5 })
	expItems := []int{5, 1, 2}
	assert.Equals(t, cPage.items, expItems)
}

func TestCPagePop(t *testing.T) {
	setupPage()
	r, val := cPage.Pop()

	assert.Equals(t, r, &cPageRefs[0])
	assert.Equals(t, val, 2)
}

func TestCPageRemove(t *testing.T) {
	setupPage()

	assert.Equals(t, cPage.Remove(cPageRefs[1]), 0)
	assert.Equals(t, cPage.Get(cPageRefs[2]), 1)
	assert.Equals(t, cPage.Get(cPageRefs[0]), 2)
	assert.Equals(t, cPage.Len(), 2)
}

func TestCPageSwap(t *testing.T) {
	setupPage()

	cPage.swap(cPageRefs[0], cPageRefs[1])
	assert.Equals(t, cPage.Get(cPageRefs[0]), 2)
	assert.Equals(t, cPage.Get(cPageRefs[1]), 0)
	assert.Equals(t, cPageRefs[0], 0)
	assert.Equals(t, cPageRefs[1], 2)
}

func TestCPageRefPreservation(t *testing.T) {
	setupPage()
	// TODO Finish
	cPage.swap(cPageRefs[0], cPageRefs[1])
}
