package fast

import stable "go-fast/math"

type heapItem[V any] struct {
	priority int
	v        V
}

// Heap holds items with periorities. Implicitly tracks the "lastPriority"
type Heap[V any] struct {
	CPage[heapItem[V]]
	lastPriority int
}

// Add a value with a reference and the next available priority.
func (h *Heap[V]) Add(ref *int, v V) {
	h.lastPriority += 1
	h.CPage.Add(ref, heapItem[V]{priority: h.lastPriority, v: v})
	h.up(h.Len() - 1)
}

// AddRelative to the priority of the last value.
func (h *Heap[V]) AddRelative(ref *int, v V, rel int) {
	h.CPage.Add(ref, heapItem[V]{priority: h.lastPriority + rel, v: v})
	// Increase lastPriority if priority is positive.
	h.lastPriority += stable.IntMaxIndex(0, rel) * rel
	h.up(h.Len() - 1)
}

// Get the value that a reference points to.
func (h *Heap[V]) Get(ref int) V {
	return h.CPage.Get(ref).v
}

// Modify the value a reference points to in place by applying f.
func (h *Heap[V]) Modify(ref int, f func(*V)) {
	f(&h.items[ref].v)
}

// Modify the priority of the value a reference points to.
// This will potentially change the order of Pop.
func (h *Heap[V]) ModifyPriority(ref int, priority int) {
	h.items[ref].priority += priority
	// Increase lastPriority if priority is positive.
	h.lastPriority += stable.IntMaxIndex(0, priority) * priority
	rp := h.refs[ref]
	h.up(ref)
	h.down(*rp)
}

// Pop the value with the lowest priority (earliest Add if AddRelative is not used).
func (h *Heap[V]) Pop() (*int, V) {
	h.swap(0, h.Len()-1)
	ref, item := h.CPage.Pop()

	h.down(0)
	return ref, item.v
}

// Remove the value pointed to by the reference.
func (h *Heap[V]) Remove(ref int) V {
	h.swap(ref, h.Len()-1)
	_, item := h.CPage.Pop()
	rp := h.refs[ref]
	h.up(ref)
	h.down(*rp)

	return item.v
}

// ToEnd sets the priority of the value a reference points to to the next available priority.
func (h *Heap[V]) ToEnd(ref int) {
	h.lastPriority += 1
	h.items[ref].priority = h.lastPriority
	h.down(ref)
}

// Factory can be used to get Ref(s) to this structure.
func (h *Heap[V]) Factory() RefFactory[V] {
	return RefFactory[V]{
		add:    h.Add,
		get:    h.Get,
		modify: h.Modify,
		remove: h.Remove,
	}
}

// VisitAll applies the function f to all values in the heap.
// Note: Visitor is not exposed for heaps since removals can result in unexpected orders
// in the underlying page.
func (h *Heap[V]) VisitAll(f func(*V)) {
	v := h.CPage.Visitor(func(h *heapItem[V]) { f(&h.v) })
	v.VisitAll()
}

// Down moves a value down the heap based on its priority.
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

// Priority of a value in the heap that handles possible int overflows.
func (h *Heap[V]) priority(ref int) int {
	return h.CPage.Get(ref).priority - h.lastPriority
}

// Up moves a value up the heap based on its priority.
func (h *Heap[V]) up(ref int) {
	for next := 0; ref > 0; ref = next {
		next = (ref - 1) >> 1
		if h.priority(ref) >= h.priority(next) {
			return
		}
		h.swap(ref, next)
	}
}
