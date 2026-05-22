package routing

// Ref: @bedrock/code-0331
// In the upstream framework the factory hand-builds Response/JsonResponse instances; the Go
// port provides a struct-based equivalent that the bedrock httpx layer can
// translate into real HTTP responses. The methods cover the public surface
// callers depend on; the internal "prepare" plumbing lives in M11.
type ResponseFactory struct {
	view       ViewFactory
	redirector *Redirector
}

// NewResponseFactory constructs a factory wired to the view layer and the
// redirector.

// HTTPResponse is the value-object form a response factory produces. The
// real httpx.Response will replace this in M11.
type HTTPResponse struct {
	Body    any
	Status  int
	Headers map[string][]string
}

func NewResponseFactory(view ViewFactory, redirector *Redirector) *ResponseFactory {
	return &ResponseFactory{view: view, redirector: redirector}
}

// Make builds a generic Response wrapping content.
func (f *ResponseFactory) Make(content any, status int, headers map[string][]string) *HTTPResponse {
	if status == 0 {
		status = 200
	}

	return &HTTPResponse{Body: content, Status: status, Headers: headers}
}

// NoContent returns a 204 (or supplied status) response with no body.
func (f *ResponseFactory) NoContent(status int, headers map[string][]string) *HTTPResponse {
	if status == 0 {
		status = 204
	}

	return &HTTPResponse{Status: status, Headers: headers}
}

// View renders a named view via the bound [ViewFactory] and wraps it in a
// response.
func (f *ResponseFactory) View(view string, data map[string]any, status int, headers map[string][]string) *HTTPResponse {
	if status == 0 {
		status = 200
	}

	body := any(nil)

	if f.view != nil {
		body = f.view.Make(view, data)
	}

	return &HTTPResponse{Body: body, Status: status, Headers: headers}
}

// JSON wraps payload in a JSON response. Encoding lives in httpx; this
// factory just packages the value.
func (f *ResponseFactory) JSON(data any, status int, headers map[string][]string) *HTTPResponse {
	if status == 0 {
		status = 200
	}

	if headers == nil {
		headers = map[string][]string{}
	}

	headers["Content-Type"] = []string{"application/json"}

	return &HTTPResponse{Body: data, Status: status, Headers: headers}
}

// RedirectTo delegates to the bound [Redirector].
func (f *ResponseFactory) RedirectTo(path string, status int, headers map[string][]string, secure *bool) *RedirectResponse {
	return f.redirector.To(path, status, headers, secure)
}

// RedirectToRoute delegates to the bound [Redirector].
func (f *ResponseFactory) RedirectToRoute(name string, parameters map[string]any, status int, headers map[string][]string) (*RedirectResponse, error) {
	return f.redirector.Route(name, parameters, status, headers)
}

// RedirectGuest sets the intended URL on the session and redirects to path.
func (f *ResponseFactory) RedirectGuest(path string, status int, headers map[string][]string, secure *bool) *RedirectResponse {
	if f.redirector.session != nil {
		f.redirector.session.Put("url.intended", "")
	}

	return f.redirector.To(path, status, headers, secure)
}

// RedirectToIntended redirects to the intended URL or default fallback.
func (f *ResponseFactory) RedirectToIntended(defaultPath string, status int, headers map[string][]string, secure *bool) *RedirectResponse {
	return f.redirector.Intended(defaultPath, status, headers, secure)
}
