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
	nodes []fNode[K]
	// Points to child nodes
	page CPage[fItem[K, V]]
	root int
}

// NewFTreeMap seeds FTreeMap
func NewFTreeMap[K, V any]() FTreeMap[K, V] {
	f := FTreeMap[K, V]{root: -1}
	return f
}

// AddAdj adds a node to the right of a leaf node pointed to by its parent + index (0 or 1).
func (f *FTreeMap[K, V]) AddAdj(parent, index int, newRef *int, k K, v V) {
	newParent := -1
	if f.page.Len() > 0 {
		newParent = len(f.nodes)
		sibling := -1
		// Set root if not set
		f.root = max(f.root, 0)
		if f.page.Len() > 1 {
			// Parent exists, get sibling
			sibling = f.nodes[parent].relatives[index]
			// Update grandparent to point to new parent
			f.nodes[parent].relatives[index] = newParent
		}
		// Update sibling to point to new parent
		f.page.items[itemRef(sibling)].parent = newParent
		// Add new parent
		f.nodes = append(
			f.nodes,
			fNode[K]{relatives: [3]int{sibling, itemRef(f.page.Len()), parent}},
		)
	}
	// Add new item
	f.page.Add(newRef, fItem[K, V]{k: k, v: v, parent: newParent})
}

// Key of the node pointed to by the reference.
func (f *FTreeMap[K, V]) Key(ref int) *K {
	if ref < 0 {
		return &f.page.items[itemRef(ref)].k
	}
	return &f.nodes[ref].k
}

// Len (gth) of the map. (Ie. # of key val pairs)
func (f *FTreeMap[K, V]) Len() int {
	return f.page.Len()
}

// Parent returns the parent reference to a non-leaf node or -1 if ref is the root.
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
	root := f.root
	for _, i := range indexes {
		root = f.Rel(root)[i]
	}
	return root
}

// Rel returns the relatives for a non-leaf node.
func (f *FTreeMap[K, V]) Rel(ref int) [3]int {
	return f.nodes[ref].relatives
}

// RemoveI takes in a negative, leaf reference, removing it.
// Returns the key value associated with that reference and the sibling if it exists.
func (f *FTreeMap[K, V]) RemoveI(iRef int) (K, V, int) {
	return f.Remove(itemRef(iRef))
}

// Remove takes in a positive, leaf reference, removing it.
// Returns the key value associated with that reference and the sibling if it exists.
func (f *FTreeMap[K, V]) Remove(ref int) (K, V, int) {
	lastItem := f.page.Len() - 1
	swapParent := f.page.items[lastItem].parent
	removed := f.page.Remove(ref)
	sibling := 0

	// If there are parents
	lastNode := len(f.nodes) - 1
	if lastNode >= 0 {
		// Fix parent of swapped item
		iRef := itemRef(ref)
		f.nodes[swapParent].replace(itemRef(lastItem), iRef)

		// Remove parent node of removed item
		rel := f.nodes[removed.parent].relatives
		sibling = rel[0] ^ rel[1] ^ iRef
		*f.parent(sibling) = rel[2]
		if rel[2] == -1 {
			f.root = sibling
		} else {
			f.nodes[rel[2]].replace(removed.parent, sibling)
		}

		// Swap parent node to delete
		if removed.parent != lastNode {
			f.nodes[removed.parent], f.nodes[lastNode] = f.nodes[lastNode], f.nodes[removed.parent]
			f.fixNode(removed.parent, lastNode)
		}
		f.nodes = f.nodes[:lastNode]
	}
	return removed.k, removed.v, sibling
}

func (f *FTreeMap[K, V]) fixNode(ref, old int) {
	*f.parent(f.nodes[ref].relatives[0]) = ref
	*f.parent(f.nodes[ref].relatives[1]) = ref
	parent := f.nodes[ref].relatives[2]
	if parent == -1 {
		f.root = ref
	} else {
		f.nodes[parent].replace(old, ref)
	}
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

	f.nodes[*parent0].replace(node0, node1)
	f.nodes[*parent1].replace(node1, node0)

	*parent0, *parent1 = *parent1, *parent0
}

// Val (ue) pointed to by the negative, leaf reference.
func (f *FTreeMap[K, V]) Val(ref int) *V {
	return &f.page.items[itemRef(ref)].v
}

// VisitAll keys and values stored inside the FTreeMap.
// Does not iterate over children/values of keys for which k returns false.
func (f *FTreeMap[K, V]) VisitAll(k func(*K) bool, kv func(*K, *V)) {
	if f.Len() == 0 {
		return
	}

	for next, prev, relI := f.root, 0, 0; next != -1 || relI != 2; {
		curr := next
		if curr >= 0 {
			rel := f.Rel(curr)
			// If we didn't backtrack
			if relI < 2 {
				if !k(&f.nodes[curr].k) {
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
			if k(&f.page.items[iRef].k) {
				kv(&f.page.items[iRef].k, &f.page.items[iRef].v)
			}
			relI = 2
			next = f.page.items[iRef].parent
		}
		prev = curr
	}
}
