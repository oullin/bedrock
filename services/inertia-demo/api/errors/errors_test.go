package errors

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia"
	"github.com/bedrock/packages/inertia/assert"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/wayfinder"
	"github.com/bedrock/services/inertia-demo/api/internal/testutil"
)

func newErrorContainer(t *testing.T) Container {
	t.Helper()

	i, err := inertia.New(testutil.TestTemplate, inertia.WithVersion("test"))

	if err != nil {
		t.Fatalf("inertia.New() error = %v", err)
	}

	requireAuthFn := func(next http.Handler) http.Handler { return next }
	renderFn := func(w http.ResponseWriter, r *http.Request, component string, props protocol.Props) {
		if strings.TrimSpace(r.RequestURI) == "" {
			r.RequestURI = r.URL.RequestURI()
		}

		if err := i.Render(w, r, component, props); err != nil {
			t.Fatalf("Render(%s) error = %v", component, err)
		}
	}

	return Container{
		RequireAuth: requireAuthFn,
		Render:      renderFn,
	}
}

func TestRegisterRoutesAndRenderHandlers(t *testing.T) {
	t.Parallel()

	container := newErrorContainer(t)
	routes := wayfinder.New()

	routes.Group("features.errors", "/features/errors", func(g *wayfinder.Group) {
		g.Add("http-error", "GET", "/http-error")
		g.Add("network-errors", "GET", "/network-errors")
	})

	mux := http.NewServeMux()

	if err := RegisterRoutes(routes, mux, container); err != nil {
		t.Fatalf("RegisterRoutes() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/features/errors/http-error", nil)

	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	page := assert.AssertFromBytes(t, w.Body.Bytes())

	page.AssertComponent(t, "Features/Errors/HttpError")

	req = httptest.NewRequest(http.MethodGet, "/features/errors/network-errors", nil)

	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	page = assert.AssertFromBytes(t, w.Body.Bytes())

	page.AssertComponent(t, "Features/Errors/NetworkErrors")
}

func TestRegisterRoutes_ValidateError(t *testing.T) {
	t.Parallel()

	err := RegisterRoutes(wayfinder.New(), http.NewServeMux(), Container{})

	if err == nil {
		t.Fatal("RegisterRoutes() error = nil, want validation error")
	}
}

func TestHTTPErrorTriggerHandler(t *testing.T) {
	t.Parallel()

	handler := httpErrorTriggerHandler()

	tests := []struct {
		code string
		want int
	}{
		{code: "403", want: http.StatusForbidden},
		{code: "404", want: http.StatusNotFound},
		{code: "500", want: http.StatusInternalServerError},
		{code: "unhandled", want: http.StatusTeapot},
		{code: "missing", want: http.StatusNotFound},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, "/features/errors/http-error/"+tt.code, nil)

		req.SetPathValue("code", tt.code)

		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != tt.want {
			t.Fatalf("code %s: status = %d, want %d", tt.code, w.Code, tt.want)
		}
	}
}
