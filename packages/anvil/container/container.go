package container

import (
	"fmt"
	"sync"
)

// Container is an IoC service container that manages bindings and resolves
// dependencies. It is safe for concurrent use.
type Container struct {
	mu                       sync.RWMutex
	bindings                 map[string]*binding
	instances                map[string]any
	aliases                  map[string]string
	tags                     map[string][]string
	contextual               map[string]map[string]any
	resolved                 map[string]bool
	extenders                map[string][]func(any, *Container) any
	reboundCallbacks         map[string][]func(any)
	beforeResolvingCallbacks []func(string, *Container)
	resolvingCallbacks       []func(string, any)
	afterResolvingCallbacks  []func(string, any)
}

// New creates an empty container.
func New() *Container {
	return &Container{
		bindings:         make(map[string]*binding),
		instances:        make(map[string]any),
		aliases:          make(map[string]string),
		tags:             make(map[string][]string),
		contextual:       make(map[string]map[string]any),
		resolved:         make(map[string]bool),
		extenders:        make(map[string][]func(any, *Container) any),
		reboundCallbacks: make(map[string][]func(any)),
	}
}

// BindIf registers a transient binding only if abstract is not already bound.
func (c *Container) BindIf(abstract string, factory Factory) {
	if !c.Bound(abstract) {
		c.Bind(abstract, factory)
	}
}

// Bind registers a transient binding. Each call to Make returns a new value.
func (c *Container) Bind(abstract string, factory Factory) {
	c.mu.Lock()

	delete(c.instances, abstract)
	c.bindings[abstract] = &binding{factory: factory, shared: false}
	callbacks := c.reboundCallbacks[abstract]

	c.mu.Unlock()

	if len(callbacks) > 0 && c.Resolved(abstract) {
		value, err := c.Make(abstract)
		if err == nil {
			for _, cb := range callbacks {
				cb(value)
			}
		}
	}
}

// SingletonIf registers a shared binding only if abstract is not already bound.
func (c *Container) SingletonIf(abstract string, factory Factory) {
	if !c.Bound(abstract) {
		c.Singleton(abstract, factory)
	}
}

// Singleton registers a shared binding. The factory executes once; subsequent
// calls to Make return the cached result.
func (c *Container) Singleton(abstract string, factory Factory) {
	c.mu.Lock()

	delete(c.instances, abstract)
	c.bindings[abstract] = &binding{factory: factory, shared: true}
	callbacks := c.reboundCallbacks[abstract]

	c.mu.Unlock()

	if len(callbacks) > 0 && c.Resolved(abstract) {
		value, err := c.Make(abstract)
		if err == nil {
			for _, cb := range callbacks {
				cb(value)
			}
		}
	}
}

// Instance binds a pre-built value. Make returns this exact value.
func (c *Container) Instance(abstract string, value any) {
	c.mu.Lock()

	delete(c.bindings, abstract)
	_, wasBound := c.instances[abstract]
	c.instances[abstract] = value
	callbacks := c.reboundCallbacks[abstract]
	wasResolved := c.resolved[abstract] || wasBound

	c.mu.Unlock()

	if len(callbacks) > 0 && wasResolved {
		for _, cb := range callbacks {
			cb(value)
		}
	}
}

// Scoped registers a shared binding that is reset when ForgetScopedInstances
// is called. This is useful for request-scoped services.
func (c *Container) Scoped(abstract string, factory Factory) {
	c.mu.Lock()

	delete(c.instances, abstract)
	c.bindings[abstract] = &binding{factory: factory, shared: true, scoped: true}
	callbacks := c.reboundCallbacks[abstract]

	c.mu.Unlock()

	if len(callbacks) > 0 && c.Resolved(abstract) {
		value, err := c.Make(abstract)
		if err == nil {
			for _, cb := range callbacks {
				cb(value)
			}
		}
	}
}

// ScopedIf registers a scoped binding only if abstract is not already bound.
func (c *Container) ScopedIf(abstract string, factory Factory) {
	if !c.Bound(abstract) {
		c.Scoped(abstract, factory)
	}
}

