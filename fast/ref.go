package fast

// RefFactory creates CRef references that allow operations on values in various
// data structures
type RefFactory[V any] struct {
	add    func(*int, V)
	get    func(int) V
	modify func(int, func(*V))
	remove func(int) V
	count  int
}

// Ref creates an empty reference to a value.
func (r *RefFactory[V]) Ref() Ref[V] {
	r.count += 1
	return Ref[V]{factory: r, ref: -1}
}

// Len(gth) returns the count of active references.
func (r *RefFactory[V]) Len() int {
	return r.count
}

// Ref provides methods for manipulating a single value.
type Ref[V any] struct {
	factory *RefFactory[V]
	ref     int
}

// Get the value pointed to by the reference.
func (r *Ref[V]) Get() V {
	return r.factory.get(r.ref)
}

// Modify the value pointed to by the reference in place by applying f.
// Does not check if the reference is set.
func (r *Ref[V]) Modify(f func(*V)) {
	r.factory.modify(r.ref, f)
}

// Set a value for the reference. Does not check if the reference is already set.
func (r *Ref[V]) Set(v V) {
	r.factory.add(&r.ref, v)
}

// Unset the value for the reference. Does not check if reference is not set.
// Does not affect the count of active references in the RefFactory,
func (c *Ref[V]) Unset() V {
	return c.factory.remove(c.ref)
}

// IsSet is true iff a value is set for the reference in the underlying data structure.
func (c *Ref[V]) IsSet() bool {
	return c.ref >= 0
}

// Destroy decreases decreases the active reference count in RefFactory and
// makes the CRef reference unusable.
func (c *Ref[V]) Destroy() {
	if c.IsSet() {
		c.Unset()
	}
	c.factory.count -= 1
	c.factory = nil
}

// CRefCached creates an unset cached reference.
func (c Ref[V]) WCache(v V) RefCached[V] {
	return RefCached[V]{Ref: c, v: v}
}

// RefCached holds a value that can be set/unset.
// Any modifications of the value while set will be stored.
type RefCached[V any] struct {
	Ref[V]
	v V
}

// Set the value.
func (c *RefCached[V]) Set() {
	c.Ref.Set(c.v)
}

// Get the value (works even if not set).
func (c *RefCached[V]) Get() V {
	return c.v
}

// Unset the value
func (c *RefCached[V]) Unset() {
	c.v = c.Ref.Unset()
}
