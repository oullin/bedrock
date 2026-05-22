package container

import (
	"fmt"
	"slices"
	"sync"
)

// Factory creates a value from the container.
type Factory func(c *Container) (any, error)

// BindingCallback is called during resolving lifecycle events.
type BindingCallback func(instance any, c *Container)

// BeforeResolvingCallback is called before resolution begins.
type BeforeResolvingCallback func(abstract string, parameters map[string]any, c *Container)

// ExtenderFunc modifies a resolved instance.
type ExtenderFunc func(instance any, c *Container) (any, error)

// MethodCallable represents a function that can be called with dependency
// injection via the container.
type MethodCallable func(c *Container, params map[string]any) (any, error)

// Binding holds a registered factory and its lifecycle metadata.
type Binding struct {
	factory Factory
	shared  bool
	scoped  bool
}

// Container is a inversion-of-control container. It manages
// service bindings, resolution, contextual bindings, tagging, extension,
// lifecycle callbacks, and method invocation. All methods are safe for
// concurrent use.
type Container struct {
	mu sync.RWMutex

	bindings        map[string]Binding
	instances       map[string]any
	aliases         map[string]string
	abstractAliases map[string][]string
	resolved        map[string]bool
	buildStack      []string
	with            []map[string]any
	contextual      map[string]map[string]any
	tags            map[string][]string
	extenders       map[string][]ExtenderFunc
	reboundCbs      map[string][]BindingCallback
	methodBindings  map[string]MethodCallable

	globalBeforeCbs []BeforeResolvingCallback
	beforeCbs       map[string][]BeforeResolvingCallback
	globalResolvCbs []BindingCallback
	resolvCbs       map[string][]BindingCallback
	globalAfterCbs  []BindingCallback
	afterCbs        map[string][]BindingCallback
}

// New creates an empty, fully initialized Container.
func New() *Container {
	return &Container{
		bindings:        make(map[string]Binding),
		instances:       make(map[string]any),
		aliases:         make(map[string]string),
		abstractAliases: make(map[string][]string),
		resolved:        make(map[string]bool),
		contextual:      make(map[string]map[string]any),
		tags:            make(map[string][]string),
		extenders:       make(map[string][]ExtenderFunc),
		reboundCbs:      make(map[string][]BindingCallback),
		methodBindings:  make(map[string]MethodCallable),
		beforeCbs:       make(map[string][]BeforeResolvingCallback),
		resolvCbs:       make(map[string][]BindingCallback),
		afterCbs:        make(map[string][]BindingCallback),
	}
}

// ---------- Binding ----------

// Bind registers a factory for the given abstract. When shared is true the
// factory is only called once and the result is cached.
func (c *Container) Bind(abstract string, factory Factory, shared bool) {
	c.mu.Lock()

	c.dropStale(abstract)

	c.bindings[abstract] = Binding{factory: factory, shared: shared}

	wasResolved := c.resolved[abstract]

	c.mu.Unlock()

	if wasResolved {
		c.rebound(abstract)
	}
}

// BindIf registers a binding only if the abstract is not already bound.
func (c *Container) BindIf(abstract string, factory Factory, shared bool) {
	if !c.Bound(abstract) {
		c.Bind(abstract, factory, shared)
	}
}

// Singleton registers a shared binding. The factory is called at most once.
func (c *Container) Singleton(abstract string, factory Factory) {
	c.Bind(abstract, factory, true)
}

// SingletonIf registers a singleton only if the abstract is not already bound.
func (c *Container) SingletonIf(abstract string, factory Factory) {
	if !c.Bound(abstract) {
		c.Singleton(abstract, factory)
	}
}

// Scoped registers a scoped binding. Scoped bindings behave like singletons
// but can be flushed independently via ForgetScopedInstances.
func (c *Container) Scoped(abstract string, factory Factory) {
	c.mu.Lock()

	c.dropStale(abstract)

	c.bindings[abstract] = Binding{factory: factory, shared: true, scoped: true}

	wasResolved := c.resolved[abstract]

	c.mu.Unlock()

	if wasResolved {
		c.rebound(abstract)
	}
}

