package protocol

// Inertia.js protocol HTTP headers.
const (
	HeaderInertia          = "X-Inertia"
	HeaderVersion          = "X-Inertia-Version"
	HeaderPartialComponent = "X-Inertia-Partial-Component"
	HeaderPartialData      = "X-Inertia-Partial-Data"
	HeaderPartialExcept    = "X-Inertia-Partial-Except"
	HeaderExceptOnceProps  = "X-Inertia-Except-Once-Props"
	HeaderInfiniteScroll   = "X-Inertia-Infinite-Scroll-Merge-Intent"
	HeaderLocation         = "X-Inertia-Location"
	HeaderRedirect         = "X-Inertia-Redirect"
	HeaderErrorBag         = "X-Inertia-Error-Bag"
	HeaderReset            = "X-Inertia-Reset"

	// HTTPPreview protocol headers.
	HeaderHTTPPreview        = "HTTPPreview"
	HeaderHTTPPreviewSuccess = HeaderHTTPPreview
	HeaderValidateOnly        = "Validate-Only"
)
