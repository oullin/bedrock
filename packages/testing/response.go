package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	collectionpkg "github.com/bedrock/packages/collection/collection"
)

type helperT interface {
	Helper()
	Fatalf(format string, args ...any)
}

// Assertions wraps an httptest.ResponseRecorder with Upstream-style response
// assertions.
type Assertions struct {
	t             helperT
	rec           *httptest.ResponseRecorder
	session       map[string]any
	viewName      string
	viewData      map[string]any
	routeResolver func(name string, params map[string]string) string
}

// AssertResponse creates a fluent assertion helper for a response recorder.

// WithSession attaches a session snapshot for session assertions.

// WithView attaches a view snapshot for view assertions.

// WithRouteResolver attaches a resolver used by route-based redirect
// assertions.

// Status asserts the exact status code.

// Ok asserts a 200 response.

// Created asserts a 201 response.

// Accepted asserts a 202 response.

// NoContent asserts a 204 response or the supplied alternate status.

// BadRequest asserts a 400 response.

// NotFound asserts a 404 response.

// MethodNotAllowed asserts a 405 response.

// NotAcceptable asserts a 406 response.

// Forbidden asserts a 403 response.

// Unauthorized asserts a 401 response.

// RequestTimeout asserts a 408 response.

// PaymentRequired asserts a 402 response.

// MovedPermanently asserts a 301 response.

// Found asserts a 302 response.

// NotModified asserts a 304 response.

// TemporaryRedirect asserts a 307 response.

// PermanentRedirect asserts a 308 response.

// Conflict asserts a 409 response.

// Gone asserts a 410 response.

// Unprocessable asserts a 422 response.

// TooManyRequests asserts a 429 response.

// FailedDependency asserts a 424 response.

// ClientError asserts any 4xx response.

// Header asserts that a header has the expected value.

// HeaderContains asserts that a header contains a substring.

// HasHeader asserts that a header is present.

// MissingHeader asserts that a header is absent.

// Location asserts the redirect location header.

// Redirect asserts that the response is a redirect to the expected location.

// RedirectToAction asserts a redirect to a resolved action URL.

// RedirectToRoute asserts a redirect to a resolved route URL.

// RedirectToSignedRoute asserts a redirect to a resolved signed route URL.
//
// This is a narrow Go adaptation: the caller supplies the final resolved URL
// either through the optional route resolver or by passing a literal path.

// RedirectToTemporarySignedRoute asserts a redirect to a resolved temporary
// signed route URL.

// RedirectContains asserts that the redirect location contains a substring.

// RedirectBack asserts a redirect to the request referrer when present.

// BodyEquals asserts the response body exactly matches the expected string.

// BodyContains asserts the response body contains the expected string.

// See asserts that the response body contains the expected string.

// DontSee asserts that the response body does not contain the expected string.

// SeeEscaped asserts that the response body contains the escaped string.

// DontSeeEscaped asserts that the response body does not contain the escaped
// string.

// SeeHtml asserts that the response body contains the raw HTML fragment.

// DontSeeHtml asserts that the response body does not contain the raw HTML
// fragment.

// SeeText asserts that the rendered text contains the expected text after
// stripping HTML tags and normalizing whitespace.

// DontSeeText asserts that the rendered text does not contain the expected
// text after stripping HTML tags and normalizing whitespace.

// SeeInOrder asserts that every string appears in the response body in order.

// SeeHtmlInOrder asserts that every HTML fragment appears in order.

// SeeTextInOrder asserts that every text fragment appears in order after HTML
// stripping and whitespace normalization.

// JSON asserts that the response body matches the expected JSON structure.

// JSONStructure is a narrow alias for JSON assertions in the Go test surface.

// SimilarJSON is an alias for JSON equality after normalising values.

// ExactJSON asserts that the decoded body exactly matches the expected JSON.

// JSONPath asserts that a JSON path exists with the expected value.

// JSONPathCanonicalizing asserts that a JSON path matches after canonicalizing
// array order recursively.

// JSONFragment asserts that the response body contains the fragment anywhere
// within the decoded JSON structure.

// JSONMissingPath asserts that a JSON path is absent.

// JSONIsArray asserts that the decoded JSON body is an array.

// JSONIsObject asserts that the decoded JSON body is an object.

// ViewIs asserts that the attached view name matches.

// ViewHas asserts that the attached view data contains the expected key and
// optional value.

// ViewHasAll asserts that the attached view data contains all the provided
// keys.

// ViewMissing asserts that the attached view data is missing a key.

// ViewData returns a resolved value from the attached view data.

// FluentJSON returns an AssertableJSON helper for fluent assertions.

// Cookie asserts that the response contains a cookie with the expected value.

// PlainCookie is an alias for Cookie.

// CookieMissing asserts that a cookie is absent.

// CookieExpired asserts that a cookie has expired.

// CookieNotExpired asserts that a cookie is still valid.

// SessionHas asserts that the attached session contains a key/value pair.

// SessionHasInput asserts that the attached session contains an input value.

// SessionHasAll asserts that the attached session contains all provided keys.

// SessionMissing asserts that the attached session is missing a key.

// SessionDoesntHaveErrors is an alias for SessionHasNoErrors.

// SessionHasErrors asserts that the attached session contains validation
// errors under the conventional "errors" key.

// SessionHasErrorsAt asserts that the attached session contains validation
// errors under the supplied key path.

// SessionHasNoErrors asserts that the session does not contain validation
// errors.

// SessionHasNoErrorsAt asserts that the supplied session key path does not
// contain validation errors.

// SessionMissingValue asserts that the attached session contains a key whose
// value does not match the supplied expectation. A missing key also passes.

// JSONValidationErrors asserts that the response body contains validation
// errors under an "errors" JSON payload key.

// JSONValidationErrorsAt asserts that the response body contains validation
// errors under the supplied key path.

// JSONMissingValidationErrors asserts that the response body does not contain
// validation errors.

// JSONMissingValidationErrorsAt asserts that the response body does not
// contain validation errors under the supplied key path.

// JSONValue returns the decoded JSON body.

// JSONCollection returns the decoded JSON body as a collection when the body
// is a JSON array.

