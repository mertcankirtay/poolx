package poolx

import "sync"

// Single is a generic memory pool for individual objects, such as structs.
// It pools pointers to the objects (*T) to prevent interface boxing allocations.
type Single[T any] struct {
	internal sync.Pool
}

// NewSingle creates a new zero-allocation pool for single objects of type T.
func NewSingle[T any]() (p *Single[T]) {
	p = new(Single[T]{
		internal: sync.Pool{
			New: func() any {
				return new(T)
			},
		},
	})

	return
}

// Get retrieves a pointer to an object from the pool.
// The returned object is guaranteed to be set to its zero value.
func (p *Single[T]) Get() (ptr *T) {
	ptr = p.internal.Get().(*T)

	return ptr
}

// Recycle safely returns an object to the pool.
// It automatically resets the object to its zero value, ensuring that any
// nested pointers (like slices or maps inside a struct) are nilled out,
// preventing background memory leaks.
func (p *Single[T]) Recycle(ptr *T) {
	if ptr == nil {
		return
	}

	var z T
	*ptr = z

	p.internal.Put(ptr)
}
