package assert_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/inertia"
	"github.com/bedrock/packages/inertia/assert"
	"github.com/bedrock/packages/inertia/protocol"
)

// --- Failure paths ---

// mockTB captures assertion failures without failing the real test.

type mockTB struct {
	testing.TB
	failed bool
}

type failReader struct{}

type readError struct{}

const testTemplate = `<!DOCTYPE html>
<html>
<head>{{ .inertiaHead }}</head>
<body>{{ .inertia }}</body>
</html>`

func TestAssertFromBytes(t *testing.T) {
	t.Parallel()

	body := []byte(`{"component":"Users/Index","props":{"name":"alice"},"url":"/users","version":"v1"}`)

	a := assert.AssertFromBytes(t, body)

	a.AssertComponent(t, "Users/Index")

	a.AssertURL(t, "/users")

	a.AssertVersion(t, "v1")

	a.AssertHasProp(t, "name")

	a.AssertPropEquals(t, "name", "alice")

	a.AssertMissingProp(t, "nonexistent")
}

func TestAssertFromHandler(t *testing.T) {
	t.Parallel()

	i, err := inertia.New(testTemplate, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i.Render(w, r, "Dashboard", protocol.Props{
			"title": "Test Dashboard",
		})
	})

	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	r.RequestURI = "/dashboard"

	a := assert.AssertFromHandler(t, i, handler, r)

	a.AssertComponent(t, "Dashboard")

	a.AssertURL(t, "/dashboard")

	a.AssertVersion(t, "v1")

	a.AssertPropEquals(t, "title", "Test Dashboard")
}

func (m *mockTB) Helper()                           {}
func (m *mockTB) Errorf(format string, args ...any) { m.failed = true }
func (m *mockTB) Fatalf(format string, args ...any) { m.failed = true }

func newAssertable() assert.AssertableInertia {
	body := []byte(`{"component":"Page","props":{"name":"alice"},"url":"/","version":"v1"}`)

	return assert.AssertFromBytes(&mockTB{}, body)
}

func TestAssertComponent_Failure(t *testing.T) {
	t.Parallel()

	m := &mockTB{}
	a := newAssertable()

	a.AssertComponent(m, "WrongComponent")

	if !m.failed {
		t.Error("expected assertion failure for wrong component")
	}
}

func TestAssertURL_Failure(t *testing.T) {
	t.Parallel()

	m := &mockTB{}
	a := newAssertable()

	a.AssertURL(m, "/wrong")

	if !m.failed {
		t.Error("expected assertion failure for wrong URL")
	}
}

func TestAssertVersion_Failure(t *testing.T) {
	t.Parallel()

	m := &mockTB{}
	a := newAssertable()

	a.AssertVersion(m, "wrong")

	if !m.failed {
		t.Error("expected assertion failure for wrong version")
	}
}

func TestAssertHasProp_Failure(t *testing.T) {
	t.Parallel()

	m := &mockTB{}
	a := newAssertable()

	a.AssertHasProp(m, "nonexistent")

	if !m.failed {
		t.Error("expected assertion failure for missing prop")
	}
}

func TestAssertPropEquals_Failure(t *testing.T) {
	t.Parallel()

	m := &mockTB{}
	a := newAssertable()

	a.AssertPropEquals(m, "name", "wrong-value")

	if !m.failed {
		t.Error("expected assertion failure for wrong prop value")
	}
}

func TestAssertPropEquals_MissingKey(t *testing.T) {
	t.Parallel()

	m := &mockTB{}
	a := newAssertable()

	a.AssertPropEquals(m, "missing", "any")

	if !m.failed {
		t.Error("expected assertion failure for missing key")
	}
}

func TestAssertMissingProp_Failure(t *testing.T) {
	t.Parallel()

	m := &mockTB{}
	a := newAssertable()

	a.AssertMissingProp(m, "name")

	if !m.failed {
		t.Error("expected assertion failure when prop exists")
	}
}

func TestAssertFromBytes_InvalidJSON(t *testing.T) {
	t.Parallel()

	m := &mockTB{}

	assert.AssertFromBytes(m, []byte(`not json`))

	if !m.failed {
		t.Error("expected failure for invalid JSON")
	}
}

func TestAssertFromReader_ReadError(t *testing.T) {
	t.Parallel()

	m := &mockTB{}

	assert.AssertFromReader(m, &failReader{})

	if !m.failed {
		t.Error("expected failure for read error")
	}
}

