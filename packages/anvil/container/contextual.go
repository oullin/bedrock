package container

import "fmt"

// ContextualBindingBuilder provides a fluent API for contextual bindings:
//
//	container.When("service").Needs("dep").Give(value)
type ContextualBindingBuilder struct {
	container *Container
	concrete  string
	need      string
}

// When begins a contextual binding declaration for concrete.
func (c *Container) When(concrete string) *ContextualBindingBuilder {
	return &ContextualBindingBuilder{container: c, concrete: concrete}
}

// Needs specifies the abstract that the concrete needs resolved differently.
func (b *ContextualBindingBuilder) Needs(abstract string) *ContextualBindingBuilder {
	b.need = abstract
	return b
}

// Give provides the value or Factory to use when concrete needs abstract.
// If value is a Factory, it is called during resolution; otherwise the value
// is returned directly.
func (b *ContextualBindingBuilder) Give(value any) {
	c := b.container

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.contextual[b.concrete] == nil {
		c.contextual[b.concrete] = make(map[string]any)
	}

	c.contextual[b.concrete][b.need] = value
}

// MakeFor resolves abstract with contextual overrides for concrete. If no
// contextual binding exists, it falls back to the normal Make resolution.
func (c *Container) MakeFor(concrete, abstract string) (any, error) {
	c.mu.RLock()
	abstract = c.resolveAlias(abstract)

	overrides := c.contextual[concrete]
	var override any
	var hasOverride bool

	if overrides != nil {
		override, hasOverride = overrides[abstract]
	}
	c.mu.RUnlock()

	if !hasOverride {
		return c.Make(abstract)
	}

	if factory, ok := override.(Factory); ok {
		value, err := factory(c)
		if err != nil {
			return nil, fmt.Errorf("%w: contextual resolving %q for %q: %v", ErrResolve, abstract, concrete, err)
		}
		return value, nil
	}

	return override, nil
}