// Tap invokes the callback with the current assertion helper and returns the
// helper for chaining.

// HTTPPreviewSuccessful asserts that the response indicates httppreview
// success.

// Streamed asserts that the response looks like a streamed response.

// NotStreamed asserts that the response is not marked as streamed.

// StreamedContent asserts the body of a streamed response.

// StreamedJSONContent asserts the body of a streamed JSON response.

// Download asserts that the response offers a file download with the expected
// filename. Pass an empty filename to only assert that a download is offered.

// StreamedBinaryFile asserts that a streamed response offers a download.

// StreamedJSONFile asserts that a streamed JSON response also offers a
// download.

// AssertableJSON provides fluent JSON assertions and interaction tracking.
type AssertableJSON struct {
	t        helperT
	root     any
	topLevel map[string]struct{}
	seen     map[string]struct{}
	skipped  bool
}

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

func AssertResponse(t helperT, rec *httptest.ResponseRecorder) *Assertions {
	t.Helper()

	return &Assertions{t: t, rec: rec}
}

func (a *Assertions) WithSession(session map[string]any) *Assertions {
	a.t.Helper()

	a.session = session

	return a
}

func (a *Assertions) WithView(name string, data map[string]any) *Assertions {
	a.t.Helper()

	a.viewName = name
	a.viewData = data

	return a
}

func (a *Assertions) WithRouteResolver(resolver func(name string, params map[string]string) string) *Assertions {
	a.t.Helper()

	a.routeResolver = resolver

	return a
}

func (a *Assertions) status(expected int) *Assertions {
	a.t.Helper()

	if a.rec.Code != expected {
		a.t.Fatalf("expected status %d, got %d", expected, a.rec.Code)
	}

	return a
}

func (a *Assertions) Status(expected int) *Assertions { return a.status(expected) }

func (a *Assertions) Ok() *Assertions { return a.status(http.StatusOK) }

func (a *Assertions) Created() *Assertions { return a.status(http.StatusCreated) }

func (a *Assertions) Accepted() *Assertions { return a.status(http.StatusAccepted) }

func (a *Assertions) NoContent(expected ...int) *Assertions {
	if len(expected) > 0 {
		return a.status(expected[0])
	}

	return a.status(http.StatusNoContent)
}

func (a *Assertions) BadRequest() *Assertions { return a.status(http.StatusBadRequest) }

func (a *Assertions) NotFound() *Assertions { return a.status(http.StatusNotFound) }

func (a *Assertions) MethodNotAllowed() *Assertions { return a.status(http.StatusMethodNotAllowed) }

func (a *Assertions) NotAcceptable() *Assertions { return a.status(http.StatusNotAcceptable) }

func (a *Assertions) Forbidden() *Assertions { return a.status(http.StatusForbidden) }

func (a *Assertions) Unauthorized() *Assertions { return a.status(http.StatusUnauthorized) }

func (a *Assertions) RequestTimeout() *Assertions { return a.status(http.StatusRequestTimeout) }

func (a *Assertions) PaymentRequired() *Assertions { return a.status(http.StatusPaymentRequired) }

func (a *Assertions) MovedPermanently() *Assertions { return a.status(http.StatusMovedPermanently) }

func (a *Assertions) Found() *Assertions { return a.status(http.StatusFound) }

func (a *Assertions) NotModified() *Assertions { return a.status(http.StatusNotModified) }

func (a *Assertions) TemporaryRedirect() *Assertions { return a.status(http.StatusTemporaryRedirect) }

func (a *Assertions) PermanentRedirect() *Assertions { return a.status(http.StatusPermanentRedirect) }

func (a *Assertions) Conflict() *Assertions { return a.status(http.StatusConflict) }

func (a *Assertions) Gone() *Assertions { return a.status(http.StatusGone) }

func (a *Assertions) Unprocessable() *Assertions { return a.status(http.StatusUnprocessableEntity) }

func (a *Assertions) TooManyRequests() *Assertions { return a.status(http.StatusTooManyRequests) }

func (a *Assertions) FailedDependency() *Assertions { return a.status(http.StatusFailedDependency) }

func (a *Assertions) ClientError() *Assertions {
	a.t.Helper()

	if a.rec.Code < 400 || a.rec.Code > 499 {
		a.t.Fatalf("expected client error status, got %d", a.rec.Code)
	}

	return a
}

func (a *Assertions) Header(key, expected string) *Assertions {
	a.t.Helper()

	got := a.rec.Header().Get(key)

	if got != expected {
		a.t.Fatalf("expected header %s=%q, got %q", key, expected, got)
	}

	return a
}

func (a *Assertions) HeaderContains(key, expected string) *Assertions {
	a.t.Helper()

	got := a.rec.Header().Get(key)

	if !strings.Contains(got, expected) {
		a.t.Fatalf("expected header %s to contain %q, got %q", key, expected, got)
	}

	return a
}

func (a *Assertions) HasHeader(key string) *Assertions {
	a.t.Helper()

	if a.rec.Header().Get(key) == "" {
		a.t.Fatalf("expected header %s to be present", key)
	}

	return a
}

func (a *Assertions) MissingHeader(key string) *Assertions {
	a.t.Helper()

	if a.rec.Header().Get(key) != "" {
		a.t.Fatalf("expected header %s to be absent", key)
	}

	return a
}

func (a *Assertions) Location(expected string) *Assertions {
	a.t.Helper()

	return a.Header("Location", expected)
}

func (a *Assertions) Redirect(expected string) *Assertions {
	a.t.Helper()

	if a.rec.Code < 300 || a.rec.Code > 399 {
		a.t.Fatalf("expected redirect status, got %d", a.rec.Code)
	}

	return a.Location(expected)
}

func (a *Assertions) RedirectToAction(action string, params ...map[string]string) *Assertions {
	a.t.Helper()

	return a.Redirect(a.resolveRouteLocation(action, params...))
}

func (a *Assertions) RedirectToRoute(route string, params ...map[string]string) *Assertions {
	a.t.Helper()

	return a.Redirect(a.resolveRouteLocation(route, params...))
}

