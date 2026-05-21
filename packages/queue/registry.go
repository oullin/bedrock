package queue

import (
	"fmt"
	"sync"
)

// HandlerRegistry is the Go replacement for the upstream queue handler
// resolution — the path that takes a serialised payload, reads its
// `job` / `displayName` field, and asks the IoC container for a matching
// handler. Go has no container; we register handlers by name up front
// and look them up by the same key written into the payload.
//
// A registry entry binds three things together:
//
//   - the display name written into Payload.DisplayName / Payload.Job,
//   - the Handler that should execute when the worker pops a job with
//     that name,
//   - the JobOptions parsed from the registered sample's struct tags
//     (see options.go), so the Manager can stamp max_tries / timeout /
//     backoff / queue / connection onto the payload at push time.
//
// Registries are safe for concurrent Register/Resolve/Names calls.
type HandlerRegistry struct {
	mu      sync.RWMutex
	entries map[string]HandlerRegistryEntry
}

// HandlerRegistryEntry is the value stored under each registered name.
type HandlerRegistryEntry struct {
	// Name is the display name under which the handler is registered.
	// It populates Payload.DisplayName / Payload.Job on the wire.
	Name string
	// Handler is the function that runs when a job with this name is
	// popped off a queue.
	Handler Handler
	// Options are the parsed struct-tag defaults for this job type.
	// They are merged into the dispatch-time options by Manager.Push.
	Options JobOptions
}

// NewHandlerRegistry constructs an empty registry.
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{entries: make(map[string]HandlerRegistryEntry)}
}

// Register binds handler to the given name with the supplied static
// options. An empty name is rejected. Registering the same name twice
// returns an error to protect against silent overrides — call
// Replace explicitly if an override is intentional.
func (r *HandlerRegistry) Register(name string, handler Handler, opts JobOptions) error {
	if name == "" {
		return fmt.Errorf("queue: HandlerRegistry.Register: empty name")
	}

	if handler == nil {
		return fmt.Errorf("queue: HandlerRegistry.Register: nil handler for %q", name)
	}

	r.mu.Lock()

	defer r.mu.Unlock()

	if _, exists := r.entries[name]; exists {
		return fmt.Errorf("queue: HandlerRegistry.Register: %q already registered", name)
	}

	r.entries[name] = HandlerRegistryEntry{Name: name, Handler: handler, Options: opts}

	return nil
}

// RegisterFor is a convenience that parses the struct tags off sample
// and registers the handler under DisplayName(sample). This is the
// typical call path for application code:
//
//	registry.RegisterFor(SendEmail{}, sendEmailHandler)
func (r *HandlerRegistry) RegisterFor(sample any, handler Handler) error {
	name := DisplayName(sample)

	if name == "" {
		return fmt.Errorf("queue: HandlerRegistry.RegisterFor: cannot derive name from nil sample")
	}

	opts, err := ParseJobOptions(sample)

	if err != nil {
		return fmt.Errorf("queue: HandlerRegistry.RegisterFor: %w", err)
	}

	return r.Register(name, handler, opts)
}

// Replace binds handler to name, overwriting any existing entry. Use
// sparingly — Register is the safer default.
func (r *HandlerRegistry) Replace(name string, handler Handler, opts JobOptions) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.entries[name] = HandlerRegistryEntry{Name: name, Handler: handler, Options: opts}
}

// Resolve returns the registered entry for name. The second return is
// false if no entry is registered.
func (r *HandlerRegistry) Resolve(name string) (HandlerRegistryEntry, bool) {
	r.mu.RLock()

	defer r.mu.RUnlock()

	e, ok := r.entries[name]

	return e, ok
}

// Names returns the list of currently registered names. Order is
// unspecified; callers should sort if they care.
func (r *HandlerRegistry) Names() []string {
	r.mu.RLock()

	defer r.mu.RUnlock()

	out := make([]string, 0, len(r.entries))

	for k := range r.entries {
		out = append(out, k)
	}

	return out
}

// Forget removes the entry for name. It is a no-op if name is not
// registered.
func (r *HandlerRegistry) Forget(name string) {
	r.mu.Lock()

	defer r.mu.Unlock()

	delete(r.entries, name)
}
