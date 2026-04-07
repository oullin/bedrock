package container

import "sync"

// Factory is a constructor function that receives the container so it can
// resolve its own dependencies during construction.
type Factory func(c *Container) (any, error)

type binding struct {
	factory  Factory
	shared   bool
	scoped   bool
	once     sync.Once
	instance any
	err      error
}

func (b *binding) resolve(c *Container) (any, error) {
	if !b.shared {
		return b.factory(c)
	}

	b.once.Do(func() {
		b.instance, b.err = b.factory(c)
	})

	return b.instance, b.err
}
