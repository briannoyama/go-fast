package threadless

type RefFactory[V any] struct {
	add    func(*int, V)
	get    func(int) V
	modify func(int) *V
	remove func(int) V
	count  int
}

func (r *RefFactory[V]) Ref() CRef[V] {
	r.count += 1
	return CRef[V]{factory: r}
}

func (r *RefFactory[V]) Length() int {
	return r.count
}

type CRef[V any] struct {
	factory *RefFactory[V]
	ref     int
}

func (c *CRef[V]) Get() V {
	return c.factory.get(c.ref)
}

func (c *CRef[V]) Modify() *V {
	return c.factory.modify(c.ref)
}

func (c *CRef[V]) Set(v V) {
	c.factory.add(&c.ref, v)
}

func (c *CRef[V]) Unset() V {
	return c.factory.remove(c.ref)
}

func (c *CRef[V]) Destroy() {
	c.factory.count -= 1
	c.factory = nil
}

func (c CRef[V]) WCache(v V) CRefCached[V] {
	return CRefCached[V]{CRef: c, v: v}
}

type CRefCached[V any] struct {
	CRef[V]
	v V
}

func (c *CRefCached[V]) Set() {
	c.CRef.Set(c.v)
}

func (c *CRefCached[V]) Unset() {
	c.v = c.CRef.Unset()
}
