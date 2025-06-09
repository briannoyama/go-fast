package stable

import (
	"testing"

	"github.com/briannoyama/go-fast/assert"
)

func TestIntMaxIndex(t *testing.T) {
	assert.Equals(t, IntMaxIndex(-1, -1), 0)
	assert.Equals(t, IntMaxIndex(-3, -1), 1)
	assert.Equals(t, IntMaxIndex(0, 0), 0)
	assert.Equals(t, IntMaxIndex(3, 2), 0)
	assert.Equals(t, IntMaxIndex(2, 3), 1)
	assert.Equals(t, IntMaxIndex(-5, 30), 1)
}

func TestIntZeroIfEqual(t *testing.T) {
	assert.Equals(t, IntZeroIfEqual(0, -1), 1)
	assert.Equals(t, IntZeroIfEqual(1, -1), 1)
	assert.Equals(t, IntZeroIfEqual(5, 5), 0)
}

func TestIntZeroIfZero(t *testing.T) {
	assert.Equals(t, IntZeroIfZero(0), 0)
	assert.Equals(t, IntZeroIfZero(-569), 1)
	assert.Equals(t, IntZeroIfZero(1024), 1)
}
