package fast

import (
	"testing"

	"github.com/briannoyama/go-fast/assert"
)

var treeMap FTreeMap[int, int]

func setupTreeMap() {
	treeMap = NewFTreeMap[int, int]()
	cPageRefs = make([]int, 4)
	treeMap.AddAdj(-1, -1, &cPageRefs[0], 0, 0)
	treeMap.AddAdj(-1, -1, &cPageRefs[1], 1, 1)
	treeMap.AddAdj(treeMap.Root(), 0, &cPageRefs[2], 2, 2)
	treeMap.AddAdj(treeMap.Path(0), 0, &cPageRefs[3], 3, 3)
}

func TestFTreeMapAddAdjacent(t *testing.T) {
	setupTreeMap()
	ref := 0
	rel := treeMap.Rel(0)
	treeMap.AddAdj(treeMap.Root(), 1, &ref, 4, 4)
	rel = treeMap.Rel(treeMap.Path(1))
	// The default child
	assert.Equals(t, *treeMap.Key(rel[0]), 1)
	assert.Equals(t, *treeMap.Val(rel[0]), 1)
	// The newly added right child
	assert.Equals(t, *treeMap.Key(rel[1]), 4)
	assert.Equals(t, *treeMap.Val(rel[1]), 4)
}

func TestFTreeMapLen(t *testing.T) {
	setupTreeMap()
	assert.Equals(t, treeMap.Len(), 4)
}

func TestFTreeMapParent(t *testing.T) {
	setupTreeMap()

	rel := treeMap.Path(0, 0, 0)
	parent := treeMap.Parent(rel)
	assert.Equals(t, parent, 2)
	parent = treeMap.Parent(parent)
	assert.Equals(t, parent, 1)
	parent = treeMap.Parent(parent)
	assert.Equals(t, parent, 0)
	parent = treeMap.Parent(parent)
	assert.Equals(t, parent, -1)
}

func TestFTreeMapSwap(t *testing.T) {
	setupTreeMap()

	assert.Equals(t, *treeMap.Key(treeMap.Path(1)), 1)
	assert.Equals(t, *treeMap.Key(treeMap.Path(0, 0, 0)), 0)
	assert.Equals(t, *treeMap.Key(treeMap.Path(0, 0, 1)), 3)

	treeMap.Swap(treeMap.Path(1), treeMap.Path(0, 0))
	assert.Equals(t, *treeMap.Key(treeMap.Path(0, 0)), 1)
	assert.Equals(t, *treeMap.Key(treeMap.Path(1, 0)), 0)
	assert.Equals(t, *treeMap.Key(treeMap.Path(1, 1)), 3)
}

func TestFTreeMapRemove(t *testing.T) {
	setupTreeMap()

	assert.Equals(t, *treeMap.Key(treeMap.Path(0, 0, 0)), 0)

	k, _, s := treeMap.RemoveI(treeMap.Path(0, 0, 1))
	assert.Equals(t, k, 3)
	assert.Equals(t, s, -1)
	assert.Equals(t, *treeMap.Key(treeMap.Path(0, 0)), 0)

	// Empty the tree
	setupTreeMap()
	k, _, s = treeMap.Remove(cPageRefs[1])
	assert.Equals(t, k, 1)
	assert.Equals(t, s, 1)
	assert.Equals(t, *treeMap.Key(treeMap.Path(1)), 2)
	assert.Equals(t, *treeMap.Key(treeMap.Path(0, 1)), 3)

	k, _, s = treeMap.Remove(cPageRefs[2])
	assert.Equals(t, k, 2)
	assert.Equals(t, s, 0)
	assert.Equals(t, *treeMap.Key(treeMap.Path(0)), 0)
	assert.Equals(t, *treeMap.Key(treeMap.Path(1)), 3)

	k, _, s = treeMap.RemoveI(treeMap.Path(0))
	assert.Equals(t, k, 0)
	assert.Equals(t, s, -1)

	k, _, s = treeMap.RemoveI(treeMap.Root())
	assert.Equals(t, k, 3)
	assert.Equals(t, s, 0)
}

func TestFTreeMapVisitAll(t *testing.T) {
	setupTreeMap()
	for i := range 3 {
		*treeMap.Key(i) = i + 4
	}

	keys := []int{}
	treeMap.VisitAll(
		func(k *int) bool { keys = append(keys, *k); return true },
		func(k, v *int) bool { keys = append(keys, *k); return true })
	assert.Equals(t, keys, []int{4, 5, 6, 0, 0, 3, 3, 2, 2, 1, 1})

	// Edge case prune
	keys = []int{}
	treeMap.VisitAll(
		func(k *int) bool { keys = append(keys, *k); return *k != 5 && *k != 2 },
		func(k, v *int) bool { keys = append(keys, *k); return true })
	assert.Equals(t, keys, []int{4, 5, 1, 1})

	// Edge case exit early
	keys = []int{}
	treeMap.VisitAll(
		func(k *int) bool { keys = append(keys, *k); return true },
		func(k, v *int) bool { keys = append(keys, *k); return *k != 3 })
	assert.Equals(t, keys, []int{4, 5, 6, 0, 0, 3, 3})

	// Edge case tree with 2 children
	treeMap.Remove(cPageRefs[0])
	treeMap.Remove(cPageRefs[1])
	keys = []int{}
	treeMap.VisitAll(
		func(k *int) bool { keys = append(keys, *k); return true },
		func(k, v *int) bool { keys = append(keys, *k); return true })
	assert.Equals(t, keys, []int{5, 3, 3, 2, 2})

	// Edge case tree with 1 child
	treeMap.Remove(cPageRefs[2])
	keys = []int{}
	treeMap.VisitAll(
		func(k *int) bool { keys = append(keys, *k); return true },
		func(k, v *int) bool { keys = append(keys, *k); return true })
	assert.Equals(t, keys, []int{3, 3})

	// Edge case tree with 0 children
	treeMap.Remove(cPageRefs[3])
	keys = []int{}
	treeMap.VisitAll(
		func(k *int) bool { keys = append(keys, *k); return true },
		func(k, v *int) bool { keys = append(keys, *k); return true })
	assert.Equals(t, keys, []int{})
}
