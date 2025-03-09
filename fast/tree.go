package fast

import stable "go-fast/math"

type FNode[K any] struct {
	// 0 = left child, 1 = right child, 2 = parent
	relatives [3]int
	k         K
}

type FItem[K, V any] struct {
	parent int
	k      K
	v      V
}

func itemRef(ref int) int {
	return ref ^ -1
}

type FBTreeMap[K, V any] struct {
	// Top nodes for "branches"
	nodes []FNode[K]
	// Points to child nodes
	page CPage[FItem[K, V]]
}

func NewFBTreeMap[K, V any]() FBTreeMap[K, V] {
	// Seed FBTreeMap with a few nothing nodes. Helps maintain non-empty variant for AddAdjacent+Remove
	f := FBTreeMap[K, V]{
		nodes: []FNode[K]{{relatives: [3]int{itemRef(0), itemRef(1), -1}}},
		page:  CPage[FItem[K, V]]{},
	}
	ref := 0
	f.page.Add(&ref, FItem[K, V]{parent: 0})
	f.page.Add(&ref, FItem[K, V]{parent: 0})
	return f
}

// AddAdjacent adds a node to the right
func (f *FBTreeMap[K, V]) AddAdjacent(to int, ref *int, item FItem[K, V]) {
	f.page.Add(ref, item)

	grandparent := f.page.Get(itemRef(to)).parent
	f.changeRef(to, to, len(f.nodes))

	f.nodes = append(f.nodes, FNode[K]{relatives: [3]int{to, itemRef(*ref), grandparent}})
}

// GetKey, GetChild
func (f *FBTreeMap[K, V]) GetKey(ref int) FItem[K, V] {
	return f.page.Get(itemRef(ref))
}

func (f *FBTreeMap[K, V]) GetItem(ref int) FItem[K, V] {
	return f.page.Get(itemRef(ref))
}

func (f *FBTreeMap[K, V]) GetNode(node int) FNode[K] {
	return f.nodes[node]
}

func (f *FBTreeMap[K, V]) GetParent(node int) int {
	if node < 0 {
		return f.page.Get(itemRef(node)).parent
	} else {
		return f.nodes[node].relatives[2]
	}
}

func (f *FBTreeMap[K, V]) Swap(node0, node1 int) {
	parent0 := f.GetParent(node0)
	parent1 := f.GetParent(node1)

	f.nodes[parent0].relatives[f.getRefIndex(parent0, node0)] = node1
	f.nodes[parent1].relatives[f.getRefIndex(parent1, node1)] = node0

	f.nodes[node0].relatives[2] = parent1
	f.nodes[node1].relatives[2] = parent0
}

func (f *FBTreeMap[K, V]) Remove(ref int) FItem[K, V] {
	// TODO Add itemRef
	// Remove node and update it's parent
	removed := f.page.Remove(itemRef(ref))
	f.nodeRemove(removed.parent, ref)
	// Look at node that was swapped in and update it's parent
	f.changeRef(ref, f.page.Len(), ref)
	return removed
}

func (f *FBTreeMap[K, V]) changeRef(ref, oldValue, newValue int) {
	parent := f.page.Get(itemRef(ref)).parent
	// Check if left child was the one moved/swapped in.
	refIndex := f.getRefIndex(parent, oldValue)
	// Update to new reference
	f.nodes[parent].relatives[refIndex] = newValue
}

func (f *FBTreeMap[K, V]) getRefIndex(parent, ref int) int {
	return stable.IntZeroIfEqual(f.nodes[parent].relatives[0], ref)
}

func (f *FBTreeMap[K, V]) nodeRemove(parent, ref int) {
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
