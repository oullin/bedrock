package lazy

import (
	"cmp"
	"math/rand/v2"
	"sort"
)

// Reverse returns a lazy collection that yields items in reverse order.
func (lc *Collection[T]) Reverse() *Collection[T] {
	return New(func(yield func(T) bool) {
		items := lc.All()

		for i := len(items) - 1; i >= 0; i-- {
			if !yield(items[i]) {
				return
			}
		}
	})
}

// Sort returns a lazy collection sorted by the provided less function.
func (lc *Collection[T]) Sort(less func(a, b T) bool) *Collection[T] {
	return New(func(yield func(T) bool) {
		items := lc.All()

		sort.SliceStable(items, func(i, j int) bool {
			return less(items[i], items[j])
		})

		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	})
}

// SortDesc returns a lazy collection sorted descending by the provided less function.
func (lc *Collection[T]) SortDesc(less func(a, b T) bool) *Collection[T] {
	return lc.Sort(func(a, b T) bool {
		return less(b, a)
	})
}

// SortBy returns a lazy collection sorted by the extracted key.
func SortBy[T any, K cmp.Ordered](lc *Collection[T], keyFunc func(T) K) *Collection[T] {
	return lc.Sort(func(a, b T) bool {
		return keyFunc(a) < keyFunc(b)
	})
}

// SortByDesc returns a lazy collection sorted descending by the extracted key.
func SortByDesc[T any, K cmp.Ordered](lc *Collection[T], keyFunc func(T) K) *Collection[T] {
	return lc.Sort(func(a, b T) bool {
		return keyFunc(a) > keyFunc(b)
	})
}

// Shuffle returns a lazy collection that yields items in random order.
func (lc *Collection[T]) Shuffle() *Collection[T] {
	return New(func(yield func(T) bool) {
		items := lc.All()

		rand.Shuffle(len(items), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})

		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	})
}

// Random returns a lazy collection with up to count randomly selected items.
func (lc *Collection[T]) Random(counts ...int) *Collection[T] {
	count := 1

	if len(counts) > 0 {
		count = counts[0]
	}

	return New(func(yield func(T) bool) {
		items := lc.Shuffle().All()

		if count > len(items) {
			count = len(items)
		}

		for _, item := range items[:count] {
			if !yield(item) {
				return
			}
		}
	})
}
