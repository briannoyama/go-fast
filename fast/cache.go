package fast

type keyVal[K, V any] struct {
	k K
	v V
}

type Cache[K comparable, V any] struct {
	heap Heap[keyVal[K, V]]
	m    map[K]*int
}

func MakeCache[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{m: make(map[K]*int)}
}

func (c *Cache[K, V]) Add(k K, v V) {
	var ref int
	c.heap.Add(&ref, keyVal[K, V]{k: k, v: v})
	c.m[k] = &ref
}

func (c *Cache[K, V]) Has(k K) bool {
	_, exists := c.m[k]
	return exists
}

func (c *Cache[K, V]) Get(k K) V {
	ref := c.m[k]
	return c.heap.Get(*ref).v
}

func (c *Cache[K, V]) Pop() (K, V) {
	_, item := c.heap.Pop()
	delete(c.m, item.k)
	return item.k, item.v
}

func (c *Cache[K, V]) Remove(k K) V {
	ref := *c.m[k]
	delete(c.m, k)
	return c.heap.Remove(ref).v
}

func (c *Cache[K, V]) Hit(k K) {
	ref := *c.m[k]
	c.heap.ToEnd(ref)
}

func (c *Cache[K, V]) Len() int {
	return len(c.m)
}

// TODO wrap in an iterator
func (c *Cache[K, V]) VisitAll(f func(*V)) {
	c.heap.VisitAll(func(k *keyVal[K, V]) { f(&k.v) })
}
