package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// BodyFormat describes the request body encoding.
type BodyFormat int

const (
	BodyJSON BodyFormat = iota
	BodyForm
	BodyMultipart
	BodyRaw
)

// PendingRequest is a fluent builder for outbound HTTP requests.
type PendingRequest struct {
	factory    *Factory
	httpClient *http.Client
	baseURL    string
	bodyFormat BodyFormat
	headers    http.Header
	cookies    []*http.Cookie
	timeout    time.Duration
	retries    int
	retryDelay time.Duration
	retryWhen  func(error, *Response) bool
	middleware []Middleware
	body       io.Reader
	bodyBytes  []byte
	ctx        context.Context
}

func newPendingRequest(f *Factory) *PendingRequest {
	return &PendingRequest{
		factory:    f,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		headers:    make(http.Header),
		bodyFormat: BodyJSON,
		timeout:    30 * time.Second,
		retries:    1,
		retryDelay: 100 * time.Millisecond,
		ctx:        context.Background(),
	}
}

// BaseURL sets the base URL prepended to relative paths.
func (p *PendingRequest) BaseURL(url string) *PendingRequest {
	p.baseURL = strings.TrimRight(url, "/")

	return p
}

// AsJSON sets the body format to JSON.
func (p *PendingRequest) AsJSON() *PendingRequest {
	p.bodyFormat = BodyJSON
	p.headers.Set("Content-Type", "application/json")

	return p
}

// AsForm sets the body format to form-encoded.
func (p *PendingRequest) AsForm() *PendingRequest {
	p.bodyFormat = BodyForm
	p.headers.Set("Content-Type", "application/x-www-form-urlencoded")

	return p
}

// AsMultipart sets the body format to multipart.
func (p *PendingRequest) AsMultipart() *PendingRequest {
	p.bodyFormat = BodyMultipart

	return p
}

// WithBody sets a raw body.
func (p *PendingRequest) WithBody(body string, contentType string) *PendingRequest {
	p.bodyFormat = BodyRaw
	p.bodyBytes = []byte(body)

	if contentType != "" {
		p.headers.Set("Content-Type", contentType)
	}

	return p
}

// WithHeaders sets multiple headers.
func (p *PendingRequest) WithHeaders(headers map[string]string) *PendingRequest {
	for k, v := range headers {
		p.headers.Set(k, v)
	}

	return p
}

// WithHeader sets a single header.
func (p *PendingRequest) WithHeader(key, value string) *PendingRequest {
	p.headers.Set(key, value)

	return p
}

// Accept sets the Accept header.
func (p *PendingRequest) Accept(contentType string) *PendingRequest {
	p.headers.Set("Accept", contentType)

	return p
}

// AcceptJSON sets the Accept header to application/json.
func (p *PendingRequest) AcceptJSON() *PendingRequest {
	return p.Accept("application/json")
}

// WithToken sets a Bearer token.
func (p *PendingRequest) WithToken(token string, tokenType ...string) *PendingRequest {
	tt := "Bearer"

	if len(tokenType) > 0 {
		tt = tokenType[0]
	}

	p.headers.Set("Authorization", tt+" "+token)

	return p
}

// WithBasicAuth sets Basic authentication.
func (p *PendingRequest) WithBasicAuth(user, password string) *PendingRequest {
	p.headers.Set("Authorization", "Basic "+basicAuth(user, password))

	return p
}

// WithCookies adds cookies to the request.
func (p *PendingRequest) WithCookies(cookies []*http.Cookie) *PendingRequest {
	p.cookies = append(p.cookies, cookies...)

	return p
}

// Timeout sets the request timeout.
func (p *PendingRequest) Timeout(d time.Duration) *PendingRequest {
	p.timeout = d
	p.httpClient.Timeout = d

	return p
}

// Retry configures automatic retries.
func (p *PendingRequest) Retry(times int, sleep time.Duration, when ...func(error, *Response) bool) *PendingRequest {
	p.retries = times
	p.retryDelay = sleep

	if len(when) > 0 {
		p.retryWhen = when[0]
	}

	return p
}

// WithMiddleware adds client middleware.
func (p *PendingRequest) WithMiddleware(mw ...Middleware) *PendingRequest {
	p.middleware = append(p.middleware, mw...)

	return p
}

// WithContext sets the request context.
func (p *PendingRequest) WithContext(ctx context.Context) *PendingRequest {
	p.ctx = ctx

	return p
}

// WithoutRedirecting disables automatic following of redirects.
func (p *PendingRequest) WithoutRedirecting() *PendingRequest {
	p.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return p
}

// MaxRedirects limits the number of redirects to follow.
func (p *PendingRequest) MaxRedirects(max int) *PendingRequest {
	p.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= max {
			return http.ErrUseLastResponse
		}

		return nil
	}

	return p
}

// Get sends a GET request.
func (p *PendingRequest) Get(url string, query ...map[string]string) (*Response, error) {
	if len(query) > 0 {
		url = p.appendQuery(url, query[0])
	}

	return p.send(http.MethodGet, url, nil)
}

// Head sends a HEAD request.
func (p *PendingRequest) Head(url string) (*Response, error) {
	return p.send(http.MethodHead, url, nil)
}

// Post sends a POST request.
func (p *PendingRequest) Post(url string, data ...any) (*Response, error) {
	return p.send(http.MethodPost, url, firstOrNil(data))
}

