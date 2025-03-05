package fast

type CBuffer[V any] struct {
	pages        []*CBufPage[V]
	pageSizeBits int
	bitmask      int
}

func MakeCBuffer[V any](pageSizeBits int) CBuffer[V] {
	return CBuffer[V]{
		pageSizeBits: pageSizeBits,
		bitmask:      (1 << pageSizeBits) - 1,
	}
}

func (c *CBuffer[V]) Add(ref *int, v V) {
	// Check if there's enough space
	if c.pages[len(c.pages)-1].Len()>>c.pageSizeBits > 0 {
		c.pages = append(c.pages, &CBufPage[V]{})
	}

	c.pages[len(c.pages)-1].Add(ref, v)
}

func (c *CBuffer[V]) getPage(ref int) *CBufPage[V] {
	return c.pages[ref>>c.pageSizeBits]
}

func (c *CBuffer[V]) Get(ref int) V {
	return c.getPage(ref).Get(ref & c.bitmask)
}

func (c *CBuffer[V]) Modify(ref int) *V {
	return c.getPage(ref).Modify(ref & c.bitmask)
}

func (c *CBuffer[V]) Pages() []*CBufPage[V] {
	return c.pages
}

// Remove the vertices pointed to by the Reference
func (c *CBuffer[V]) Remove(ref int) V {
	lastPage := len(c.pages) - 1
	lastRef, v := c.pages[lastPage].pop()
	// If we did not just pop the ref
	if ref != *lastRef {
		page := c.getPage(ref)
		page.refs[ref&c.bitmask] = lastRef
		page.items[ref&c.bitmask] = v
		*lastRef = ref
		// Move the vertices and reference to the end of the page where the removal happened.
		page.pushToEnd(ref)
	}
	// Check if the lastPage is empty. If so, destroy + remove.
	if c.pages[lastPage].Len() == 0 {
		if lastPage > 0 {
			c.pages = c.pages[0:lastPage]
		}
	}
	return v
}

func (c *CBuffer[V]) Factory() RefFactory[V] {
	return RefFactory[V]{
		add:    c.Add,
		get:    c.Get,
		modify: c.Modify,
		remove: c.Remove,
	}
}
