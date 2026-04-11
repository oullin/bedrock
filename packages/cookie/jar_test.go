package cookie_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/bedrock/packages/cookie"
)

func defaultOpts() cookie.Options {
	return cookie.Options{Path: "/", HTTPOnly: true, SameSite: cookie.SameSiteLax}
}

// ---------------------------------------------------------------------------
// Standalone factory functions
// ---------------------------------------------------------------------------

func TestMake(t *testing.T) {
	t.Parallel()

	c := cookie.Make("session", "abc", defaultOpts())

	if c.Name != "session" || c.Value != "abc" {
		t.Fatalf("unexpected cookie: %+v", c)
	}

	if !c.HttpOnly {
		t.Fatal("expected HttpOnly")
	}

	if c.Path != "/" {
		t.Fatalf("expected path '/', got %q", c.Path)
	}

	if c.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSiteLax, got %v", c.SameSite)
	}
}

func TestMakeWithAllOptions(t *testing.T) {
	t.Parallel()

	opts := cookie.Options{
		Path:     "/path",
		Domain:   "example.com",
		MaxAge:   600,
		Secure:   true,
		HTTPOnly: false,
		SameSite: cookie.SameSiteStrict,
	}

	c := cookie.Make("color", "blue", opts)

	if c.Value != "blue" {
		t.Fatalf("expected value 'blue', got %q", c.Value)
	}

	if c.Path != "/path" {
		t.Fatalf("expected path '/path', got %q", c.Path)
	}

	if c.Domain != "example.com" {
		t.Fatalf("expected domain 'example.com', got %q", c.Domain)
	}

	if !c.Secure {
		t.Fatal("expected Secure")
	}

	if c.HttpOnly {
		t.Fatal("expected HttpOnly=false")
	}

	if c.SameSite != http.SameSiteStrictMode {
		t.Fatalf("expected SameSiteStrict, got %v", c.SameSite)
	}
}

func TestForever(t *testing.T) {
	t.Parallel()

	c := cookie.Forever("remember", "xyz", defaultOpts())

	expected := int((400 * 24 * time.Hour).Seconds())

	if c.MaxAge != expected {
		t.Fatalf("expected %d, got %d", expected, c.MaxAge)
	}

	if c.Value != "xyz" {
		t.Fatalf("expected value 'xyz', got %q", c.Value)
	}
}

func TestForget(t *testing.T) {
	t.Parallel()

	c := cookie.Forget("session", defaultOpts())

	if c.MaxAge != -1 {
		t.Fatalf("expected MaxAge=-1, got %d", c.MaxAge)
	}

	if c.Value != "" {
		t.Fatalf("expected empty value, got %q", c.Value)
	}
}

func TestMakeRaw(t *testing.T) {
	t.Parallel()

	opts := defaultOpts()
	opts.Raw = true

	c := cookie.Make("token", "abc=123", opts)

	if c.Raw != "token=abc=123" {
		t.Fatalf("expected raw string, got %q", c.Raw)
	}
}

// ---------------------------------------------------------------------------
// Jar — defaults
// ---------------------------------------------------------------------------

func TestJarSetDefaultsRoundTrip(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())

	newOpts := cookie.Options{Path: "/app", Domain: "test.com", Secure: true}
	j.SetDefaults(newOpts)

	got := j.Defaults()

	if got.Path != "/app" || got.Domain != "test.com" || !got.Secure {
		t.Fatalf("defaults not updated: %+v", got)
	}
}

func TestJarMakeInheritsDefaults(t *testing.T) {
	t.Parallel()

	opts := defaultOpts()
	opts.Secure = true
	j := cookie.NewJar(opts)

	c := j.Make("x", "y", cookie.Options{})

	if !c.Secure {
		t.Fatal("expected Secure from defaults")
	}

	if c.Path != "/" {
		t.Fatalf("expected path '/', got %q", c.Path)
	}

	if !c.HttpOnly {
		t.Fatal("expected HttpOnly from defaults")
	}
}

func TestJarMakeOverridesPath(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())

	c := j.Make("x", "y", cookie.Options{Path: "/custom"})

	if c.Path != "/custom" {
		t.Fatalf("expected path '/custom', got %q", c.Path)
	}
}

func TestJarMakeOverridesDomain(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(cookie.Options{Path: "/", Domain: "default.com"})

	c := j.Make("x", "y", cookie.Options{Domain: "override.com"})

	if c.Domain != "override.com" {
		t.Fatalf("expected domain 'override.com', got %q", c.Domain)
	}
}

func TestJarMakeSecureOrLogic(t *testing.T) {
	t.Parallel()

	// When default is true, Secure stays true even if opts.Secure=false.
	opts := defaultOpts()
	opts.Secure = true
	j := cookie.NewJar(opts)

	c := j.Make("x", "y", cookie.Options{Secure: false})

	if !c.Secure {
		t.Fatal("expected Secure=true (OR with default)")
	}
}

