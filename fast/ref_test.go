package fast

import (
	"testing"

	"github.com/briannoyama/go-fast/assert"
)

var (
	refFactory RefFactory[int]
	ref        Ref[int]
)

func setupFactory() {
	ival := -1
	refFactory = RefFactory[int]{
		add: func(ref *int, v int) {
			ival = v
		},
		get: func(ref int) int {
			return ival
		},
		modify: func(ref int, f func(v *int)) {
			f(&ival)
		},
		remove: func(ref int) int {
			rval := ival
			ival = -1
			return rval
		},
	}
}

func setupRef() {
	setupFactory()
	ref = refFactory.Ref()
}

func TestRefFactory(t *testing.T) {
	setupFactory()

	r := refFactory.Ref()
	assert.Equals(t, refFactory.Len(), 1)
	refFactory.Ref()
	assert.Equals(t, refFactory.Len(), 2)
	r.Destroy()
	assert.Equals(t, refFactory.Len(), 1)
}

func TestRef(t *testing.T) {
	setupRef()

	ref.Set(5)
	assert.Equals(t, ref.Get(), 5)

	ref.Modify(func(v *int) { *v *= 5 })
	assert.Equals(t, ref.Get(), 25)

	ref.Unset()
	assert.Equals(t, ref.Get(), -1)
}

func TestRefCached(t *testing.T) {
	setupRef()

	rc := ref.WCache(5)
	assert.Equals(t, ref.Get(), -1)

	rc.Set()
	assert.Equals(t, ref.Get(), 5)

	rc.Modify(func(v *int) { *v *= 5 })
	rc.Unset()
	assert.Equals(t, ref.Get(), -1)
	assert.Equals(t, rc.Get(), 25)

	rc.Set()
	assert.Equals(t, ref.Get(), 25)
}
