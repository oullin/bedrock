package features

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
)

func TestRegisterRoutes_MountsHandlers(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)
	mux := http.NewServeMux()

	requireAuthFn := func(next http.Handler) http.Handler { return next }
	renderFn := func(w http.ResponseWriter, r *http.Request, component string, props protocol.Props) {
		if err := h.inertia.Render(w, r, component, props); err != nil {
			t.Fatalf("Render(%s) error = %v", component, err)
		}
	}

	redirectFn := func(w http.ResponseWriter, r *http.Request, url string) {
		h.inertia.Redirect(w, r, url)
	}

	locationFn := func(w http.ResponseWriter, r *http.Request, url string) {
		h.inertia.Location(w, r, url)
	}

	setFlashFn := func(w http.ResponseWriter, msg flash.Message) error {
		return nil
	}

	err := RegisterRoutes(h.registry, mux, Container{
		DB:          h.db,
		RequireAuth: requireAuthFn,
		Render:      renderFn,

		Redirect: redirectFn,

		Location: locationFn,

		RouteURL: h.registry.URL,
		SetFlash: setFlashFn,
	})

	if err != nil {
		t.Fatalf("RegisterRoutes() error = %v", err)
	}

	req := h.request(http.MethodGet, "/features/forms/use-form", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	req = h.request(http.MethodPost, "/features/http/use-http/api", nil)
	w = httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("use-http api status = %d, want %d", w.Code, http.StatusOK)
	}

	req = h.request(http.MethodPost, "/features/navigation/redirects/external", nil)
	w = httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("redirect external status = %d, want %d", w.Code, http.StatusConflict)
	}
}
