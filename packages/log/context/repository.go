package context

import (
	"fmt"
	"sync"
)

// Repository provides a scoped key-value store for structured log context.
// It supports public data, hidden (sensitive) data, stack operations, counters,
// scoped execution, and serialization hooks.
type Repository struct {
	mu          sync.RWMutex
	data        map[string]any
	hidden      map[string]any
	dehydrating []func(*Repository)
	hydrated    []func(*Repository)
}

// New creates an empty Repository.
func New() *Repository {
	return &Repository{
		data:   make(map[string]any),
		hidden: make(map[string]any),
	}
}

// Has reports whether all of the given keys are set in the public context.
func (r *Repository) Has(keys ...string) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	for _, key := range keys {
		if _, ok := r.data[key]; !ok {
			return false
		}
	}

	return true
}

// Missing reports whether the given key is not set.
func (r *Repository) Missing(key string) bool {
	return !r.Has(key)
}

// HasHidden reports whether all of the given keys are set in the hidden context.
func (r *Repository) HasHidden(keys ...string) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	for _, key := range keys {
		if _, ok := r.hidden[key]; !ok {
			return false
		}
	}

	return true
}

// MissingHidden reports whether the given key is not set in the hidden context.
func (r *Repository) MissingHidden(key string) bool {
	return !r.HasHidden(key)
}

// Get returns the value for the given key, or the fallback if not set.
func (r *Repository) Get(key string, fallback ...any) any {
	r.mu.RLock()

	defer r.mu.RUnlock()

	if v, ok := r.data[key]; ok {
		return v
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return nil
}

// GetHidden returns the value for the given key from the hidden context.
func (r *Repository) GetHidden(key string, fallback ...any) any {
	r.mu.RLock()

	defer r.mu.RUnlock()

	if v, ok := r.hidden[key]; ok {
		return v
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return nil
}

// All returns a copy of all public context data.
func (r *Repository) All() map[string]any {
	r.mu.RLock()

	defer r.mu.RUnlock()

	return r.copyMap(r.data)
}

// AllHidden returns a copy of all hidden context data.
func (r *Repository) AllHidden() map[string]any {
	r.mu.RLock()

	defer r.mu.RUnlock()

	return r.copyMap(r.hidden)
}

// Only returns a map containing only the specified keys.
func (r *Repository) Only(keys ...string) map[string]any {
	r.mu.RLock()

	defer r.mu.RUnlock()

	result := make(map[string]any, len(keys))

	for _, key := range keys {
		if v, ok := r.data[key]; ok {
			result[key] = v
		}
	}

	return result
}

// OnlyHidden returns a map of only the specified hidden keys.
func (r *Repository) OnlyHidden(keys ...string) map[string]any {
	r.mu.RLock()

	defer r.mu.RUnlock()

	result := make(map[string]any, len(keys))

	for _, key := range keys {
		if v, ok := r.hidden[key]; ok {
			result[key] = v
		}
	}

	return result
}

// Except returns all public data except the specified keys.
func (r *Repository) Except(keys ...string) map[string]any {
	r.mu.RLock()

	defer r.mu.RUnlock()

	exclude := make(map[string]struct{}, len(keys))

	for _, key := range keys {
		exclude[key] = struct{}{}
	}

	result := make(map[string]any, len(r.data))

	for k, v := range r.data {
		if _, skip := exclude[k]; !skip {
			result[k] = v
		}
	}

	return result
}

// ExceptHidden returns all hidden data except the specified keys.
func (r *Repository) ExceptHidden(keys ...string) map[string]any {
	r.mu.RLock()

	defer r.mu.RUnlock()

	exclude := make(map[string]struct{}, len(keys))

	for _, key := range keys {
		exclude[key] = struct{}{}
	}

	result := make(map[string]any, len(r.hidden))

	for k, v := range r.hidden {
		if _, skip := exclude[k]; !skip {
			result[k] = v
		}
	}

	return result
}

// Add sets a key-value pair in the public context.
func (r *Repository) Add(key string, value any) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.data[key] = value

	return r
}

// AddHidden sets a key-value pair in the hidden context.
func (r *Repository) AddHidden(key string, value any) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.hidden[key] = value

	return r
}

// AddIf sets a key-value pair only if the key is not already present.
func (r *Repository) AddIf(key string, value any) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	if _, ok := r.data[key]; !ok {
		r.data[key] = value
	}

	return r
}

// AddHiddenIf sets a hidden key-value pair only if the key is not already present.
func (r *Repository) AddHiddenIf(key string, value any) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	if _, ok := r.hidden[key]; !ok {
		r.hidden[key] = value
	}

	return r
}

