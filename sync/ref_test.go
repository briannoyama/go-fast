package sync

import (
	"reflect"
	"testing"
)

func TestVertBufferAdd(t *testing.T) {
	c := CPage[int]{}
	f := c.Factory()

	r := f.Ref()
	r.Set(1)

	data := []int{1}

	if reflect.DeepEqual(c.items, data) {
		t.Fail()
	}
}