// Put sends a PUT request.
func (p *PendingRequest) Put(url string, data ...any) (*Response, error) {
	return p.send(http.MethodPut, url, firstOrNil(data))
}

// Patch sends a PATCH request.
func (p *PendingRequest) Patch(url string, data ...any) (*Response, error) {
	return p.send(http.MethodPatch, url, firstOrNil(data))
}

// Delete sends a DELETE request.
func (p *PendingRequest) Delete(url string, data ...any) (*Response, error) {
	return p.send(http.MethodDelete, url, firstOrNil(data))
}

// Options sends an OPTIONS request.
func (p *PendingRequest) Options(url string) (*Response, error) {
	return p.send(http.MethodOptions, url, nil)
}

func (p *PendingRequest) send(method, requestURL string, data any) (*Response, error) {
	fullURL := p.buildURL(requestURL)
	body, contentType, err := p.encodeBody(data)

	if err != nil {
		return nil, err
	}

	if contentType != "" && p.headers.Get("Content-Type") == "" {
		p.headers.Set("Content-Type", contentType)
	}

	var resp *Response
	var lastErr error

	for attempt := 0; attempt < p.retries; attempt++ {
		if attempt > 0 {
			time.Sleep(p.retryDelay)
		}

		var bodyReader io.Reader

		if body != nil {
			bodyReader = bytes.NewReader(body)
		}

		req, err := http.NewRequestWithContext(p.ctx, method, fullURL, bodyReader)

		if err != nil {
			return nil, err
		}

		// Set headers.
		for k, vals := range p.headers {
			for _, v := range vals {
				req.Header.Add(k, v)
			}
		}

		// Set cookies.
		for _, c := range p.cookies {
			req.AddCookie(c)
		}

		// Check if factory has a fake/stub.
		if p.factory != nil && p.factory.isFaking() {
			resp, lastErr = p.factory.handleFake(req, body)
		} else {
			// Execute through middleware chain.
			transport := p.buildTransport()
			rawResp, err := transport(req)

			if err != nil {
				lastErr = &ConnectionError{URL: fullURL, Err: err}

				if p.retryWhen != nil && !p.retryWhen(lastErr, nil) {
					return nil, lastErr
				}

				continue
			}

			resp = NewResponse(rawResp)
		}

		if lastErr != nil {
			if p.retryWhen != nil && !p.retryWhen(lastErr, resp) {
				return resp, lastErr
			}

			continue
		}

		if resp != nil && resp.Failed() && attempt < p.retries-1 {
			if p.retryWhen != nil && !p.retryWhen(nil, resp) {
				return resp, nil
			}

			lastErr = resp.Throw()

			continue
		}

		return resp, nil
	}

	return resp, lastErr
}

func (p *PendingRequest) buildURL(requestURL string) string {
	if strings.HasPrefix(requestURL, "http://") || strings.HasPrefix(requestURL, "https://") {
		return requestURL
	}

	if p.baseURL != "" {
		if !strings.HasPrefix(requestURL, "/") {
			requestURL = "/" + requestURL
		}

		return p.baseURL + requestURL
	}

	return requestURL
}

func (p *PendingRequest) appendQuery(u string, params map[string]string) string {
	parsed, err := url.Parse(u)

	if err != nil {
		return u
	}

	q := parsed.Query()

	for k, v := range params {
		q.Set(k, v)
	}

	parsed.RawQuery = q.Encode()

	return parsed.String()
}

func (p *PendingRequest) encodeBody(data any) ([]byte, string, error) {
	if data == nil {
		return p.bodyBytes, "", nil
	}

	switch p.bodyFormat {
	case BodyJSON:
		b, err := json.Marshal(data)

		return b, "application/json", err
	case BodyForm:
		switch v := data.(type) {
		case map[string]string:
			vals := url.Values{}

			for k, val := range v {
				vals.Set(k, val)
			}

			return []byte(vals.Encode()), "application/x-www-form-urlencoded", nil
		case url.Values:
			return []byte(v.Encode()), "application/x-www-form-urlencoded", nil
		default:
			return nil, "", fmt.Errorf("client: form body must be map[string]string or url.Values")
		}
	case BodyMultipart:
		return p.encodeMultipart(data)
	case BodyRaw:
		if b, ok := data.([]byte); ok {
			return b, "", nil
		}

		if s, ok := data.(string); ok {
			return []byte(s), "", nil
		}

		return nil, "", fmt.Errorf("client: raw body must be []byte or string")
	}

	return nil, "", nil
}

func (p *PendingRequest) encodeMultipart(data any) ([]byte, string, error) {
	fields, ok := data.(map[string]string)

	if !ok {
		return nil, "", fmt.Errorf("client: multipart body must be map[string]string")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	writer.Close()

	return buf.Bytes(), writer.FormDataContentType(), nil
}

func (p *PendingRequest) buildTransport() RoundTripFunc {
	base := RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		return p.httpClient.Do(req)
	})

	// Apply middleware in reverse order so the first middleware wraps outermost.
	chain := base

	for i := len(p.middleware) - 1; i >= 0; i-- {
		mw := p.middleware[i]
		next := chain

		chain = func(req *http.Request) (*http.Response, error) {
			return mw(req, next)
		}
	}

	return chain
}

func basicAuth(user, password string) string {
	credentials := user + ":" + password

	return base64.StdEncoding.EncodeToString([]byte(credentials))
}

func firstOrNil(data []any) any {
	if len(data) > 0 {
		return data[0]
	}

	return nil
}
