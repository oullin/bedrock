package wayfinder

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	route, ok := reg.Lookup("login")

	if !ok {
		t.Fatal("expected to find route 'login'")
	}

	if route.Method != "GET" {
		t.Errorf("expected method GET, got %s", route.Method)
	}

	if route.Pattern != "/login" {
		t.Errorf("expected pattern /login, got %s", route.Pattern)
	}
}

func TestAddNormalizesMethod(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "post", "/login")

	route, _ := reg.Lookup("login")

	if route.Method != "POST" {
		t.Errorf("expected method POST, got %s", route.Method)
	}
}

func TestAddChaining(t *testing.T) {
	t.Parallel()

	reg := New().
		Add("login", "GET", "/login").
		Add("logout", "POST", "/logout")

	if _, ok := reg.Lookup("login"); !ok {
		t.Error("expected to find 'login'")
	}

	if _, ok := reg.Lookup("logout"); !ok {
		t.Error("expected to find 'logout'")
	}
}

func TestGroup(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Group("contacts", "/contacts", func(g *Group) {
		g.Add("index", "GET", "")
		g.Add("show", "GET", "/{contact}")
		g.Add("store", "POST", "")
	})

	tests := []struct {
		name    string
		method  string
		pattern string
	}{
		{"contacts.index", "GET", "/contacts"},
		{"contacts.show", "GET", "/contacts/{contact}"},
		{"contacts.store", "POST", "/contacts"},
	}

	for _, tt := range tests {
		route, ok := reg.Lookup(tt.name)

		if !ok {
			t.Errorf("expected to find route %q", tt.name)

			continue
		}

		if route.Method != tt.method {
			t.Errorf("%s: expected method %s, got %s", tt.name, tt.method, route.Method)
		}

		if route.Pattern != tt.pattern {
			t.Errorf("%s: expected pattern %s, got %s", tt.name, tt.pattern, route.Pattern)
		}
	}
}

func TestNestedGroup(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Group("contacts", "/contacts", func(g *Group) {
		g.Add("index", "GET", "")

		g.Group("notes", "", func(ng *Group) {
			ng.Add("store", "POST", "/{contact}/notes")
		})
	})

	route, ok := reg.Lookup("contacts.notes.store")

	if !ok {
		t.Fatal("expected to find route 'contacts.notes.store'")
	}

	if route.Pattern != "/contacts/{contact}/notes" {
		t.Errorf("expected pattern /contacts/{contact}/notes, got %s", route.Pattern)
	}
}

func TestURL(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	url := reg.URL("contacts.show", map[string]string{"contact": "42"})

	if url != "/contacts/42" {
		t.Errorf("expected /contacts/42, got %s", url)
	}
}

func TestURLEscapesParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	got := reg.URL("contacts.show", map[string]string{"contact": "foo/bar?baz#qux"})
	want := "/contacts/foo%2Fbar%3Fbaz%23qux"

	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}

func TestURLUnknownRoute(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)

	defer log.SetOutput(os.Stderr)

	reg := New()

	url := reg.URL("nonexistent", nil)

	if url != "#!wayfinder:unknown-route" {
		t.Errorf("expected #!wayfinder:unknown-route, got %s", url)
	}

	if !strings.Contains(buf.String(), `wayfinder: unknown route "nonexistent"`) {
		t.Errorf("expected warning log, got %q", buf.String())
	}
}

func TestURLNoParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("dashboard", "GET", "/dashboard")

	url := reg.URL("dashboard", nil)

	if url != "/dashboard" {
		t.Errorf("expected /dashboard, got %s", url)
	}
}

func TestURLEncodesParams(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	url := reg.URL("contacts.show", map[string]string{"contact": "hello world/foo"})

	if url != "/contacts/hello%20world%2Ffoo" {
		t.Errorf("expected /contacts/hello%%20world%%2Ffoo, got %s", url)
	}
}

func TestManifest(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")
	reg.Add("contacts.show", "GET", "/contacts/{contact}")

	m := reg.Manifest()

	if m["login"] != "/login" {
		t.Errorf("expected /login, got %s", m["login"])
	}

	if m["contacts.show"] != "/contacts/{contact}" {
		t.Errorf("expected /contacts/{contact}, got %s", m["contacts.show"])
	}
}

func TestManifestProps(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	props := reg.ManifestProps()

	if props["login"] != "/login" {
		t.Errorf("expected /login, got %v", props["login"])
	}
}

func TestExportOrder(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("c", "GET", "/c")
	reg.Add("a", "GET", "/a")
	reg.Add("b", "GET", "/b")

	routes := reg.Export()

	if len(routes) != 3 {
		t.Fatalf("expected 3 routes, got %d", len(routes))
	}

	if routes[0].Name != "c" || routes[1].Name != "a" || routes[2].Name != "b" {
		t.Error("expected insertion order c, a, b")
	}
}

func TestToJSON(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	data, err := reg.ToJSON()

	if err != nil {
		t.Fatal(err)
	}

	var routes []Route

	if err := json.Unmarshal(data, &routes); err != nil {
		t.Fatal(err)
	}

	if len(routes) != 1 || routes[0].Name != "login" {
		t.Errorf("unexpected JSON output: %s", string(data))
	}
}

func TestLookupMiss(t *testing.T) {
	t.Parallel()

	reg := New()

	_, ok := reg.Lookup("missing")

	if ok {
		t.Error("expected ok=false for missing route")
	}
}

func TestRouteParams(t *testing.T) {
	t.Parallel()

	route := Route{Pattern: "/contacts/{contact}/notes/{note}"}
	params := route.Params()

	if len(params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(params))
	}

	if params[0] != "contact" || params[1] != "note" {
		t.Errorf("expected [contact, note], got %v", params)
	}
}

func TestRouteParamsNone(t *testing.T) {
	t.Parallel()

	route := Route{Pattern: "/dashboard"}
	params := route.Params()

	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestAddOverwrite(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")
	reg.Add("login", "POST", "/auth/login")

	route, _ := reg.Lookup("login")

	if route.Method != "POST" || route.Pattern != "/auth/login" {
		t.Errorf("expected overwritten route, got %+v", route)
	}

	// Should not duplicate in order slice.
	routes := reg.Export()

	if len(routes) != 1 {
		t.Errorf("expected 1 route after overwrite, got %d", len(routes))
	}
}

func TestHandle_RegistersOnMux(t *testing.T) {
	t.Parallel()

	reg := New()

	reg.Add("login", "GET", "/login")

	mux := http.NewServeMux()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("login page"))
	})

	reg.Handle("login", handler, mux)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if rec.Body.String() != "login page" {
		t.Errorf("expected 'login page', got %q", rec.Body.String())
	}
}

func TestHandle_UnknownRoute_Skips(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)

	defer log.SetOutput(os.Stderr)

	reg := New()
	mux := http.NewServeMux()

	reg.Handle("nonexistent", http.NotFoundHandler(), mux)

	if !strings.Contains(buf.String(), "wayfinder: Handle: unknown route") {
		t.Errorf("expected warning log, got %q", buf.String())
	}
}

func TestGroupHandle(t *testing.T) {
	t.Parallel()

	reg := New()
	mux := http.NewServeMux()

	reg.Group("contacts", "/contacts", func(g *Group) {
		g.Add("index", "GET", "")

		g.Handle("index", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			w.Write([]byte("contacts list"))
		}), mux)
	})

	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if rec.Body.String() != "contacts list" {
		t.Errorf("expected 'contacts list', got %q", rec.Body.String())
	}
}