// ForgetScopedInstances resets all scoped bindings so they are re-resolved
// on the next Make call.
func (c *Container) ForgetScopedInstances() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for abstract, b := range c.bindings {
		if b.scoped && b.instance != nil {
			c.bindings[abstract] = &binding{
				factory: b.factory,
				shared:  true,
				scoped:  true,
			}
		}
	}
}

// Bound reports whether abstract has a binding, instance, or alias.
func (c *Container) Bound(abstract string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	abstract = c.resolveAlias(abstract)

	if _, ok := c.bindings[abstract]; ok {
		return true
	}

	if _, ok := c.instances[abstract]; ok {
		return true
	}

	return false
}

// Make resolves abstract from the container. It checks instances first, then
// bindings. Shared bindings are resolved once and cached.
func (c *Container) Make(abstract string) (any, error) {
	c.mu.RLock()
	abstract = c.resolveAlias(abstract)

	if value, ok := c.instances[abstract]; ok {
		c.mu.RUnlock()
		c.fireResolvingCallbacks(abstract, value)
		return value, nil
	}

	b, ok := c.bindings[abstract]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrNotBound, abstract)
	}

	c.fireBeforeResolvingCallbacks(abstract)

	value, err := b.resolve(c)
	if err != nil {
		return nil, fmt.Errorf("%w: resolving %q: %v", ErrResolve, abstract, err)
	}

	value = c.applyExtenders(abstract, value)

	c.mu.Lock()
	c.resolved[abstract] = true
	c.mu.Unlock()

	c.fireResolvingCallbacks(abstract, value)

	return value, nil
}

// MustMake resolves abstract or panics.
func (c *Container) MustMake(abstract string) any {
	value, err := c.Make(abstract)
	if err != nil {
		panic(err)
	}

	return value
}

// Alias creates an alternative name for an abstract. It panics if abstract
// and alias are the same string.
func (c *Container) Alias(abstract, alias string) {
	if abstract == alias {
		panic(fmt.Sprintf("container: %q is aliased to itself", abstract))
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.aliases[alias] = abstract
}

// IsAlias reports whether name is a registered alias.
func (c *Container) IsAlias(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.aliases[name]
	return ok
}

// GetAlias returns the canonical abstract for name. If name is not an alias,
// it returns name unchanged.
func (c *Container) GetAlias(name string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.resolveAlias(name)
}

// Has reports whether abstract is registered (binding, instance, or alias).
// This is an alias for Bound.
func (c *Container) Has(abstract string) bool {
	return c.Bound(abstract)
}

// FactoryFunc returns a closure that resolves abstract from the container each
// time it is called.
func (c *Container) FactoryFunc(abstract string) func() (any, error) {
	return func() (any, error) {
		return c.Make(abstract)
	}
}

// Tag assigns a tag to a group of abstracts.
func (c *Container) Tag(abstracts []string, tag string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.tags[tag] = append(c.tags[tag], abstracts...)
}

// Tagged resolves all abstracts registered under tag.
func (c *Container) Tagged(tag string) ([]any, error) {
	c.mu.RLock()
	abstracts := make([]string, len(c.tags[tag]))
	copy(abstracts, c.tags[tag])
	c.mu.RUnlock()

	results := make([]any, 0, len(abstracts))
	for _, abstract := range abstracts {
		value, err := c.Make(abstract)
		if err != nil {
			return nil, err
		}
		results = append(results, value)
	}

	return results, nil
}

// Resolved reports whether abstract has been resolved at least once.
func (c *Container) Resolved(abstract string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	abstract = c.resolveAlias(abstract)
	return c.resolved[abstract]
}

// Flush resets the container, clearing all bindings, instances, aliases,
// tags, and resolved state.
func (c *Container) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bindings = make(map[string]*binding)
	c.instances = make(map[string]any)
	c.aliases = make(map[string]string)
	c.tags = make(map[string][]string)
	c.contextual = make(map[string]map[string]any)
	c.resolved = make(map[string]bool)
	c.extenders = make(map[string][]func(any, *Container) any)
	c.reboundCallbacks = make(map[string][]func(any))
	c.beforeResolvingCallbacks = nil
	c.resolvingCallbacks = nil
	c.afterResolvingCallbacks = nil
}

