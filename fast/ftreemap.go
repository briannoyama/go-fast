package fast

import (
	stable "go-fast/math"
)

type fNode[K any] struct {
	// 0 = left child, 1 = right child, 2 = parent
	relatives [3]int
	k         K
}

func (f *fNode[K]) replace(old, new int) {
	f.relatives[stable.IntZeroIfEqual(f.relatives[0], old)] = new
}

type fItem[K, V any] struct {
	parent int
	k      K
	v      V
}

func itemRef(ref int) int {
	return ref ^ -1
}

// FTreeMap or Full Tree Map creates a tree where internal nodes, referenced by positive ints,
// always have a key and two children.
// Leaf nodes, referenced by negative ints, contain a key and a value.
type FTreeMap[K, V any] struct {
	// Top nodes for "branches"
	nodes []fNode[K]
	// Points to child nodes
	page CPage[fItem[K, V]]
	seed [2]int
}

// NewFTreeMap seeds a non-empty FTreeMap
func NewFTreeMap[K, V any]() FTreeMap[K, V] {
	// Seed FBTreeMap with a few nothing nodes. Helps maintain non-empty variant for AddAdjacent+Remove
	f := FTreeMap[K, V]{
		nodes: []fNode[K]{{relatives: [3]int{itemRef(0), itemRef(1), -1}}},
		page:  CPage[fItem[K, V]]{},
	}
	f.page.Add(&f.seed[0], fItem[K, V]{parent: 0})
	f.page.Add(&f.seed[1], fItem[K, V]{parent: 0})
	return f
}

// AddAdj adds a node to the right of a leaf node pointed to by its parent + index (0 or 1).
func (f *FTreeMap[K, V]) AddAdj(parent, index int, newRef *int, k K, v V) int {
	newParent := len(f.nodes)
	sibling := f.nodes[parent].relatives[index]
	// Update grandparent to point to new parent
	f.nodes[parent].relatives[index] = newParent
	// Update sibling to point to new parent
	f.page.items[itemRef(sibling)].parent = newParent
	// Add new item
	f.page.Add(newRef, fItem[K, V]{k: k, v: v, parent: len(f.nodes)})
	// Add new parent
	f.nodes = append(f.nodes, fNode[K]{relatives: [3]int{sibling, itemRef(*newRef), parent}})
	return newParent
}

// Key of the node pointed to by the reference.
func (f *FTreeMap[K, V]) Key(ref int) *K {
	if ref < 0 {
		return &f.page.items[itemRef(ref)].k
	}
	return &f.nodes[ref].k
}

// Parent returns the parent reference to a non-leaf node or -1 if there is no parent.
func (f *FTreeMap[K, V]) Parent(ref int) int {
	return *f.parent(ref)
}

func (f *FTreeMap[K, V]) parent(ref int) *int {
	if ref < 0 {
		return &f.page.items[itemRef(ref)].parent
	} else {
		return &f.nodes[ref].relatives[2]
	}
}

// Path returns the reference (non-leafs positive, leafs negative) to a node.
func (f *FTreeMap[K, V]) Path(indexes ...int) int {
	root := 0
	for _, i := range indexes {
		root = f.Rel(root)[i]
	}
	return root
}

// Rel returns the relatives for a non-leaf node.
func (f *FTreeMap[K, V]) Rel(ref int) [3]int {
	return f.nodes[ref].relatives
}

// Remove takes in a positive, leaf reference, removing it.
// Returns the key value associated with that reference.
func (f *FTreeMap[K, V]) Remove(ref int) (K, V) {
	return f.RemoveRef(itemRef(ref))
}

// RemoveRef takes in a negative, leaf reference, removing it.
// Returns the key value associated with that reference.
func (f *FTreeMap[K, V]) RemoveRef(ref int) (K, V) {
	// Get nodes involved in removal
	iRef := itemRef(ref)
	parentRef := f.page.Get(ref).parent
	parent := &f.nodes[parentRef]
	sibling := parent.relatives[0] ^ parent.relatives[1] ^ iRef
	grandparent := parent.relatives[2]
	lastNode := len(f.nodes) - 1

	// Swap parent w/last node + remove
	*parent = f.nodes[lastNode]
	*f.parent(parent.relatives[0]) = parentRef
	*f.parent(parent.relatives[1]) = parentRef
	if parent.relatives[2] > -1 {
		f.nodes[parent.relatives[2]].replace(lastNode, parentRef)
	}
	// Update pointer of parent of last item that will be swapped in.
	f.nodes[f.page.items[lastNode+1].parent].replace(itemRef(lastNode+1), iRef)

	// Connect grandparent with sibling
	f.nodes[grandparent].replace(parentRef, sibling)
	*f.parent(sibling) = grandparent

	// Finish removal
	f.nodes = f.nodes[:lastNode]
	removed := f.page.Remove(ref)
	return removed.k, removed.v
}

// Swap two nodes in the FTreeMap
// Can result in unreachable nodes if one of the argumetns is a descendant of the other.
func (f *FTreeMap[K, V]) Swap(node0, node1 int) {
	parent0 := *f.parent(node0)
	parent1 := *f.parent(node1)

	f.nodes[parent0].replace(node0, node1)
	f.nodes[parent1].replace(node1, node0)

	f.nodes[node0].relatives[2] = parent1
	f.nodes[node1].relatives[2] = parent0
}

// Val(ue) pointed to by the negative, leaf reference.
func (f *FTreeMap[K, V]) Val(ref int) *V {
	return &f.page.items[itemRef(ref)].v
}

// VisitAllKeys stored inside the FTreeMap.
// Note: this will also visit the seeded values.
func (f *FTreeMap[K, V]) VisitAllKeys(k func(*K)) {
	for next, prev, relI := 0, -1, 0; next != -1 || relI != 2; {
		curr := next
		if curr >= 0 {
			// If we didn't backtrack
			if relI < 2 {
				k(&f.nodes[curr].k)
			}
			rel := f.Rel(curr)
			// Go to next node
			eq0 := stable.IntZeroIfEqual(rel[0], prev)
			eq1 := stable.IntZeroIfEqual(rel[1], prev)
			relI = (eq0 ^ eq1) << eq0
			next = rel[relI]
		} else {
			iRef := itemRef(curr)
			k(&f.page.items[iRef].k)
			relI = 2
			next = f.page.items[iRef].parent
		}
		prev = curr
	}
}

// VisitAllValues added to the FTreeMap.
func (f *FTreeMap[K, V]) VisitAllValues(v func(*V)) {
	maxSeed := stable.IntMaxIndex(f.seed[1], f.seed[0])
	for _, r := range [3][2]int{
		{0, f.seed[0^maxSeed]},
		{f.seed[0^maxSeed], f.seed[0^maxSeed]},
		{f.seed[0^maxSeed], f.page.Len()},
	} {
		for i := r[0]; i < r[1]; i++ {
			v(&f.page.items[i].v)
		}
	}
}
