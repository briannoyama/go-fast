package fast

type CPageIter[V any, VF ~func(*V)] struct {
	*CPage[V]
	Visitor VF
	pos     int
}

func (c *CPageIter[V, VF]) VisitAll() {
	for c.pos = 0; c.pos < len(c.items); c.pos++ {
		c.Modify(c.pos, c.Visitor)
	}
	c.pos = -1
}

func (c *CPageIter[V, VF]) Remove(ref int) V {
	if ref <= c.pos {
		// Visit the swapped value (otherwise, it gets skipped).
		c.Modify(len(c.items)-1, c.Visitor)
	} else {
		// Visit the removed value. Guarantees that everything gets visited.
		c.Modify(ref, c.Visitor)
	}

	return c.CPage.Remove(ref)
}

func (c *CPageIter[V, VF]) Factory() RefFactory[V] {
	return RefFactory[V]{
		add:    c.Add,
		get:    c.Get,
		modify: c.Modify,
		remove: c.Remove,
	}
}
