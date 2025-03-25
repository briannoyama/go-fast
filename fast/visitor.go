package fast

// CVisitor
type CVisitor[V any] struct {
	*CPage[V]
	Visitor func(*V)
	pos     int
}

func (c *CVisitor[V]) VisitAll() {
	for c.pos = 0; c.pos < len(c.items); c.pos++ {
		c.Modify(c.pos, c.Visitor)
	}
	c.pos = -1
}

func (c *CVisitor[V]) Remove(ref int) V {
	v := c.CPage.Remove(ref)
	if ref < c.pos {
		// Visit the swapped value (if we didn't do this, it would get skipped).
		c.Modify(ref, c.Visitor)
	} else if ref > c.pos {
		// Visit the removed value. Guarantees that everything gets visited.
		c.Visitor(&v)
	} else {
		// Otherwise we're removing the currently visited value.
		c.pos -= 1
	}
	return v
}

func (c *CVisitor[V]) Factory() RefFactory[V] {
	return RefFactory[V]{
		add:    c.Add,
		get:    c.Get,
		modify: c.Modify,
		remove: c.Remove,
	}
}
