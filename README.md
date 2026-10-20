# poolx

[![Go Reference](https://pkg.go.dev/badge/github.com/mertcankirtay/poolx.svg)](https://pkg.go.dev/github.com/mertcankirtay/poolx)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A generic, bounded `sync.Pool` wrapper for Go. `poolx` is designed to safely reuse memory while preventing common pooling pitfalls such as memory leaks and dirty reads from stale data.

## Features

* **Zero Allocations**: Pools pointers (`*[]T` and `*T`) to avoid interface boxing allocations on `Get` and `Recycle`.
* **Bounded Capacity**: `Slice` pools enforce a `maxCap` limit. Oversized slices are dropped to the garbage collector, preventing long-term memory bloat.
* **Automatic Sanitization**: 
  * `Slice` pools automatically extend slices to their full capacity and `clear()` them before recycling, preventing memory leaks from nested pointers in shortened slices.
  * `Single` pools automatically reset structs to their zero values upon recycling.
* **Type-Safe**: Fully utilizes Go generics.

## Installation

```bash
go get github.com/mertcankirtay/poolx
```

## Usage

### Pooling Slices

The `Slice` pool is ideal for high-throughput scenarios where you frequently allocate and discard slices, such as reading network buffers or building batches.

```go
package main

import (
	"fmt"
	"github.com/mertcankirtay/poolx"
)

var pool *poolx.Slice[int]

func main() {
	// Init the pool for []int.
	// maxCap: 1024 (slices larger than this are sent to the GC)
	// initialCap: 64 (newly allocated slices start with this capacity)
	pool = poolx.NewSlice[int](1024, 64)

	// 1. Get a slice pointer from the pool
	s := pool.Get()

	// 2. Use it (dereference the pointer)
	*s = append(*s, 1, 2, 3)
	fmt.Println(*s) // Output: [1 2 3]

	// 3. Recycle it
	// The pool automatically clears the underlying array and resets length to 0
	pool.Recycle(s)
}
```

### Pooling Single Objects (Structs)

The `Single` pool is perfect for complex structs, preventing the allocation overhead of frequently creating and destroying them.

```go
package main

import (
	"fmt"
	"github.com/mertcankirtay/poolx"
)

type User struct {
	ID   int
	Name string
}

var pool *poolx.Single[User]

func main() {
	// Init the pool for User structs
	pool = poolx.NewSingle[User]()

	// 1. Get a pointer to a User
	user := pool.Get()

	// 2. Use it
	user.ID = 1
	user.Name = "Alice"
	fmt.Printf("%+v\n", *user)

	// 3. Recycle it
	// The pool automatically resets the struct to its zero value ({ID: 0, Name: ""})
	pool.Recycle(user)
}
```

## License

MIT License