// ---------------------------------------------------------------------------
// Jar — Forever and Forget via jar
// ---------------------------------------------------------------------------

func TestJarForever(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())

	c := j.Forever("remember", "val", cookie.Options{})

	expected := int((400 * 24 * time.Hour).Seconds())

	if c.MaxAge != expected {
		t.Fatalf("expected %d, got %d", expected, c.MaxAge)
	}

	if c.Path != "/" {
		t.Fatalf("expected path from defaults, got %q", c.Path)
	}
}

func TestJarForget(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())

	c := j.Forget("session", cookie.Options{})

	if c.MaxAge != -1 {
		t.Fatalf("expected MaxAge=-1, got %d", c.MaxAge)
	}

	if c.Value != "" {
		t.Fatalf("expected empty value, got %q", c.Value)
	}

	if c.Path != "/" {
		t.Fatalf("expected path from defaults, got %q", c.Path)
	}
}

// ---------------------------------------------------------------------------
// Jar — Queue, HasQueued, Queued
// ---------------------------------------------------------------------------

func TestJarQueue(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())

	if len(j.GetQueued()) != 0 {
		t.Fatal("expected empty queue")
	}

	j.Queue(&http.Cookie{Name: "foo", Value: "bar"})

	if !j.HasQueued("foo") {
		t.Fatal("expected cookie to be queued")
	}

	c := j.Queued("foo")

	if c == nil || c.Value != "bar" {
		t.Fatalf("unexpected queued cookie: %v", c)
	}
}

func TestJarQueueReplacesExisting(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar", Path: "/"})
	j.Queue(&http.Cookie{Name: "foo", Value: "newBar", Path: "/"})

	c := j.Queued("foo", "/")

	if c == nil || c.Value != "newBar" {
		t.Fatalf("expected replaced value 'newBar', got %v", c)
	}

	// Only one cookie should be queued.
	if len(j.GetQueued()) != 1 {
		t.Fatalf("expected 1 queued cookie, got %d", len(j.GetQueued()))
	}
}

func TestJarQueueEmptyValue(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: ""})

	if !j.HasQueued("foo") {
		t.Fatal("expected cookie with empty value to be queued")
	}

	c := j.Queued("foo")

	if c.Value != "" {
		t.Fatalf("expected empty value, got %q", c.Value)
	}
}

func TestJarQueuedWithPath(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	c1 := &http.Cookie{Name: "foo", Value: "bar", Path: "/path"}
	c2 := &http.Cookie{Name: "foo", Value: "rab", Path: "/"}

	j.Queue(c1)
	j.Queue(c2)

	got1 := j.Queued("foo", "/path")

	if got1 == nil || got1.Value != "bar" {
		t.Fatalf("expected 'bar' at /path, got %v", got1)
	}

	got2 := j.Queued("foo", "/")

	if got2 == nil || got2.Value != "rab" {
		t.Fatalf("expected 'rab' at /, got %v", got2)
	}
}

func TestJarQueuedWithoutPathReturnsLast(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar", Path: "/path"})
	j.Queue(&http.Cookie{Name: "foo", Value: "rab", Path: "/"})

	// Without path, should return a cookie (any entry).
	c := j.Queued("foo")

	if c == nil {
		t.Fatal("expected a cookie when called without path")
	}
}

func TestJarQueuedReturnsNilForMissing(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())

	if j.Queued("nonexistent") != nil {
		t.Fatal("expected nil for missing cookie")
	}
}

func TestJarQueuedWithWrongPathReturnsNil(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar", Path: "/path"})

	if j.Queued("foo", "/wrong") != nil {
		t.Fatal("expected nil for wrong path")
	}
}

// ---------------------------------------------------------------------------
// Jar — HasQueued
// ---------------------------------------------------------------------------

func TestJarHasQueuedEmpty(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())

	if j.HasQueued("foo") {
		t.Fatal("expected false for empty queue")
	}
}

func TestJarHasQueuedWithPath(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar", Path: "/path"})
	j.Queue(&http.Cookie{Name: "foo", Value: "rab", Path: "/"})

	if !j.HasQueued("foo", "/path") {
		t.Fatal("expected true for /path")
	}

	if !j.HasQueued("foo", "/") {
		t.Fatal("expected true for /")
	}

	if j.HasQueued("foo", "/wrongPath") {
		t.Fatal("expected false for /wrongPath")
	}
}

func TestJarHasQueuedNonexistent(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar"})

	if j.HasQueued("nonexistent") {
		t.Fatal("expected false for nonexistent cookie")
	}
}

// ---------------------------------------------------------------------------
// Jar — Expire
// ---------------------------------------------------------------------------

