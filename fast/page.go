package fast

type CPage[V any] struct {
	items []V
	refs  []*int
}

func (c *CPage[V]) Add(ref *int, v V) {
	c.items = append(c.items, v)
	c.refs = append(c.refs, ref)
	*ref = len(c.items) - 1
}

func (c *CPage[V]) Get(ref int) V {
	return c.items[ref]
}

func (c *CPage[V]) Len() int {
	return len(c.items)
}

func (c *CPage[V]) Modify(ref int, f func(*V)) {
	f(&c.items[ref])
}

func (c *CPage[V]) Pop() (*int, V) {
	ref := c.refs[len(c.items)-1]
	item := c.Remove(len(c.items) - 1)

	*ref = -1
	return ref, item
}

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

func (c *CPage[V]) swap(ref0, ref1 int) {
	c.items[ref0], c.items[ref1] = c.items[ref1], c.items[ref0]
	c.refs[ref0], c.refs[ref1] = c.refs[ref1], c.refs[ref0]
	*c.refs[ref0] = ref0
	*c.refs[ref1] = ref1
}

func (c *CPage[V]) Factory() RefFactory[V] {
	return RefFactory[V]{
		add:    c.Add,
		get:    c.Get,
		modify: c.Modify,
		remove: c.Remove,
	}
}
