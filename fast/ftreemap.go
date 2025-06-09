package fast

import "github.com/briannoyama/go-fast/stable"

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
// Returns the key value associated with that reference.
func (f *FTreeMap[K, V]) RemoveI(iRef int) (K, V) {
	return f.Remove(itemRef(iRef))
}

// Remove takes in a positive, leaf reference, removing it.
// Returns the key value associated with that reference.
func (f *FTreeMap[K, V]) Remove(ref int) (K, V) {
	// If there is a parent
	lastNode := len(f.nodes)
	if lastNode > 0 {

		// Swap lastnode and ref
		iRef := itemRef(ref)
		iLastNode := itemRef(lastNode)
		f.page.swap(ref, lastNode)
		parentRef := f.page.items[lastNode].parent

		f.nodes[f.page.items[ref].parent].replace(iLastNode, iRef)
		f.nodes[parentRef].replace(iRef, iLastNode)

		ref, iRef = lastNode, iLastNode

		// Swap last parent and parent
		lastNode -= 1
		f.nodes[parentRef], f.nodes[lastNode] = f.nodes[lastNode], f.nodes[parentRef]
		f.fixNode(parentRef, lastNode)
		f.fixNode(lastNode, parentRef)

		// Connect sibling to grandparent
		lastRel := f.nodes[lastNode].relatives
		gparent := lastRel[2]
		sibling := lastRel[0] ^ lastRel[1] ^ iRef
		if gparent == -1 {
			f.root = sibling
		} else {
			f.nodes[gparent].replace(lastNode, sibling)
		}
		*f.parent(sibling) = gparent

		// Remove last node
		f.nodes = f.nodes[:lastNode]
	}
	removed := f.page.Remove(ref)
	return removed.k, removed.v
}

func (f *FTreeMap[K, V]) fixNode(ref, old int) {
	*f.parent(f.nodes[ref].relatives[0]) = ref
	*f.parent(f.nodes[ref].relatives[1]) = ref
	parent := f.nodes[ref].relatives[2]
	if parent != -1 {
		f.nodes[parent].replace(old, ref)
	}
}

// Root returns reference to root node
func (f *FTreeMap[K, V]) Root() int {
	return f.root
}

// Swap two nodes in the FTreeMap
// Can result in unreachable nodes if one of the argumetns is a descendant of the other.
func (f *FTreeMap[K, V]) Swap(node0, node1 int) {
	parent0 := f.parent(node0)
	parent1 := f.parent(node1)

	f.nodes[*parent0].replace(node0, node1)
	f.nodes[*parent1].replace(node1, node0)

	*parent0, *parent1 = *parent1, *parent0
}

// Val(ue) pointed to by the negative, leaf reference.
func (f *FTreeMap[K, V]) Val(ref int) *V {
	return &f.page.items[itemRef(ref)].v
}

// VisitAllKeys stored inside the FTreeMap.
func (f *FTreeMap[K, V]) VisitAllKeys(k func(*K)) {
	for next, prev, relI := f.root, -1, 0; next != -1 || relI != 2; {
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
	visitor := f.page.Visitor(func(i *fItem[K, V]) { v(&i.v) })
	visitor.VisitAll()
}
