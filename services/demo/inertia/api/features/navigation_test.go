package features

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
)

func TestNavigationHandlers_RenderBranches(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	tests := []struct {
		name      string
		component string
		target    string
		handler   func(http.ResponseWriter, *http.Request)
	}{
		{name: "preserve state", component: "Features/Navigation/PreserveState", target: "/features/navigation/preserve-state", handler: h.app.preserveStateHandler},
		{name: "preserve scroll", component: "Features/Navigation/PreserveScroll", target: "/features/navigation/preserve-scroll", handler: h.app.preserveScrollHandler},
		{name: "view transitions", component: "Features/Navigation/ViewTransitions", target: "/features/navigation/view-transitions", handler: h.app.viewTransitionsHandler},
		{name: "async requests", component: "Features/Navigation/AsyncRequests", target: "/features/navigation/async-requests", handler: h.app.asyncRequestsHandler},
		{name: "manual visits", component: "Features/Navigation/ManualVisits", target: "/features/navigation/manual-visits", handler: h.app.manualVisitsHandler},
		{name: "scroll management", component: "Features/Navigation/ScrollManagement", target: "/features/navigation/scroll-management", handler: h.app.scrollManagementHandler},
		{name: "instant visits", component: "Features/Navigation/InstantVisits", target: "/features/navigation/instant-visits", handler: h.app.instantVisitsHandler},
		{name: "instant visit target", component: "Features/Navigation/InstantVisitTarget", target: "/features/navigation/instant-visit-target", handler: h.app.instantVisitTargetHandler},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := h.request(http.MethodGet, tt.target, nil)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			page := h.page(t, w)

			page.AssertComponent(t, tt.component)
		})
	}
}

func TestNavigationHandlers_PostBranches(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	tests := []struct {
		name     string
		target   string
		location string
		handler  func(http.ResponseWriter, *http.Request)
	}{
		{name: "links", target: "/features/navigation/links", location: "/features/navigation/links", handler: h.app.linksHandler},
		{name: "history management", target: "/features/navigation/history-management", location: "/features/navigation/history-management", handler: h.app.historyManagementHandler},
		{name: "redirects", target: "/features/navigation/redirects", location: "/features/navigation/redirects", handler: h.app.redirectsHandler},
		{name: "url fragments", target: "/features/navigation/url-fragments", location: "/features/navigation/url-fragments", handler: h.app.urlFragmentsHandler},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := h.request(http.MethodGet, tt.target, nil)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			page := h.page(t, w)

			if strings.TrimSpace(page.Component) == "" {
				t.Fatal("expected rendered component")
			}

			req = h.request(http.MethodPost, tt.target, nil)
			w = httptest.NewRecorder()

			tt.handler(w, req)

			if w.Code != http.StatusFound || w.Header().Get("Location") != tt.location {
				t.Fatalf("status = %d, location = %q", w.Code, w.Header().Get("Location"))
			}

			req = h.request(http.MethodDelete, tt.target, nil)
			w = httptest.NewRecorder()

			tt.handler(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestRedirectsAndFragmentsActionHandlers(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	for _, action := range []string{"back", "to-route"} {
		req := h.request(http.MethodPost, "/features/navigation/redirects/"+action, nil)

		req.SetPathValue("action", action)

		w := httptest.NewRecorder()

		h.app.redirectsActionHandler(w, req)

		if w.Code != http.StatusFound || w.Header().Get("Location") != "/features/navigation/redirects" {
			t.Fatalf("action %s: status = %d, location = %q", action, w.Code, w.Header().Get("Location"))
		}
	}

	req := h.request(http.MethodPost, "/features/navigation/redirects/external", nil)

	req.SetPathValue("action", "external")

	w := httptest.NewRecorder()

	h.app.redirectsActionHandler(w, req)

	if w.Code != http.StatusConflict || w.Header().Get(protocol.HeaderLocation) != "https://oullin.io" {
		t.Fatalf("external redirect: status = %d, location = %q", w.Code, w.Header().Get(protocol.HeaderLocation))
	}

	req = h.request(http.MethodPost, "/features/navigation/redirects/missing", nil)

	req.SetPathValue("action", "missing")

	w = httptest.NewRecorder()

	h.app.redirectsActionHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}

	for _, tt := range []struct {
		action       string
		wantLocation string
	}{
		{"redirect-with-hash", "/features/navigation/url-fragments#section-2"},
		{"preserve-fragment", "/features/navigation/url-fragments"},
	} {
		req = h.request(http.MethodPost, "/features/navigation/url-fragments/"+tt.action, nil)

		req.SetPathValue("action", tt.action)

		w = httptest.NewRecorder()

		h.app.urlFragmentsActionHandler(w, req)

		if w.Code != http.StatusFound || w.Header().Get("Location") != tt.wantLocation {
			t.Fatalf("action %s: status = %d, location = %q, want %d %q", tt.action, w.Code, w.Header().Get("Location"), http.StatusFound, tt.wantLocation)
		}
	}

	req = h.request(http.MethodPost, "/features/navigation/url-fragments/missing", nil)

	req.SetPathValue("action", "missing")

	w = httptest.NewRecorder()

	h.app.urlFragmentsActionHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestAsyncSlowHandler_Branches(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)
	req := h.request(http.MethodGet, "/features/navigation/async-slow", nil)
	w := httptest.NewRecorder()

	h.app.asyncSlowHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/Navigation/AsyncRequests")

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	req = h.request(http.MethodGet, "/features/navigation/async-slow", nil).WithContext(ctx)
	w = httptest.NewRecorder()

	h.app.asyncSlowHandler(w, req)

	if w.Body.Len() != 0 {
		t.Fatalf("cancelled body = %q, want empty", w.Body.String())
	}
}
