package support

import "strings"

// NamespacedItemResolver parses strings like "namespace::group.item".
type NamespacedItemResolver struct {
	cache map[string][3]string
}

// NewNamespacedItemResolver creates a resolver with an empty parse cache.
func NewNamespacedItemResolver() *NamespacedItemResolver {
	return &NamespacedItemResolver{cache: make(map[string][3]string)}
}

// Parse returns namespace, group, and item components.
func (r *NamespacedItemResolver) Parse(key string) (namespace, group, item string) {
	if cached, ok := r.cache[key]; ok {
		return cached[0], cached[1], cached[2]
	}

	target := key
	if before, after, ok := strings.Cut(key, "::"); ok {
		namespace = before
		target = after
	}

	group = target
	if before, after, ok := strings.Cut(target, "."); ok {
		group = before
		item = after
	}

	r.cache[key] = [3]string{namespace, group, item}

	return namespace, group, item
}

// Flush clears the parse cache.
func (r *NamespacedItemResolver) Flush() {
	r.cache = make(map[string][3]string)
}

// CacheSize returns the number of cached parses.
func (r *NamespacedItemResolver) CacheSize() int {
	return len(r.cache)
}
