package lazy

import (
	"fmt"
	"strings"
	"time"

	"github.com/gollin/packages/support/collection/arr"
)

// Each iterates over items, calling the callback for each one.
// If the callback returns false, iteration stops early.
func (lc *Collection[T]) Each(callback func(T, int) bool) *Collection[T] {
	idx := 0
	lc.source(func(item T) bool {
		cont := callback(item, idx)
		idx++

		return cont
	})

	return lc
}

// Tap passes the lazy collection to the given callback and returns it unchanged.
func (lc *Collection[T]) Tap(callback func(*Collection[T])) *Collection[T] {
	callback(lc)

	return lc
}

// TapEach returns a new lazy collection that calls the callback on each item as it passes through.
func (lc *Collection[T]) TapEach(callback func(T, int)) *Collection[T] {
	return New(func(yield func(T) bool) {
		idx := 0
		lc.source(func(item T) bool {
			callback(item, idx)
			idx++

			return yield(item)
		})
	})
}

// Throttle returns a new lazy collection that inserts a delay between each yielded item.
func (lc *Collection[T]) Throttle(delay time.Duration) *Collection[T] {
	return New(func(yield func(T) bool) {
		first := true
		lc.source(func(item T) bool {
			if !first {
				time.Sleep(delay)
			}

			first = false

			return yield(item)
		})
	})
}

// Remember returns a lazy collection that caches items on the first iteration,
// so later iterations reuse the cached values.
func (lc *Collection[T]) Remember() *Collection[T] {
	var cache []T
	cached := false

	return New(func(yield func(T) bool) {
		if cached {
			for _, item := range cache {
				if !yield(item) {
					return
				}
			}

			return
		}

		cache = make([]T, 0)
		lc.source(func(item T) bool {
			cache = append(cache, item)

			return yield(item)
		})
		cached = true
	})
}

// Implode joins all items into a single string separated by the given glue.
func (lc *Collection[T]) Implode(glue string) string {
	parts := make([]string, 0)
	lc.source(func(item T) bool {
		parts = append(parts, fmt.Sprint(item))

		return true
	})

	return strings.Join(parts, glue)
}

// Join joins all items into a string with the given separator.
// An optional final separator can be provided for the last element.
func (lc *Collection[T]) Join(glue string, finalGlues ...string) string {
	parts := make([]string, 0)
	lc.source(func(item T) bool {
		parts = append(parts, fmt.Sprint(item))

		return true
	})

	return arr.Join(parts, glue, finalGlues...)
}

// When applies the callback if the condition is true; otherwise applies the optional default.
func (lc *Collection[T]) When(condition bool, callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	if condition {
		return callback(lc)
	}

	if len(defaults) > 0 {
		return defaults[0](lc)
	}

	return lc
}

// WhenEmpty applies the callback if the lazy collection is empty.
func (lc *Collection[T]) WhenEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return lc.When(lc.IsEmpty(), callback, defaults...)
}

// WhenNotEmpty applies the callback if the lazy collection is not empty.
func (lc *Collection[T]) WhenNotEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return lc.When(lc.IsNotEmpty(), callback, defaults...)
}

// Unless applies the callback unless the condition is true.
func (lc *Collection[T]) Unless(condition bool, callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return lc.When(!condition, callback, defaults...)
}

// Dump prints the lazy collection items to stdout for debugging and returns a new
// lazy collection backed by the evaluated items.
func (lc *Collection[T]) Dump() *Collection[T] {
	items := lc.All()
	fmt.Printf("%v\n", items)

	return From(items)
}