// ScopedIf registers a scoped binding only if the abstract is not already bound.
func (c *Container) ScopedIf(abstract string, factory Factory) {
	if !c.Bound(abstract) {
		c.Scoped(abstract, factory)
	}
}

// Instance registers a pre-existing value in the container. Returns the
// instance for convenience.
func (c *Container) Instance(abstract string, instance any) any {
	c.mu.Lock()

	c.removeAlias(abstract)

	wasBound := c.isBound(abstract)
	c.instances[abstract] = instance

	c.mu.Unlock()

	if wasBound {
		c.rebound(abstract)
	}

	return instance
}

// ---------- Resolution ----------

// Make resolves the given abstract from the container.
func (c *Container) Make(abstract string) (any, error) {
	return c.resolve(abstract, nil)
}

// MakeWith resolves the given abstract, passing parameters to the factory via
// the parameter override stack.
func (c *Container) MakeWith(abstract string, parameters map[string]any) (any, error) {
	return c.resolve(abstract, parameters)
}

// Build executes the given factory directly. It is useful for one-off
// instantiation without registering a binding.
func (c *Container) Build(factory Factory) (any, error) {
	return factory(c)
}

// Get resolves the abstract. If the abstract is not bound, it returns
// ErrNotBound (PSR-11 parity).
func (c *Container) Get(abstract string) (any, error) {
	if c.Has(abstract) {
		return c.resolve(abstract, nil)
	}

	return nil, fmt.Errorf("%w: %q", ErrNotBound, abstract)
}

// FactoryFunc returns a closure that resolves the abstract each time it is
// called.
func (c *Container) FactoryFunc(abstract string) func() (any, error) {
	return func() (any, error) {
		return c.Make(abstract)
	}
}

// resolve is the core resolution engine.
func (c *Container) resolve(abstract string, parameters map[string]any) (any, error) {
	c.mu.Lock()

	original := abstract
	abstract = c.getAlias(abstract)

	// Fire before-resolving callbacks.
	beforeGlobal := slices.Clone(c.globalBeforeCbs)
	beforeSpecific := slices.Clone(c.beforeCbs[abstract])

	// Check contextual binding first (takes precedence over cached instances).
	concrete := c.getContextualConcrete(abstract)

	if concrete == nil && original != abstract {
		concrete = c.getContextualConcrete(original)
	}

	// Check for cached instance when no parameters are given and no
	// contextual override exists.
	if concrete == nil && len(parameters) == 0 {
		if inst, ok := c.instances[abstract]; ok {
			c.mu.Unlock()

			fireBeforeCallbacks(beforeGlobal, abstract, parameters, c)
			fireBeforeCallbacks(beforeSpecific, abstract, parameters, c)

			return inst, nil
		}
	}

	// Determine the concrete factory.
	var factory Factory

	if concrete != nil {
		if f, ok := concrete.(Factory); ok {
			factory = f
		} else {
			val := concrete
			factory = func(_ *Container) (any, error) { return val, nil }
		}
	} else if b, ok := c.bindings[abstract]; ok {
		factory = b.factory
	}

	if factory == nil {
		c.mu.Unlock()

		return nil, fmt.Errorf("%w: %q", ErrNotBound, abstract)
	}

	// Circular dependency detection.
	if slices.Contains(c.buildStack, abstract) {
		// Snapshot the stack while still holding the lock; otherwise
		// the deferred fmt.Errorf read would race with concurrent
		// resolve() calls mutating c.buildStack at lines 265/283.
		stackSnapshot := slices.Clone(c.buildStack)
		c.mu.Unlock()

		return nil, fmt.Errorf("%w: %q (build stack: %v)", ErrCircularDependency, abstract, stackSnapshot)
	}

	// Capture binding metadata before unlocking.
	b := c.bindings[abstract]
	extenders := slices.Clone(c.extenders[abstract])
	resolvGlobal := slices.Clone(c.globalResolvCbs)
	resolvSpecific := slices.Clone(c.resolvCbs[abstract])
	afterGlobal := slices.Clone(c.globalAfterCbs)
	afterSpecific := slices.Clone(c.afterCbs[abstract])

	c.buildStack = append(c.buildStack, abstract)

	// Always push parameters (even nil) so nested Make calls get their own
	// scope and don't inherit the parent's parameters.
	c.with = append(c.with, parameters)

	c.mu.Unlock()

	fireBeforeCallbacks(beforeGlobal, abstract, parameters, c)
	fireBeforeCallbacks(beforeSpecific, abstract, parameters, c)

	// Execute factory.
	instance, err := factory(c)

	c.mu.Lock()

	// Pop build stack.
	if len(c.buildStack) > 0 {
		c.buildStack = c.buildStack[:len(c.buildStack)-1]
	}

	if len(c.with) > 0 {
		c.with = c.with[:len(c.with)-1]
	}

	c.mu.Unlock()

	if err != nil {
		return nil, err
	}

	// Apply extenders.
	for _, ext := range extenders {
		instance, err = ext(instance, c)

		if err != nil {
			return nil, err
		}
	}

	// Cache shared instances.
	if b.shared && len(parameters) == 0 {
		c.mu.Lock()
		c.instances[abstract] = instance
		c.mu.Unlock()
	}

	c.mu.Lock()
	c.resolved[abstract] = true
	c.mu.Unlock()

	// Fire resolving callbacks.
	fireCallbacks(resolvGlobal, instance, c)
	fireCallbacks(resolvSpecific, instance, c)

	// Fire after-resolving callbacks.
	fireCallbacks(afterGlobal, instance, c)
	fireCallbacks(afterSpecific, instance, c)

	return instance, nil
}

