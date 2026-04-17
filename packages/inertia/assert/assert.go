package assert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia"

	"github.com/bedrock/packages/inertia/protocol"
)

// AssertableInertia holds a decoded Inertia page object for test
// assertions.
type AssertableInertia struct {
	Component      string              `json:"component"`
	Props          protocol.Props         `json:"props"`
	URL            string              `json:"url"`
	Version        string              `json:"version"`
	EncryptHistory bool                `json:"encryptHistory"`
	ClearHistory   bool                `json:"clearHistory"`
	MergeProps     []string            `json:"mergeProps"`
	DeepMergeProps []string            `json:"deepMergeProps"`
	DeferredProps  map[string][]string `json:"deferredProps"`
	ScrollProps    map[string]struct {
		PageName     string `json:"pageName"`
		PreviousPage any    `json:"previousPage"`
		NextPage     any    `json:"nextPage"`
		CurrentPage  any    `json:"currentPage"`
		Reset        bool   `json:"reset"`
	} `json:"scrollProps"`
	OnceProps map[string]struct {
		Prop      string `json:"prop"`
		ExpiresAt *int64 `json:"expiresAt"`
	} `json:"onceProps"`
}

// AssertFromBytes decodes a JSON body into an AssertableInertia.
func AssertFromBytes(t testing.TB, body []byte) AssertableInertia {
	t.Helper()

	var a AssertableInertia

	if err := json.Unmarshal(body, &a); err != nil {
		t.Fatalf("inertia/testing: unmarshal page: %v", err)
	}

	return a
}

// AssertFromReader decodes a JSON body from an io.Reader.
func AssertFromReader(t testing.TB, body io.Reader) AssertableInertia {
	t.Helper()

	data, err := io.ReadAll(body)

	if err != nil {
		t.Fatalf("inertia/testing: read body: %v", err)
	}

	return AssertFromBytes(t, data)
}

// AssertFromHandler runs handler with r through the given Inertia
// instance's middleware and decodes the JSON response. The request
// must have X-Inertia: true set.
func AssertFromHandler(t testing.TB, i *inertia.Inertia, handler http.HandlerFunc, r *http.Request) AssertableInertia {
	t.Helper()

	if strings.TrimSpace(r.Header.Get(protocol.HeaderInertia)) == "" {
		r.Header.Set(protocol.HeaderInertia, "true")
	}

	w := httptest.NewRecorder()

	i.Middleware(handler).ServeHTTP(w, r)

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("inertia/testing: expected 200, got %d", resp.StatusCode)
	}

	return AssertFromReader(t, resp.Body)
}

// AssertComponent checks that the page component matches want.
func (a AssertableInertia) AssertComponent(t testing.TB, want string) {
	t.Helper()

	if a.Component != want {
		t.Errorf("component: got %q, want %q", a.Component, want)
	}
}

// AssertURL checks that the page URL matches want.
func (a AssertableInertia) AssertURL(t testing.TB, want string) {
	t.Helper()

	if a.URL != want {
		t.Errorf("url: got %q, want %q", a.URL, want)
	}
}

// AssertVersion checks that the page version matches want.
func (a AssertableInertia) AssertVersion(t testing.TB, want string) {
	t.Helper()

	if a.Version != want {
		t.Errorf("version: got %q, want %q", a.Version, want)
	}
}

// AssertHasProp checks that the given key exists in the page props.
func (a AssertableInertia) AssertHasProp(t testing.TB, key string) {
	t.Helper()

	if _, ok := a.Props[key]; !ok {
		t.Errorf("prop %q: not found", key)
	}
}

// AssertPropEquals checks that the given key in props equals want
// using reflect.DeepEqual.
func (a AssertableInertia) AssertPropEquals(t testing.TB, key string, want any) {
	t.Helper()

	got, ok := a.Props[key]

	if !ok {
		t.Fatalf("prop %q: not found", key)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("prop %q: got %v, want %v", key, got, want)
	}
}

// AssertMissingProp checks that the given key does not exist in props.
func (a AssertableInertia) AssertMissingProp(t testing.TB, key string) {
	t.Helper()

	if _, ok := a.Props[key]; ok {
		t.Errorf("prop %q: expected absent, but found", key)
	}
}
