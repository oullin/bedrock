package container

import "fmt"

// Make resolves abstract from the container and asserts the result to type T.
func Make[T any](c *Container, abstract string) (T, error) {
	var zero T

	value, err := c.Make(abstract)
	if err != nil {
		return zero, err
	}

	typed, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("%w: wanted %T, got %T", ErrTypeMismatch, zero, value)
	}

	return typed, nil
}

// MustMake resolves abstract and asserts to type T, or panics.
func MustMake[T any](c *Container, abstract string) T {
	value, err := Make[T](c, abstract)
	if err != nil {
		panic(err)
	}

	return value
}
