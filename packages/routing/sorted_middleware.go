package routing

import (
	"reflect"
	"sort"
)

// MiddlewarePriority manages priority-based sorting of middleware functions.
// Lower priority values execute first.
type MiddlewarePriority struct {
	priority map[uintptr]int
}

// NewMiddlewarePriority creates a new middleware priority manager.
func NewMiddlewarePriority() *MiddlewarePriority {
	return &MiddlewarePriority{
		priority: make(map[uintptr]int),
	}
}

// Set assigns a priority to a middleware function. Lower values run first.
func (mp *MiddlewarePriority) Set(mw MiddlewareFunc, priority int) {
	ptr := reflect.ValueOf(mw).Pointer()
	mp.priority[ptr] = priority
}

// Sort returns a new slice of middleware sorted by priority. Middleware without
// an assigned priority retains its relative position after all prioritized
// middleware.
func (mp *MiddlewarePriority) Sort(middleware []MiddlewareFunc) []MiddlewareFunc {
	type entry struct {
		mw       MiddlewareFunc
		priority int
		index    int
		hasPrio  bool
	}

	entries := make([]entry, len(middleware))

	for i, mw := range middleware {
		ptr := reflect.ValueOf(mw).Pointer()
		prio, ok := mp.priority[ptr]
		entries[i] = entry{mw: mw, priority: prio, index: i, hasPrio: ok}
	}

	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].hasPrio && entries[j].hasPrio {
			return entries[i].priority < entries[j].priority
		}

		if entries[i].hasPrio && !entries[j].hasPrio {
			return true
		}

		if !entries[i].hasPrio && entries[j].hasPrio {
			return false
		}

		return entries[i].index < entries[j].index
	})

	result := make([]MiddlewareFunc, len(entries))

	for i, e := range entries {
		result[i] = e.mw
	}

	return result
}
