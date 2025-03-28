package fast

// CBuffer holds a collection of pages, pushing modifications/removals to the
// right of each page.
type CBuffer[V any] struct {
	pages        []CBufPage[V]
	pageSizeBits int
	bitmask      int
}

// MakeCBuffer makes a CBuffer with pages of size 2^pageSizeBits
func MakeCBuffer[V any](pageSizeBits int) CBuffer[V] {
	return CBuffer[V]{
		pages:        make([]CBufPage[V], 1),
		pageSizeBits: pageSizeBits,
		bitmask:      (1 << pageSizeBits) - 1,
	}
}

// Add a value by adding a reference and a value.
// The reference will be set to the position of the value.
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

// getPage returns the page a reference belongs to.
func (c *CBuffer[V]) getPage(ref int) *CBufPage[V] {
	return &c.pages[ref>>c.pageSizeBits]
}

// Get the value that a reference points to.
func (c *CBuffer[V]) Get(ref int) V {
	return c.getPage(ref).Get(ref & c.bitmask)
}

// Len(gth) or number of values held in the buffer.
func (c *CBuffer[V]) Len() int {
	return (len(c.pages)-1)*(c.bitmask+1) + c.pages[len(c.pages)-1].Len()
}

// MarkClean resets the CleanLen of each page.
func (c *CBuffer[V]) MarkClean() {
	for i := range c.pages {
		c.pages[i].MarkClean()
	}
}

// Modify the value a reference points to in place by applying f.
func (c *CBuffer[V]) Modify(ref int, f func(*V)) {
	c.getPage(ref).Modify(ref&c.bitmask, f)
}

// Pages in the buffer.
func (c *CBuffer[V]) Pages() []CBufPage[V] {
	return c.pages
}

// Remove the value pointed to by the reference.
func (c *CBuffer[V]) Remove(ref int) V {
	bufEnd := len(c.pages) - 1
	lPage := &c.pages[bufEnd]
	lIndex := lPage.Len() - 1
	rPage := c.getPage(ref)
	rIndex := ref & c.bitmask

	// Swap item with last and then pop to remove.
	lPage.refs[lIndex], rPage.refs[rIndex] = rPage.refs[rIndex], lPage.refs[lIndex]
	lPage.items[lIndex], rPage.items[rIndex] = rPage.items[rIndex], lPage.items[lIndex]
	*rPage.refs[rIndex] = ref
	_, v := lPage.Pop()

	// Push change to the end
	rPage.pushToEnd(ref & c.bitmask)

	// Check if the lastPage is empty. If so, destroy + remove.
	if lPage.Len() == 0 && bufEnd > 0 {
		c.pages = c.pages[0:bufEnd]
	}
	return v
}

// Factory can be used to get Ref(s) to this structure.
func (c *CBuffer[V]) Factory() RefFactory[V] {
	return RefFactory[V]{
		add:    c.Add,
		get:    c.Get,
		modify: c.Modify,
		remove: c.Remove,
	}
}
