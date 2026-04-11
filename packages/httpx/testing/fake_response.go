package testing

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

// ResponseAssertions provides test assertion helpers for httptest.ResponseRecorder.
type ResponseAssertions struct {
	t   *testing.T
	rec *httptest.ResponseRecorder
}

// AssertResponse creates a ResponseAssertions wrapper for fluent assertions.
func AssertResponse(t *testing.T, rec *httptest.ResponseRecorder) *ResponseAssertions {
	t.Helper()

	return &ResponseAssertions{t: t, rec: rec}
}

// Status asserts the response has the expected status code.
func (a *ResponseAssertions) Status(expected int) *ResponseAssertions {
	a.t.Helper()

	if a.rec.Code != expected {
		a.t.Fatalf("expected status %d, got %d", expected, a.rec.Code)
	}

	return a
}

// Ok asserts the response has a 200 status.
func (a *ResponseAssertions) Ok() *ResponseAssertions {
	return a.Status(200)
}

// Created asserts the response has a 201 status.
func (a *ResponseAssertions) Created() *ResponseAssertions {
	return a.Status(201)
}

// NoContent asserts the response has a 204 status.
func (a *ResponseAssertions) NoContent() *ResponseAssertions {
	return a.Status(204)
}

// Header asserts a response header has the expected value.
func (a *ResponseAssertions) Header(key, expected string) *ResponseAssertions {
	a.t.Helper()

	got := a.rec.Header().Get(key)

	if got != expected {
		a.t.Fatalf("expected header %s=%q, got %q", key, expected, got)
	}

	return a
}

// HasHeader asserts a response header is present.
func (a *ResponseAssertions) HasHeader(key string) *ResponseAssertions {
	a.t.Helper()

	if a.rec.Header().Get(key) == "" {
		a.t.Fatalf("expected header %s to be present", key)
	}

	return a
}

// MissingHeader asserts a response header is absent.
func (a *ResponseAssertions) MissingHeader(key string) *ResponseAssertions {
	a.t.Helper()

	if a.rec.Header().Get(key) != "" {
		a.t.Fatalf("expected header %s to be absent", key)
	}

	return a
}

// BodyContains asserts the response body contains the expected string.
func (a *ResponseAssertions) BodyContains(expected string) *ResponseAssertions {
	a.t.Helper()

	if !strings.Contains(a.rec.Body.String(), expected) {
		a.t.Fatalf("expected body to contain %q, got %q", expected, a.rec.Body.String())
	}

	return a
}

// BodyEquals asserts the response body exactly matches the expected string.
func (a *ResponseAssertions) BodyEquals(expected string) *ResponseAssertions {
	a.t.Helper()

	if a.rec.Body.String() != expected {
		a.t.Fatalf("expected body %q, got %q", expected, a.rec.Body.String())
	}

	return a
}

// JSONPath asserts that a JSON response contains the given key with the
// expected value (compared as JSON).
func (a *ResponseAssertions) JSONPath(key string, expected any) *ResponseAssertions {
	a.t.Helper()

	var data map[string]any

	if err := json.Unmarshal(a.rec.Body.Bytes(), &data); err != nil {
		a.t.Fatalf("failed to parse JSON body: %v", err)
	}

	got, ok := data[key]

	if !ok {
		a.t.Fatalf("expected JSON key %q to exist", key)
	}

	gotJSON, _ := json.Marshal(got)
	expectedJSON, _ := json.Marshal(expected)

	if string(gotJSON) != string(expectedJSON) {
		a.t.Fatalf("expected JSON %s=%s, got %s", key, string(expectedJSON), string(gotJSON))
	}

	return a
}

// HasJSON asserts the body is valid JSON and contains the given key.
func (a *ResponseAssertions) HasJSON(key string) *ResponseAssertions {
	a.t.Helper()

	var data map[string]any

	if err := json.Unmarshal(a.rec.Body.Bytes(), &data); err != nil {
		a.t.Fatalf("failed to parse JSON body: %v", err)
	}

	if _, ok := data[key]; !ok {
		a.t.Fatalf("expected JSON key %q to exist", key)
	}

	return a
}
