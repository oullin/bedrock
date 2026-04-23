package config_test

import (
	"errors"
	"reflect"
	"testing"

	collection "github.com/bedrock/packages/collection/collection"
	"github.com/bedrock/packages/config"
)

func newUpstreamRepository() (*config.Repository, map[string]any) {
	items := map[string]any{
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
		"a.b":       "c",
		"a":         map[string]any{"b.c": "d"},
	}

	return config.New(items), items
}

func assertSame[T comparable](t *testing.T, want T, got T) {
	t.Helper()

	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func assertDeepSame(t *testing.T, want any, got any) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %#v, got %#v", want, got)
	}
}

// Ports of Upstream tests/Config/RepositoryTest.php read behavior:
// RepositoryTest::testGetValueWhenKeyContainDot
// RepositoryTest::testGetBooleanValue
// RepositoryTest::testGetNullValue
// RepositoryTest::testConstruct
// RepositoryTest::testHasIsTrue
// RepositoryTest::testHasIsFalse
// RepositoryTest::testGet
// RepositoryTest::testGetWithDefault
// RepositoryTest::testAll
func TestUpstreamRepositoryReadBehavior(t *testing.T) {
	t.Parallel()

	repo, items := newUpstreamRepository()

	if repo == nil {
		t.Fatal("expected repository")
	}

	assertSame(t, "c", repo.Get("a.b"))

	if got := repo.Get("a.b.c"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	if got := repo.Get("x.y.z"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	if got := repo.Get("."); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	assertSame(t, true, repo.Get("boolean"))

	if !repo.Has("null") {
		t.Fatal("expected explicit nil key to be present")
	}

	if got := repo.Get("null"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	if !repo.Has("foo") {
		t.Fatal("expected foo to exist")
	}

	if repo.Has("not-exist") {
		t.Fatal("expected not-exist to be missing")
	}

	assertSame(t, "bar", repo.Get("foo"))
	assertSame(t, "default", repo.Get("not-exist", "default"))
	assertDeepSame(t, items, repo.All())
}

// Ports of Upstream tests/Config/RepositoryTest.php multi-get behavior:
// RepositoryTest::testGetWithArrayOfKeys
// RepositoryTest::testGetMany
func TestUpstreamRepositoryGetMany(t *testing.T) {
	t.Parallel()

	repo, _ := newUpstreamRepository()

	assertDeepSame(t, map[string]any{
		"foo":  "bar",
		"bar":  "baz",
		"none": nil,
	}, repo.GetMany([]string{"foo", "bar", "none"}))

	assertDeepSame(t, map[string]any{
		"x.y": "default",
		"x.z": "zoo",
		"bar": "baz",
		"baz": "bat",
	}, repo.GetMany([]string{"x.y", "x.z", "bar", "baz"}, map[string]any{
		"x.y": "default",
		"x.z": "default",
		"bar": "default",
	}))
}

// Ports of Upstream tests/Config/RepositoryTest.php write behavior:
// RepositoryTest::testSet
// RepositoryTest::testSetArray
// RepositoryTest::testPrepend
// RepositoryTest::testPush
// RepositoryTest::testPrependWithNewKey
// RepositoryTest::testPushWithNewKey
func TestUpstreamRepositoryWriteBehavior(t *testing.T) {
	t.Parallel()

	repo, _ := newUpstreamRepository()

	repo.Set("key", "value")
	assertSame(t, "value", repo.Get("key"))

	repo.SetMany(map[string]any{
		"key1": "value1",
		"key2": "value2",
		"key3": nil,
		"key4": map[string]any{
			"foo": "bar",
			"bar": map[string]any{"foo": "bar"},
		},
	})

	assertSame(t, "value1", repo.Get("key1"))
	assertSame(t, "value2", repo.Get("key2"))

	if got := repo.Get("key3"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	assertSame(t, "bar", repo.Get("key4.foo"))
	assertSame(t, "bar", repo.Get("key4.bar.foo"))

	if got := repo.Get("key5"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	repo.Prepend("array", "xxx")
	assertSame(t, "xxx", repo.Get("array.0"))
	assertSame(t, "aaa", repo.Get("array.1"))
	assertSame(t, "zzz", repo.Get("array.2"))

	if got := repo.Get("array.3"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	assertDeepSame(t, []any{"xxx", "aaa", "zzz"}, repo.Get("array"))

	repo.Push("array", "yyy")
	assertDeepSame(t, []any{"xxx", "aaa", "zzz", "yyy"}, repo.Get("array"))

	repo.Prepend("new_prepend_key", "xxx")
	assertDeepSame(t, []any{"xxx"}, repo.Get("new_prepend_key"))

	repo.Push("new_push_key", "xxx")
	assertDeepSame(t, []any{"xxx"}, repo.Get("new_push_key"))
}

// Ports of Upstream tests/Config/RepositoryTest.php ArrayAccess behavior through
// explicit Go repository methods:
// RepositoryTest::testOffsetExists
// RepositoryTest::testOffsetGet
// RepositoryTest::testOffsetSet
// RepositoryTest::testOffsetUnset
func TestUpstreamRepositoryExplicitOffsetAdaptation(t *testing.T) {
	t.Parallel()

	repo, items := newUpstreamRepository()

	repo.SetMany(map[string]any{
		"null_value":    nil,
		"empty_string":  "",
		"numeric_value": 123,
	})

	if !repo.Has("foo") || !repo.Has("null_value") || !repo.Has("empty_string") || !repo.Has("numeric_value") {
		t.Fatal("expected explicit keys to exist")
	}

	if repo.Has("not-exist") {
		t.Fatal("expected missing key")
	}

	assertSame(t, "bar", repo.Get("foo"))
	assertDeepSame(t, map[string]any{"x": "xxx", "y": "yyy"}, repo.Get("associate"))

	if got := repo.Get("not-exist"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	repo.Set("offset_key", "value")
	assertSame(t, "value", repo.Get("offset_key"))

	repo.Set("offset_key", "new_value")
	assertSame(t, "new_value", repo.Get("offset_key"))

	repo.Set("new_key", nil)

	if !repo.Has("new_key") {
		t.Fatal("expected nil key to exist")
	}

	if got := repo.Get("new_key"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}

	repo.Set("", "value")
	assertSame(t, "value", repo.Get(""))

	if _, ok := repo.All()["associate"]; !ok {
		t.Fatal("expected associate key before unset")
	}

	assertDeepSame(t, items["associate"], repo.Get("associate"))

	repo.Unset("associate")

	if _, ok := repo.All()["associate"]; !ok {
		t.Fatal("expected associate key to remain present after unset")
	}

	if got := repo.Get("associate"); got != nil {
		t.Fatalf("expected nil after unset, got %v", got)
	}
}

// Ports of Upstream tests/Config/RepositoryTest.php typed getter behavior:
// RepositoryTest::testItGetsAsString
// RepositoryTest::testItThrowsAnExceptionWhenTryingToGetNonStringValueAsString
// RepositoryTest::testItGetsAsArray
// RepositoryTest::testItThrowsAnExceptionWhenTryingToGetNonArrayValueAsArray
// RepositoryTest::testItGetsAsCollection
// RepositoryTest::testItGetsAsBoolean
// RepositoryTest::testItThrowsAnExceptionWhenTryingToGetNonBooleanValueAsBoolean
// RepositoryTest::testItGetsAsInteger
// RepositoryTest::testItThrowsAnExceptionWhenTryingToGetNonIntegerValueAsInteger
// RepositoryTest::testItGetsAsFloat
// RepositoryTest::testItThrowsAnExceptionWhenTryingToGetNonFloatValueAsFloat
func TestUpstreamRepositoryTypedGetters(t *testing.T) {
	t.Parallel()

	repo, _ := newUpstreamRepository()

	s, err := repo.String("a.b")

	if err != nil {
		t.Fatal(err)
	}

	assertSame(t, "c", s)

	if _, err := repo.String("a"); !errors.Is(err, config.ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}

	a, err := repo.Array("array")

	if err != nil {
		t.Fatal(err)
	}

	assertDeepSame(t, []any{"aaa", "zzz"}, a)

	if _, err := repo.Array("a.b"); !errors.Is(err, config.ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}

	c, err := repo.Collection("array")

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := any(c).(*collection.Collection[any]); !ok {
		t.Fatalf("expected collection, got %T", c)
	}

	assertDeepSame(t, []any{"aaa", "zzz"}, c.All())

	b, err := repo.Boolean("boolean")

	if err != nil {
		t.Fatal(err)
	}

	assertSame(t, true, b)

	if _, err := repo.Boolean("a.b"); !errors.Is(err, config.ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}

	i, err := repo.Integer("integer")

	if err != nil {
		t.Fatal(err)
	}

	assertSame(t, 1, i)

	if _, err := repo.Integer("a.b"); !errors.Is(err, config.ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}

	f, err := repo.Float("float")

	if err != nil {
		t.Fatal(err)
	}

	assertSame(t, 1.1, f)

	if _, err := repo.Float("a.b"); !errors.Is(err, config.ErrInvalidType) {
		t.Fatalf("expected ErrInvalidType, got %v", err)
	}
}
