package fast

// SPage creates a page of values with static references
type SPage[V any] struct {
	values []V
	free   []int
}

func NewSPage[V any]() SPage[V] {
	return SPage[V]{}
}

func NewPreAllocSPage[V any](len int) SPage[V] {
	return SPage[V]{values: make([]V, 0, len)}
}

// Add a value. Returns a static reference set to the position of the value.
// Will write over values marked for removal
func (s *SPage[V]) Add(v V) int {
	if len(s.free) == 0 {
		s.values = append(s.values, v)
		return len(s.values) - 1
	} else {
		pop := len(s.free) - 1
		i := s.free[pop]
		s.free = s.free[:pop]
		s.values[i] = v
		return i
	}
}

// Get the value that a reference points to.
func (s *SPage[V]) Get(ref int) V {
	return s.values[ref]
}

// Len (gth) or number of values held in the page.
func (s *SPage[V]) Len() int {
	return len(s.values) - len(s.free)
}

// Modify the value a reference points to in place by applying f.
func (c *SPage[V]) Modify(ref int, f func(*V)) {
	f(&c.values[ref])
}

// Remove marks an item for removal.
// Does not actually remove or check if the item can be removed.
func (s *SPage[V]) Remove(ref int) V {
	s.free = append(s.free, ref)
	return s.values[ref]
}
