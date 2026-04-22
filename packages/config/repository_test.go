package config_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/config"
	"github.com/spf13/viper"
)

func newRepository() *config.Repository {
	return config.New(map[string]any{
		"foo":       "bar",
		"bar":       "baz",
		"baz":       "bat",
		"null":      nil,
		"boolean":   true,
		"integer":   1,
		"float":     1.1,
		"associate": map[string]any{"x": "xxx", "y": "yyy"},
		"array":     []any{"aaa", "zzz"},
		"x":         map[string]any{"z": "zoo"},
	})
}

// ---------------------------------------------------------------------------
// Constructor
// ---------------------------------------------------------------------------

func TestNew(t *testing.T) {
	t.Parallel()

	repo := config.New(map[string]any{"key": "value"})

	if repo == nil {
		t.Fatal("expected non-nil repository")
	}
}

func TestNewFromViper(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("key", "value")

	repo := config.NewFromViper(v)

	if got := repo.Get("key"); got != "value" {
		t.Fatalf("expected %q, got %v", "value", got)
	}
}

func TestViper(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	if repo.Viper() == nil {
		t.Fatal("expected non-nil viper instance")
	}
}

// ---------------------------------------------------------------------------
// Has
// ---------------------------------------------------------------------------

func TestHasIsTrue(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	if !repo.Has("foo") {
		t.Fatal("expected Has(foo) to be true")
	}
}

func TestHasIsFalse(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	if repo.Has("not-exist") {
		t.Fatal("expected Has(not-exist) to be false")
	}
}

