package fast

import (
	"testing"

	"go-fast/assert"
)

var bufPage CBufPage[int]

func setupBufPage() {
	bufPage = CBufPage[int]{}
	cPageRefs = make([]int, 5)
	for i := 0; i < len(cPageRefs); i++ {
		bufPage.Add(&cPageRefs[i], i)
	}
	bufPage.MarkClean()
}

func TestBufPagePop(t *testing.T) {
	setupBufPage()

	ref, v := bufPage.Pop()
	assert.Equals(t, ref, &cPageRefs[4])
	assert.Equals(t, v, 4)

	ref, v = bufPage.Pop()
	assert.Equals(t, ref, &cPageRefs[3])
	assert.Equals(t, v, 3)

	assert.Equals(t, bufPage.CleanLen(), bufPage.Len())
}

func TestBufPageModify(t *testing.T) {
	setupBufPage()

	modified := []int{}
	bufPage.Modify(cPageRefs[0], func(i *int) { modified = append(modified, *i) })
	bufPage.Modify(cPageRefs[1], func(i *int) { modified = append(modified, *i) })
	bufPage.Modify(cPageRefs[2], func(i *int) { modified = append(modified, *i) })
	assert.Equals(t, modified, []int{0, 1, 2})

	// First modified value should be at the end.
	_, v := bufPage.Pop()
	assert.Equals(t, v, 0)
	assert.Equals(t, bufPage.CleanLen(), 0)

	cleanData := bufPage.Data()[1]
	bufPage.MarkClean()

	modified = []int{}
	bufPage.Modify(cPageRefs[1], func(i *int) { modified = append(modified, *i) })
	bufPage.Modify(cPageRefs[2], func(i *int) { modified = append(modified, *i) })
	assert.Equals(t, modified, []int{1, 2})

	// First modified value should be at the end.
	_, v = bufPage.Pop()
	assert.Equals(t, v, 1)
	// First two elements should not been touched.
	assert.Equals(t, bufPage.CleanLen(), 2)
	assert.Equals(t, bufPage.Data()[1], cleanData)
}

func TestBufPageRemove(t *testing.T) {
	setupBufPage()

	assert.Equals(t, bufPage.Remove(cPageRefs[2]), 2)
	assert.Equals(t, bufPage.CleanLen(), 2)
	// Modify value after removal. Should not affect clean portion.
	bufPage.Modify(cPageRefs[3], func(*int) {})
	assert.Equals(t, bufPage.CleanLen(), 2)

	bufPage.MarkClean()
	bufPage.Modify(cPageRefs[4], func(*int) {})
	assert.Equals(t, bufPage.Remove(cPageRefs[1]), 1)
	// Expect 0 at beginning (never moved/modified), and 4 at the end (most recently modified)
	assert.Equals(t, bufPage.Data(), []int{0, 3, 4})
	assert.Equals(t, bufPage.CleanLen(), 1)
}
