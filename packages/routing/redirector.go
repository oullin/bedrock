package routing

import "fmt"

// RedirectResponse is a minimal redirect-response value type used by the
// routing port. The bedrock httpx package supplies a richer type that the
// service provider in M11 wires in via [Redirector.SetResponseFactory].
//
// Ref: @bedrock/code-0223
type RedirectResponse struct {
	URL     string
	Status  int
	Headers map[string][]string
	Session SessionStore
	with    map[string]any
}

// With attaches a flash data entry to the response for the next request.

// Send applies any pending flash data to the session and returns the URL.
// In real use the bedrock httpx layer drives the actual HTTP write — the
// Send call here is a parity shim.

// SessionStore is the minimum session surface Redirector touches.
// httpx.SessionStore narrowly.
type SessionStore interface {
	Flash(key string, value any)
	GetOldInput(key, fallback string) string
	HasOldInput(key string) bool
	FlashInput(input map[string]any)
	Get(key string, fallback any) any
	Put(key string, value any)
}

// Ref: @bedrock/code-0329
// [RedirectResponse] values for the various redirect verbs (To, Away, Secure,
// Back, Refresh, Route, Action).
type Redirector struct {
	generator   *UrlGenerator
	session     SessionStore
	intendedUrl string
}

func (r *RedirectResponse) With(key string, value any) *RedirectResponse {
	if r.with == nil {
		r.with = map[string]any{}
	}

	r.with[key] = value

	return r
}

func (r *RedirectResponse) Send() error {
	if r.Session != nil && r.with != nil {
		for k, v := range r.with {
			r.Session.Flash(k, v)
		}
	}

	return nil
}

// NewRedirector constructs a redirector bound to a URL generator.
func NewRedirector(generator *UrlGenerator) *Redirector {
	return &Redirector{generator: generator}
}

// SetSession wires the session store used by [Redirector.Back],
// [Redirector.Intended] and flash data attached via [RedirectResponse.With].
func (r *Redirector) SetSession(s SessionStore) { r.session = s }

// GetUrlGenerator returns the underlying URL generator.
func (r *Redirector) GetUrlGenerator() *UrlGenerator { return r.generator }

// To redirects to the given path.
func (r *Redirector) To(path string, status int, headers map[string][]string, secure *bool) *RedirectResponse {
	if status == 0 {
		status = 302
	}

	return r.createResponse(r.generator.To(path, nil, secure), status, headers)
}

// Away redirects to an arbitrary external URL without validation.
func (r *Redirector) Away(path string, status int, headers map[string][]string) *RedirectResponse {
	if status == 0 {
		status = 302
	}

	return r.createResponse(path, status, headers)
}

// Secure redirects to the given path forcing HTTPS.
func (r *Redirector) Secure(path string, status int, headers map[string][]string) *RedirectResponse {
	t := true

	return r.To(path, status, headers, &t)
}

// Back redirects to the URL recorded in the session as the previous URL.
// fallback is used when the session has no previous URL.
func (r *Redirector) Back(status int, headers map[string][]string, fallback string) *RedirectResponse {
	if status == 0 {
		status = 302
	}

	target := fallback

	if r.session != nil {
		if v := r.session.Get("_previous.url", nil); v != nil {
			if s, ok := v.(string); ok && s != "" {
				target = s
			}
		}
	}

	if target == "" {
		target = "/"
	}

	return r.createResponse(target, status, headers)
}

// Refresh redirects to the request's current URL (used after form submissions).
func (r *Redirector) Refresh(status int, headers map[string][]string) *RedirectResponse {
	if status == 0 {
		status = 302
	}

	target := ""

	if r.generator != nil {
		target = r.generator.Current()
	}

	return r.createResponse(target, status, headers)
}

// Intended redirects to the URL the user attempted to access before
// authenticating, falling back to default.
func (r *Redirector) Intended(defaultPath string, status int, headers map[string][]string, secure *bool) *RedirectResponse {
	target := defaultPath

	if r.intendedUrl != "" {
		target = r.intendedUrl
		r.intendedUrl = ""
	}

	return r.To(target, status, headers, secure)
}

// SetIntendedUrl records the URL the user attempted before being redirected
// to login.
func (r *Redirector) SetIntendedUrl(url string) { r.intendedUrl = url }

// GetIntendedUrl returns the previously stored intended URL.
func (r *Redirector) GetIntendedUrl() string { return r.intendedUrl }

// Route redirects to a named route.
func (r *Redirector) Route(name string, parameters map[string]any, status int, headers map[string][]string) (*RedirectResponse, error) {
	if status == 0 {
		status = 302
	}

	target, err := r.generator.Route(name, parameters, true)

	if err != nil {
		return nil, err
	}

	return r.createResponse(target, status, headers), nil
}

// SignedRoute redirects to a signed named route.
func (r *Redirector) SignedRoute(name string, parameters map[string]any, expiration int64, status int, headers map[string][]string) (*RedirectResponse, error) {
	if status == 0 {
		status = 302
	}

	target, err := r.generator.SignedRoute(name, parameters, expiration, true)

	if err != nil {
		return nil, err
	}

	return r.createResponse(target, status, headers), nil
}

// TemporarySignedRoute is the parity-named two-arg shortcut.
func (r *Redirector) TemporarySignedRoute(name string, expiration int64, parameters map[string]any, status int, headers map[string][]string) (*RedirectResponse, error) {
	return r.SignedRoute(name, parameters, expiration, status, headers)
}

// Action redirects to a controller action — the parity entry point used by
// "Controller@method" style references.
func (r *Redirector) Action(action string, parameters map[string]any, status int, headers map[string][]string) (*RedirectResponse, error) {
	return nil, fmt.Errorf("redirect to action [%s]: action lookup not yet implemented", action)
}

func (r *Redirector) createResponse(url string, status int, headers map[string][]string) *RedirectResponse {
	return &RedirectResponse{
		URL:     url,
		Status:  status,
		Headers: headers,
		Session: r.session,
	}
}