// Parameters returns the current parameter override map from the top of the
// with stack. Factories can call this to access parameters passed via MakeWith.
func (c *Container) Parameters() map[string]any {
	c.mu.RLock()

	defer c.mu.RUnlock()

	if len(c.with) == 0 {
		return nil
	}

	return c.with[len(c.with)-1]
}

// ---------- Aliases ----------

// Alias creates an alias that resolves to the given abstract.
func (c *Container) Alias(abstract, alias string) {
	if abstract == alias {
		panic(fmt.Sprintf("%s: %q", ErrSelfAlias.Error(), alias))
	}

	c.mu.Lock()

	defer c.mu.Unlock()

	c.aliases[alias] = abstract
	c.abstractAliases[abstract] = append(c.abstractAliases[abstract], alias)
}

// GetAlias resolves an alias chain to the actual abstract name.
func (c *Container) GetAlias(abstract string) string {
	c.mu.RLock()

	defer c.mu.RUnlock()

	return c.getAlias(abstract)
}

// IsAlias reports whether the given name is a registered alias.
func (c *Container) IsAlias(name string) bool {
	c.mu.RLock()

	defer c.mu.RUnlock()

	_, ok := c.aliases[name]

	return ok
}

// getAlias resolves the full alias chain without locking. Caller must hold the
// lock.
func (c *Container) getAlias(abstract string) string {
	for {
		target, ok := c.aliases[abstract]

		if !ok {
			return abstract
		}

		abstract = target
	}
}

// ---------- Queries ----------

// Bound reports whether an abstract has a binding, instance, or alias.
func (c *Container) Bound(abstract string) bool {
	c.mu.RLock()

	defer c.mu.RUnlock()

	return c.isBound(abstract)
}

// isBound is the lock-free version of Bound. Caller must hold the lock.
func (c *Container) isBound(abstract string) bool {
	_, hasBind := c.bindings[abstract]
	_, hasInst := c.instances[abstract]
	_, hasAlias := c.aliases[abstract]

	return hasBind || hasInst || hasAlias
}

