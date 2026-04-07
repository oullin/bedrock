// Package cache provides a Laravel-inspired cache store abstraction.
//
// The Store interface defines get/put/forget/flush semantics with TTL
// support. ArrayStore is an in-memory implementation suitable for
// testing and single-process use. It is safe for concurrent use.
//
//	s := cache.New()
//	s.Put(ctx, "key", "value", 5*time.Minute)
//	v, _ := s.Get(ctx, "key")
package cache
