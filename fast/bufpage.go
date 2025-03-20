package fast

type CBufPage[V any] struct {
	CPage[V]
	lastSwap, lastMod int
}

func (c *CBufPage[V]) Pop() (*int, V) {
	ref, item := c.CPage.Pop()

	// Ensure the lastSwap and Mod are not past the length of the vertices
	c.lastSwap = max(len(c.items), c.lastSwap)
	c.lastMod = max(len(c.items), c.lastMod)

	return ref, item
}

func (c *CBufPage[V]) pushToEnd(ref int) int {
	if ref < c.lastMod {
		if ref < c.lastSwap {
			c.lastSwap = ref
		}
		// Swap with last non modified reference (pushes modifications towards end of page)
		c.lastMod -= 1
		c.swap(c.lastMod, ref)
		*c.refs[ref] = ref
		ref = c.lastMod
	}
	return ref
}

func (c *CBufPage[V]) Modify(ref int, f func(*V)) {
	c.CPage.Modify(c.pushToEnd(ref), f)
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
