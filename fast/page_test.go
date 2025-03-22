package fast

import (
	"fmt"
	"testing"

	"go-fast/assert"
)

var (
	c   CPage[int]
	ref []int
)

func setup() {
	c = CPage[int]{}
	ref = make([]int, 3)
	c.Add(&ref[1], 0)
	c.Add(&ref[2], 1)
	c.Add(&ref[0], 2)
}

func TestCPageAdd(t *testing.T) {
	setup()

	expItems := []int{0, 1, 2}
	assert.Equals(t, c.items, expItems)

	expRefs := []int{2, 0, 1}
	assert.Equals(t, ref, expRefs)
}

func TestCPageGet(t *testing.T) {
	setup()

	assert.Equals(t, c.Get(ref[1]), 0)
	assert.Equals(t, c.Get(ref[2]), 1)
	assert.Equals(t, c.Get(ref[0]), 2)
}

func TestCPageLen(t *testing.T) {
	setup()

	assert.Equals(t, c.Len(), 3)
}

func TestCPageModify(t *testing.T) {
	setup()

	c.Modify(0, func(i *int) { *i = 5 })
	expItems := []int{5, 1, 2}
	assert.Equals(t, c.items, expItems)
}

func TestCPagePop(t *testing.T) {
	setup()
	r, val := c.Pop()

	assert.Equals(t, r, &ref[0])
	assert.Equals(t, val, 2)
}

func TestCPageRemove(t *testing.T) {
	setup()

	assert.Equals(t, c.Remove(ref[1]), 0)
	fmt.Printf("%v", ref)
	assert.Equals(t, c.Get(ref[2]), 1)
	assert.Equals(t, c.Get(ref[0]), 2)
	assert.Equals(t, c.Len(), 2)
}

func TestCPageSwap(t *testing.T) {
	setup()

	c.swap(ref[0], ref[1])
	assert.Equals(t, c.Get(ref[0]), 2)
	assert.Equals(t, c.Get(ref[1]), 0)
	assert.Equals(t, ref[0], 0)
	assert.Equals(t, ref[1], 2)
}