// Has is an alias for Bound (PSR-11 parity).
func (c *Container) Has(abstract string) bool {
	return c.Bound(abstract)
}

// Resolved reports whether the given abstract has been resolved at least once.
func (c *Container) Resolved(abstract string) bool {
	c.mu.RLock()

	defer c.mu.RUnlock()

	abs := c.getAlias(abstract)

	if _, ok := c.instances[abs]; ok {
		return true
	}

	return c.resolved[abs]
}

// IsShared reports whether the given abstract is a singleton or scoped binding.
func (c *Container) IsShared(abstract string) bool {
	c.mu.RLock()

	defer c.mu.RUnlock()

	abs := c.getAlias(abstract)

	if _, ok := c.instances[abs]; ok {
		return true
	}

	if b, ok := c.bindings[abs]; ok {
		return b.shared
	}

	return false
}

// CurrentlyResolving returns the abstract at the top of the build stack, or an
// empty string if nothing is being resolved.
func (c *Container) CurrentlyResolving() string {
	c.mu.RLock()

	defer c.mu.RUnlock()

	if len(c.buildStack) == 0 {
		return ""
	}

	return c.buildStack[len(c.buildStack)-1]
}

// ---------- Tagging ----------

// Tag assigns one or more tags to the given abstracts.
func (c *Container) Tag(abstracts []string, tags ...string) {
	c.mu.Lock()

	defer c.mu.Unlock()

	for _, tag := range tags {
		c.tags[tag] = append(c.tags[tag], abstracts...)
	}
}

// Tagged resolves all abstracts registered under the given tag.
func (c *Container) Tagged(tag string) []any {
	c.mu.RLock()
	abstracts := slices.Clone(c.tags[tag])
	c.mu.RUnlock()

	var results []any

	for _, abstract := range abstracts {
		instance, err := c.Make(abstract)

		if err == nil {
			results = append(results, instance)
		}
	}

	return results
}

// ---------- Extension ----------

// Extend registers a callback that modifies a resolved instance. If the
// abstract already has a cached instance, the extender is applied immediately.
func (c *Container) Extend(abstract string, extender ExtenderFunc) {
	c.mu.Lock()

	abs := c.getAlias(abstract)

	if inst, ok := c.instances[abs]; ok {
		c.mu.Unlock()

		newInst, err := extender(inst, c)

		if err == nil {
			c.mu.Lock()
			c.instances[abs] = newInst
			c.mu.Unlock()

			c.rebound(abs)
		}

		return
	}

	c.extenders[abs] = append(c.extenders[abs], extender)

	c.mu.Unlock()

	if c.Resolved(abs) {
		c.rebound(abs)
	}
}

// ForgetExtenders removes all extension callbacks for the given abstract.
func (c *Container) ForgetExtenders(abstract string) {
	c.mu.Lock()

	defer c.mu.Unlock()

	delete(c.extenders, c.getAlias(abstract))
}

// ---------- Lifecycle Callbacks ----------

// BeforeResolving registers a callback that fires before the given abstract is
// resolved.
func (c *Container) BeforeResolving(abstract string, callback BeforeResolvingCallback) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.beforeCbs[abstract] = append(c.beforeCbs[abstract], callback)
}

// BeforeResolvingAny registers a global before-resolving callback.
func (c *Container) BeforeResolvingAny(callback BeforeResolvingCallback) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.globalBeforeCbs = append(c.globalBeforeCbs, callback)
}

// Resolving registers a callback that fires when the given abstract is being
// resolved.
func (c *Container) Resolving(abstract string, callback BindingCallback) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.resolvCbs[abstract] = append(c.resolvCbs[abstract], callback)
}

// ResolvingAny registers a global resolving callback.
func (c *Container) ResolvingAny(callback BindingCallback) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.globalResolvCbs = append(c.globalResolvCbs, callback)
}

