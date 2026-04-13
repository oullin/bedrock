package jsonx_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/jsonx"
)

// unknownType is a custom type that implements SchemaType but is not recognized by the serializer.
type unknownType struct{}

func (u *unknownType) ToMap() map[string]any { return nil }
func (u *unknownType) String() string        { return "" }

func TestSerializeUnknownType(t *testing.T) {
	t.Parallel()

	_, err := jsonx.Serialize(&unknownType{})

	if err == nil {
		t.Fatal("expected error for unknown type, got nil")
	}

	if !errors.Is(err, jsonx.ErrUnknownType) {
		t.Fatalf("expected ErrUnknownType, got %v", err)
	}
}
