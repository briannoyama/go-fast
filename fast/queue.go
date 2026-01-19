package fast

type Queue[V any] struct {
	values     []V
	start, end int
}

func NewQueue[V any]() Queue[V] {
	return Queue[V]{
		values: make([]V, 2),
		start:  0,
		end:    0,
	}
}

// Add a value to the end of the queue
func (q *Queue[V]) Add(v V) {
	q.values[q.end] = v
	end := q.end + 1
	q.end = end % cap(q.values)
	if q.start == q.end {
		items := make([]V, cap(q.values)<<1)
		copy(items[:q.end], q.values[:q.end])
		copy(items[len(items)-len(q.values)+q.start:], q.values[q.start:])
		q.start += len(q.values)
		q.values = items
	}
}

// Len of the queue
func (q *Queue[V]) Len() int {
	return (q.end - q.start + cap(q.values)) % cap(q.values)
}

// Pop the value from the beginning of the queue
func (q *Queue[V]) Pop() V {
	item := q.values[q.start]
	if q.start != q.end {
		q.start = (q.start + 1) % cap(q.values)
	}
	return item
}
