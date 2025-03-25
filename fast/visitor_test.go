package fast

import (
	"testing"

	"go-fast/assert"
)

func TestPageIter(t *testing.T) {
	cPageRefs = make([]int, 4)
	var visitor CVisitor[func() int]

	// Some nodes removed ahead, some behind
	cPage := CPage[func() int]{}
	cPage.Add(&cPageRefs[0], func() int {
		visitor.Remove(cPageRefs[0])
		return 0
	})
	cPage.Add(&cPageRefs[1], func() int {
		visitor.Remove(cPageRefs[2])
		return 1
	})
	cPage.Add(&cPageRefs[2], func() int { return 2 })
	cPage.Add(&cPageRefs[3], func() int { return 3 })

	visited := map[int]bool{}
	visitor = cPage.Visitor(func(f *func() int) {
		visited[(*f)()] = true
	})

	expected := map[int]bool{0: true, 1: true, 2: true, 3: true}
	visitor.VisitAll()
	assert.Equals(t, visited, expected)

	// Everything removes itself
	cPage = CPage[func() int]{}
	cPage.Add(&cPageRefs[0], func() int {
		visitor.Remove(cPageRefs[0])
		return 0
	})
	cPage.Add(&cPageRefs[1], func() int {
		visitor.Remove(cPageRefs[1])
		return 1
	})
	cPage.Add(&cPageRefs[2], func() int {
		visitor.Remove(cPageRefs[2])
		return 2
	})
	cPage.Add(&cPageRefs[3], func() int {
		visitor.Remove(cPageRefs[3])
		return 3
	})

	visited = map[int]bool{}
	visitor.VisitAll()
	assert.Equals(t, visited, expected)
}
