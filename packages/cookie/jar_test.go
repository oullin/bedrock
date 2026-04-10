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

func TestMake(t *testing.T) {
	t.Parallel()

	c := cookie.Make("session", "abc", defaultOpts())

	if c.Name != "session" || c.Value != "abc" {
		t.Fatalf("unexpected cookie: %+v", c)
	}
	if !c.HttpOnly {
		t.Fatal("expected HttpOnly")
	}
}

func TestForever(t *testing.T) {
	t.Parallel()

	c := cookie.Forever("remember", "xyz", defaultOpts())

	if c.MaxAge <= 0 {
		t.Fatalf("expected positive max-age, got %d", c.MaxAge)
	}
	// 400 days in seconds
	expected := int((400 * 24 * time.Hour).Seconds())
	if c.MaxAge != expected {
		t.Fatalf("expected %d, got %d", expected, c.MaxAge)
	}
}

func TestForget(t *testing.T) {
	t.Parallel()

	c := cookie.Forget("session", defaultOpts())

	if c.MaxAge != -1 {
		t.Fatalf("expected MaxAge=-1, got %d", c.MaxAge)
	}
}

func TestJarQueue(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar"})

	if !j.HasQueued("foo") {
		t.Fatal("expected cookie to be queued")
	}

	c := j.Queued("foo")
	if c == nil || c.Value != "bar" {
		t.Fatalf("unexpected queued cookie: %v", c)
	}
}

func TestJarUnqueue(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.Queue(&http.Cookie{Name: "foo", Value: "bar"})
	j.Unqueue("foo")

	if j.HasQueued("foo") {
		t.Fatal("expected cookie to be removed from queue")
	}
}

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
}

func TestJarQueueMake(t *testing.T) {
	t.Parallel()

	j := cookie.NewJar(defaultOpts())
	j.QueueMake("tok", "abc123", cookie.Options{})

	if !j.HasQueued("tok") {
		t.Fatal("expected tok to be queued")
	}
}
