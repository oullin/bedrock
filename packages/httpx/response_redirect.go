package httpx

import (
	"net/http"
	"net/url"
	"strings"
)

// RedirectResponse sends an HTTP redirect with optional session flash data,
// input persistence, and fragment handling.
type RedirectResponse struct {
	writer      http.ResponseWriter
	request     *http.Request
	statusCode  int
	headers     http.Header
	cookies     []*http.Cookie
	targetURL   string
	session     SessionStore
	exception   error
	flashData   map[string]any
	withInput   bool
	inputOnly   []string
	inputExcept []string
	fragment    string
}

// NewRedirectResponse creates a redirect to the given URL.
func NewRedirectResponse(w http.ResponseWriter, r *http.Request, target string, status int) *RedirectResponse {
	return &RedirectResponse{
		writer:     w,
		request:    r,
		statusCode: status,
		headers:    make(http.Header),
		targetURL:  target,
		flashData:  make(map[string]any),
	}
}

// Status overrides the HTTP redirect status code.
func (r *RedirectResponse) Status(code int) *RedirectResponse {
	r.statusCode = code

	return r
}

// GetStatusCode returns the current status code.
func (r *RedirectResponse) GetStatusCode() int {
	return r.statusCode
}

// GetTargetURL returns the redirect target URL.
func (r *RedirectResponse) GetTargetURL() string {
	target := r.targetURL

	if r.fragment != "" {
		target += "#" + r.fragment
	}

	return target
}

// Header sets a single response header.
func (r *RedirectResponse) Header(key, value string) *RedirectResponse {
	r.headers.Set(key, value)

	return r
}

// WithHeaders sets multiple response headers.
func (r *RedirectResponse) WithHeaders(headers map[string]string) *RedirectResponse {
	for k, v := range headers {
		r.headers.Set(k, v)
	}

	return r
}

// Cookie appends a cookie to the response.
func (r *RedirectResponse) Cookie(c *http.Cookie) *RedirectResponse {
	r.cookies = append(r.cookies, c)

	return r
}

// With flashes a key/value pair to the session.
func (r *RedirectResponse) With(key string, value any) *RedirectResponse {
	r.flashData[key] = value

	return r
}

// WithData flashes multiple key/value pairs to the session.
func (r *RedirectResponse) WithData(data map[string]any) *RedirectResponse {
	for k, v := range data {
		r.flashData[k] = v
	}

	return r
}

// WithInput flashes all current request input to the session.
func (r *RedirectResponse) WithInput() *RedirectResponse {
	r.withInput = true

	return r
}

// OnlyInput flashes only the specified input keys.
func (r *RedirectResponse) OnlyInput(keys ...string) *RedirectResponse {
	r.withInput = true
	r.inputOnly = keys

	return r
}

// ExceptInput flashes all input except the specified keys.
func (r *RedirectResponse) ExceptInput(keys ...string) *RedirectResponse {
	r.withInput = true
	r.inputExcept = keys

	return r
}

// WithErrors flashes validation errors to the session under the "errors" key.
func (r *RedirectResponse) WithErrors(errors map[string][]string) *RedirectResponse {
	r.flashData["errors"] = errors

	return r
}

// WithFragment appends a URL fragment (#section) to the redirect target.
func (r *RedirectResponse) WithFragment(fragment string) *RedirectResponse {
	r.fragment = fragment

	return r
}

// WithoutFragment removes any previously set fragment.
func (r *RedirectResponse) WithoutFragment() *RedirectResponse {
	r.fragment = ""

	return r
}

// SetSession attaches a session store for flash data persistence.
func (r *RedirectResponse) SetSession(s SessionStore) *RedirectResponse {
	r.session = s

	return r
}

// WithException attaches an error for later inspection.
func (r *RedirectResponse) WithException(err error) *RedirectResponse {
	r.exception = err

	return r
}

// EnforceSameOrigin checks that the redirect target's origin matches the
// request origin. If it does not, an ErrOriginMismatch error is returned and
// no redirect is performed.
func (r *RedirectResponse) EnforceSameOrigin() error {
	parsed, err := url.Parse(r.targetURL)

	if err != nil {
		return ErrMalformedURL
	}

	// Relative URLs are always same-origin.
	if parsed.Host == "" {
		return nil
	}

	if r.request != nil && parsed.Host != r.request.Host {
		return ErrOriginMismatch
	}

	return nil
}

// Away creates a redirect to an external URL without same-origin checks.
func (r *RedirectResponse) Away(target string, status ...int) *RedirectResponse {
	r.targetURL = target

	if len(status) > 0 {
		r.statusCode = status[0]
	}

	return r
}

// Send writes the redirect response including flash data.
func (r *RedirectResponse) Send() error {
	// Flash session data.
	if r.session != nil {
		for k, v := range r.flashData {
			r.session.Flash(k, v)
		}

		if r.withInput {
			if r.request != nil {
				req := NewRequest(r.request)

				var input map[string]any

				switch {
				case len(r.inputOnly) > 0:
					input = req.Only(r.inputOnly...)
				case len(r.inputExcept) > 0:
					input = req.Except(r.inputExcept...)
				default:
					input = req.All()
				}

				r.session.FlashInput(input)
			}
		}
	}

	// Build the target URL.
	target := r.GetTargetURL()

	// Flush headers.
	dest := r.writer.Header()

	for k, vals := range r.headers {
		for _, v := range vals {
			dest.Add(k, v)
		}
	}

	for _, c := range r.cookies {
		http.SetCookie(r.writer, c)
	}

	dest.Set("Location", target)
	r.writer.WriteHeader(r.statusCode)

	return nil
}

// Secure creates a redirect to the HTTPS version of the given path.
func Secure(w http.ResponseWriter, r *http.Request, path string, status ...int) *RedirectResponse {
	code := http.StatusFound

	if len(status) > 0 {
		code = status[0]
	}

	scheme := "https"
	host := r.Host

	target := scheme + "://" + host

	if !strings.HasPrefix(path, "/") {
		target += "/"
	}

	target += path

	return NewRedirectResponse(w, r, target, code)
}
