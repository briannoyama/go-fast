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

func Contains[V any](t *testing.T, vS []V, vI ...V) {
	t.Helper()
	for _, i := range vI {
		notIn := true
		for _, s := range vS {
			if reflect.DeepEqual(s, i) {
				notIn = false
			}
		}
		if notIn {
			t.Errorf("%v not in %v", i, vS)
		}
	}
}
