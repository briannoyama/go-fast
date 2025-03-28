package fast

import stable "go-fast/math"

type heapItem[V any] struct {
	priority int
	v        V
}

type Heap[V any] struct {
	CPage[heapItem[V]]
	lastPriority int
}

func (h *Heap[V]) Add(ref *int, v V) {
	h.lastPriority += 1
	h.CPage.Add(ref, heapItem[V]{priority: h.lastPriority, v: v})
	h.up(h.Len() - 1)
}

// AddRelative to the priority of the last value.
// Positive priorities are not guaranteed to work due to the possibility of int overflow.
// Positive priorities allow adding items past the "end", such that calles to "ToEnd" moves
// a value to before the relatively added value.
func (h *Heap[V]) AddRelative(ref *int, v V, priority int) {
	h.CPage.Add(ref, heapItem[V]{priority: h.lastPriority + priority, v: v})
	h.up(h.Len() - 1)
}

func (h *Heap[V]) Get(ref int) V {
	return h.CPage.Get(ref).v
}

func (h *Heap[V]) Modify(ref int, f func(*V)) {
	f(&h.items[ref].v)
}

func (h *Heap[V]) ModifyPriority(ref int, priority int) {
	h.items[ref].priority += priority
	rp := h.refs[ref]
	h.up(ref)
	h.down(*rp)
}

func (h *Heap[V]) Pop() (*int, V) {
	h.swap(0, h.Len()-1)
	ref, item := h.CPage.Pop()

	h.down(0)
	return ref, item.v
}

func (h *Heap[V]) Remove(ref int) V {
	h.swap(ref, h.Len()-1)
	_, item := h.CPage.Pop()
	rp := h.refs[ref]
	h.up(ref)
	h.down(*rp)

	return item.v
}

func (h *Heap[V]) ToEnd(ref int) {
	h.lastPriority += 1
	h.items[ref].priority = h.lastPriority
	h.down(ref)
}

func (h *Heap[V]) Factory() RefFactory[V] {
	return RefFactory[V]{
		add:    h.Add,
		get:    h.Get,
		modify: h.Modify,
		remove: h.Remove,
	}
}

func (h *Heap[V]) VisitAll(f func(*V)) {
	v := h.CPage.Visitor(func(h *heapItem[V]) { f(&h.v) })
	v.VisitAll()
}

func (h *Heap[V]) down(ref int) {
	l := h.Len()
	left := (ref << 1) + 1
	for ; left+1 < l; left = (ref << 1) + 1 {
		k := stable.IntMaxIndex(h.priority(left+1), h.priority(left))

		// left can be the left or right child dependent on k
		left += k
		if h.priority(ref) <= h.priority(left) {
			// Heap property maintained
			return
		}
		h.swap(ref, left)
		ref = left
	}
	// Chance there's one more left node
	if left < l && h.priority(ref) > h.priority(left) {
		h.swap(ref, left)
	}
}

func (h *Heap[V]) priority(ref int) int {
	return h.CPage.Get(ref).priority - h.lastPriority
}

func (h *Heap[V]) up(ref int) {
	for next := 0; ref > 0; ref = next {
		next = (ref - 1) >> 1
		if h.priority(ref) >= h.priority(next) {
			return
		}
		h.swap(ref, next)
	}
}