// Forget removes the given keys from the public context.
func (r *Repository) Forget(keys ...string) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, key := range keys {
		delete(r.data, key)
	}

	return r
}

// ForgetHidden removes the given keys from the hidden context.
func (r *Repository) ForgetHidden(keys ...string) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, key := range keys {
		delete(r.hidden, key)
	}

	return r
}

// Pull retrieves a value and removes it from the public context.
func (r *Repository) Pull(key string, fallback ...any) any {
	r.mu.Lock()

	defer r.mu.Unlock()

	if v, ok := r.data[key]; ok {
		delete(r.data, key)

		return v
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return nil
}

// PullHidden retrieves a value and removes it from the hidden context.
func (r *Repository) PullHidden(key string, fallback ...any) any {
	r.mu.Lock()

	defer r.mu.Unlock()

	if v, ok := r.hidden[key]; ok {
		delete(r.hidden, key)

		return v
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return nil
}

// Remember returns the value for the given key. If the key is not set, fn is
// called, the result is stored, and then returned.
func (r *Repository) Remember(key string, fn func() any) any {
	r.mu.Lock()

	defer r.mu.Unlock()

	if v, ok := r.data[key]; ok {
		return v
	}

	v := fn()
	r.data[key] = v

	return v
}

// RememberHidden is like Remember but for the hidden context.
func (r *Repository) RememberHidden(key string, fn func() any) any {
	r.mu.Lock()

	defer r.mu.Unlock()

	if v, ok := r.hidden[key]; ok {
		return v
	}

	v := fn()
	r.hidden[key] = v

	return v
}

// Push appends values to a slice stored at the given key.
func (r *Repository) Push(key string, values ...any) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	existing, ok := r.data[key]

	if !ok {
		r.data[key] = values

		return r
	}

	slice, ok := existing.([]any)

	if !ok {
		panic(fmt.Sprintf("context: cannot push to non-slice value at key %q", key))
	}

	if !isList(slice) {
		panic(fmt.Sprintf("context: cannot push to non-list array at key %q", key))
	}

	r.data[key] = append(slice, values...)

	return r
}

// PushHidden appends values to a slice stored at the given hidden key.
func (r *Repository) PushHidden(key string, values ...any) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	existing, ok := r.hidden[key]

	if !ok {
		r.hidden[key] = values

		return r
	}

	slice, ok := existing.([]any)

	if !ok {
		panic(fmt.Sprintf("context: cannot push to non-slice hidden value at key %q", key))
	}

	if !isList(slice) {
		panic(fmt.Sprintf("context: cannot push to non-list hidden array at key %q", key))
	}

	r.hidden[key] = append(slice, values...)

	return r
}

// Pop removes and returns the last element from the slice at the given key.
func (r *Repository) Pop(key string) any {
	r.mu.Lock()

	defer r.mu.Unlock()

	existing, ok := r.data[key]

	if !ok {
		panic(fmt.Sprintf("context: cannot pop from empty stack at key %q", key))
	}

	slice, ok := existing.([]any)

	if !ok {
		panic(fmt.Sprintf("context: cannot pop from non-slice value at key %q", key))
	}

	if len(slice) == 0 {
		panic(fmt.Sprintf("context: cannot pop from empty stack at key %q", key))
	}

	if !isList(slice) {
		panic(fmt.Sprintf("context: cannot pop from non-list array at key %q", key))
	}

	last := slice[len(slice)-1]
	r.data[key] = slice[:len(slice)-1]

	return last
}

// PopHidden removes and returns the last element from the hidden slice at the
// given key.
func (r *Repository) PopHidden(key string) any {
	r.mu.Lock()

	defer r.mu.Unlock()

	existing, ok := r.hidden[key]

	if !ok {
		panic(fmt.Sprintf("context: cannot pop from empty hidden stack at key %q", key))
	}

	slice, ok := existing.([]any)

	if !ok {
		panic(fmt.Sprintf("context: cannot pop from non-slice hidden value at key %q", key))
	}

	if len(slice) == 0 {
		panic(fmt.Sprintf("context: cannot pop from empty hidden stack at key %q", key))
	}

	if !isList(slice) {
		panic(fmt.Sprintf("context: cannot pop from non-list hidden array at key %q", key))
	}

	last := slice[len(slice)-1]
	r.hidden[key] = slice[:len(slice)-1]

	return last
}

