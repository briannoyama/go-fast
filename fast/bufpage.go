package fast

// CBufPage or buffer optimized page swaps/removes items in a way that moves "clean" or
// untouched values to the front.
type CBufPage[V any] struct {
	CPage[V]
	lastSwap, lastMod int
}

// CleanLen is length of data that is "clean". Clean data equals c.Data()[:c.CleanLen()].
func (c *CBufPage[V]) CleanLen() int {
	return c.lastSwap
}

// Data held in this page buffer
func (c *CBufPage[V]) Data() []V {
	return c.items
}

// MarkClean resets the CleanLen
func (c *CBufPage[V]) MarkClean() {
	c.lastMod = c.Len()
	c.lastSwap = c.Len()
}

// Modify the value a reference points to in place by applying f.
func (c *CBufPage[V]) Modify(ref int, f func(*V)) {
	c.CPage.Modify(c.pushToEnd(ref), f)
}

// Pop an item out of the page. Returns its reference and value.
// The item returned is not guaranteed to be the last item added.
// Ie. CBufPage is not a stack.
func (c *CBufPage[V]) Pop() (*int, V) {
	ref, item := c.CPage.Pop()

	// Ensure the lastSwap and Mod are not past the length of the vertices
	c.lastSwap = min(len(c.items), c.lastSwap)
	c.lastMod = min(len(c.items), c.lastMod)

	return ref, item
}

// Remove the value pointed to by the reference.
func (c *CBufPage[V]) Remove(ref int) V {
	return c.CPage.Remove(c.pushToEnd(ref))
}

// pushToEnd of data the value pointed to by ref. Returns the ref's new position
func (c *CBufPage[V]) pushToEnd(ref int) int {
	if ref < c.lastMod {
		if ref < c.lastSwap {
			c.lastSwap = ref
		}
		// Swap with last non modified reference (pushes modifications towards end of page)
		c.lastMod -= 1
		c.swap(c.lastMod, ref)
		ref = c.lastMod
	}
	return ref
}

// Factory can be used to get Ref(s) to this structure.
func (c *CBufPage[V]) Factory() RefFactory[V] {
	return RefFactory[V]{
		add:    c.Add,
		get:    c.Get,
		modify: c.Modify,
		remove: c.Remove,
	}
}
