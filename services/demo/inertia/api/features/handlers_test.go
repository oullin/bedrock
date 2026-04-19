package features

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
)

func TestStateHandlers(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	for _, tt := range []struct {
		component string
		target    string
		handler   func(http.ResponseWriter, *http.Request)
	}{
		{component: "Features/State/Remember", target: "/features/state/remember", handler: h.app.rememberHandler},
		{component: "Features/State/SharedProps", target: "/features/state/shared-props", handler: h.app.sharedPropsHandler},
	} {
		req := h.request(http.MethodGet, tt.target, nil)
		w := httptest.NewRecorder()

		tt.handler(w, req)

		page := h.page(t, w)

		page.AssertComponent(t, tt.component)
	}

	req := h.request(http.MethodGet, "/features/state/flash-data", nil)
	w := httptest.NewRecorder()

	h.app.flashDataHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/State/FlashData")

	for _, kind := range []string{"", "warning", "error"} {
		req = h.request(http.MethodPost, "/features/state/flash-data", []byte(url.Values{
			"kind": {kind},
		}.Encode()))

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w = httptest.NewRecorder()

		h.app.flashDataHandler(w, req)

		if w.Code != http.StatusFound || w.Header().Get("Location") != "/features/state/flash-data" {
			t.Fatalf("kind %q: status = %d, location = %q", kind, w.Code, w.Header().Get("Location"))
		}
	}

	req = h.request(http.MethodPut, "/features/state/flash-data", nil)
	w = httptest.NewRecorder()

	h.app.flashDataHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	req = httptest.NewRequest(http.MethodPost, "/features/state/flash-data", strings.NewReader("bad"))

	req.Header.Set("Content-Type", "multipart/form-data; boundary=bad")
	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.flashDataHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("flashDataHandler status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	for _, action := range []string{"error", "warning", "info"} {
		req = h.request(http.MethodPost, "/features/state/flash-data/"+action, nil)

		req.SetPathValue("action", action)

		w = httptest.NewRecorder()

		h.app.flashDataActionHandler(w, req)

		if w.Code != http.StatusFound {
			t.Fatalf("action %q: status = %d, want %d", action, w.Code, http.StatusFound)
		}
	}
}

func TestLayoutPrefetchEventAndHTTPHandlers(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	for _, tt := range []struct {
		component string
		target    string
		handler   func(http.ResponseWriter, *http.Request)
	}{
		{component: "Features/Layouts/PersistentLayouts", target: "/features/layouts/persistent-layouts", handler: h.app.persistentLayoutsHandler},
		{component: "Features/Layouts/PersistentLayoutsPageTwo", target: "/features/layouts/persistent-layouts/page-2", handler: h.app.persistentLayoutsPage2Handler},
		{component: "Features/Layouts/NestedLayouts", target: "/features/layouts/nested-layouts", handler: h.app.nestedLayoutsHandler},
		{component: "Features/Layouts/Head", target: "/features/layouts/head", handler: h.app.headHandler},
		{component: "Features/Layouts/LayoutProps", target: "/features/layouts/layout-props", handler: h.app.layoutPropsHandler},
		{component: "Features/Prefetching/LinkPrefetch", target: "/features/prefetching/link-prefetch", handler: h.app.linkPrefetchHandler},
		{component: "Features/Prefetching/StaleWhileRevalidate", target: "/features/prefetching/stale-while-revalidate", handler: h.app.staleWhileRevalidateHandler},
		{component: "Features/Prefetching/ManualPrefetch", target: "/features/prefetching/manual-prefetch", handler: h.app.manualPrefetchHandler},
		{component: "Features/Prefetching/CacheManagement", target: "/features/prefetching/cache-management", handler: h.app.cacheManagementHandler},
		{component: "Features/Events/Progress", target: "/features/events/progress", handler: h.app.progressHandler},
		{component: "Features/Http/UseHttp", target: "/features/http/use-http", handler: h.app.useHttpHandler},
	} {
		req := h.request(http.MethodGet, tt.target, nil)
		w := httptest.NewRecorder()

		tt.handler(w, req)

		page := h.page(t, w)

		page.AssertComponent(t, tt.component)
	}

	req := h.request(http.MethodGet, "/features/events/global-events", nil)
	w := httptest.NewRecorder()

	h.app.globalEventsHandler(w, req)

	h.page(t, w).AssertComponent(t, "Features/Events/GlobalEvents")

	req = h.request(http.MethodPost, "/features/events/global-events", nil)
	w = httptest.NewRecorder()

	h.app.globalEventsHandler(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("globalEventsHandler status = %d, want %d", w.Code, http.StatusFound)
	}

	req = h.request(http.MethodDelete, "/features/events/global-events", nil)
	w = httptest.NewRecorder()

	h.app.globalEventsHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("globalEventsHandler delete status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	req = h.request(http.MethodGet, "/features/events/visit-callbacks", nil)
	w = httptest.NewRecorder()

	h.app.visitCallbacksHandler(w, req)

	h.page(t, w).AssertComponent(t, "Features/Events/VisitCallbacks")

	req = h.request(http.MethodPost, "/features/events/visit-callbacks", nil)
	w = httptest.NewRecorder()

	h.app.visitCallbacksHandler(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("visitCallbacksHandler status = %d, want %d", w.Code, http.StatusFound)
	}

	req = h.request(http.MethodDelete, "/features/events/visit-callbacks", nil)
	w = httptest.NewRecorder()

	h.app.visitCallbacksHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("visitCallbacksHandler delete status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestProgressSlowAndUseHttpAPIHandlers(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)
	req := h.request(http.MethodGet, "/features/events/progress/slow", nil)
	w := httptest.NewRecorder()

	h.app.progressSlowHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/Events/Progress")

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	req = h.request(http.MethodGet, "/features/events/progress/slow", nil).WithContext(ctx)
	w = httptest.NewRecorder()

	h.app.progressSlowHandler(w, req)

	if w.Body.Len() != 0 {
		t.Fatalf("cancelled body = %q, want empty", w.Body.String())
	}

	req = h.request(http.MethodPost, "/features/http/use-http/api", []byte(url.Values{
		"name": {"Ada"},
	}.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()

	h.app.useHttpApiHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var payload map[string]any

	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if payload["greeting"] != "Hello, Ada!" {
		t.Fatalf("greeting = %#v, want Hello, Ada!", payload["greeting"])
	}

	req = h.request(http.MethodPost, "/features/http/use-http/api", nil)
	w = httptest.NewRecorder()

	h.app.useHttpApiHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() default error = %v", err)
	}

	if payload["greeting"] != "Hello, World!" {
		t.Fatalf("default greeting = %#v, want Hello, World!", payload["greeting"])
	}

	req = h.request(http.MethodGet, "/features/http/use-http/api", nil)
	w = httptest.NewRecorder()

	h.app.useHttpApiHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	req = httptest.NewRequest(http.MethodPost, "/features/http/use-http/api", strings.NewReader("bad"))

	req.Header.Set("Content-Type", "multipart/form-data; boundary=bad")
	req.Header.Set(protocol.HeaderInertia, "true")

	req.RequestURI = req.URL.RequestURI()
	w = httptest.NewRecorder()

	h.app.useHttpApiHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("useHttpApiHandler status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
