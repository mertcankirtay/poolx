package poolx

import (
	"testing"
)

type mockStruct struct {
	ID   int
	Data string
}

func TestSlicePool(t *testing.T) {
	pool := NewSlice[int](10, 5)

	s1 := pool.Get()
	if cap(*s1) != 5 {
		t.Errorf("expected initial capacity 5, got %d", cap(*s1))
	}

	if len(*s1) != 0 {
		t.Errorf("expected length 0, got %d", len(*s1))
	}

	*s1 = append(*s1, 1, 2, 3)
	pool.Recycle(s1)

	s2 := pool.Get()
	if len(*s2) != 0 {
		t.Errorf("expected length 0 after recycle, got %d", len(*s2))
	}

	extended := (*s2)[:3]
	if extended[0] != 0 || extended[1] != 0 || extended[2] != 0 {
		t.Errorf("expected slice to be cleared, got %v", extended)
	}

	*s2 = make([]int, 0, 15)
	pool.Recycle(s2)

	s3 := pool.Get()
	if cap(*s3) == 15 {
		t.Errorf("slice exceeding maxCap was pooled and returned")
	}
}

func TestSinglePool(t *testing.T) {
	pool := NewSingle[mockStruct]()

	obj1 := pool.Get()
	if obj1.ID != 0 || obj1.Data != "" {
		t.Errorf("expected zeroed struct, got %+v", obj1)
	}

	obj1.ID = 42
	obj1.Data = "poolx"
	pool.Recycle(obj1)

	obj2 := pool.Get()
	if obj2.ID != 0 || obj2.Data != "" {
		t.Errorf("expected struct to be zeroed on recycle, got %+v", obj2)
	}
}

func BenchmarkSlicePool_GetPut(b *testing.B) {
	pool := NewSlice[int](1024, 64)
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		s := pool.Get()
		pool.Recycle(s)
	}
}

func BenchmarkSinglePool_GetPut(b *testing.B) {
	pool := NewSingle[mockStruct]()
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		obj := pool.Get()
		pool.Recycle(obj)
	}
}
