package support

import "sync"

type onceCache struct {
	mu       sync.Mutex
	enabled  bool
	disabled int
	values   map[string]any
}

var onceStore = &onceCache{values: make(map[string]any)}

// Once memoizes a callback result by an explicit key.
func Once[T any](key string, fn func() T) T {
	onceStore.mu.Lock()
	enabled := onceStore.disabled == 0

	if enabled {
		if value, ok := onceStore.values[key]; ok {
			onceStore.mu.Unlock()

			return value.(T)
		}
	}

	onceStore.mu.Unlock()

	value := fn()

	if !enabled {
		return value
	}

	onceStore.mu.Lock()
	onceStore.values[key] = value
	onceStore.mu.Unlock()

	return value
}

// FlushOnce clears all memoized Once results.
func FlushOnce() {
	onceStore.mu.Lock()

	defer onceStore.mu.Unlock()

	onceStore.values = make(map[string]any)
}

// DisableOnce disables memoization until EnableOnce is called.
func DisableOnce() {
	onceStore.mu.Lock()

	defer onceStore.mu.Unlock()

	onceStore.disabled++
}

// EnableOnce reverses one DisableOnce call.
func EnableOnce() {
	onceStore.mu.Lock()

	defer onceStore.mu.Unlock()

	if onceStore.disabled > 0 {
		onceStore.disabled--
	}
}

// WithoutOnce runs fn while memoization is temporarily disabled.
func WithoutOnce(fn func()) {
	DisableOnce()

	defer EnableOnce()

	fn()
}
