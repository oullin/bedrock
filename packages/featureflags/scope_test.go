package featureflags_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bedrock/packages/featureflags"
)

// scopeableUser implements Scopeable.
type scopeableUser struct{ id string }

// stringerUser implements fmt.Stringer.
type stringerUser struct{ name string }

// structWithID has an exported ID field.
type structWithID struct {
	ID   int
	Name string
}

// structWithLowercaseId has exported Id field.
type structWithLowercaseId struct {
	Id string
}

// structWithNoID has no ID field.
type structWithNoID struct{ Name string }

// ptrStructWithID is tested via pointer.
type ptrStructWithID struct{ ID int64 }

func (u scopeableUser) FeatureScopeIdentifier() string { return "user:" + u.id }

func (u stringerUser) String() string { return "stringer:" + u.name }

func TestSerializeScope_Nil(t *testing.T) {
	t.Parallel()

	got, err := featureflags.SerializeScope(nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != featureflags.NullScope {
		t.Fatalf("expected %q, got %q", featureflags.NullScope, got)
	}
}

func TestSerializeScope_String(t *testing.T) {
	t.Parallel()

	got, err := featureflags.SerializeScope("hello-world")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "hello-world" {
		t.Fatalf("expected %q, got %q", "hello-world", got)
	}
}

func TestSerializeScope_EmptyString(t *testing.T) {
	t.Parallel()

	got, err := featureflags.SerializeScope("")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestSerializeScope_BoolTrue(t *testing.T) {
	t.Parallel()

	got, err := featureflags.SerializeScope(true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "true" {
		t.Fatalf("expected %q, got %q", "true", got)
	}
}

func TestSerializeScope_BoolFalse(t *testing.T) {
	t.Parallel()

	got, err := featureflags.SerializeScope(false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "false" {
		t.Fatalf("expected %q, got %q", "false", got)
	}
}

func TestSerializeScope_Integers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    any
		expected string
	}{
		{int(42), "42"},
		{int8(8), "8"},
		{int16(16), "16"},
		{int32(32), "32"},
		{int64(64), "64"},
		{uint(10), "10"},
		{uint8(8), "8"},
		{uint16(16), "16"},
		{uint32(32), "32"},
		{uint64(64), "64"},
		{int(-1), "-1"},
	}

	for _, c := range cases {
		c := c

		t.Run(fmt.Sprintf("%T(%v)", c.input, c.input), func(t *testing.T) {
			t.Parallel()

			got, err := featureflags.SerializeScope(c.input)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != c.expected {
				t.Fatalf("expected %q, got %q", c.expected, got)
			}
		})
	}
}

func TestSerializeScope_Scopeable(t *testing.T) {
	t.Parallel()

	u := scopeableUser{id: "abc123"}
	got, err := featureflags.SerializeScope(u)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "user:abc123" {
		t.Fatalf("expected %q, got %q", "user:abc123", got)
	}
}

func TestSerializeScope_Stringer(t *testing.T) {
	t.Parallel()

	u := stringerUser{name: "alice"}
	got, err := featureflags.SerializeScope(u)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "stringer:alice" {
		t.Fatalf("expected %q, got %q", "stringer:alice", got)
	}
}

func TestSerializeScope_StructWithID(t *testing.T) {
	t.Parallel()

	u := structWithID{ID: 99, Name: "test"}
	got, err := featureflags.SerializeScope(u)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == "" {
		t.Fatal("expected non-empty key")
	}

	// Must contain the type name and the ID value.
	expected := fmt.Sprintf("github.com/bedrock/packages/featureflags_test.structWithID|99")

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestSerializeScope_StructWithLowercaseId(t *testing.T) {
	t.Parallel()

	u := structWithLowercaseId{Id: "xyz"}
	got, err := featureflags.SerializeScope(u)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := fmt.Sprintf("github.com/bedrock/packages/featureflags_test.structWithLowercaseId|xyz")

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestSerializeScope_PointerToStruct(t *testing.T) {
	t.Parallel()

	u := &ptrStructWithID{ID: 7}
	got, err := featureflags.SerializeScope(u)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "github.com/bedrock/packages/featureflags_test.ptrStructWithID|7"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestSerializeScope_NilPointer(t *testing.T) {
	t.Parallel()

	var u *ptrStructWithID

	got, err := featureflags.SerializeScope(u)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != featureflags.NullScope {
		t.Fatalf("expected %q, got %q", featureflags.NullScope, got)
	}
}

func TestSerializeScope_StructWithNoID_Error(t *testing.T) {
	t.Parallel()

	u := structWithNoID{Name: "test"}
	_, err := featureflags.SerializeScope(u)

	if !errors.Is(err, featureflags.ErrUnserializableScope) {
		t.Fatalf("expected ErrUnserializableScope, got %v", err)
	}
}

func TestSerializeScope_UnsupportedType_Error(t *testing.T) {
	t.Parallel()

	_, err := featureflags.SerializeScope([]int{1, 2, 3})

	if !errors.Is(err, featureflags.ErrUnserializableScope) {
		t.Fatalf("expected ErrUnserializableScope, got %v", err)
	}
}

func TestSerializeScope_MapType_Error(t *testing.T) {
	t.Parallel()

	_, err := featureflags.SerializeScope(map[string]any{"key": "val"})

	if !errors.Is(err, featureflags.ErrUnserializableScope) {
		t.Fatalf("expected ErrUnserializableScope, got %v", err)
	}
}
