package fast

import (
	"testing"

	"github.com/briannoyama/go-fast/v2/assert"
)

func TestQueue(t *testing.T) {
	queue := NewQueue[int]()
	queue.Add(0)
	queue.Add(1)
	queue.Add(2)
	assert.Equals(t, queue.Len(), 3)

	v := queue.Pop()
	assert.Equals(t, v, 0)
	v = queue.Pop()
	assert.Equals(t, v, 1)
	assert.Equals(t, queue.Len(), 1)

	queue.Add(3)
	queue.Add(4)
	queue.Add(5)
	queue.Add(6)
	queue.Add(7)
	assert.Equals(t, queue.Len(), 6)

	v = queue.Pop()
	assert.Equals(t, v, 2)

	v = queue.Pop()
	assert.Equals(t, v, 3)
	queue.Pop()
	assert.Equals(t, queue.Len(), 3)
	queue.Pop()
	queue.Pop()
	queue.Pop()
	queue.Pop()
	assert.Equals(t, queue.Len(), 0)
}