// ForgetInstance removes a single resolved instance.
func (c *Container) ForgetInstance(abstract string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.instances, abstract)
}

// ForgetInstances removes all resolved instances.
func (c *Container) ForgetInstances() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.instances = make(map[string]any)
}

// Extend registers a decorator for abstract. Each time abstract is resolved,
// the extender receives the resolved value and may return a replacement.
// Multiple extenders are applied in registration order.
func (c *Container) Extend(abstract string, extender func(any, *Container) any) {
	c.mu.Lock()

	abstract = c.resolveAlias(abstract)

	if value, ok := c.instances[abstract]; ok {
		c.instances[abstract] = extender(value, c)
		callbacks := c.reboundCallbacks[abstract]
		newValue := c.instances[abstract]

		c.mu.Unlock()

		for _, cb := range callbacks {
			cb(newValue)
		}

		return
	}

	c.extenders[abstract] = append(c.extenders[abstract], extender)
	wasResolved := c.resolved[abstract]

	c.mu.Unlock()

	if wasResolved {
		value, err := c.Make(abstract)
		if err == nil {
			c.mu.RLock()
			callbacks := make([]func(any), len(c.reboundCallbacks[abstract]))
			copy(callbacks, c.reboundCallbacks[abstract])
			c.mu.RUnlock()

			for _, cb := range callbacks {
				cb(value)
			}
		}
	}
}

// ForgetExtenders removes all extenders for abstract.
func (c *Container) ForgetExtenders(abstract string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	abstract = c.resolveAlias(abstract)
	delete(c.extenders, abstract)
}

// Rebinding registers a callback that fires when abstract is re-bound.
func (c *Container) Rebinding(abstract string, callback func(any)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	abstract = c.resolveAlias(abstract)
	c.reboundCallbacks[abstract] = append(c.reboundCallbacks[abstract], callback)
}

// BeforeResolving registers a callback that fires before each resolution begins.
func (c *Container) BeforeResolving(callback func(string, *Container)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.beforeResolvingCallbacks = append(c.beforeResolvingCallbacks, callback)
}

// Resolving registers a callback that fires each time any abstract is resolved.
func (c *Container) Resolving(callback func(string, any)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.resolvingCallbacks = append(c.resolvingCallbacks, callback)
}

// AfterResolving registers a callback that fires after each resolution.
func (c *Container) AfterResolving(callback func(string, any)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.afterResolvingCallbacks = append(c.afterResolvingCallbacks, callback)
}

// resolveAlias follows the alias chain to the canonical abstract. The caller
// must hold at least a read lock.
func (c *Container) resolveAlias(name string) string {
	visited := map[string]bool{}

	for {
		target, ok := c.aliases[name]
		if !ok {
			return name
		}

		if visited[name] {
			return name
		}

		visited[name] = true
		name = target
	}
}

func (c *Container) fireBeforeResolvingCallbacks(abstract string) {
	c.mu.RLock()
	callbacks := make([]func(string, *Container), len(c.beforeResolvingCallbacks))
	copy(callbacks, c.beforeResolvingCallbacks)
	c.mu.RUnlock()

	for _, cb := range callbacks {
		cb(abstract, c)
	}
}

func (c *Container) applyExtenders(abstract string, value any) any {
	c.mu.RLock()
	exts := make([]func(any, *Container) any, len(c.extenders[abstract]))
	copy(exts, c.extenders[abstract])
	c.mu.RUnlock()

	for _, ext := range exts {
		value = ext(value, c)
	}

	return value
}

func (c *Container) fireResolvingCallbacks(abstract string, value any) {
	c.mu.RLock()
	resolving := make([]func(string, any), len(c.resolvingCallbacks))
	copy(resolving, c.resolvingCallbacks)
	after := make([]func(string, any), len(c.afterResolvingCallbacks))
	copy(after, c.afterResolvingCallbacks)
	c.mu.RUnlock()

	for _, cb := range resolving {
		cb(abstract, value)
	}

	for _, cb := range after {
		cb(abstract, value)
	}
}
