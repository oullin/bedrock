// Package rules contains the built-in validation rule functions and the
// global rule registry.
package rules

import "sync"

// RuleFunc is the signature every built-in (and user-registered) rule must
// satisfy.  It returns true when the value is valid.
//
// Parameters:
//   - attribute: the dot-notation field name being validated
//   - value:     the field value (may be nil when the field is absent)
//   - params:    the colon-separated parameters, e.g. ["255"] for max:255
//   - ctx:       a RuleContext that provides access to the full data set,
//     other fields, and the ability to add custom failure messages
type RuleFunc func(attribute string, value any, params []string, ctx RuleContext) bool

// RuleContext is passed to every RuleFunc and gives it access to the
// surrounding validation state without importing the validator package
// (which would create a circular dependency).
type RuleContext interface {
	// GetValue returns the value for another attribute in the data set.
	GetValue(attribute string) any

	// GetData returns the full data set as a flat dot-notation map.
	GetData() map[string]any

	// GetOriginalData returns the full original (nested) data set.
	GetOriginalData() map[string]any

	// IsSometimes reports whether the "sometimes" modifier is active for attr.
	IsSometimes(attribute string) bool

	// IsPresent reports whether the attribute key exists in the data.
	IsPresent(attribute string) bool

	// SetMessage overrides the failure message for the current rule invocation.
	SetMessage(msg string)

	// GetPresenceVerifier returns the database presence verifier, or nil.
	GetPresenceVerifier() PresenceVerifier
}

// PresenceVerifier checks database-level constraints (exists, unique).
type PresenceVerifier interface {
	GetCount(table, column, value string, excludeID *string, idColumn *string, extras map[string]any) int
	GetMultiCount(table, column string, values []string, extras map[string]any) int
}

// registry is the global rule store.
type registry struct {
	mu       sync.RWMutex
	rules    map[string]RuleFunc
	implicit map[string]bool // rules that run even when the field is absent
}

var global = &registry{
	rules:    make(map[string]RuleFunc),
	implicit: make(map[string]bool),
}

// Register adds a rule under name.  Panics if name is empty.
func Register(name string, fn RuleFunc) {
	if name == "" {
		panic("validation/rules: rule name must not be empty")
	}

	global.mu.Lock()

	defer global.mu.Unlock()

	global.rules[name] = fn
}

// RegisterImplicit registers an implicit rule — one that runs even when the
// field is absent or its value is empty.
func RegisterImplicit(name string, fn RuleFunc) {
	Register(name, fn)

	global.mu.Lock()

	defer global.mu.Unlock()

	global.implicit[name] = true
}

// Lookup returns the RuleFunc for name, or (nil, false) if not found.
func Lookup(name string) (RuleFunc, bool) {
	global.mu.RLock()

	defer global.mu.RUnlock()

	fn, ok := global.rules[name]

	return fn, ok
}

// IsImplicit reports whether name is an implicit rule.
func IsImplicit(name string) bool {
	global.mu.RLock()

	defer global.mu.RUnlock()

	return global.implicit[name]
}

// All returns a copy of all registered rules (name → fn).
func All() map[string]RuleFunc {
	global.mu.RLock()

	defer global.mu.RUnlock()

	out := make(map[string]RuleFunc, len(global.rules))

	for k, v := range global.rules {
		out[k] = v
	}

	return out
}

func init() {
	registerAll()
}
