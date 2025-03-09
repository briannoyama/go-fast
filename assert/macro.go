package assert

import (
	"reflect"
	"testing"
)

func Equals[V any](t *testing.T, v0, v1 V) {
	t.Helper()
	if !reflect.DeepEqual(v0, v1) {
		t.Errorf("%v != %v", v0, v1)
	}
}
