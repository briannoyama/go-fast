package fast

import (
	"testing"

	"github.com/briannoyama/go-fast/v2/assert"
)

func TestSPage(t *testing.T) {
	spage := NewSPage[int]()
	refs := []int{spage.Add(5), spage.Add(3), spage.Add(2), spage.Add(7)}
	assert.Equals(t, spage.Get(refs[0]), 5)
	assert.Equals(t, spage.Remove(refs[1]), 3)
	assert.Equals(t, spage.Get(refs[2]), 2)
	assert.Equals(t, spage.Remove(refs[3]), 7)
	refs = []int{spage.Add(4), spage.Add(6)}
	assert.Equals(t, refs[0], 3)
	assert.Equals(t, refs[1], 1)
	assert.Equals(t, spage.Get(refs[0]), 4)
	assert.Equals(t, spage.Get(refs[1]), 6)
}
