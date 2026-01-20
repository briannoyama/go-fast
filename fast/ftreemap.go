package fast

import (
	"github.com/briannoyama/go-fast/stable"
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
	nodes SPage[fNode[K]]
	// Points to key/value nodes
	items SPage[fItem[K, V]]
	root  int
}

// NewFTreeMap seeds FTreeMap
func NewFTreeMap[K, V any]() FTreeMap[K, V] {
	f := FTreeMap[K, V]{
		nodes: NewSPage[fNode[K]](),
		items: NewSPage[fItem[K, V]](),
	}
	return f
}

// Add0 adds a node to the tree. Use only when lenth is 0.
// Returns the static position of the key value pair.
func (f *FTreeMap[K, V]) Add0(k K, v V) int {
	f.root = itemRef(f.items.Add(fItem[K, V]{k: k, v: v, parent: -1}))
	return itemRef(f.root)
}

// Add1 adds a node to the tree. Use only when lenth is 1.
// Returns the static position of the key value pair.
func (f *FTreeMap[K, V]) Add1(k K, v V) int {
	ref := f.items.Add(fItem[K, V]{k: k, v: v})
	sibling := f.root
	f.root = f.nodes.Add(fNode[K]{relatives: [3]int{sibling, itemRef(ref), -1}})
	f.items.values[itemRef(sibling)].parent = f.root
	f.items.values[ref].parent = f.root
	return ref
}

// AddAdj adds a node to the right of a leaf node pointed to by its parent + index (0 or 1).
// Returns the static position of the key value pair.
func (f *FTreeMap[K, V]) AddAdj(parent, index int, k K, v V) int {
	ref := f.items.Add(fItem[K, V]{k: k, v: v})
	// Create a new parent underneath current parent
	sibling := f.nodes.values[parent].relatives[index]
	newParent := f.nodes.Add(fNode[K]{relatives: [3]int{sibling, itemRef(ref), parent}})
	f.nodes.values[parent].relatives[index] = newParent
	// Update sibling to point to new parent
	f.items.values[itemRef(sibling)].parent = newParent
	f.items.values[ref].parent = newParent
	return ref
}

// Key of the node pointed to by the reference.
func (f *FTreeMap[K, V]) Key(ref int) *K {
	if ref < 0 {
		return &f.items.values[itemRef(ref)].k
	}
	return &f.nodes.values[ref].k
}

// Len (gth) of the map. (Ie. # of key val pairs)
func (f *FTreeMap[K, V]) Len() int {
	return f.items.Len()
}

// Parent returns the parent reference to a non-leaf node or -1 if ref is the root.
func (f *FTreeMap[K, V]) Parent(ref int) int {
	return *f.parent(ref)
}

func (f *FTreeMap[K, V]) parent(ref int) *int {
	if ref < 0 {
		return &f.items.values[itemRef(ref)].parent
	} else {
		return &f.nodes.values[ref].relatives[2]
	}
}

// Path returns the reference (non-leafs positive, leafs negative) to a node.
func (f *FTreeMap[K, V]) Path(indexes ...int) int {
	curr := f.root
	for _, i := range indexes {
		curr = f.Rel(curr)[i]
	}
	return curr
}

// Rel returns the relatives for a non-leaf node.
func (f *FTreeMap[K, V]) Rel(ref int) [3]int {
	return f.nodes.values[ref].relatives
}

// RemoveI0 takes in a negative, leaf reference, removing it. Use only when Len is 1.
// Returns the key, value of the reference.
func (f *FTreeMap[K, V]) RemoveI0(iRef int) (K, V) {
	return f.Remove0(itemRef(iRef))
}

// RemoveI takes in a negative, leaf reference, removing it.
// Returns the key, value and sibling of the reference.
func (f *FTreeMap[K, V]) RemoveI(iRef int) (K, V, int) {
	return f.Remove(itemRef(iRef))
}

// Remove0 removes the referenced item. Use only when Len is 1.
// Returns the key, value of the reference.
func (f *FTreeMap[K, V]) Remove0(ref int) (K, V) {
	removed := f.items.Remove(ref)
	f.root = -1
	return removed.k, removed.v
}

// Remove takes in a positive, leaf reference, removing it.
// Returns the key, value and sibling of the reference.
func (f *FTreeMap[K, V]) Remove(ref int) (K, V, int) {
	removed := f.items.Remove(ref)
	rparent := f.nodes.Remove(removed.parent)
	rel := rparent.relatives
	sibling := rel[0] ^ itemRef(ref) ^ rel[1]
	if rel[2] != -1 {
		gparent := &f.nodes.values[rel[2]]
		pindex := stable.IntZeroIfEqual(gparent.relatives[0], removed.parent)
		gparent.relatives[pindex] = sibling
		*f.parent(sibling) = rel[2]
	} else {
		f.root = sibling
		*f.parent(sibling) = -1
	}

	return removed.k, removed.v, sibling
}

// Root returns reference to root node
func (f *FTreeMap[K, V]) Root() int {
	return f.root
}

// Swap two nodes in the FTreeMap
// Can result in unreachable nodes if one of the arguments is a descendant of the other.
func (f *FTreeMap[K, V]) Swap(node0, node1 int) {
	parent0 := f.parent(node0)
	parent1 := f.parent(node1)

	f.nodes.values[*parent0].replace(node0, node1)
	f.nodes.values[*parent1].replace(node1, node0)

	*parent0, *parent1 = *parent1, *parent0
}

// Val (ue) pointed to by the negative, leaf reference.
func (f *FTreeMap[K, V]) Val(ref int) *V {
	return &f.items.values[itemRef(ref)].v
}

// VisitAll keys and values stored inside the FTreeMap.
// Does not iterate over children/values of keys for which k returns false.
// Exits early if kv returns false
func (f *FTreeMap[K, V]) VisitAll(k func(*K) bool, kv func(*K, *V) bool) {
	if f.Len() == 0 {
		return
	}

	for next, prev, relI := f.root, 0, 0; next != -1 || relI != 2; {
		curr := next
		if curr >= 0 {
			rel := f.Rel(curr)
			// If we didn't backtrack
			if relI < 2 {
				if !k(&f.nodes.values[curr].k) {
					// Skip and go to parent
					prev = rel[1]
				}
			}
			// Go to next node
			eq0 := stable.IntZeroIfEqual(rel[0], prev)
			eq1 := stable.IntZeroIfEqual(rel[1], prev)
			relI = (eq0 ^ eq1) << eq0
			next = rel[relI]
		} else {
			iRef := itemRef(curr)
			if k(&f.items.values[iRef].k) {
				if !kv(&f.items.values[iRef].k, &f.items.values[iRef].v) {
					// For iter yield functions, return early when yield is false
					return
				}
			}
			relI = 2
			next = f.items.values[iRef].parent
		}
		prev = curr
	}
}
