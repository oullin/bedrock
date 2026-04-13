package container

// ContextualBindingBuilder provides a fluent interface for defining contextual
// bindings. It is created by Container.When and stores bindings that resolve
// differently based on the consuming concrete type.
type ContextualBindingBuilder struct {
	container *Container
	concrete  []string
	needs     string
}

// Needs specifies the abstract dependency that the concrete type(s) require.
func (b *ContextualBindingBuilder) Needs(abstract string) *ContextualBindingBuilder {
	b.needs = abstract

	return b
}

// Give provides the implementation to use when the concrete type needs the
// abstract. The implementation can be a Factory or a raw value.
func (b *ContextualBindingBuilder) Give(implementation any) {
	for _, concrete := range b.concrete {
		b.container.AddContextualBinding(concrete, b.needs, implementation)
	}
}

// GiveTagged provides all services tagged with the given tag as the
// implementation. The tagged services are resolved into a slice.
func (b *ContextualBindingBuilder) GiveTagged(tag string) {
	b.Give(Factory(func(c *Container) (any, error) {
		return c.Tagged(tag), nil
	}))
}

// GiveConfig provides a configuration value as the implementation. The
// container must have a "config" binding that implements a Get method.
func (b *ContextualBindingBuilder) GiveConfig(key string, fallback ...any) {
	b.Give(Factory(func(c *Container) (any, error) {
		cfg, err := c.Make("config")

		if err != nil {
			if len(fallback) > 0 {
				return fallback[0], nil
			}

			return nil, err
		}

		type configGetter interface {
			Get(key string, fallback ...any) any
		}

		if getter, ok := cfg.(configGetter); ok {
			return getter.Get(key, fallback...), nil
		}

		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return nil, nil
	}))
}