func (a *Assertions) RedirectToSignedRoute(route string, params ...map[string]string) *Assertions {
	a.t.Helper()

	return a.Redirect(a.resolveRouteLocation(route, params...))
}

func (a *Assertions) RedirectToTemporarySignedRoute(route string, params ...map[string]string) *Assertions {
	a.t.Helper()

	return a.Redirect(a.resolveRouteLocation(route, params...))
}

func (a *Assertions) RedirectContains(expected string) *Assertions {
	a.t.Helper()

	if a.rec.Code < 300 || a.rec.Code > 399 {
		a.t.Fatalf("expected redirect status, got %d", a.rec.Code)
	}

	return a.HeaderContains("Location", expected)
}

func (a *Assertions) RedirectBack() *Assertions {
	a.t.Helper()

	if referrer := a.rec.Header().Get("Referer"); referrer != "" {
		return a.Redirect(referrer)
	}

	if referrer := a.rec.Header().Get("Referrer"); referrer != "" {
		return a.Redirect(referrer)
	}

	return a.Redirect("/")
}

func (a *Assertions) BodyEquals(expected string) *Assertions {
	a.t.Helper()

	if a.rec.Body.String() != expected {
		a.t.Fatalf("expected body %q, got %q", expected, a.rec.Body.String())
	}

	return a
}

func (a *Assertions) BodyContains(expected string) *Assertions {
	a.t.Helper()

	if !strings.Contains(a.rec.Body.String(), expected) {
		a.t.Fatalf("expected body to contain %q, got %q", expected, a.rec.Body.String())
	}

	return a
}

func (a *Assertions) See(expected string) *Assertions {
	a.t.Helper()

	return a.BodyContains(expected)
}

func (a *Assertions) DontSee(expected string) *Assertions {
	a.t.Helper()

	if strings.Contains(a.rec.Body.String(), expected) {
		a.t.Fatalf("expected body to not contain %q, got %q", expected, a.rec.Body.String())
	}

	return a
}

func (a *Assertions) SeeEscaped(expected string) *Assertions {
	a.t.Helper()

	return a.See(html.EscapeString(expected))
}

func (a *Assertions) DontSeeEscaped(expected string) *Assertions {
	a.t.Helper()

	return a.DontSee(html.EscapeString(expected))
}

func (a *Assertions) SeeHtml(expected string) *Assertions {
	a.t.Helper()

	return a.See(expected)
}

func (a *Assertions) DontSeeHtml(expected string) *Assertions {
	a.t.Helper()

	return a.DontSee(expected)
}

func (a *Assertions) SeeText(expected string) *Assertions {
	a.t.Helper()

	body := normalizeRenderedText(a.rec.Body.String())
	want := normalizeRenderedText(expected)

	if !strings.Contains(body, want) {
		a.t.Fatalf("expected body text to contain %q, got %q", want, body)
	}

	return a
}

func (a *Assertions) DontSeeText(expected string) *Assertions {
	a.t.Helper()

	body := normalizeRenderedText(a.rec.Body.String())
	want := normalizeRenderedText(expected)

	if strings.Contains(body, want) {
		a.t.Fatalf("expected body text to not contain %q, got %q", want, body)
	}

	return a
}

func (a *Assertions) SeeInOrder(expected ...string) *Assertions {
	a.t.Helper()

	return a.seeInOrder(a.rec.Body.String(), expected...)
}

func (a *Assertions) SeeHtmlInOrder(expected ...string) *Assertions {
	a.t.Helper()

	return a.SeeInOrder(expected...)
}

func (a *Assertions) SeeTextInOrder(expected ...string) *Assertions {
	a.t.Helper()

	body := normalizeRenderedText(a.rec.Body.String())
	normalized := make([]string, 0, len(expected))

	for _, item := range expected {
		normalized = append(normalized, normalizeRenderedText(item))
	}

	return a.seeInOrder(body, normalized...)
}

func (a *Assertions) seeInOrder(body string, expected ...string) *Assertions {
	cursor := 0

	for _, item := range expected {
		idx := strings.Index(body[cursor:], item)

		if idx < 0 {
			a.t.Fatalf("expected body to contain %q in order, got %q", item, body)
		}

		cursor += idx + len(item)
	}

	return a
}

func normalizeRenderedText(body string) string {
	body = htmlTagRe.ReplaceAllString(body, " ")
	body = html.UnescapeString(body)

	return strings.Join(strings.Fields(body), " ")
}

func (a *Assertions) JSON(expected any) *Assertions {
	a.t.Helper()

	got := a.mustJSONBody()
	want := mustNormalizedJSON(a.t, expected)

	if !reflect.DeepEqual(got, want) {
		a.t.Fatalf("expected JSON %v, got %v", want, got)
	}

	return a
}

func (a *Assertions) JSONStructure(expected any) *Assertions {
	return a.JSON(expected)
}

func (a *Assertions) SimilarJSON(expected any) *Assertions {
	return a.JSON(expected)
}

func (a *Assertions) ExactJSON(expected any) *Assertions {
	return a.JSON(expected)
}

func (a *Assertions) JSONPath(key string, expected any) *Assertions {
	a.t.Helper()

	got := a.mustJSONBody()
	value, ok := lookupJSON(got, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)

		return a
	}

	if !valueMatches(a.t, value, expected) {
		a.t.Fatalf("expected JSON path %q to match %v, got %v", key, expected, value)
	}

	return a
}

func (a *Assertions) JSONPathCanonicalizing(key string, expected any) *Assertions {
	a.t.Helper()

	got := a.mustJSONBody()
	value, ok := lookupJSON(got, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)

		return a
	}

	want := canonicalizeJSONValue(a.t, expected)
	have := canonicalizeJSONValue(a.t, value)

	if !reflect.DeepEqual(have, want) {
		a.t.Fatalf("expected JSON path %q=%v, got %v", key, want, value)
	}

	return a
}

func (a *Assertions) JSONFragment(fragment any) *Assertions {
	a.t.Helper()

	got := a.mustJSONBody()
	want := mustNormalizedJSON(a.t, fragment)

	if !jsonContains(got, want) {
		a.t.Fatalf("expected JSON body to contain fragment %v", want)
	}

	return a
}

