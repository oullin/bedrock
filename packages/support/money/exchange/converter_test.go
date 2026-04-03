package exchange

import (
	"errors"
	"testing"

	"github.com/gollin/packages/support/money/exception"
)

func TestNewConverter(t *testing.T) {
	t.Parallel()
	e := NewExchange()
	c, err := NewConverter(e)

	if err != nil {
		t.Fatalf("NewConverter unexpected error: %v", err)
	}

	if c == nil {
		t.Fatal("NewConverter returned nil")
	}

	ex, err := c.GetExchange()

	if err != nil || ex != e {
		t.Error("GetExchange did not return the expected exchange instance")
	}

	c2, err2 := NewConverter(nil)

	if err2 == nil {
		t.Error("Expected error when exchange is nil, got nil")
	}

	if c2 != nil {
		t.Error("Expected nil converter when exchange is nil")
	}
}

func TestConverterGetExchange(t *testing.T) {
	t.Parallel()

	var c *Converter

	ex, err := c.GetExchange()

	if err == nil {
		t.Error("Expected error from nil converter")
	}

	if !errors.Is(err, exception.ErrNoConverterProvided) {
		t.Errorf("Expected ErrNoConverterProvided, got %v", err)
	}

	if ex != nil {
		t.Error("Expected nil from nil converter")
	}

	e := NewExchange()
	c, err = NewConverter(e)

	if err != nil {
		t.Fatalf("NewConverter unexpected error: %v", err)
	}

	ex, err = c.GetExchange()

	if err != nil || ex != e {
		t.Error("GetExchange did not return the expected exchange instance")
	}
}
