package routing

import "strings"

// SortedMiddleware is the deduplicated, priority-ordered middleware list.
//
// In the upstream framework this is a Collection subclass; in Go it is a plain slice with
// a constructor that performs the sort. Slice elements are typed as any
// because middleware may be either a string class name (with optional ":args"
// suffix) or a closure.
//
// Ref: @bedrock/code-0345
type SortedMiddleware []any

// NewSortedMiddleware returns the middleware list reordered so that any
// middleware appearing in priorityMap respects its relative position there,
// then deduplicated with the same logic as Router::uniqueMiddleware.
func NewSortedMiddleware(priorityMap []string, middleware []any) SortedMiddleware {
	return uniqueMiddleware(sortMiddleware(priorityMap, middleware))
}

func sortMiddleware(priorityMap []string, middleware []any) []any {
	work := make([]any, len(middleware))
	copy(work, middleware)

	for {
		lastIndex := -1
		lastPriorityIndex := -1
		moved := false

		for i, m := range work {
			s, ok := m.(string)

			if !ok {
				continue
			}

			pi := priorityMapIndex(priorityMap, s)

			if pi < 0 {
				continue
			}

			if lastPriorityIndex >= 0 && pi < lastPriorityIndex {
				work = moveMiddleware(work, i, lastIndex)
				moved = true

				break
			}

			lastIndex = i
			lastPriorityIndex = pi
		}

		if !moved {
			return work
		}
	}
}

func priorityMapIndex(priorityMap []string, middleware string) int {
	stripped := middleware

	if idx := strings.Index(middleware, ":"); idx >= 0 {
		stripped = middleware[:idx]
	}

	for i, name := range priorityMap {
		if name == stripped {
			return i
		}
	}

	return -1
}

func moveMiddleware(in []any, from, to int) []any {
	if from == to {
		return in
	}

	if to < 0 {
		to = 0
	}

	out := make([]any, 0, len(in))
	out = append(out, in[:to]...)
	out = append(out, in[from])
	out = append(out, in[to:from]...)
	out = append(out, in[from+1:]...)

	return out
}

// occurrence and discards later duplicates. Only string entries are eligible
// for dedup; closures are kept verbatim.
func uniqueMiddleware(middleware []any) SortedMiddleware {
	seen := map[string]struct{}{}
	out := make(SortedMiddleware, 0, len(middleware))

	for _, m := range middleware {
		if s, ok := m.(string); ok {
			if _, dup := seen[s]; dup {
				continue
			}

			seen[s] = struct{}{}
		}

		out = append(out, m)
	}

	return out
}
