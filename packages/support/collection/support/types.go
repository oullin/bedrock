package support

// Numeric is a constraint for numeric types.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Pair represents a key-value pair.
type Pair[K any, V any] struct {
	Key   K
	Value V
}
