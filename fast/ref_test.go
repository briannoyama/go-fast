package fast

import (
	"testing"

	"go-fast/assert"
)

func TestRefSet(t *testing.T) {
	c := CPage[int]{}
	f := c.Factory()

	r := f.Ref()
	r.Set(1)

	data := []int{1}

	assert.Equals(t, c.items, data)
}
