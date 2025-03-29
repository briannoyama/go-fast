package fast

import (
	"math"
	"testing"

	"go-fast/assert"
)

var heap Heap[int]

func setupHeap() {
	heap = Heap[int]{}
	cPageRefs = make([]int, 9)
	for i := 0; i < len(cPageRefs); i++ {
		heap.Add(&cPageRefs[i], i)
	}
}

func getHeapData() []int {
	values := []int{}
	heap.VisitAll(func(v *int) { values = append(values, *v) })
	return values
}

func emptyHeap() []int {
	values := []int{}
	for heap.Len() > 0 {
		_, v := heap.Pop()
		values = append(values, v)
	}
	return values
}

func TestHeapAdd(t *testing.T) {
	setupHeap()

	assert.Equals(t, getHeapData(), []int{0, 1, 2, 3, 4, 5, 6, 7, 8})
}

func TestHeapAddRelative(t *testing.T) {
	setupHeap()

	refs := make([]int, 4)
	heap.AddRelative(&refs[0], 9, -5)
	heap.AddRelative(&refs[1], 10, -9)
	heap.AddRelative(&refs[2], 11, 2)
	heap.AddRelative(&refs[3], 12, -1)
	assert.Contains(t, getHeapData(), 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	assert.Equals(t, emptyHeap(), []int{10, 0, 1, 2, 3, 9, 4, 5, 6, 7, 8, 12, 11})
}

func TestHeapGet(t *testing.T) {
	setupHeap()

	assert.Equals(t, heap.Get(cPageRefs[0]), 0)
	assert.Equals(t, heap.Get(cPageRefs[3]), 3)
	assert.Equals(t, heap.Get(cPageRefs[6]), 6)
	// Shift elements around, check that references still work
	ref := 0
	heap.AddRelative(&ref, 10, -9)
	assert.Equals(t, heap.Get(cPageRefs[0]), 0)
	assert.Equals(t, heap.Get(cPageRefs[3]), 3)
	assert.Equals(t, heap.Get(cPageRefs[6]), 6)
}

func TestHeapHighPriority(t *testing.T) {
	// Test resilience against int overflow
	heap = Heap[int]{}
	heap.lastPriority = math.MaxInt - 3
	cPageRefs = make([]int, 9)
	for i := 0; i < len(cPageRefs); i++ {
		heap.Add(&cPageRefs[i], i)
	}

	refs := make([]int, 4)
	heap.AddRelative(&refs[0], 9, -5)
	heap.AddRelative(&refs[1], 10, -9)
	heap.Add(&refs[2], 12)
	heap.Add(&refs[3], 11)
	assert.Contains(t, getHeapData(), 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	assert.Equals(t, emptyHeap(), []int{10, 0, 1, 2, 3, 9, 4, 5, 6, 7, 8, 12, 11})
}

func TestHeapModify(t *testing.T) {
	setupHeap()

	heap.Modify(cPageRefs[4], func(i *int) { *i *= 5 })
	assert.Equals(t, emptyHeap(), []int{0, 1, 2, 3, 20, 5, 6, 7, 8})
}

func TestHeapModifyPriority(t *testing.T) {
	setupHeap()
	heap.ModifyPriority(cPageRefs[4], 10)
	assert.Equals(t, emptyHeap(), []int{0, 1, 2, 3, 5, 6, 7, 8, 4})

	setupHeap()
	heap.ModifyPriority(cPageRefs[4], 10)
	heap.ModifyPriority(cPageRefs[6], -3)
	assert.Equals(t, emptyHeap(), []int{0, 1, 2, 3, 6, 5, 7, 8, 4})
}

func TestHeapRemove(t *testing.T) {
	setupHeap()

	heap.Remove(cPageRefs[4])
	assert.Equals(t, emptyHeap(), []int{0, 1, 2, 3, 5, 6, 7, 8})
}

func TestHeapToEnd(t *testing.T) {
	setupHeap()

	heap.ToEnd(cPageRefs[4])
	heap.ToEnd(cPageRefs[0])
	assert.Equals(t, emptyHeap(), []int{1, 2, 3, 5, 6, 7, 8, 4, 0})

	setupHeap()
	heap.ToEnd(cPageRefs[4])
	heap.ToEnd(cPageRefs[0])
	heap.ToEnd(cPageRefs[4])
	assert.Equals(t, emptyHeap(), []int{1, 2, 3, 5, 6, 7, 8, 0, 4})
}
