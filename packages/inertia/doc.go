// Package inertia is the server-side Go adapter for the Inertia.js protocol.
// It renders Inertia pages (JSON on XHR visits, HTML on the first request),
// merges shared and per-request props, manages the head (title, meta, links),
// and integrates with CSRF, i18n, precognition, and flash middleware.
//
// This package is ported from github.com/oullin/inertia-go (MIT) and adapted
// for the Bedrock monorepo. The upstream httpx sub-package is renamed here to
// protocol to avoid colliding with github.com/bedrock/packages/httpx, which is
// a HTTP Request/Response wrapper for a different purpose.
package inertia
