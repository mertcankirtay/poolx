package poolx

import "sync"

// Slice is a generic, bounded memory pool for slices.
// It pools pointers to slices (*[]T) to achieve zero-allocation performance
// during Get and Recycle operations.
type Slice[T any] struct {
	internal sync.Pool
	maxCap   uint64
}

// NewSlice creates a new bounded pool for slices.
// maxCap defines the maximum capacity a slice can have to be eligible for recycling.
// Slices exceeding maxCap are dropped and left for the garbage collector.
// initialCap determines the pre-allocated capacity of newly created slices.
func NewSlice[T any](maxCap, initialCap uint64) (p *Slice[T]) {
	p = new(Slice[T]{
		maxCap: maxCap,
		internal: sync.Pool{
			New: func() any {
				s := make([]T, 0, initialCap)

				return &s
			},
		},
	})

	return
}

// Get retrieves a pointer to a slice from the pool.
// The returned slice is guaranteed to have a length of 0 and all underlying
// elements up to its capacity are cleared.
func (p *Slice[T]) Get() (ptr *[]T) {
	ptr = p.internal.Get().(*[]T)

	return ptr
}

// Recycle safely returns a slice to the pool.
// It enforces the maxCap limit, fully clears all elements up to the slice's
// capacity to prevent memory leaks from nested pointers, and resets the length to 0.
func (p *Slice[T]) Recycle(ptr *[]T) {
	if uint64(cap(*ptr)) > p.maxCap {
		return
	}

	*ptr = (*ptr)[:cap(*ptr)]
	clear(*ptr)

	*ptr = (*ptr)[:0]

	p.internal.Put(ptr)
}