func TestHasNestedKey(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	if !repo.Has("x.z") {
		t.Fatal("expected Has(x.z) to be true")
	}

	if repo.Has("x.missing") {
		t.Fatal("expected Has(x.missing) to be false")
	}
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

func TestGet(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	if got := repo.Get("foo"); got != "bar" {
		t.Fatalf("expected %q, got %v", "bar", got)
	}
}

func TestGetWithDefault(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	if got := repo.Get("not-exist", "default"); got != "default" {
		t.Fatalf("expected %q, got %v", "default", got)
	}
}

func TestGetMissingReturnsNil(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	if got := repo.Get("not-exist"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestGetBooleanValue(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	if got := repo.Get("boolean"); got != true {
		t.Fatalf("expected true, got %v", got)
	}
}

func TestGetNestedKey(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	if got := repo.Get("x.z"); got != "zoo" {
		t.Fatalf("expected %q, got %v", "zoo", got)
	}

	if got := repo.Get("associate.x"); got != "xxx" {
		t.Fatalf("expected %q, got %v", "xxx", got)
	}
}

func TestGetNestedKeyMissing(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	if got := repo.Get("x.y.z"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// GetMany
// ---------------------------------------------------------------------------

func TestGetMany(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	result := repo.GetMany([]string{"foo", "bar", "none"})

	if result["foo"] != "bar" {
		t.Fatalf("expected foo=bar, got %v", result["foo"])
	}

	if result["bar"] != "baz" {
		t.Fatalf("expected bar=baz, got %v", result["bar"])
	}

	if result["none"] != nil {
		t.Fatalf("expected none=nil, got %v", result["none"])
	}
}

// ---------------------------------------------------------------------------
// Set
// ---------------------------------------------------------------------------

func TestSet(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	repo.Set("key", "value")

	if got := repo.Get("key"); got != "value" {
		t.Fatalf("expected %q, got %v", "value", got)
	}
}

func TestSetOverwritesExisting(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	repo.Set("foo", "new")

	if got := repo.Get("foo"); got != "new" {
		t.Fatalf("expected %q, got %v", "new", got)
	}
}

func TestSetNested(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	repo.Set("new.nested.key", "value")

	if got := repo.Get("new.nested.key"); got != "value" {
		t.Fatalf("expected %q, got %v", "value", got)
	}
}

func TestSetMany(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	repo.SetMany(map[string]any{
		"key1": "value1",
		"key2": "value2",
	})

	if got := repo.Get("key1"); got != "value1" {
		t.Fatalf("expected %q, got %v", "value1", got)
	}

	if got := repo.Get("key2"); got != "value2" {
		t.Fatalf("expected %q, got %v", "value2", got)
	}
}

func TestSetArrayThenGetNestedKey(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	repo.Set("key4", map[string]any{"foo": "bar"})

	if got := repo.Get("key4.foo"); got != "bar" {
		t.Fatalf("expected %q, got %v", "bar", got)
	}
}

// ---------------------------------------------------------------------------
// Prepend
// ---------------------------------------------------------------------------

func TestPrepend(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	repo.Prepend("array", "xxx")

	got, err := repo.Array("array")

	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 3 {
		t.Fatalf("expected length 3, got %d", len(got))
	}

	if got[0] != "xxx" {
		t.Fatalf("expected first element %q, got %v", "xxx", got[0])
	}

	if got[1] != "aaa" {
		t.Fatalf("expected second element %q, got %v", "aaa", got[1])
	}

	if got[2] != "zzz" {
		t.Fatalf("expected third element %q, got %v", "zzz", got[2])
	}
}

func TestPrependWithNewKey(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	repo.Prepend("new_key", "xxx")

	got, err := repo.Array("new_key")

	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("expected length 1, got %d", len(got))
	}

	if got[0] != "xxx" {
		t.Fatalf("expected %q, got %v", "xxx", got[0])
	}
}

// ---------------------------------------------------------------------------
// Push
// ---------------------------------------------------------------------------

func TestPush(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	repo.Push("array", "xxx")

	got, err := repo.Array("array")

	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 3 {
		t.Fatalf("expected length 3, got %d", len(got))
	}

	if got[0] != "aaa" {
		t.Fatalf("expected first element %q, got %v", "aaa", got[0])
	}

	if got[1] != "zzz" {
		t.Fatalf("expected second element %q, got %v", "zzz", got[1])
	}

	if got[2] != "xxx" {
		t.Fatalf("expected third element %q, got %v", "xxx", got[2])
	}
}

func TestPushWithNewKey(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	repo.Push("new_key", "xxx")

	got, err := repo.Array("new_key")

	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("expected length 1, got %d", len(got))
	}

	if got[0] != "xxx" {
		t.Fatalf("expected %q, got %v", "xxx", got[0])
	}
}

// ---------------------------------------------------------------------------
// All
// ---------------------------------------------------------------------------

func TestAll(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	all := repo.All()

	if all["foo"] != "bar" {
		t.Fatalf("expected foo=bar, got %v", all["foo"])
	}

	if all["bar"] != "baz" {
		t.Fatalf("expected bar=baz, got %v", all["bar"])
	}
}

func TestAllDeepMergesExplicitNestedItems(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("app.env", "local")

	repo := config.NewFromViper(v)
	repo.Set("app.name", "val")

	all := repo.All()
	app, ok := all["app"].(map[string]any)

	if !ok {
		t.Fatalf("expected app to be a map, got %T", all["app"])
	}

	if app["env"] != "local" {
		t.Fatalf("expected app.env=local, got %v", app["env"])
	}

	if app["name"] != "val" {
		t.Fatalf("expected app.name=val, got %v", app["name"])
	}
}

// ---------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------

func TestString(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	got, err := repo.String("foo")

	if err != nil {
		t.Fatal(err)
	}

	if got != "bar" {
		t.Fatalf("expected %q, got %q", "bar", got)
	}
}

func TestStringError(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	_, err := repo.String("integer")

	if err == nil {
		t.Fatal("expected error for non-string value")
	}

	if !errors.Is(err, config.ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}
}

func TestStringDefault(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	got, err := repo.String("missing", "fallback")

	if err != nil {
		t.Fatal(err)
	}

	if got != "fallback" {
		t.Fatalf("expected %q, got %q", "fallback", got)
	}
}

// ---------------------------------------------------------------------------
// Integer
// ---------------------------------------------------------------------------

func TestInteger(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	got, err := repo.Integer("integer")

	if err != nil {
		t.Fatal(err)
	}

	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestIntegerError(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	_, err := repo.Integer("foo")

	if err == nil {
		t.Fatal("expected error for non-int value")
	}

	if !errors.Is(err, config.ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}
}

func TestIntegerDefault(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	got, err := repo.Integer("missing", 42)

	if err != nil {
		t.Fatal(err)
	}

	if got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
}

// ---------------------------------------------------------------------------
// Float
// ---------------------------------------------------------------------------

func TestFloat(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	got, err := repo.Float("float")

	if err != nil {
		t.Fatal(err)
	}

	if got != 1.1 {
		t.Fatalf("expected 1.1, got %f", got)
	}
}

func TestFloatError(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	_, err := repo.Float("foo")

	if err == nil {
		t.Fatal("expected error for non-float value")
	}

	if !errors.Is(err, config.ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}
}

func TestFloatDefault(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	got, err := repo.Float("missing", 3.14)

	if err != nil {
		t.Fatal(err)
	}

	if got != 3.14 {
		t.Fatalf("expected 3.14, got %f", got)
	}
}

// ---------------------------------------------------------------------------
// Boolean
// ---------------------------------------------------------------------------

func TestBoolean(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	got, err := repo.Boolean("boolean")

	if err != nil {
		t.Fatal(err)
	}

	if got != true {
		t.Fatalf("expected true, got %v", got)
	}
}

func TestBooleanError(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	_, err := repo.Boolean("foo")

	if err == nil {
		t.Fatal("expected error for non-bool value")
	}

	if !errors.Is(err, config.ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}
}

func TestBooleanDefault(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	got, err := repo.Boolean("missing", true)

	if err != nil {
		t.Fatal(err)
	}

	if got != true {
		t.Fatalf("expected true, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// Array
// ---------------------------------------------------------------------------

func TestArray(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	got, err := repo.Array("array")

	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("expected length 2, got %d", len(got))
	}

	if got[0] != "aaa" {
		t.Fatalf("expected %q, got %v", "aaa", got[0])
	}

	if got[1] != "zzz" {
		t.Fatalf("expected %q, got %v", "zzz", got[1])
	}
}

func TestArrayError(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	_, err := repo.Array("foo")

	if err == nil {
		t.Fatal("expected error for non-array value")
	}

	if !errors.Is(err, config.ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}
}

func TestArrayDefault(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	fb := []any{"x"}

	got, err := repo.Array("missing", fb)

	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 || got[0] != "x" {
		t.Fatalf("expected [x], got %v", got)
	}
}

// ---------------------------------------------------------------------------
// Type getter error format
// ---------------------------------------------------------------------------

func TestTypeGetterErrorFormat(t *testing.T) {
	t.Parallel()

	repo := newRepository()

	_, err := repo.Integer("foo")

	if err == nil {
		t.Fatal("expected error")
	}

	expected := "config: invalid type: value for key [foo] must be an int, string given"

	if err.Error() != expected {
		t.Fatalf("expected error %q, got %q", expected, err.Error())
	}
}
