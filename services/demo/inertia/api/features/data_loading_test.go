package features

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
)

func TestDeferredPropsHandler_InitialAndPartial(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)
	req := h.request(http.MethodGet, "/features/data-loading/deferred-props", nil)
	w := httptest.NewRecorder()

	h.app.deferredPropsHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/DataLoading/DeferredProps")

	page.AssertPropEquals(t, "quickStat", float64(42))

	if len(page.DeferredProps["slow"]) == 0 || len(page.DeferredProps["heavy"]) == 0 {
		t.Fatalf("DeferredProps = %#v, want slow and heavy groups", page.DeferredProps)
	}

	req = h.partialRequest(http.MethodGet, "/features/data-loading/deferred-props", "Features/DataLoading/DeferredProps", "slowStats,heavyData")
	w = httptest.NewRecorder()

	h.app.deferredPropsHandler(w, req)

	page = h.page(t, w)

	page.AssertHasProp(t, "slowStats")
	page.AssertHasProp(t, "heavyData")
}

func TestPartialReloadsInfiniteScrollAndPolling(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	req := h.request(http.MethodGet, "/features/data-loading/partial-reloads", nil)
	w := httptest.NewRecorder()

	h.app.partialReloadsHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/DataLoading/PartialReloads")

	page.AssertHasProp(t, "users")
	page.AssertHasProp(t, "stats")

	req = h.request(http.MethodGet, "/features/data-loading/infinite-scroll", nil)
	w = httptest.NewRecorder()

	h.app.infiniteScrollHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Features/DataLoading/InfiniteScroll")

	page.AssertHasProp(t, "contacts")

	contacts, ok := page.Props["contacts"].(map[string]any)

	if !ok || contacts["next_cursor"] == nil {
		t.Fatalf("contacts = %#v, want next cursor", page.Props["contacts"])
	}

	cursor, _ := contacts["next_cursor"].(string)
	req = h.request(http.MethodGet, "/features/data-loading/infinite-scroll?cursor="+cursor, nil)
	w = httptest.NewRecorder()

	h.app.infiniteScrollHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Features/DataLoading/InfiniteScroll")

	req = h.request(http.MethodGet, "/features/data-loading/polling", nil)
	w = httptest.NewRecorder()

	h.app.pollingHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Features/DataLoading/Polling")

	page.AssertHasProp(t, "contactCount")

	h.db.Close()

	req = h.request(http.MethodGet, "/features/data-loading/partial-reloads", nil)
	w = httptest.NewRecorder()

	h.app.partialReloadsHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("partialReloadsHandler status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	req = h.request(http.MethodGet, "/features/data-loading/infinite-scroll", nil)
	w = httptest.NewRecorder()

	h.app.infiniteScrollHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("infiniteScrollHandler status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestWhenVisibleOptionalAndOnceProps(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)

	req := h.request(http.MethodGet, "/features/data-loading/when-visible", nil)
	w := httptest.NewRecorder()

	h.app.whenVisibleHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/DataLoading/WhenVisible")

	page.AssertMissingProp(t, "section1")

	req = h.partialRequest(http.MethodGet, "/features/data-loading/when-visible", "Features/DataLoading/WhenVisible", "section1,section2,section3")
	w = httptest.NewRecorder()

	h.app.whenVisibleHandler(w, req)

	page = h.page(t, w)

	page.AssertHasProp(t, "section1")
	page.AssertHasProp(t, "section2")
	page.AssertHasProp(t, "section3")

	req = h.request(http.MethodGet, "/features/data-loading/optional-props", nil)
	w = httptest.NewRecorder()

	h.app.optionalPropsHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Features/DataLoading/OptionalProps")

	page.AssertHasProp(t, "regularData")

	page.AssertMissingProp(t, "optionalData")

	req = h.partialRequest(http.MethodGet, "/features/data-loading/optional-props", "Features/DataLoading/OptionalProps", "optionalData,deferredData")
	w = httptest.NewRecorder()

	h.app.optionalPropsHandler(w, req)

	page = h.page(t, w)

	page.AssertHasProp(t, "optionalData")
	page.AssertHasProp(t, "deferredData")

	req = h.request(http.MethodGet, "/features/data-loading/once-props/0", nil)

	req.SetPathValue("page", "0")

	w = httptest.NewRecorder()

	h.app.oncePropsHandler(w, req)

	page = h.page(t, w)

	page.AssertComponent(t, "Features/DataLoading/OnceProps")

	page.AssertPropEquals(t, "page", float64(1))
}

func TestPropMergingAndSlowCancellation(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)
	req := h.request(http.MethodGet, "/features/data-loading/prop-merging", nil)
	w := httptest.NewRecorder()

	h.app.propMergingHandler(w, req)

	page := h.page(t, w)

	page.AssertComponent(t, "Features/DataLoading/PropMerging")

	if len(page.MergeProps) == 0 {
		t.Fatalf("MergeProps = %#v, want merge props", page.MergeProps)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	req = h.partialRequest(http.MethodGet, "/features/data-loading/deferred-props", "Features/DataLoading/DeferredProps", "slowStats")
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()

	h.app.deferredPropsHandler(w, req)

	page = h.page(t, w)

	if page.Props["slowStats"] != nil {
		t.Fatalf("slowStats = %#v, want nil on cancelled context", page.Props["slowStats"])
	}

	req = h.partialRequest(http.MethodGet, "/features/data-loading/optional-props", "Features/DataLoading/OptionalProps", "deferredData")
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()

	h.app.optionalPropsHandler(w, req)

	page = h.page(t, w)

	if page.Props["deferredData"] != nil {
		t.Fatalf("deferredData = %#v, want nil on cancelled context", page.Props["deferredData"])
	}
}

func TestOncePropsHandler_ValidPage(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)
	req := h.request(http.MethodGet, "/features/data-loading/once-props/3", nil)

	req.SetPathValue("page", "3")

	w := httptest.NewRecorder()

	h.app.oncePropsHandler(w, req)

	page := h.page(t, w)

	page.AssertPropEquals(t, "page", float64(3))
}

func TestPartialReloadHandlersSupportInertiaPartials(t *testing.T) {
	t.Parallel()

	h := newFeatureHarness(t)
	req := h.partialRequest(http.MethodGet, "/features/data-loading/partial-reloads", "Features/DataLoading/PartialReloads", "stats")
	w := httptest.NewRecorder()

	h.app.partialReloadsHandler(w, req)

	page := h.page(t, w)

	page.AssertHasProp(t, "stats")

	req = h.partialRequest(http.MethodGet, "/features/data-loading/polling", "Features/DataLoading/Polling", "contactCount")
	w = httptest.NewRecorder()

	h.app.pollingHandler(w, req)

	page = h.page(t, w)

	page.AssertHasProp(t, "contactCount")

	if strings.TrimSpace(req.Header.Get(protocol.HeaderPartialComponent)) == "" {
		t.Fatal("partial headers should be set")
	}
}
