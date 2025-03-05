package fast

type CBufPage[V any] struct {
	CPage[V]
	lastSwap, lastMod int
}

func (c *CBufPage[V]) pop() (*int, V) {
	ref := c.refs[len(c.items)-1]
	item := c.CPage.Remove(len(c.items) - 1)

	// Ensure the lastSwap and Mod are not past the length of the vertices
	c.lastSwap = max(len(c.items), c.lastSwap)
	c.lastMod = max(len(c.items), c.lastMod)

	*ref = -1
	return ref, item
}

func (c *CBufPage[V]) pushToEnd(ref int) int {
	if ref < c.lastMod {
		if ref < c.lastSwap {
			c.lastSwap = ref
		}
		// Swap with last non modified reference (pushes modifications towards end of page)
		c.lastMod -= 1
		c.items[c.lastMod], c.items[ref] = c.items[ref], c.items[c.lastMod]
		c.refs[c.lastMod], c.refs[ref] = c.refs[ref], c.refs[c.lastMod]
		*c.refs[c.lastMod] = c.lastMod
		*c.refs[ref] = ref
		ref = c.lastMod
	}
	return ref
}

func (c *CBufPage[V]) Modify(ref int) *V {
	return c.CPage.Modify(c.pushToEnd(ref))
}

func (c *CBufPage[V]) Remove(ref int) V {
	return c.CPage.Remove(c.pushToEnd(ref))
}

func (c *CBufPage[V]) LastMod() int {
	return c.lastMod
}

func (c *CBufPage[V]) LastSwap() int {
	return c.lastSwap
}

func (c *CBufPage[V]) Data() []V {
	return c.items
}