func (a *Assertions) JSONMissingPath(key string) *Assertions {
	a.t.Helper()

	if _, ok := lookupJSON(a.mustJSONBody(), key); ok {
		a.t.Fatalf("expected JSON path %q to be absent", key)
	}

	return a
}

func (a *Assertions) JSONIsArray() *Assertions {
	a.t.Helper()

	if _, ok := a.mustJSONBody().([]any); !ok {
		a.t.Fatalf("expected JSON array body")
	}

	return a
}

func (a *Assertions) JSONIsObject() *Assertions {
	a.t.Helper()

	if _, ok := a.mustJSONBody().(map[string]any); !ok {
		a.t.Fatalf("expected JSON object body")
	}

	return a
}

func (a *Assertions) ViewIs(expected string) *Assertions {
	a.t.Helper()

	if a.viewName != expected {
		a.t.Fatalf("expected view %q, got %q", expected, a.viewName)
	}

	return a
}

func (a *Assertions) ViewHas(key string, expected ...any) *Assertions {
	a.t.Helper()

	if a.viewData == nil {
		a.t.Fatalf("expected view data to be attached")
	}

	value, ok := lookupMap(a.viewData, key)

	if !ok {
		a.t.Fatalf("expected view key %q to be present", key)
	}

	if len(expected) > 0 {
		if !valueMatches(a.t, value, expected[0]) {
			a.t.Fatalf("expected view %s=%v, got %v", key, expected[0], value)
		}
	}

	return a
}

func (a *Assertions) ViewHasAll(keys ...string) *Assertions {
	a.t.Helper()

	for _, key := range keys {
		a.ViewHas(key)
	}

	return a
}

func (a *Assertions) ViewMissing(key string) *Assertions {
	a.t.Helper()

	if a.viewData == nil {
		a.t.Fatalf("expected view data to be attached")
	}

	if _, ok := lookupMap(a.viewData, key); ok {
		a.t.Fatalf("expected view key %q to be absent", key)
	}

	return a
}

func (a *Assertions) ViewData(key string) (any, bool) {
	a.t.Helper()

	if a.viewData == nil {
		return nil, false
	}

	return lookupMap(a.viewData, key)
}

func (a *Assertions) FluentJSON() *AssertableJSON {
	a.t.Helper()

	return newAssertableJSON(a.t, a.mustJSONBody())
}

func (a *Assertions) Cookie(name string, expected ...string) *Assertions {
	a.t.Helper()

	cookie, ok := a.findCookie(name)

	if !ok {
		a.t.Fatalf("expected cookie %s to be present", name)
	}

	if len(expected) > 0 && cookie.Value != expected[0] {
		a.t.Fatalf("expected cookie %s=%q, got %q", name, expected[0], cookie.Value)
	}

	return a
}

func (a *Assertions) PlainCookie(name string, expected ...string) *Assertions {
	return a.Cookie(name, expected...)
}

func (a *Assertions) CookieMissing(name string) *Assertions {
	a.t.Helper()

	if _, ok := a.findCookie(name); ok {
		a.t.Fatalf("expected cookie %s to be absent", name)
	}

	return a
}

func (a *Assertions) CookieExpired(name string) *Assertions {
	a.t.Helper()

	cookie, ok := a.findCookie(name)

	if !ok {
		a.t.Fatalf("expected cookie %s to be present", name)
	}

	if !cookieExpired(cookie) {
		a.t.Fatalf("expected cookie %s to be expired", name)
	}

	return a
}

func (a *Assertions) CookieNotExpired(name string) *Assertions {
	a.t.Helper()

	cookie, ok := a.findCookie(name)

	if !ok {
		a.t.Fatalf("expected cookie %s to be present", name)
	}

	if cookieExpired(cookie) {
		a.t.Fatalf("expected cookie %s to be unexpired", name)
	}

	return a
}

func (a *Assertions) SessionHas(key string, expected any) *Assertions {
	a.t.Helper()

	if a.session == nil {
		a.t.Fatalf("expected session snapshot to be attached")
	}

	value, ok := lookupMap(a.session, key)

	if !ok {
		a.t.Fatalf("expected session key %q to be present", key)
	}

	if !valueMatches(a.t, value, expected) {
		a.t.Fatalf("expected session %s=%v, got %v", key, expected, value)
	}

	return a
}

func (a *Assertions) SessionHasInput(key string, expected any) *Assertions {
	return a.SessionHas(key, expected)
}

func (a *Assertions) SessionHasAll(keys ...string) *Assertions {
	a.t.Helper()

	for _, key := range keys {
		if a.session == nil {
			a.t.Fatalf("expected session snapshot to be attached")
		}

		if _, ok := lookupMap(a.session, key); !ok {
			a.t.Fatalf("expected session key %q to be present", key)
		}
	}

	return a
}

func (a *Assertions) SessionMissing(key string) *Assertions {
	a.t.Helper()

	if a.session == nil {
		a.t.Fatalf("expected session snapshot to be attached")
	}

	if _, ok := lookupMap(a.session, key); ok {
		a.t.Fatalf("expected session key %q to be absent", key)
	}

	return a
}

func (a *Assertions) SessionDoesntHaveErrors() *Assertions {
	return a.SessionHasNoErrors()
}

func (a *Assertions) SessionHasErrors(expected map[string][]string) *Assertions {
	return a.SessionHasErrorsAt("errors", expected)
}

func (a *Assertions) SessionHasErrorsAt(key string, expected map[string][]string) *Assertions {
	a.t.Helper()

	if a.session == nil {
		a.t.Fatalf("expected session snapshot to be attached")
	}

	got, ok := lookupMap(a.session, key)

	if !ok {
		a.t.Fatalf("expected session validation errors")

		return a
	}

	if raw, ok := got.(string); ok {
		var decoded any

		if err := json.Unmarshal([]byte(raw), &decoded); err == nil {
			got = decoded
		}
	}

	if !reflect.DeepEqual(mustNormalizedJSON(a.t, got), mustNormalizedJSON(a.t, expected)) {
		a.t.Fatalf("expected session errors %v, got %v", expected, got)
	}

	return a
}

