// Package registry holds multiple named workflows and resolves the right one
// for a given subject via SupportStrategy predicates.
package registry

import (
	"fmt"
	"sync"

	"github.com/bedrock/packages/workflow"
)

// SupportStrategy reports whether a workflow applies to a subject.
type SupportStrategy[T any] func(T) bool

// Entry registers a workflow with an optional name and support predicate.
type Entry[T any] struct {
	Name     string
	Workflow workflow.Engine[T]
	Supports SupportStrategy[T]
}

// Registry is a thread-safe set of workflow Entries.
type Registry[T any] struct {
	mu      sync.RWMutex
	entries []Entry[T]
}

func New[T any]() *Registry[T] {
	return &Registry[T]{}
}

// Add registers a workflow entry. Nil-Workflow entries are ignored.
func (r *Registry[T]) Add(entry Entry[T]) {
	if entry.Workflow == nil {
		return
	}

	r.mu.Lock()

	defer r.mu.Unlock()

	r.entries = append(r.entries, entry)
}

// Get returns the first workflow whose name matches (if provided) and whose
// SupportStrategy accepts the subject.
func (r *Registry[T]) Get(subject T, name string) (workflow.Engine[T], error) {
	r.mu.RLock()

	defer r.mu.RUnlock()

	for _, entry := range r.entries {
		if name != "" && entry.Name != name {
			continue
		}

		if entry.Supports != nil && !entry.Supports(subject) {
			continue
		}

		return entry.Workflow, nil
	}

	if name != "" {
		return nil, fmt.Errorf("workflow %q not found", name)
	}

	return nil, fmt.Errorf("no workflow matched subject")
}
