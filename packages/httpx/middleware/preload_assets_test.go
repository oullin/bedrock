package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx/middleware"
)

func TestPreloadAssetsLinkHeader(t *testing.T) {
	t.Parallel()

	handler := middleware.NewAddLinkHeadersForPreloadedAssets(
		middleware.PreloadAsset{URI: "/css/app.css", As: "style"},
		middleware.PreloadAsset{URI: "/js/app.js", As: "script"},
	).Wrap(okHandler)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	link := rec.Header().Get("Link")

	if !strings.Contains(link, "</css/app.css>; rel=preload; as=style") {
		t.Fatalf("expected CSS preload link, got %s", link)
	}

	if !strings.Contains(link, "</js/app.js>; rel=preload; as=script") {
		t.Fatalf("expected JS preload link, got %s", link)
	}
}

func TestPreloadAssetsNoAssets(t *testing.T) {
	t.Parallel()

	handler := middleware.NewAddLinkHeadersForPreloadedAssets().Wrap(okHandler)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Header().Get("Link") != "" {
		t.Fatal("expected no Link header when no assets configured")
	}
}

func TestPreloadAssetsWithoutAs(t *testing.T) {
	t.Parallel()

	handler := middleware.NewAddLinkHeadersForPreloadedAssets(
		middleware.PreloadAsset{URI: "/font.woff2"},
	).Wrap(okHandler)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	link := rec.Header().Get("Link")

	if !strings.Contains(link, "</font.woff2>; rel=preload") {
		t.Fatalf("expected font preload link, got %s", link)
	}

	if strings.Contains(link, "as=") {
		t.Fatal("expected no 'as' parameter when not specified")
	}
}