func TestAssertFromHandler_Non200(t *testing.T) {
	t.Parallel()

	i, err := inertia.New(testTemplate, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	m := &mockTB{}
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RequestURI = "/"

	assert.AssertFromHandler(m, i, handler, r)

	if !m.failed {
		t.Error("expected failure for non-200 response")
	}
}

func (f *failReader) Read(p []byte) (int, error) {
	return 0, &readError{}
}

func (e *readError) Error() string { return "read failed" }

// --- Metadata decoding ---

func TestAssert_DeferredPropsDecoded(t *testing.T) {
	t.Parallel()

	body := []byte(`{"component":"Page","props":{},"url":"/","version":"v1","deferredProps":{"sidebar":["stats","forecast"]}}`)

	a := assert.AssertFromBytes(t, body)

	if len(a.DeferredProps["sidebar"]) != 2 {
		t.Errorf("DeferredProps[sidebar] = %v, want 2 items", a.DeferredProps["sidebar"])
	}
}

func TestAssert_MergePropsDecoded(t *testing.T) {
	t.Parallel()

	body := []byte(`{"component":"Page","props":{},"url":"/","version":"v1","mergeProps":["items"]}`)

	a := assert.AssertFromBytes(t, body)

	if len(a.MergeProps) != 1 || a.MergeProps[0] != "items" {
		t.Errorf("MergeProps = %v, want [items]", a.MergeProps)
	}
}

func TestAssert_DeepMergePropsDecoded(t *testing.T) {
	t.Parallel()

	body := []byte(`{"component":"Page","props":{},"url":"/","version":"v1","deepMergeProps":["data"]}`)

	a := assert.AssertFromBytes(t, body)

	if len(a.DeepMergeProps) != 1 || a.DeepMergeProps[0] != "data" {
		t.Errorf("DeepMergeProps = %v, want [data]", a.DeepMergeProps)
	}
}

func TestAssert_ScrollPropsDecoded(t *testing.T) {
	t.Parallel()

	body := []byte(`{"component":"Page","props":{},"url":"/","version":"v1","scrollProps":{"feed":{"pageName":"feedPage","previousPage":null,"nextPage":2,"currentPage":1,"reset":false}}}`)

	a := assert.AssertFromBytes(t, body)

	scroll, ok := a.ScrollProps["feed"]

	if !ok {
		t.Fatal("missing scrollProps[feed]")
	}

	if scroll.PageName != "feedPage" {
		t.Errorf("pageName = %q", scroll.PageName)
	}
}

func TestAssert_OncePropsDecoded(t *testing.T) {
	t.Parallel()

	body := []byte(`{"component":"Page","props":{},"url":"/","version":"v1","onceProps":{"notes":{"prop":"notes","expiresAt":1700000000}}}`)

	a := assert.AssertFromBytes(t, body)

	once, ok := a.OnceProps["notes"]

	if !ok {
		t.Fatal("missing onceProps[notes]")
	}

	if once.Prop != "notes" {
		t.Errorf("prop = %q", once.Prop)
	}

	if once.ExpiresAt == nil || *once.ExpiresAt != 1700000000 {
		t.Errorf("expiresAt = %v", once.ExpiresAt)
	}
}

func TestAssert_EncryptHistoryDecoded(t *testing.T) {
	t.Parallel()

	body := []byte(`{"component":"Page","props":{},"url":"/","version":"v1","encryptHistory":true}`)

	a := assert.AssertFromBytes(t, body)

	if !a.EncryptHistory {
		t.Error("encryptHistory should be true")
	}
}

func TestAssert_ClearHistoryDecoded(t *testing.T) {
	t.Parallel()

	body := []byte(`{"component":"Page","props":{},"url":"/","version":"v1","clearHistory":true}`)

	a := assert.AssertFromBytes(t, body)

	if !a.ClearHistory {
		t.Error("clearHistory should be true")
	}
}

func TestAssertFromHandler_WithSharedProps(t *testing.T) {
	t.Parallel()

	i, err := inertia.New(testTemplate, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	i.ShareProp("app_name", "TestApp")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i.Render(w, r, "Page")
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RequestURI = "/"

	a := assert.AssertFromHandler(t, i, handler, r)

	a.AssertHasProp(t, "app_name")

	a.AssertPropEquals(t, "app_name", "TestApp")
}

func TestAssertFromHandler_WithContextProps(t *testing.T) {
	t.Parallel()

	i, err := inertia.New(testTemplate, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := inertia.SetProp(r.Context(), "user", "alice")

		i.Render(w, r.WithContext(ctx), "Page")
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RequestURI = "/"

	a := assert.AssertFromHandler(t, i, handler, r)

	a.AssertHasProp(t, "user")

	a.AssertPropEquals(t, "user", "alice")
}

func TestAssertFromHandler_AutoSetsInertiaHeader(t *testing.T) {
	t.Parallel()

	i, err := inertia.New(testTemplate, inertia.WithVersion("v1"))

	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the header was auto-set.
		if r.Header.Get(protocol.HeaderInertia) != "true" {
			t.Error("X-Inertia header should be auto-set by AssertFromHandler")
		}

		i.Render(w, r, "Page")
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RequestURI = "/"
	// Deliberately NOT setting X-Inertia header.

	assert.AssertFromHandler(t, i, handler, r)
}