// StackContains reports whether the slice at the given key contains the value.
func (r *Repository) StackContains(key string, value any) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	existing, ok := r.data[key]

	if !ok {
		return false
	}

	slice, ok := existing.([]any)

	if !ok {
		return false
	}

	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

// StackContainsFunc reports whether the slice at the given key contains an
// element satisfying fn.
func (r *Repository) StackContainsFunc(key string, fn func(any) bool) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	existing, ok := r.data[key]

	if !ok {
		return false
	}

	slice, ok := existing.([]any)

	if !ok {
		return false
	}

	for _, item := range slice {
		if fn(item) {
			return true
		}
	}

	return false
}

// HiddenStackContains reports whether the hidden slice at the given key
// contains the value.
func (r *Repository) HiddenStackContains(key string, value any) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	existing, ok := r.hidden[key]

	if !ok {
		return false
	}

	slice, ok := existing.([]any)

	if !ok {
		return false
	}

	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

// HiddenStackContainsFunc reports whether the hidden slice at the given key
// contains an element satisfying fn.
func (r *Repository) HiddenStackContainsFunc(key string, fn func(any) bool) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	existing, ok := r.hidden[key]

	if !ok {
		return false
	}

	slice, ok := existing.([]any)

	if !ok {
		return false
	}

	for _, item := range slice {
		if fn(item) {
			return true
		}
	}

	return false
}

// Increment increments an integer counter at the given key.
func (r *Repository) Increment(key string, by ...int) *Repository {
	amount := 1

	if len(by) > 0 {
		amount = by[0]
	}

	r.mu.Lock()

	defer r.mu.Unlock()

	current, ok := r.data[key]

	if !ok {
		r.data[key] = amount

		return r
	}

	if v, ok := current.(int); ok {
		r.data[key] = v + amount
	}

	return r
}

// Decrement decrements an integer counter at the given key.
func (r *Repository) Decrement(key string, by ...int) *Repository {
	amount := 1

	if len(by) > 0 {
		amount = by[0]
	}

	r.mu.Lock()

	defer r.mu.Unlock()

	current, ok := r.data[key]

	if !ok {
		r.data[key] = -amount

		return r
	}

	if v, ok := current.(int); ok {
		r.data[key] = v - amount
	}

	return r
}

// Scope executes fn with a copy of the current data, discarding any changes
// made within fn after it returns.
func (r *Repository) Scope(fn func(*Repository), data ...map[string]any) *Repository {
	r.mu.RLock()
	dataCopy := r.copyMap(r.data)
	hiddenCopy := r.copyMap(r.hidden)
	r.mu.RUnlock()

	if len(data) > 0 && data[0] != nil {
		for k, v := range data[0] {
			dataCopy[k] = v
		}
	}

	scoped := &Repository{
		data:   dataCopy,
		hidden: hiddenCopy,
	}

	fn(scoped)

	return r
}

// Dehydrate serializes the repository's public data for transport. Returns nil
// if the repository is empty.
func (r *Repository) Dehydrate() map[string]any {
	r.mu.RLock()

	defer r.mu.RUnlock()

	for _, cb := range r.dehydrating {
		cb(r)
	}

	if len(r.data) == 0 {
		return nil
	}

	return r.copyMap(r.data)
}

// Hydrate restores repository data from a previously dehydrated state.
func (r *Repository) Hydrate(data map[string]any) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	if data != nil {
		for k, v := range data {
			r.data[k] = v
		}
	}

	for _, cb := range r.hydrated {
		cb(r)
	}

	return r
}

// Dehydrating registers a callback to run before dehydration.
func (r *Repository) Dehydrating(fn func(*Repository)) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.dehydrating = append(r.dehydrating, fn)

	return r
}

// Hydrated registers a callback to run after hydration.
func (r *Repository) Hydrated(fn func(*Repository)) *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.hydrated = append(r.hydrated, fn)

	return r
}

// Flush clears all data, hidden data, and callbacks.
func (r *Repository) Flush() *Repository {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.data = make(map[string]any)
	r.hidden = make(map[string]any)

	return r
}

// IsEmpty reports whether the repository has no public or hidden data.
func (r *Repository) IsEmpty() bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	return len(r.data) == 0 && len(r.hidden) == 0
}

func (r *Repository) copyMap(m map[string]any) map[string]any {
	cp := make(map[string]any, len(m))

	for k, v := range m {
		cp[k] = v
	}

	return cp
}

// isList checks if a slice is a sequential list (not associative).
// In Go, all []any slices are lists, so this always returns true.
// This exists for API parity with Laravel's array checks.
func isList(_ []any) bool {
	return true
}
