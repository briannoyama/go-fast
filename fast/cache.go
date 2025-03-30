package fast

type keyVal[K, V any] struct {
	k K
	v V
}

// Cache tracks the age of Keys/Values allowing for removal of values based on age.
type Cache[K comparable, V any] struct {
	heap Heap[keyVal[K, V]]
	m    map[K]*int
}

// MakeCache initializes a Cache
func MakeCache[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{m: make(map[K]*int)}
}

// Add a Key/Calue pair.
func (c *Cache[K, V]) Add(k K, v V) {
	var ref int
	c.heap.Add(&ref, keyVal[K, V]{k: k, v: v})
	c.m[k] = &ref
}

// Has returns true iff the key exists.
func (c *Cache[K, V]) Has(k K) bool {
	_, exists := c.m[k]
	return exists
}

// Get a value for the key. Panics if key does not exist.
func (c *Cache[K, V]) Get(k K) V {
	ref := c.m[k]
	return c.heap.Get(*ref).v
}

// Pop the oldest key/value from the cache.
func (c *Cache[K, V]) Pop() (K, V) {
	_, item := c.heap.Pop()
	delete(c.m, item.k)
	return item.k, item.v
}

// Remove a key/value using the key.
func (c *Cache[K, V]) Remove(k K) V {
	ref := *c.m[k]
	delete(c.m, k)
	return c.heap.Remove(ref).v
}

// Hit refreshes the age of the key/value pair matching the key.
// Hit will mark the item as the most recent.
func (c *Cache[K, V]) Hit(k K) {
	ref := *c.m[k]
	c.heap.ToEnd(ref)
}

// Len(gth) returns the number of items in the Cache.
func (c *Cache[K, V]) Len() int {
	return len(c.m)
}

// VisitAll applies the funtion f to all values in the Cache.
func (c *Cache[K, V]) VisitAll(f func(*V)) {
	c.heap.VisitAll(func(k *keyVal[K, V]) { f(&k.v) })
}
