package fast

type CBuffer[V any] struct {
	pages        []CBufPage[V]
	pageSizeBits int
	bitmask      int
}

func MakeCBuffer[V any](pageSizeBits int) CBuffer[V] {
	return CBuffer[V]{
		pages:        make([]CBufPage[V], 1),
		pageSizeBits: pageSizeBits,
		bitmask:      (1 << pageSizeBits) - 1,
	}
}

func (c *CBuffer[V]) Add(ref *int, v V) {
	// Check if there's enough space
	if c.pages[len(c.pages)-1].Len()>>c.pageSizeBits > 0 {
		c.pages = append(c.pages, CBufPage[V]{})
	}

	// Add page number to reference
	page := len(c.pages) - 1
	c.pages[page].Add(ref, v)
	*ref |= page << c.pageSizeBits
}

func (c *CBuffer[V]) getPage(ref int) *CBufPage[V] {
	return &c.pages[ref>>c.pageSizeBits]
}

func (c *CBuffer[V]) Get(ref int) V {
	return c.getPage(ref).Get(ref & c.bitmask)
}

func (c *CBuffer[V]) MarkClean() {
	for i := range c.pages {
		c.pages[i].MarkClean()
	}
}

func (c *CBuffer[V]) Modify(ref int, f func(*V)) {
	c.getPage(ref).Modify(ref&c.bitmask, f)
}

func (c *CBuffer[V]) Pages() []CBufPage[V] {
	return c.pages
}

// Remove the vertices pointed to by the Reference
func (c *CBuffer[V]) Remove(ref int) V {
	// Get last item in buffer.
	lastIndex := len(c.pages) - 1
	lastPage := &c.pages[lastIndex]
	lastRef := lastPage.refs[lastPage.Len()-1]
	v := lastPage.Get(*lastRef)

	// Overwrite item to remove with last item.
	page := c.getPage(ref)
	page.refs[ref&c.bitmask] = lastRef
	page.items[ref&c.bitmask] = v
	*lastRef = ref
	page.pushToEnd(ref)

	// Remove redundant copy of last item.
	lastPage.Pop()

	// Check if the lastPage is empty. If so, destroy + remove.
	if lastPage.Len() == 0 && lastIndex > 0 {
		c.pages = c.pages[0:lastIndex]
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
