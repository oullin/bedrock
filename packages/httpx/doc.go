// Package httpx provides HTTP request and response primitives
// built on Go's net/http. It wraps *http.Request with rich input, content-type,
// flash-data and precognitive helpers, offers fluent Response, JsonResponse,
// RedirectResponse and StreamedEvent writers, and includes file-upload handling
// with pluggable storage. Sub-packages supply HTTP middleware, an outbound HTTP
// client with testing fakes, JSON API resources, and test-double utilities.
package httpx
