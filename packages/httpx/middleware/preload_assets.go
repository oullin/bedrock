package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// PreloadAsset describes a resource to preload via the Link header.
type PreloadAsset struct {
	URI string
	As  string // "script", "style", "font", "image", etc.
}

// AddLinkHeadersForPreloadedAssets appends Link headers for assets that should
// be preloaded by the browser.
type AddLinkHeadersForPreloadedAssets struct {
	assets []PreloadAsset
}

// NewAddLinkHeadersForPreloadedAssets creates the middleware with the given
// assets to preload.
func NewAddLinkHeadersForPreloadedAssets(assets ...PreloadAsset) *AddLinkHeadersForPreloadedAssets {
	return &AddLinkHeadersForPreloadedAssets{assets: assets}
}

// Wrap returns an http.Handler that sets Link headers before delegating.
func (m *AddLinkHeadersForPreloadedAssets) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var links []string

		for _, asset := range m.assets {
			link := fmt.Sprintf("<%s>; rel=preload", asset.URI)

			if asset.As != "" {
				link += fmt.Sprintf("; as=%s", asset.As)
			}

			links = append(links, link)
		}

		if len(links) > 0 {
			existing := w.Header().Get("Link")

			if existing != "" {
				w.Header().Set("Link", existing+", "+strings.Join(links, ", "))
			} else {
				w.Header().Set("Link", strings.Join(links, ", "))
			}
		}

		next.ServeHTTP(w, r)
	})
}
