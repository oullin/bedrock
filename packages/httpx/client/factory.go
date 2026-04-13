package client

import (
	"net/http"
	"strings"
	"sync"
)

// Factory is the central hub for creating PendingRequest instances and managing
// test fakes, stubs and request recording.
type Factory struct {
	mu               sync.Mutex
	baseURL          string
	globalMiddleware []Middleware
	recording        bool
	recorded         []RecordedRequest
	faking           bool
	stubCallbacks    []StubCallback
	sequences        map[string]*ResponseSequence
	preventStray     bool
	dispatcher       *EventDispatcher
}

// RecordedRequest holds a recorded outbound request and its response.
type RecordedRequest struct {
	Request  *Request
	Response *Response
}

// StubCallback maps a request to a stubbed response. Return nil to pass
// through to the next stub or to the real transport.
type StubCallback func(req *http.Request) *http.Response

// NewFactory creates a new Factory.
func NewFactory() *Factory {
	return &Factory{
		sequences: make(map[string]*ResponseSequence),
	}
}

// BaseURL sets a base URL for all requests created by this factory.
func (f *Factory) BaseURL(url string) *Factory {
	f.baseURL = strings.TrimRight(url, "/")

	return f
}

// WithDispatcher sets the event dispatcher for the factory. Events like
// RequestSending, ResponseReceived, and ConnectionFailed will be dispatched
// through it.
func (f *Factory) WithDispatcher(d *EventDispatcher) *Factory {
	f.dispatcher = d

	return f
}

// GlobalMiddleware adds middleware applied to every request.
func (f *Factory) GlobalMiddleware(mw ...Middleware) *Factory {
	f.globalMiddleware = append(f.globalMiddleware, mw...)

	return f
}

// PendingRequest creates a new PendingRequest.
func (f *Factory) PendingRequest() *PendingRequest {
	p := newPendingRequest(f)

	if f.baseURL != "" {
		p.BaseURL(f.baseURL)
	}

	p.middleware = append(p.middleware, f.globalMiddleware...)

	return p
}

// Fake puts the factory in fake mode. All requests will be intercepted.
func (f *Factory) Fake(stubs ...StubCallback) *Factory {
	f.mu.Lock()

	defer f.mu.Unlock()

	f.faking = true
	f.recording = true
	f.stubCallbacks = append(f.stubCallbacks, stubs...)

	return f
}

// FakeSequence registers a response sequence for a URL pattern.
func (f *Factory) FakeSequence(urlPattern string, seq *ResponseSequence) *Factory {
	f.mu.Lock()

	defer f.mu.Unlock()

	f.faking = true
	f.recording = true
	f.sequences[urlPattern] = seq

	return f
}

// PreventStrayRequests causes the factory to return an error for any request
// that doesn't match a stub.
func (f *Factory) PreventStrayRequests() *Factory {
	f.mu.Lock()

	defer f.mu.Unlock()

	f.preventStray = true

	return f
}

// isFaking returns true when the factory is in fake mode.
func (f *Factory) isFaking() bool {
	f.mu.Lock()

	defer f.mu.Unlock()

	return f.faking
}

// handleFake processes a request through the stub/sequence chain.
func (f *Factory) handleFake(req *http.Request, body []byte) (*Response, error) {
	f.mu.Lock()

	defer f.mu.Unlock()

	// Check sequences first.
	for pattern, seq := range f.sequences {
		if matchesURLPattern(pattern, req.URL.String()) {
			raw := seq.Next()
			resp := NewResponse(raw)

			f.recorded = append(f.recorded, RecordedRequest{
				Request:  NewRequest(req, body),
				Response: resp,
			})

			return resp, nil
		}
	}

	// Check stub callbacks.
	for _, stub := range f.stubCallbacks {
		if raw := stub(req); raw != nil {
			resp := NewResponse(raw)

			f.recorded = append(f.recorded, RecordedRequest{
				Request:  NewRequest(req, body),
				Response: resp,
			})

			return resp, nil
		}
	}

	// Default: return 200 empty or error if preventing stray.
	if f.preventStray {
		return nil, &ConnectionError{
			URL: req.URL.String(),
			Err: ErrStrayRequest,
		}
	}

	raw := stubToResponse(ResponseStub{Status: http.StatusOK})
	resp := NewResponse(raw)

	f.recorded = append(f.recorded, RecordedRequest{
		Request:  NewRequest(req, body),
		Response: resp,
	})

	return resp, nil
}

// Recorded returns all recorded requests.
func (f *Factory) Recorded(filter ...func(RecordedRequest) bool) []RecordedRequest {
	f.mu.Lock()

	defer f.mu.Unlock()

	if len(filter) == 0 {
		return f.recorded
	}

	var result []RecordedRequest

	for _, r := range f.recorded {
		if filter[0](r) {
			result = append(result, r)
		}
	}

	return result
}

// AssertSent asserts that at least one request matches the callback.
func (f *Factory) AssertSent(fn func(RecordedRequest) bool) bool {
	for _, r := range f.Recorded() {
		if fn(r) {
			return true
		}
	}

	return false
}

// AssertNotSent asserts that no request matches the callback.
func (f *Factory) AssertNotSent(fn func(RecordedRequest) bool) bool {
	for _, r := range f.Recorded() {
		if fn(r) {
			return false
		}
	}

	return true
}

// AssertNothingSent asserts that no requests were recorded.
func (f *Factory) AssertNothingSent() bool {
	return len(f.Recorded()) == 0
}

// AssertSentCount asserts that exactly n requests were recorded.
func (f *Factory) AssertSentCount(n int) bool {
	return len(f.Recorded()) == n
}

// AssertSentInOrder asserts that recorded requests match the given callbacks in
// order. Returns true only when every callback matches its corresponding
// recorded request.
func (f *Factory) AssertSentInOrder(callbacks []func(RecordedRequest) bool) bool {
	recorded := f.Recorded()

	if len(callbacks) != len(recorded) {
		return false
	}

	for i, cb := range callbacks {
		if !cb(recorded[i]) {
			return false
		}
	}

	return true
}

// AllowStrayRequests disables the stray request prevention, allowing requests
// that don't match any stub.
func (f *Factory) AllowStrayRequests() *Factory {
	f.mu.Lock()

	defer f.mu.Unlock()

	f.preventStray = false

	return f
}

// Response creates a stubbed *http.Response for use in test fakes.
func (f *Factory) Response(body string, status int, headers ...map[string]string) *http.Response {
	stub := ResponseStub{
		Status: status,
		Body:   body,
	}

	if len(headers) > 0 {
		stub.Headers = headers[0]
	}

	return stubToResponse(stub)
}

// Sequence creates a new ResponseSequence from the given stubs.
func (f *Factory) Sequence(stubs ...ResponseStub) *ResponseSequence {
	return NewResponseSequence(stubs...)
}

// matchesURLPattern checks if a URL matches a simple pattern. Supports "*" as
// a wildcard.
func matchesURLPattern(pattern, u string) bool {
	if pattern == "*" {
		return true
	}

	if strings.Contains(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")

		return strings.HasPrefix(u, prefix)
	}

	return strings.Contains(u, pattern)
}