func (a *Assertions) SessionHasNoErrors() *Assertions {
	return a.SessionHasNoErrorsAt("errors")
}

func (a *Assertions) SessionHasNoErrorsAt(key string) *Assertions {
	a.t.Helper()

	if a.session == nil {
		a.t.Fatalf("expected session snapshot to be attached")
	}

	if got, ok := lookupMap(a.session, key); ok && !isEmptyValue(got) {
		a.t.Fatalf("expected session to have no validation errors, got %v", got)
	}

	return a
}

func (a *Assertions) SessionMissingValue(key string, expected any) *Assertions {
	a.t.Helper()

	if a.session == nil {
		a.t.Fatalf("expected session snapshot to be attached")
	}

	got, ok := lookupMap(a.session, key)

	if !ok {
		return a
	}

	if valueMatches(a.t, got, expected) {
		a.t.Fatalf("expected session key %q to be missing value %v", key, expected)
	}

	return a
}

func (a *Assertions) JSONValidationErrors(expected map[string][]string) *Assertions {
	return a.JSONValidationErrorsAt("errors", expected)
}

func (a *Assertions) JSONValidationErrorsAt(key string, expected map[string][]string) *Assertions {
	a.t.Helper()

	value, ok := lookupJSON(a.mustJSONBody(), key)

	if !ok {
		a.t.Fatalf("expected JSON validation errors")

		return a
	}

	if !reflect.DeepEqual(mustNormalizedJSON(a.t, value), mustNormalizedJSON(a.t, expected)) {
		a.t.Fatalf("expected JSON validation errors %v, got %v", expected, value)
	}

	return a
}

func (a *Assertions) JSONMissingValidationErrors() *Assertions {
	return a.JSONMissingValidationErrorsAt("errors")
}

func (a *Assertions) JSONMissingValidationErrorsAt(key string) *Assertions {
	a.t.Helper()

	if value, ok := lookupJSON(a.mustJSONBody(), key); ok && !isEmptyValue(value) {
		a.t.Fatalf("expected JSON validation errors to be absent, got %v", value)
	}

	return a
}

func (a *Assertions) JSONValue() any {
	a.t.Helper()

	return a.mustJSONBody()
}

func (a *Assertions) JSONCollection() *collectionpkg.Collection[any] {
	a.t.Helper()

	if items, ok := a.mustJSONBody().([]any); ok {
		return collectionpkg.Collect(items)
	}

	return collectionpkg.Empty[any]()
}

func (a *Assertions) Tap(callback func(*Assertions)) *Assertions {
	a.t.Helper()

	if callback != nil {
		callback(a)
	}

	return a
}

func (a *Assertions) HTTPPreviewSuccessful() *Assertions {
	a.t.Helper()

	if got := a.rec.Header().Get("HTTPPreview-Success"); got != "true" {
		a.t.Fatalf("expected HTTPPreview-Success header to be true, got %q", got)
	}

	return a
}

func (a *Assertions) Streamed() *Assertions {
	a.t.Helper()

	if !isStreamed(a.rec) {
		a.t.Fatalf("expected streamed response")
	}

	return a
}

func (a *Assertions) NotStreamed() *Assertions {
	a.t.Helper()

	if isStreamed(a.rec) {
		a.t.Fatalf("expected non-streamed response")
	}

	return a
}

func (a *Assertions) StreamedContent(expected string) *Assertions {
	a.t.Helper()

	return a.Streamed().BodyEquals(expected)
}

func (a *Assertions) StreamedJSONContent(expected any) *Assertions {
	a.t.Helper()

	return a.Streamed().JSON(expected)
}

func (a *Assertions) Download(filename ...string) *Assertions {
	a.t.Helper()

	value := a.rec.Header().Get("Content-Disposition")

	if value == "" {
		a.t.Fatalf("expected Content-Disposition header to be present")
	}

	disposition, params, err := mime.ParseMediaType(value)

	if err != nil {
		a.t.Fatalf("expected valid Content-Disposition header: %v", err)
	}

	if disposition != "attachment" {
		a.t.Fatalf("expected attachment disposition, got %q", disposition)
	}

	if len(filename) > 0 && filename[0] != "" {
		want := filename[0]
		got := params["filename"]

		if got != want {
			a.t.Fatalf("expected download filename %q, got %q", want, got)
		}
	}

	return a
}

func (a *Assertions) StreamedBinaryFile(filename string) *Assertions {
	a.t.Helper()

	return a.Streamed().Download(filename)
}

func (a *Assertions) StreamedJSONFile(expected any, filename string) *Assertions {
	a.t.Helper()

	return a.StreamedJSONContent(expected).Download(filename)
}

func (a *Assertions) mustJSONBody() any {
	a.t.Helper()

	decoded, err := decodeJSON(a.rec.Body.Bytes())

	if err != nil {
		a.t.Fatalf("expected valid JSON body: %v", err)
	}

	return decoded
}

func (a *Assertions) findCookie(name string) (*http.Cookie, bool) {
	cookies := a.rec.Result().Cookies()

	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie, true
		}
	}

	return nil, false
}

func cookieExpired(cookie *http.Cookie) bool {
	if cookie.MaxAge < 0 {
		return true
	}

	if cookie.Expires.IsZero() {
		return false
	}

	return cookie.Expires.Before(time.Now())
}

func isStreamed(rec *httptest.ResponseRecorder) bool {
	if rec == nil {
		return false
	}

	header := rec.Header()

	if strings.Contains(strings.ToLower(header.Get("Transfer-Encoding")), "chunked") {
		return true
	}

	if strings.Contains(strings.ToLower(header.Get("Content-Type")), "text/event-stream") {
		return true
	}

	return header.Get("X-Streamed") == "true"
}

func lookupMap(data map[string]any, key string) (any, bool) {
	current := any(data)

	for _, segment := range strings.Split(key, ".") {
		switch node := current.(type) {
		case map[string]any:
			value, ok := node[segment]

			if !ok {
				return nil, false
			}

			current = value
		case []any:
			idx, err := strconv.Atoi(segment)

			if err != nil || idx < 0 || idx >= len(node) {
				return nil, false
			}

			current = node[idx]
		default:
			return nil, false
		}
	}

	return current, true
}