// AfterResolving registers a callback that fires after the given abstract is
// resolved.
func (c *Container) AfterResolving(abstract string, callback BindingCallback) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.afterCbs[abstract] = append(c.afterCbs[abstract], callback)
}

// AfterResolvingAny registers a global after-resolving callback.
func (c *Container) AfterResolvingAny(callback BindingCallback) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.globalAfterCbs = append(c.globalAfterCbs, callback)
}

// ---------- Rebinding ----------

// Rebinding registers a callback that fires when the given abstract is rebound.
// If the abstract is already bound, the callback is invoked immediately with
// the current instance.
func (c *Container) Rebinding(abstract string, callback BindingCallback) (any, error) {
	c.mu.Lock()

	abs := c.getAlias(abstract)
	c.reboundCbs[abs] = append(c.reboundCbs[abs], callback)

	wasBound := c.isBound(abs)

	c.mu.Unlock()

	if wasBound {
		instance, err := c.Make(abstract)

		if err != nil {
			return nil, err
		}

		callback(instance, c)

		return instance, nil
	}

	return nil, nil
}

// Refresh is a convenience wrapper around Rebinding. When the abstract is
// rebound the setter is called with the new instance.
func (c *Container) Refresh(abstract string, setter func(any)) {
	c.Rebinding(abstract, func(instance any, _ *Container) { //nolint:errcheck
		setter(instance)
	})
}

// rebound fires all registered rebound callbacks for the given abstract.
func (c *Container) rebound(abstract string) {
	c.mu.RLock()

	abs := c.getAlias(abstract)
	cbs := slices.Clone(c.reboundCbs[abs])

	c.mu.RUnlock()

	if len(cbs) == 0 {
		return
	}

	instance, err := c.Make(abstract)

	if err != nil {
		return
	}

	for _, cb := range cbs {
		cb(instance, c)
	}
}

// ---------- Contextual Binding ----------

// When creates a ContextualBindingBuilder for the given concrete type(s).
func (c *Container) When(concrete ...string) *ContextualBindingBuilder {
	return &ContextualBindingBuilder{
		container: c,
		concrete:  concrete,
	}
}

// AddContextualBinding registers a contextual binding. When the concrete type
// is being resolved and needs the given abstract, the implementation is used
// instead.
func (c *Container) AddContextualBinding(concrete, abstract string, implementation any) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if c.contextual[concrete] == nil {
		c.contextual[concrete] = make(map[string]any)
	}

	c.contextual[concrete][abstract] = implementation
}

// getContextualConcrete looks up a contextual binding for the given abstract
// based on the current build stack. Caller must hold the lock.
func (c *Container) getContextualConcrete(abstract string) any {
	if len(c.buildStack) == 0 {
		return nil
	}

	current := c.buildStack[len(c.buildStack)-1]

	if bindings, ok := c.contextual[current]; ok {
		if impl, ok := bindings[abstract]; ok {
			return impl
		}
	}

	return nil
}

// ---------- Method Invocation ----------

// Call invokes the given callable, passing the container and parameters.
func (c *Container) Call(callable MethodCallable, parameters map[string]any) (any, error) {
	return callable(c, parameters)
}

// Wrap returns a closure that invokes the given callable with the container and
// parameters when called.
func (c *Container) Wrap(callable MethodCallable, parameters map[string]any) func() (any, error) {
	return func() (any, error) {
		return c.Call(callable, parameters)
	}
}

// BindMethod registers a callable for a named method binding.
func (c *Container) BindMethod(method string, callback MethodCallable) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.methodBindings[method] = callback
}

// HasMethodBinding reports whether a method binding is registered.
func (c *Container) HasMethodBinding(method string) bool {
	c.mu.RLock()

	defer c.mu.RUnlock()

	_, ok := c.methodBindings[method]

	return ok
}

