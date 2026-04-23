package support

import "time"

// BenchmarkResult captures a measured callback result and elapsed time.
type BenchmarkResult[T any] struct {
	Value    T
	Duration time.Duration
}

// BenchmarkValue runs fn and returns its value with elapsed duration.
func BenchmarkValue[T any](fn func() T) BenchmarkResult[T] {
	start := time.Now()
	value := fn()

	return BenchmarkResult[T]{Value: value, Duration: time.Since(start)}
}
