package fast

// CVisitor visits all nodes in a CPage, even those removed during visit.
type CVisitor[V any] struct {
	*CPage[V]
	visitor func(*V)
	pos     int
}

// VisitAll nodes applies the visitor function.
func (c *CVisitor[V]) VisitAll() {
	for c.pos = 0; c.pos < len(c.items); c.pos++ {
		c.Modify(c.pos, c.visitor)
	}
	c.pos = -1
}

// Remove the value, calling visitor on shuffled items if necessary.
func (c *CVisitor[V]) Remove(ref int) V {
	v := c.CPage.Remove(ref)
	if ref < c.pos && c.pos < c.CPage.Len() {
		// Visit the swapped value if not already visited.
		// (if we didn't do this, it would get skipped).
		c.Modify(ref, c.visitor)
	} else if ref > c.pos {
		// Visit the removed value. Guarantees that everything gets visited.
		c.visitor(&v)
	} else {
		// Otherwise we're removing the currently visited value.
		// Current value in this position should be visited next.
		c.pos -= 1
	}
	return v
}

// Remove the current visited item if possible
func (c *CVisitor[V]) RmCurrent() {
	if c.pos != -1 {
		c.Remove(c.pos)
	}
}

// Factory can be used to get CPage ref's, CRef(s).
func (c *CVisitor[V]) Factory() RefFactory[V] {
	return RefFactory[V]{
		add:    c.Add,
		get:    c.Get,
		modify: c.Modify,
		remove: c.Remove,
	}
}
