package fast

import stable "go-fast/math"

type fNode[K any] struct {
	// 0 = left child, 1 = right child, 2 = parent
	relatives [3]int
	k         K
}

type fItem[K, V any] struct {
	parent int
	k      K
	v      V
}

func itemRef(ref int) int {
	return ref ^ -1
}

type FTreeMap[K, V any] struct {
	// Top nodes for "branches"
	nodes []fNode[K]
	// Points to child nodes
	page CPage[fItem[K, V]]
}

func NewFTreeMap[K, V any]() FTreeMap[K, V] {
	// Seed FBTreeMap with a few nothing nodes. Helps maintain non-empty variant for AddAdjacent+Remove
	f := FTreeMap[K, V]{
		nodes: []fNode[K]{{relatives: [3]int{itemRef(0), itemRef(1), -1}}},
		page:  CPage[fItem[K, V]]{},
	}
	ref := 0
	f.page.Add(&ref, fItem[K, V]{parent: 0})
	f.page.Add(&ref, fItem[K, V]{parent: 0})
	return f
}

// AddAdjacent adds a node to the right of ref
func (f *FTreeMap[K, V]) AddAdjacent(ref int, newref *int, k K, v V) {
	// Add new item
	f.page.Add(newref, fItem[K, V]{k: k, v: v})

	// Update grandparent to point to new parent
	grandparent := f.page.Get(itemRef(ref)).parent
	parent := len(f.nodes)
	f.changeRef(ref, ref, parent)

	// Add new parent
	sibling := itemRef(*newref)
	f.nodes = append(f.nodes, fNode[K]{relatives: [3]int{ref, sibling, grandparent}})

	// Update parent pointers
	f.page.items[ref].parent = parent
	f.page.items[sibling].parent = parent
}

func (f *FTreeMap[K, V]) GetKey(ref int) K {
	if ref < 0 {
		return f.page.Get(itemRef(ref)).k
	} else {
		return f.nodes[ref].k
	}
}

func (f *FTreeMap[K, V]) GetValue(ref int) V {
	return f.page.Get(itemRef(ref)).v
}

func (f *FTreeMap[K, V]) GetRelatives(node int) [3]int {
	return f.nodes[node].relatives
}

func (f *FTreeMap[K, V]) ModifyParent(ref int) *int {
	if ref < 0 {
		return &f.page.items[itemRef(ref)].parent
	} else {
		return &f.nodes[ref].relatives[2]
	}
}

func (f *FTreeMap[K, V]) Swap(node0, node1 int) {
	parent0 := *f.ModifyParent(node0)
	parent1 := *f.ModifyParent(node1)

	f.nodes[parent0].relatives[f.getRefIndex(parent0, node0)] = node1
	f.nodes[parent1].relatives[f.getRefIndex(parent1, node1)] = node0

	f.nodes[node0].relatives[2] = parent1
	f.nodes[node1].relatives[2] = parent0
}

func (f *FTreeMap[K, V]) Remove(ref int) fItem[K, V] {
	// Remove item node
	removed := f.page.Remove(itemRef(ref))
	parent := f.nodes[removed.parent]

	// Swap parent w/last node + remove
	f.nodes[removed.parent] = f.nodes[len(f.nodes)-1]
	f.nodes = f.nodes[:len(f.nodes)-1]

	// Get sibling index
	sibling := parent.relatives[0] ^ parent.relatives[1] ^ ref
	grandparent := parent.relatives[2]

	parentIndex := f.getRefIndex(grandparent, removed.parent)
	f.nodes[grandparent].relatives[parentIndex] = sibling

	// TODO Update parent of sibling
	*f.ModifyParent(sibling) = grandparent
	// f.nodeRemove(removed.parent, ref)
	// Look at node that was swapped in and update it's parent
	f.changeRef(ref, f.page.Len(), ref)
	return removed
}

func (f *FTreeMap[K, V]) changeRef(ref, oldValue, newValue int) {
	parent := f.page.Get(itemRef(ref)).parent
	// Check if left child was the one moved/swapped in.
	refIndex := f.getRefIndex(parent, oldValue)
	// Update to new reference
	f.nodes[parent].relatives[refIndex] = newValue
}

func (f *FTreeMap[K, V]) getRefIndex(parent, ref int) int {
	return stable.IntZeroIfEqual(f.nodes[parent].relatives[0], ref)
}

func (f *FTreeMap[K, V]) nodeRemove(parent, ref int) {
	removed := f.nodes[parent]
	// Swap in last node
	f.nodes[parent] = f.nodes[len(f.nodes)-1]
	f.nodes = f.nodes[:len(f.nodes)-1]
	parentSibling := removed.relatives[0] ^ removed.relatives[1] ^ ref
	grandparent := removed.relatives[2]
	// Check if left child was the parent
	parentIndex := f.getRefIndex(grandparent, parent)
	f.nodes[grandparent].relatives[parentIndex] = parentSibling
}
