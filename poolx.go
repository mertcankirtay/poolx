// Package poolx provides generic, bounded memory pools built on top of sync.Pool.
// It is designed to safely reuse memory while preventing common pooling pitfalls
// such as memory leaks and dirty reads from stale data.
package poolx
