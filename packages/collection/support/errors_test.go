package support

import (
	"testing"
)

func TestItemNotFoundErrorDefault(t *testing.T) {
	err := &ItemNotFoundError{}

	if err.Error() != "item not found" {
		t.Errorf("expected 'item not found', got '%s'", err.Error())
	}
}

func TestItemNotFoundErrorCustom(t *testing.T) {
	err := &ItemNotFoundError{Message: "user not found"}

	if err.Error() != "user not found" {
		t.Errorf("expected 'user not found', got '%s'", err.Error())
	}
}

func TestMultipleItemsFoundErrorDefault(t *testing.T) {
	err := &MultipleItemsFoundError{Count: 5}

	if err.Error() != "multiple items found: 5 items" {
		t.Errorf("expected 'multiple items found: 5 items', got '%s'", err.Error())
	}
}

func TestMultipleItemsFoundErrorCustom(t *testing.T) {
	err := &MultipleItemsFoundError{Count: 3, Message: "too many users"}

	if err.Error() != "too many users" {
		t.Errorf("expected 'too many users', got '%s'", err.Error())
	}
}
