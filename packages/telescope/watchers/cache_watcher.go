package watchers

import (
	"strings"

	"github.com/bedrock/packages/telescope"
)

// ignoredCachePrefixes contains internal cache key prefixes that should not be

// CacheWatcher monitors cache operations (hit, miss, set, forget) and records
// them as Telescope entries. It mirrors the the underlying behavior class.
//
// Options:
//   - "hidden" ([]string): key names whose values will be masked.
type CacheWatcher struct {
	telescope.BaseWatcher
}

var ignoredCachePrefixes = []string{
	"illuminate:queue:restart",
	"framework/schedule",
	"telescope:",
}

// NewCacheWatcher creates a CacheWatcher with the given options.
func NewCacheWatcher(t *telescope.Telescope, options map[string]any) *CacheWatcher {
	w := &CacheWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

// Register is a no-op for CacheWatcher; callers drive it by calling Hit,
// Missed, Written, and Forgotten directly. When integrated with the bedrock
// cache package, attach listeners to the cache event dispatcher.
func (w *CacheWatcher) Register(_ any) error { return nil }

// CacheWatcher::shouldIgnore().
func (w *CacheWatcher) ShouldIgnore(key string) bool {
	for _, prefix := range ignoredCachePrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}

	return false
}

// maskedValue returns "********" when key is in the hidden list, otherwise v.
func (w *CacheWatcher) maskedValue(key string, v any) any {
	hidden := w.StringsOption("hidden")

	for _, h := range hidden {
		if key == h {
			return "********"
		}
	}

	return v
}

// Hit records a cache-hit entry.
func (w *CacheWatcher) Hit(key string, value any) {
	if w.ShouldIgnore(key) {
		return
	}

	content := map[string]any{
		"type":  "hit",
		"key":   key,
		"value": w.maskedValue(key, value),
	}

	w.Scope().RecordCache(telescope.NewEntry(telescope.EntryTypeCache, content))
}

// Missed records a cache-miss entry.
func (w *CacheWatcher) Missed(key string) {
	if w.ShouldIgnore(key) {
		return
	}

	content := map[string]any{
		"type":  "missed",
		"key":   key,
		"value": nil,
	}

	w.Scope().RecordCache(telescope.NewEntry(telescope.EntryTypeCache, content))
}

// Written records a cache-write entry. expireSeconds is the TTL in seconds
// (0 means no expiry).
func (w *CacheWatcher) Written(key string, value any, expireSeconds int64) {
	if w.ShouldIgnore(key) {
		return
	}

	content := map[string]any{
		"type":    "set",
		"key":     key,
		"value":   w.maskedValue(key, value),
		"expires": expireSeconds,
	}

	w.Scope().RecordCache(telescope.NewEntry(telescope.EntryTypeCache, content))
}

// Forgotten records a cache-forget entry.
func (w *CacheWatcher) Forgotten(key string) {
	if w.ShouldIgnore(key) {
		return
	}

	content := map[string]any{
		"type":  "forget",
		"key":   key,
		"value": nil,
	}

	w.Scope().RecordCache(telescope.NewEntry(telescope.EntryTypeCache, content))
}