func valueMatches(t helperT, got, expected any) bool {
	t.Helper()

	switch want := expected.(type) {
	case func(any) bool:
		return want(got)
	case func(any):
		want(got)

		return true
	default:
		return reflect.DeepEqual(mustNormalizedJSON(t, got), mustNormalizedJSON(t, expected))
	}
}

func canonicalizeJSONValue(t helperT, value any) any {
	t.Helper()

	normalized := mustNormalizedJSON(t, value)

	return canonicalizeJSONNode(normalized)
}

func canonicalizeJSONNode(value any) any {
	switch current := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(current))

		for key, item := range current {
			out[key] = canonicalizeJSONNode(item)
		}

		return out
	case []any:
		out := make([]any, len(current))

		for i, item := range current {
			out[i] = canonicalizeJSONNode(item)
		}

		sort.SliceStable(out, func(i, j int) bool {
			return canonicalJSONString(out[i]) < canonicalJSONString(out[j])
		})

		return out
	default:
		return current
	}
}

func canonicalJSONString(value any) string {
	data, err := json.Marshal(value)

	if err != nil {
		return fmt.Sprintf("%T:%v", value, value)
	}

	return string(data)
}

func (a *Assertions) resolveRouteLocation(name string, params ...map[string]string) string {
	merged := make(map[string]string)

	if len(params) > 0 && params[0] != nil {
		for key, value := range params[0] {
			merged[key] = value
		}
	}

	if a.routeResolver != nil {
		if resolved := a.routeResolver(name, merged); resolved != "" {
			return resolved
		}

		return "/"
	}

	return buildRouteLocation(name, merged)
}

func buildRouteLocation(name string, params map[string]string) string {
	if name == "" {
		return "/"
	}

	if len(params) == 0 {
		return name
	}

	keys := make([]string, 0, len(params))

	for key := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	values := url.Values{}

	for _, key := range keys {
		values.Set(key, params[key])
	}

	if strings.Contains(name, "?") {
		return name + "&" + values.Encode()
	}

	return name + "?" + values.Encode()
}

func isEmptyValue(v any) bool {
	if v == nil {
		return true
	}

	switch value := v.(type) {
	case string:
		return value == ""
	case []string:
		return len(value) == 0
	case []any:
		return len(value) == 0
	case map[string]any:
		return len(value) == 0
	case map[string][]string:
		return len(value) == 0
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return rv.Len() == 0
	}

	return false
}

func normalizeJSONValue(t helperT, value any) any {
	t.Helper()

	data, err := jsonMarshal(value)

	if err != nil {
		t.Fatalf("expected JSON-compatible value %T: %v", value, err)
	}

	decoded, err := decodeJSON(data)

	if err != nil {
		t.Fatalf("expected JSON-compatible value %T: %v", value, err)
	}

	return decoded
}

func numericJSONValue(value any) (float64, bool) {
	switch current := value.(type) {
	case json.Number:
		parsed, err := current.Float64()

		return parsed, err == nil
	case float64:
		return current, true
	case float32:
		return float64(current), true
	case int:
		return float64(current), true
	case int8:
		return float64(current), true
	case int16:
		return float64(current), true
	case int32:
		return float64(current), true
	case int64:
		return float64(current), true
	case uint:
		return float64(current), true
	case uint8:
		return float64(current), true
	case uint16:
		return float64(current), true
	case uint32:
		return float64(current), true
	case uint64:
		return float64(current), true
	default:
		return 0, false
	}
}

func mustNormalizedJSON(t helperT, value any) any {
	t.Helper()

	return normalizeJSONValue(t, value)
}

func jsonContains(node any, fragment any) bool {
	switch frag := fragment.(type) {
	case map[string]any:
		current, ok := node.(map[string]any)

		if ok && mapContains(current, frag) {
			return true
		}

		for _, child := range walkJSON(node) {
			if jsonContains(child, frag) {
				return true
			}
		}
	case []any:
		current, ok := node.([]any)

		if ok && sliceContains(current, frag) {
			return true
		}

		for _, child := range walkJSON(node) {
			if jsonContains(child, frag) {
				return true
			}
		}
	default:
		if reflect.DeepEqual(node, frag) {
			return true
		}

		for _, child := range walkJSON(node) {
			if jsonContains(child, frag) {
				return true
			}
		}
	}

	return false
}

func mapContains(have, want map[string]any) bool {
	for key, expected := range want {
		got, ok := have[key]

		if !ok {
			return false
		}

		if !reflect.DeepEqual(got, expected) {
			return false
		}
	}

	return true
}