// CallMethodBinding invokes the registered method binding. The instance is
// passed in the parameters map under the key "_instance".
func (c *Container) CallMethodBinding(method string, instance any) (any, error) {
	c.mu.RLock()
	cb, ok := c.methodBindings[method]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrMethodNotBound, method)
	}

	return cb(c, map[string]any{"_instance": instance})
}

// ---------- Instance Management ----------

// ForgetInstance removes the cached instance for the given abstract.
func (c *Container) ForgetInstance(abstract string) {
	c.mu.Lock()

	defer c.mu.Unlock()

	delete(c.instances, abstract)
}

// ForgetInstances removes all cached instances.
func (c *Container) ForgetInstances() {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.instances = make(map[string]any)
}

// ForgetScopedInstances removes only the cached instances for scoped bindings.
func (c *Container) ForgetScopedInstances() {
	c.mu.Lock()

	defer c.mu.Unlock()

	for abstract, b := range c.bindings {
		if b.scoped {
			delete(c.instances, abstract)
		}
	}
}

// Flush resets the container to an empty state, clearing all bindings,
// instances, aliases, resolved flags, callbacks, and method bindings.
func (c *Container) Flush() {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.bindings = make(map[string]Binding)
	c.instances = make(map[string]any)
	c.aliases = make(map[string]string)
	c.abstractAliases = make(map[string][]string)
	c.resolved = make(map[string]bool)
	c.buildStack = nil
	c.with = nil
	c.contextual = make(map[string]map[string]any)
	c.tags = make(map[string][]string)
	c.extenders = make(map[string][]ExtenderFunc)
	c.reboundCbs = make(map[string][]BindingCallback)
	c.methodBindings = make(map[string]MethodCallable)
	c.globalBeforeCbs = nil
	c.beforeCbs = make(map[string][]BeforeResolvingCallback)
	c.globalResolvCbs = nil
	c.resolvCbs = make(map[string][]BindingCallback)
	c.globalAfterCbs = nil
	c.afterCbs = make(map[string][]BindingCallback)
}

// GetBindings returns a copy of all registered bindings.
func (c *Container) GetBindings() map[string]Binding {
	c.mu.RLock()

	defer c.mu.RUnlock()

	out := make(map[string]Binding, len(c.bindings))

	for k, v := range c.bindings {
		out[k] = v
	}

	return out
}

// ---------- Static Instance ----------

var (
	globalInstance   *Container
	globalInstanceMu sync.Mutex
)

// GetInstance returns the global container instance, creating one if needed.
func GetInstance() *Container {
	globalInstanceMu.Lock()

	defer globalInstanceMu.Unlock()

	if globalInstance == nil {
		globalInstance = New()
	}

	return globalInstance
}

// SetInstance sets or clears the global container instance.
func SetInstance(c *Container) {
	globalInstanceMu.Lock()

	defer globalInstanceMu.Unlock()

	globalInstance = c
}

// ---------- Internal Helpers ----------

// dropStale removes the cached instance and alias entries for the given
// abstract. Caller must hold the write lock.
func (c *Container) dropStale(abstract string) {
	delete(c.instances, abstract)
	delete(c.aliases, abstract)
}

// removeAlias removes the abstract from all alias mappings. Caller must hold
// the write lock.
func (c *Container) removeAlias(abstract string) {
	for abs, aliases := range c.abstractAliases {
		for i, alias := range aliases {
			if alias == abstract {
				c.abstractAliases[abs] = slices.Delete(aliases, i, i+1)

				break
			}
		}
	}

	delete(c.aliases, abstract)
}

// fireCallbacks invokes each callback with the given instance and container.
func fireCallbacks(callbacks []BindingCallback, instance any, c *Container) {
	for _, cb := range callbacks {
		cb(instance, c)
	}
}

// fireBeforeCallbacks invokes each before-resolving callback.
func fireBeforeCallbacks(callbacks []BeforeResolvingCallback, abstract string, params map[string]any, c *Container) {
	for _, cb := range callbacks {
		cb(abstract, params, c)
	}
}