func TestJarExpire(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())

	if len(j.GetQueued()) != 0 {
		t.Fatal("expected empty queue")
	}

	j.Expire("foobar", cookie.Options{Path: "/path"})

	c := j.Queued("foobar")

	if c == nil {
		t.Fatal("expected expired cookie to be queued")
	}

	if c.Name != "foobar" {
		t.Fatalf("expected name 'foobar', got %q", c.Name)
	}

	if c.Value != "" {
		t.Fatalf("expected empty value, got %q", c.Value)
	}

	if c.MaxAge != -1 {
		t.Fatalf("expected MaxAge=-1, got %d", c.MaxAge)
	}

	if c.Path != "/path" {
		t.Fatalf("expected path '/path', got %q", c.Path)
	}

	if len(j.GetQueued()) != 1 {
		t.Fatalf("expected 1 queued cookie, got %d", len(j.GetQueued()))
	}
}

// ---------------------------------------------------------------------------
// Jar — Unqueue
// ---------------------------------------------------------------------------

func TestJarUnqueue(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar"})
	j.Unqueue("foo")

	if j.HasQueued("foo") {
		t.Fatal("expected cookie to be removed from queue")
	}

	if len(j.GetQueued()) != 0 {
		t.Fatal("expected empty queue")
	}
}

func TestJarUnqueueNonexistent(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Unqueue("nonexistent") // should not panic

	if len(j.GetQueued()) != 0 {
		t.Fatal("expected empty queue")
	}
}

func TestJarUnqueueMultipleCookies(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar"})
	j.Queue(&http.Cookie{Name: "baz", Value: "qux"})
	j.Unqueue("foo")

	if j.HasQueued("foo") {
		t.Fatal("expected foo to be removed")
	}

	if !j.HasQueued("baz") {
		t.Fatal("expected baz to remain")
	}
}

func TestJarUnqueueWithPath(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	c1 := &http.Cookie{Name: "foo", Value: "bar", Path: "/path"}
	c2 := &http.Cookie{Name: "foo", Value: "rab", Path: "/"}

	j.Queue(c1)
	j.Queue(c2)

	j.Unqueue("foo", "/path")

	// Only the "/" entry should remain.
	if j.HasQueued("foo", "/path") {
		t.Fatal("expected /path to be removed")
	}

	if !j.HasQueued("foo", "/") {
		t.Fatal("expected / to remain")
	}

	got := j.Queued("foo", "/")

	if got == nil || got.Value != "rab" {
		t.Fatalf("expected 'rab', got %v", got)
	}
}

func TestJarUnqueueOnlyCookieForName(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar", Path: "/path"})
	j.Unqueue("foo", "/path")

	// The name bucket should be cleaned up entirely.
	if j.HasQueued("foo") {
		t.Fatal("expected name bucket to be cleaned up")
	}

	if len(j.GetQueued()) != 0 {
		t.Fatal("expected empty queue")
	}
}

func TestJarUnqueueWithPathNonexistentPath(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar", Path: "/path"})

	j.Unqueue("foo", "/wrong") // should not panic or remove anything

	if !j.HasQueued("foo", "/path") {
		t.Fatal("expected /path to still be present")
	}
}

// ---------------------------------------------------------------------------
// Jar — GetQueued
// ---------------------------------------------------------------------------

func TestJarGetQueued(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar", Path: "/path"})
	j.Queue(&http.Cookie{Name: "foo", Value: "rab", Path: "/"})
	j.Queue(&http.Cookie{Name: "baz", Value: "qux", Path: "/path"})

	cookies := j.GetQueued()

	if len(cookies) != 3 {
		t.Fatalf("expected 3 cookies, got %d", len(cookies))
	}
}

func TestJarGetQueuedEmpty(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())

	if len(j.GetQueued()) != 0 {
		t.Fatal("expected empty slice")
	}
}

// ---------------------------------------------------------------------------
// Jar — Flush
// ---------------------------------------------------------------------------

func TestJarFlush(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "a", Value: "1"})
	j.Queue(&http.Cookie{Name: "b", Value: "2"})
	j.Flush()

	if len(j.GetQueued()) != 0 {
		t.Fatal("expected empty queue after flush")
	}
}

func TestJarFlushWithPathCookies(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar", Path: "/a"})
	j.Queue(&http.Cookie{Name: "foo", Value: "baz", Path: "/b"})
	j.Flush()

	if j.HasQueued("foo") {
		t.Fatal("expected all cookies to be flushed")
	}
}

// ---------------------------------------------------------------------------
// Jar — QueueMake, QueueForever
// ---------------------------------------------------------------------------

func TestJarQueueMake(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.QueueMake("tok", "abc123", cookie.Options{})

	if !j.HasQueued("tok") {
		t.Fatal("expected tok to be queued")
	}

	c := j.Queued("tok")

	if c.Value != "abc123" {
		t.Fatalf("expected value 'abc123', got %q", c.Value)
	}
}

func TestJarQueueForever(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.QueueForever("remember", "token", cookie.Options{})

	if !j.HasQueued("remember") {
		t.Fatal("expected remember to be queued")
	}

	c := j.Queued("remember")
	expected := int((400 * 24 * time.Hour).Seconds())

	if c.MaxAge != expected {
		t.Fatalf("expected MaxAge=%d, got %d", expected, c.MaxAge)
	}
}