func sliceContains(have, want []any) bool {
	if len(want) == 0 {
		return false
	}

	if len(want) > len(have) {
		return false
	}

	used := make([]bool, len(have))

	for _, expected := range want {
		found := false

		for i, item := range have {
			if used[i] {
				continue
			}

			if reflect.DeepEqual(item, expected) {
				used[i] = true
				found = true

				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func walkJSON(node any) []any {
	switch value := node.(type) {
	case map[string]any:
		children := make([]any, 0, len(value))

		for _, child := range value {
			children = append(children, child)
		}

		return children
	case []any:
		return append([]any(nil), value...)
	default:
		return nil
	}
}

func decodeJSON(data []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var value any

	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}

	return value, nil
}

func jsonMarshal(value any) ([]byte, error) {
	return json.Marshal(value)
}

func lookupJSON(data any, key string) (any, bool) {
	current := data

	for _, segment := range strings.Split(key, ".") {
		switch node := current.(type) {
		case map[string]any:
			value, ok := node[segment]

			if !ok {
				return nil, false
			}

			current = value
		case []any:
			idx, err := strconv.Atoi(segment)

			if err != nil || idx < 0 || idx >= len(node) {
				return nil, false
			}

			current = node[idx]
		default:
			return nil, false
		}
	}

	return current, true
}

func newAssertableJSON(t helperT, root any) *AssertableJSON {
	t.Helper()

	a := &AssertableJSON{
		t:    t,
		root: root,
		seen: make(map[string]struct{}),
	}

	if obj, ok := root.(map[string]any); ok {
		a.topLevel = make(map[string]struct{}, len(obj))

		for key := range obj {
			a.topLevel[key] = struct{}{}
		}
	} else {
		a.skipped = true
	}

	return a
}

// Has asserts that the JSON path exists.
func (a *AssertableJSON) Has(key string) *AssertableJSON {
	a.t.Helper()

	if _, ok := lookupJSON(a.root, key); !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	a.mark(key)

	return a
}

// HasAll asserts that every provided JSON path exists.
func (a *AssertableJSON) HasAll(keys ...string) *AssertableJSON {
	a.t.Helper()

	for _, key := range keys {
		a.Has(key)
	}

	return a
}

// HasAny asserts that at least one JSON path exists.
func (a *AssertableJSON) HasAny(keys ...string) *AssertableJSON {
	a.t.Helper()

	for _, key := range keys {
		if _, ok := lookupJSON(a.root, key); ok {
			a.mark(key)

			return a
		}
	}

	a.t.Fatalf("expected at least one JSON path to exist: %v", keys)

	return a
}

// HasOnly asserts that the root object only contains the provided top-level
// keys.
func (a *AssertableJSON) HasOnly(keys ...string) *AssertableJSON {
	a.t.Helper()

	root, ok := a.root.(map[string]any)

	if !ok {
		a.t.Fatalf("expected JSON object root")
	}

	if len(root) != len(keys) {
		a.t.Fatalf("expected JSON root to contain only %v, got %v", keys, root)
	}

	allowed := make(map[string]struct{}, len(keys))

	for _, key := range keys {
		allowed[key] = struct{}{}
	}

	for key := range root {
		if _, ok := allowed[key]; !ok {
			a.t.Fatalf("expected JSON root to contain only %v, got extra key %q", keys, key)
		}
	}

	for _, key := range keys {
		a.mark(key)
	}

	return a
}

// Missing asserts that the JSON path does not exist.
func (a *AssertableJSON) Missing(key string) *AssertableJSON {
	a.t.Helper()

	if _, ok := lookupJSON(a.root, key); ok {
		a.t.Fatalf("expected JSON path %q to be absent", key)
	}

	a.mark(key)

	return a
}

// MissingAll asserts that every provided JSON path is absent.
func (a *AssertableJSON) MissingAll(keys ...string) *AssertableJSON {
	a.t.Helper()

	for _, key := range keys {
		a.Missing(key)
	}

	return a
}

// Where asserts that the JSON path has the expected value.
func (a *AssertableJSON) Where(key string, expected any) *AssertableJSON {
	a.t.Helper()

	got, ok := lookupJSON(a.root, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	if !valueMatches(a.t, got, expected) {
		a.t.Fatalf("expected JSON path %q=%v, got %v", key, expected, got)
	}

	a.mark(key)

	return a
}

// WhereAll asserts multiple JSON paths in one call.
func (a *AssertableJSON) WhereAll(values map[string]any) *AssertableJSON {
	a.t.Helper()

	for key, expected := range values {
		a.Where(key, expected)
	}

	return a
}

// WhereAllType asserts multiple JSON paths in one call.
func (a *AssertableJSON) WhereAllType(values map[string]string) *AssertableJSON {
	a.t.Helper()

	for key, wantType := range values {
		a.WhereType(key, wantType)
	}

	return a
}

// WhereNot asserts that the JSON path does not match the expected value.
func (a *AssertableJSON) WhereNot(key string, expected any) *AssertableJSON {
	a.t.Helper()

	got, ok := lookupJSON(a.root, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	if valueMatches(a.t, got, expected) {
		a.t.Fatalf("expected JSON path %q to not equal %v", key, expected)
	}

	a.mark(key)

	return a
}

// WhereNull asserts that the JSON path is null.
func (a *AssertableJSON) WhereNull(key string) *AssertableJSON {
	a.t.Helper()

	return a.Where(key, nil)
}

// WhereNotNull asserts that the JSON path is not null.
func (a *AssertableJSON) WhereNotNull(key string) *AssertableJSON {
	a.t.Helper()

	got, ok := lookupJSON(a.root, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	if got == nil {
		a.t.Fatalf("expected JSON path %q to be non-null", key)
	}

	a.mark(key)

	return a
}

// WhereContains asserts that the JSON path contains the expected value.
func (a *AssertableJSON) WhereContains(key string, expected any) *AssertableJSON {
	a.t.Helper()

	got, ok := lookupJSON(a.root, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	want := expected

	switch expected.(type) {
	case func(any) bool, func(any):
	default:
		want = mustNormalizedJSON(a.t, expected)
	}

	if !jsonValueContains(got, want) {
		a.t.Fatalf("expected JSON path %q to contain %v, got %v", key, expected, got)
	}

	a.mark(key)

	return a
}

// CountMultiple asserts that multiple JSON paths contain arrays of the
// expected lengths.
func (a *AssertableJSON) CountMultiple(values map[string]int) *AssertableJSON {
	a.t.Helper()

	for key, expected := range values {
		a.Count(key, expected)
	}

	return a
}

// WhereType asserts that the JSON path has the given JSON type.
func (a *AssertableJSON) WhereType(key, wantType string) *AssertableJSON {
	a.t.Helper()

	got, ok := lookupJSON(a.root, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	if !jsonTypeMatches(jsonTypeOf(got), wantType) {
		a.t.Fatalf("expected JSON path %q to be %s, got %s", key, wantType, jsonTypeOf(got))
	}

	a.mark(key)

	return a
}

// Count asserts that the JSON path holds an array with the expected length.
func (a *AssertableJSON) Count(key string, expected int) *AssertableJSON {
	a.t.Helper()

	got, ok := lookupJSON(a.root, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	arr, ok := got.([]any)

	if !ok {
		a.t.Fatalf("expected JSON path %q to be an array, got %T", key, got)
	}

	if len(arr) != expected {
		a.t.Fatalf("expected JSON path %q to contain %d items, got %d", key, expected, len(arr))
	}

	a.mark(key)

	return a
}

// Between asserts that the JSON path holds a numeric value within the
// inclusive range.
func (a *AssertableJSON) Between(key string, lower, upper any) *AssertableJSON {
	a.t.Helper()

	got, ok := lookupJSON(a.root, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	value, ok := numericJSONValue(got)

	if !ok {
		a.t.Fatalf("expected JSON path %q to be numeric, got %T", key, got)
	}

	min, ok := numericJSONValue(lower)

	if !ok {
		a.t.Fatalf("expected lower bound for JSON path %q to be numeric, got %T", key, lower)
	}

	max, ok := numericJSONValue(upper)

	if !ok {
		a.t.Fatalf("expected upper bound for JSON path %q to be numeric, got %T", key, upper)
	}

	if value < min || value > max {
		a.t.Fatalf("expected JSON path %q to be between %v and %v, got %v", key, lower, upper, got)
	}

	a.mark(key)

	return a
}

// Scope resolves a nested JSON path and executes the callback against the
// resolved value as a fresh fluent assertion helper.
func (a *AssertableJSON) Scope(key string, callback func(*AssertableJSON)) *AssertableJSON {
	a.t.Helper()

	value, ok := lookupJSON(a.root, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	switch value.(type) {
	case map[string]any, []any:
	default:
		a.t.Fatalf("expected JSON path %q to be an object or array, got %T", key, value)
	}

	a.mark(key)
	a.scopeValue(value, callback)

	return a
}

// FirstScope resolves the first item from an array path and scopes assertions
// to that item.
func (a *AssertableJSON) FirstScope(key string, callback func(*AssertableJSON)) *AssertableJSON {
	a.t.Helper()

	value, ok := lookupJSON(a.root, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	items, ok := value.([]any)

	if !ok {
		a.t.Fatalf("expected JSON path %q to be an array, got %T", key, value)
	}

	if len(items) == 0 {
		a.t.Fatalf("expected JSON path %q to contain at least one item", key)

		return a
	}

	a.mark(key)
	a.scopeValue(items[0], callback)

	return a
}

// EachScope resolves an array path and scopes assertions to each item.
func (a *AssertableJSON) EachScope(key string, callback func(*AssertableJSON)) *AssertableJSON {
	a.t.Helper()

	value, ok := lookupJSON(a.root, key)

	if !ok {
		a.t.Fatalf("expected JSON path %q to exist", key)
	}

	items, ok := value.([]any)

	if !ok {
		a.t.Fatalf("expected JSON path %q to be an array, got %T", key, value)
	}

	a.mark(key)

	for _, item := range items {
		a.scopeValue(item, callback)
	}

	return a
}

// Assert verifies that every top-level JSON key was interacted with.
func (a *AssertableJSON) Assert() {
	a.t.Helper()

	if a.skipped {
		return
	}

	for key := range a.topLevel {
		if _, ok := a.seen[key]; !ok {
			a.t.Fatalf("expected JSON key %q to be interacted with", key)
		}
	}
}

func (a *AssertableJSON) scopeValue(value any, callback func(*AssertableJSON)) {
	if callback == nil {
		return
	}

	scoped := newAssertableJSON(a.t, value)
	callback(scoped)
	scoped.Assert()
}

// AllowMissingInteractions disables the interaction check for the current
// fluent JSON assertion.
func (a *AssertableJSON) AllowMissingInteractions() *AssertableJSON {
	a.t.Helper()

	a.skipped = true

	return a
}

// Tap invokes the callback with the fluent JSON assertion helper.
func (a *AssertableJSON) Tap(callback func(*AssertableJSON)) *AssertableJSON {
	a.t.Helper()

	if callback != nil {
		callback(a)
	}

	return a
}

// Any returns true when the given path exists.
func (a *AssertableJSON) Any(key string) bool {
	_, ok := lookupJSON(a.root, key)

	if ok {
		a.mark(key)
	}

	return ok
}

func (a *AssertableJSON) mark(key string) {
	if a.skipped {
		return
	}

	top := key

	if idx := strings.Index(top, "."); idx >= 0 {
		top = top[:idx]
	}

	if a.topLevel == nil {
		return
	}

	if _, ok := a.topLevel[top]; ok {
		a.seen[top] = struct{}{}
	}
}

func jsonTypeOf(v any) string {
	switch v := v.(type) {
	case nil:
		return "null"
	case string:
		return "string"
	case bool:
		return "boolean"
	case json.Number:
		if strings.Contains(v.String(), ".") {
			return "double"
		}

		return "integer"
	case float64:
		if float64(int64(v)) == v {
			return "integer"
		}

		return "double"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return fmt.Sprintf("%T", v)
	}
}

func jsonTypeMatches(got, want string) bool {
	if got == want {
		return true
	}

	for _, candidate := range strings.Split(want, "|") {
		if strings.TrimSpace(candidate) == got {
			return true
		}
	}

	return false
}

func jsonValueContains(have any, want any) bool {
	switch predicate := want.(type) {
	case func(any) bool:
		if _, ok := have.([]any); ok {
			break
		}

		return predicate(have)
	case func(any):
		if _, ok := have.([]any); ok {
			break
		}

		predicate(have)

		return true
	}

	switch current := have.(type) {
	case string:
		needle, ok := want.(string)

		return ok && strings.Contains(current, needle)
	case []any:
		if expectedSlice, ok := want.([]any); ok {
			return sliceContains(current, expectedSlice)
		}

		for _, item := range current {
			switch predicate := want.(type) {
			case func(any) bool:
				if predicate(item) {
					return true
				}

				continue
			case func(any):
				predicate(item)

				return true
			}

			if reflect.DeepEqual(item, want) {
				return true
			}
		}

		return false
	case map[string]any:
		if expectedMap, ok := want.(map[string]any); ok {
			return mapContains(current, expectedMap)
		}

		return false
	default:
		return reflect.DeepEqual(current, want)
	}
}

// JSONPathValue resolves and returns the JSON path value.
func (a *AssertableJSON) JSONPathValue(key string) (any, bool) {
	return lookupJSON(a.root, key)
}

// Root returns the decoded JSON root value.
func (a *AssertableJSON) Root() any {
	return a.root
}
