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

	k, _ := treeMap.RemoveI(treeMap.Path(0, 0, 1))
	assert.Equals(t, k, 3)
	assert.Equals(t, *treeMap.Key(treeMap.Path(0, 0)), 0)

	// Empty the tree
	setupTreeMap()
	k, _ = treeMap.Remove(cPageRefs[1])
	assert.Equals(t, k, 1)
	assert.Equals(t, *treeMap.Key(treeMap.Path(1)), 2)
	assert.Equals(t, *treeMap.Key(treeMap.Path(0, 1)), 3)

	k, _ = treeMap.Remove(cPageRefs[2])
	assert.Equals(t, k, 2)
	assert.Equals(t, *treeMap.Key(treeMap.Path(0)), 0)
	assert.Equals(t, *treeMap.Key(treeMap.Path(1)), 3)

	k, _ = treeMap.RemoveI(treeMap.Path(0))
	assert.Equals(t, k, 0)

	k, _ = treeMap.RemoveI(treeMap.Root())
	assert.Equals(t, k, 3)
}

func TestFTreeMapVisitAllKeys(t *testing.T) {
	setupTreeMap()

	keys := []int{}
	treeMap.VisitAllKeys(func(k *int) { keys = append(keys, *k) })
	assert.Equals(t, keys, []int{0, 0, 0, 0, 3, 2, 1})
}

func TestFTreeMapVisitAllValues(t *testing.T) {
	setupTreeMap()

	vals := []int{}
	treeMap.VisitAllValues(func(v *int) { vals = append(vals, *v) })
	assert.Equals(t, vals, []int{0, 1, 2, 3})
}
