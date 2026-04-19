package httppreview

import (
	"net/http"
	"net/http/httptest"

	"github.com/bedrock/packages/container"
	routeContracts "github.com/bedrock/packages/routing/contracts"
)

// HandlePrecognitiveRequests is middleware that intercepts precognitive HTTP
// requests. It marks the request as precognitive, swaps route dispatchers to
// httppreview versions (when a container is provided), and manages the
// HTTPPreview and Vary response headers.
//
// Mirrors Framework\Foundation\Http\Middleware\HandlePrecognitiveRequests.
type HandlePrecognitiveRequests struct {
	container *container.Container
}

// New creates a HandlePrecognitiveRequests middleware. An optional container
// enables dispatcher swapping; without it the middleware relies on
// [AfterValidationHook] being called explicitly by handlers.
func New(c ...*container.Container) *HandlePrecognitiveRequests {
	m := &HandlePrecognitiveRequests{}

	if len(c) > 0 {
		m.container = c[0]
	}

	return m
}

// Wrap returns an http.Handler that handles precognitive requests.
func (m *HandlePrecognitiveRequests) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAttemptingHTTPPreview(r) {
			next.ServeHTTP(w, r)
			appendVaryHeader(w)

			return
		}

		// Save original dispatchers and swap to httppreview versions.
		var restoreDispatchers func()

		if m.container != nil {
			restoreDispatchers = m.prepareForHTTPPreview(r)
		}

		r = MarkPrecognitive(r)

		// Capture the downstream response so we can inject headers.
		rec := httptest.NewRecorder()
		panicked := true

		func() {
			defer func() {
				if v := recover(); v != nil {
					if _, ok := v.(SuccessResponse); ok {
						panicked = false
						WriteSuccessResponse(w)
						AddHTTPPreviewHeader(w)
						appendVaryHeader(w)

						return
					}

					// Re-panic for unexpected panics.
					panic(v)
				}

				panicked = false
			}()

			next.ServeHTTP(rec, r)
		}()

		if restoreDispatchers != nil {
			restoreDispatchers()
		}

		if panicked {
			return
		}

		// If the handler completed without a SuccessResponse panic, forward
		// the captured response with httppreview headers injected.
		if rec.Code != 0 {
			copyHeaders(w, rec)
			AddHTTPPreviewHeader(w)
			appendVaryHeader(w)
			w.WriteHeader(rec.Code)
			w.Write(rec.Body.Bytes())
		}
	})
}

// prepareForHTTPPreview swaps the route dispatchers in the container to
// httppreview versions and returns a function that restores the originals.
func (m *HandlePrecognitiveRequests) prepareForHTTPPreview(r *http.Request) func() {
	c := m.container

	// Save originals.
	origCallable, _ := c.Make("routing.callable_dispatcher")
	origController, _ := c.Make("routing.controller_dispatcher")

	// Swap to httppreview dispatchers.
	c.Instance("routing.callable_dispatcher", NewCallableDispatcher(nil))
	c.Instance("routing.controller_dispatcher", NewControllerDispatcher(nil))

	return func() {
		if origCallable != nil {
			c.Instance("routing.callable_dispatcher", origCallable)
		}

		if origController != nil {
			c.Instance("routing.controller_dispatcher", origController)
		}
	}
}

// appendVaryHeader adds "HTTPPreview" to the Vary header. This is called for
// both precognitive and non-precognitive responses, matching Upstream's
// behaviour where the Vary header is always set.
func appendVaryHeader(w http.ResponseWriter) {
	AddVaryHeader(w)
}

// copyHeaders copies all headers from a recorder to a response writer.
func copyHeaders(w http.ResponseWriter, rec *httptest.ResponseRecorder) {
	for k, vals := range rec.Header() {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}
}

// Compile-time assertion that dispatchers satisfy the contracts.
var (
	_ routeContracts.CallableDispatcher   = (*CallableDispatcher)(nil)
	_ routeContracts.ControllerDispatcher = (*ControllerDispatcher)(nil)
)
