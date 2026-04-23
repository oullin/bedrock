package cookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/cookie"
)

// CookieTest::testCookiesAreCreatedWithProperOptions
// CookieTest::testCookiesAreCreatedWithProperOptionsUsingDefaultPathAndDomain
// CookieTest::testCookiesCanSetSecureOptionUsingDefaultPathAndDomain
func TestLaravelCookieCreationInventoryEquivalents(t *testing.T) {
	t.Parallel()

	opts := defaultOpts()
	opts.Domain = "example.test"
	opts.Secure = cookie.BoolPtr(true)
	jar := cookie.NewJar(opts)

	created := jar.Make("session", "abc", cookie.Options{})

	if created.Name != "session" || created.Value != "abc" {
		t.Fatalf("unexpected cookie identity: %#v", created)
	}

	if created.Path != "/" || created.Domain != "example.test" || !created.Secure || !created.HttpOnly {
		t.Fatalf("defaults were not applied: %#v", created)
	}

	overridden := jar.Make("session", "abc", cookie.Options{Secure: cookie.BoolPtr(false)})

	if overridden.Secure {
		t.Fatal("expected explicit secure=false to override the default")
	}
}

// CookieTest::testQueuedCookiesWithoutName
// CookieTest::testQueuedCookiesWithHandlingEmptyValues
// CookieTest::testQueuedCookiesWithRepeatedValue
// CookieTest::testQueuedCookies
// CookieTest::testQueuedWithPath
// CookieTest::testQueuedWithoutPath
// CookieTest::testHasQueued
// CookieTest::testHasQueuedWithPath
// CookieTest::testExpire
// CookieTest::testUnqueue
// CookieTest::testUnqueueMultipleCookies
// CookieTest::testUnqueueWithPath
// CookieTest::testUnqueueOnlyCookieForName
// CookieTest::testQueueCookie
// CookieTest::testQueueWithCreatingNewCookie
// CookieTest::testGetQueuedCookies
// CookieTest::testFlushQueuedCookies
func TestLaravelCookieQueueInventoryEquivalents(t *testing.T) {
	t.Parallel()

	jar := cookie.NewJar(defaultOpts())

	if err := jar.Queue(&http.Cookie{}); err != cookie.ErrEmptyName {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}

	if err := jar.Queue(&http.Cookie{Name: "foo", Value: "", Path: "/"}); err != nil {
		t.Fatal(err)
	}

	if err := jar.Queue(&http.Cookie{Name: "foo", Value: "bar", Path: "/path"}); err != nil {
		t.Fatal(err)
	}

	if err := jar.Queue(jar.Make("foo", "root", cookie.Options{Path: "/"})); err != nil {
		t.Fatal(err)
	}

	if !jar.HasQueued("foo", "/path") || !jar.HasQueued("foo", "/") {
		t.Fatal("expected both path-specific cookies to be queued")
	}

	if got := jar.Queued("foo", "/"); got == nil || got.Value != "root" {
		t.Fatalf("expected root queued cookie, got %#v", got)
	}

	if err := jar.Expire("expired", cookie.Options{Path: "/"}); err != nil {
		t.Fatal(err)
	}

	if got := jar.Queued("expired", "/"); got == nil || got.MaxAge != -1 {
		t.Fatalf("expected expired cookie to be queued, got %#v", got)
	}

	jar.Unqueue("foo", "/path")

	if jar.HasQueued("foo", "/path") || !jar.HasQueued("foo", "/") {
		t.Fatal("expected only the path-specific cookie to be removed")
	}

	if len(jar.GetQueued()) != 2 {
		t.Fatalf("expected root and expired cookies to remain, got %d", len(jar.GetQueued()))
	}

	jar.Flush()

	if len(jar.GetQueued()) != 0 {
		t.Fatal("expected queued cookies to be flushed")
	}
}

// Middleware/AddQueuedCookiesToResponseTest::testHandle
func TestLaravelAddQueuedCookiesMiddlewareInventoryEquivalent(t *testing.T) {
	t.Parallel()

	jar := cookie.NewJar(defaultOpts())

	if err := jar.Queue(&http.Cookie{Name: "queued", Value: "yes"}); err != nil {
		t.Fatal(err)
	}

	mw := cookie.NewAttachQueued(jar)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	cookies := rec.Result().Cookies()

	if len(cookies) != 1 || cookies[0].Name != "queued" || cookies[0].Value != "yes" {
		t.Fatalf("expected queued cookie to be attached, got %#v", cookies)
	}

	if len(jar.GetQueued()) != 0 {
		t.Fatal("expected middleware to flush queued cookies")
	}
}

// Middleware/EncryptCookiesTest::testSetCookieEncryption
// Middleware/EncryptCookiesTest::testQueuedCookieEncryption
// Middleware/EncryptCookiesTest::testCookieDecryption
func TestLaravelEncryptCookiesMiddlewareInventoryEquivalents(t *testing.T) {
	t.Parallel()

	mw := cookie.NewEncryptCookies(stubEncrypter{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "enc:request"})

	var requestValue string

	mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ := r.Cookie("session")
		requestValue = got.Value
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "response"})
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)

	if requestValue != "request" {
		t.Fatalf("expected decrypted request cookie, got %q", requestValue)
	}

	cookies := rec.Result().Cookies()

	if len(cookies) != 1 || cookies[0].Value != "enc:response" {
		t.Fatalf("expected encrypted response cookie, got %#v", cookies)
	}
}
