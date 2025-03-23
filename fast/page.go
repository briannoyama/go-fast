package fast

// CPage holds a Contiguous "Page" or array of values pointed to by references.
// Cpage swaps the positions of values, updating references as it does so.
type CPage[V any] struct {
	items []V
	refs  []*int
}

// Add a value by adding a reference and a value.
// The reference will be set to the position of the value.
func (c *CPage[V]) Add(ref *int, v V) {
	c.items = append(c.items, v)
	c.refs = append(c.refs, ref)
	*ref = len(c.items) - 1
}

// Get the value that a reference points to.
func (c *CPage[V]) Get(ref int) V {
	return c.items[ref]
}

// Len(gth) or number of values held in the page.
func (c *CPage[V]) Len() int {
	return len(c.items)
}

// Modify the value a reference points to in place by applying f.
func (c *CPage[V]) Modify(ref int, f func(*V)) {
	f(&c.items[ref])
}

// Pop an item out of the page. Returns its reference and value.
// The item returned is not guaranteed to be the last item added.
// Ie. CPage is not a stack.
func (c *CPage[V]) Pop() (*int, V) {
	ref := c.refs[len(c.items)-1]
	item := c.Remove(len(c.items) - 1)

	*ref = -1
	return ref, item
}

// Remove the item pointed to by the reference.
func (c *CPage[V]) Remove(ref int) V {
	*c.refs[ref] = -1
	c.refs[ref] = c.refs[len(c.items)-1]
	*c.refs[ref] = ref
	c.refs = c.refs[:len(c.items)-1]

	item := c.items[ref]
	c.items[ref] = c.items[len(c.items)-1]
	c.items = c.items[:len(c.items)-1]
	return item
}

// Swap the internal positions of two values pointed to by the references.
func (c *CPage[V]) swap(ref0, ref1 int) {
	c.items[ref0], c.items[ref1] = c.items[ref1], c.items[ref0]
	c.refs[ref0], c.refs[ref1] = c.refs[ref1], c.refs[ref0]
	*c.refs[ref0] = ref0
	*c.refs[ref1] = ref1
}

// Factory can be used to get CPage ref's, CRef(s).
func (c *CPage[V]) Factory() RefFactory[V] {
	return RefFactory[V]{
		add:    c.Add,
		get:    c.Get,
		modify: c.Modify,
		remove: c.Remove,
	}
}
