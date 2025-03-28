package fast

import (
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
	heap.AddRelative(&refs[0], 10, -9)
	heap.AddRelative(&refs[0], 11, 2)
	heap.Add(&refs[0], 12)
	items := []int{}
	heap.VisitAll(func(v *int) { items = append(items, *v) })
	assert.Contains(t, getHeapData(), 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	assert.Equals(t, emptyHeap(), []int{10, 0, 1, 2, 3, 9, 4, 5, 6, 7, 8, 12, 11})
}

func TestHeapGet(t *testing.T) {
}

func TestHeapHighPriority(t *testing.T) {
}

func TestHeapModify(t *testing.T) {
}

func TestHeapModifyPriority(t *testing.T) {
}

func TestHeapRemove(t *testing.T) {
}

func TestHeapToEnd(t *testing.T) {
}
