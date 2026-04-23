package watchers

import (
	"fmt"

	"github.com/bedrock/packages/debugbar"
)

// GateResult constants mirror the two authorization outcomes recorded by
// Upstream's GateWatcher.

// GateWatcher monitors authorization gate checks and records them as DebugBar
// entries. It mirrors Upstream's GateWatcher class.
//
// Options:
//   - "ignore_abilities" ([]string): gate ability names to skip.
//   - "ignore_packages"  ([]string): packages whose abilities are skipped.
type GateWatcher struct {
	debugbar.BaseWatcher
}

const (
	GateResultAllowed = "allowed"
	GateResultDenied  = "denied"
)

// NewGateWatcher creates a GateWatcher with the given options.
func NewGateWatcher(t *debugbar.DebugBar, options map[string]any) *GateWatcher {
	w := &GateWatcher{}
	w.SetDebugBar(t)
	w.Options = options

	return w
}

// Register is a no-op for GateWatcher; callers drive it via Record.
func (w *GateWatcher) Register(_ any) error { return nil }

// ShouldIgnore reports whether the gate ability check should be skipped.
func (w *GateWatcher) ShouldIgnore(ability string) bool {
	for _, name := range w.StringsOption("ignore_abilities") {
		if name == ability {
			return true
		}
	}

	return false
}

// Record records a gate authorization check entry. ability is the gate
// ability name, allowed reports the decision, user is the authorised subject
// (optional), and arguments are the gate arguments (e.g. model instance).
func (w *GateWatcher) Record(ability string, allowed bool, user any, arguments []any) {
	w.RecordWithMessage(ability, allowed, "", user, arguments)
}

// RecordWithMessage records a gate authorization check with an optional response message.
func (w *GateWatcher) RecordWithMessage(ability string, allowed bool, message string, user any, arguments []any) {
	if w.ShouldIgnore(ability) {
		return
	}

	result := GateResultAllowed

	if !allowed {
		result = GateResultDenied
	}

	formattedArgs := make([]string, 0, len(arguments))

	for _, arg := range arguments {
		formattedArgs = append(formattedArgs, fmt.Sprintf("%v", arg))
	}

	content := map[string]any{
		"ability":   ability,
		"result":    result,
		"arguments": formattedArgs,
	}

	if user != nil {
		content["user"] = fmt.Sprintf("%v", user)
	}

	if message != "" {
		content["message"] = message
	}

	entry := debugbar.NewEntry(debugbar.EntryTypeGate, content)
	entry.AddTags(result)

	w.Scope().RecordGate(entry)
}
