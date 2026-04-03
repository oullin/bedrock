package support

import "fmt"

// ItemNotFoundError is returned when an item lookup fails.
type ItemNotFoundError struct {
	Message string
}

// MultipleItemsFoundError is returned when a single item is expected, but multiple is found.
type MultipleItemsFoundError struct {
	Count   int
	Message string
}

func (e *ItemNotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "item not found"
}

func (e *MultipleItemsFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return fmt.Sprintf("multiple items found: %d items", e.Count)
}
