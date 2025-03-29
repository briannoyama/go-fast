package fast

import (
	"testing"

	"go-fast/assert"
)

var cache Cache[int, int]

func setupCache() {
	cache = MakeCache[int, int]()
	for i := 0; i < 5; i++ {
		cache.Add(i, i)
	}
}

func TestCacheAdd(t *testing.T) {
	setupCache()

	cache.Add(20, 6)
	assert.Equals(t, cache.Get(20), 6)
}

func TestCacheHas(t *testing.T) {
	setupCache()

	assert.Equals(t, cache.Has(3), true)
	assert.Equals(t, cache.Has(20), false)
}

func TestCachePop(t *testing.T) {
	setupCache()

	k, v := cache.Pop()
	assert.Equals(t, k, 0)
	assert.Equals(t, v, 0)
	assert.Equals(t, cache.Len(), 4)
	assert.Equals(t, cache.Has(0), false)
}

func TestCacheRemove(t *testing.T) {
	setupCache()

	assert.Equals(t, cache.Remove(1), 1)
	assert.Equals(t, cache.Remove(2), 2)
	assert.Equals(t, cache.Len(), 3)
	assert.Equals(t, cache.Has(1), false)
	assert.Equals(t, cache.Has(2), false)
}

func TestCacheHit(t *testing.T) {
	setupCache()

	cache.Hit(0)
	cache.Hit(1)
	k, v := cache.Pop()
	assert.Equals(t, k, 2)
	assert.Equals(t, v, 2)
	assert.Equals(t, cache.Len(), 4)
	assert.Equals(t, cache.Has(0), true)
	assert.Equals(t, cache.Has(1), true)
}

func TestCacheVisitAll(t *testing.T) {
	setupCache()
	cache.VisitAll(func(i *int) { *i *= 5 })
	assert.Equals(t, cache.Get(1), 5)
	assert.Equals(t, cache.Get(2), 10)
}
