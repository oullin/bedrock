package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

// BodyFormat describes the request body encoding.
type BodyFormat int

// attachment represents a file attachment for multipart requests.
type attachment struct {
	Name     string
	Contents io.Reader
	Filename string
}

// PendingRequest is a fluent builder for outbound HTTP requests.
type PendingRequest struct {
	factory         *Factory
	httpClient      *http.Client
	baseURL         string
	bodyFormat      BodyFormat
	headers         http.Header
	cookies         []*http.Cookie
	timeout         time.Duration
	connectTimeout  time.Duration
	retries         int
	retryDelay      time.Duration
	retryWhen       func(error, *Response) bool
	middleware      []Middleware
	body            io.Reader
	bodyBytes       []byte
	ctx             context.Context
	queryParams     map[string]string
	urlParams       map[string]string
	skipTLSVerify   bool
	sinkWriter      io.Writer
	attachments     []attachment
	beforeCallbacks []func(*http.Request)
	afterCallbacks  []func(*Response)
	throwOnFailure  bool
	throwCallbacks  []func(*Response, error)
	stub            StubCallback
	preventStray    bool
	attributes      map[string]any
}

const (
	BodyJSON BodyFormat = iota
	BodyForm
	BodyMultipart
	BodyRaw
)

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

// WithDigestAuth sets HTTP Digest authentication credentials. The digest
// handshake is performed automatically when the server responds with a 401
// and a WWW-Authenticate: Digest challenge.
func (p *PendingRequest) WithDigestAuth(user, password string) *PendingRequest {
	p.middleware = append(p.middleware, digestMiddleware(&digestAuth{
		username: user,
		password: password,
	}))

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

// ConnectTimeout sets the connection timeout separately from the overall
// request timeout.
func (p *PendingRequest) ConnectTimeout(d time.Duration) *PendingRequest {
	p.connectTimeout = d

	return p
}

// WithQueryParameters sets query parameters that are merged into every request
// URL.
func (p *PendingRequest) WithQueryParameters(params map[string]string) *PendingRequest {
	p.queryParams = params

	return p
}

// ContentType sets the Content-Type header.
func (p *PendingRequest) ContentType(contentType string) *PendingRequest {
	p.headers.Set("Content-Type", contentType)

	return p
}

// WithUserAgent sets the User-Agent header.
func (p *PendingRequest) WithUserAgent(agent string) *PendingRequest {
	p.headers.Set("User-Agent", agent)

	return p
}

// WithUrlParameters sets URL template parameters. Placeholders like {key} in
// the URL will be replaced with the corresponding value.
func (p *PendingRequest) WithUrlParameters(params map[string]string) *PendingRequest {
	p.urlParams = params

	return p
}

// ReplaceHeaders replaces all headers with the given set.
func (p *PendingRequest) ReplaceHeaders(headers map[string]string) *PendingRequest {
	p.headers = make(http.Header)

	for k, v := range headers {
		p.headers.Set(k, v)
	}

	return p
}

// Send sends a request with the given HTTP method.
func (p *PendingRequest) Send(method, url string, data ...any) (*Response, error) {
	return p.send(method, url, firstOrNil(data))
}

// WithoutVerifying disables TLS certificate verification.
func (p *PendingRequest) WithoutVerifying() *PendingRequest {
	p.skipTLSVerify = true

	return p
}

// Sink sets a writer that the response body will be written to.
func (p *PendingRequest) Sink(w io.Writer) *PendingRequest {
	p.sinkWriter = w

	return p
}

// Attach adds a file attachment for multipart requests.
func (p *PendingRequest) Attach(name string, contents io.Reader, filename string) *PendingRequest {
	p.attachments = append(p.attachments, attachment{
		Name:     name,
		Contents: contents,
		Filename: filename,
	})
	p.bodyFormat = BodyMultipart

	return p
}

// BeforeSending registers a callback that runs before each request is sent.
func (p *PendingRequest) BeforeSending(fn func(*http.Request)) *PendingRequest {
	p.beforeCallbacks = append(p.beforeCallbacks, fn)

	return p
}

// AfterResponse registers a callback that runs after each response is received.
func (p *PendingRequest) AfterResponse(fn func(*Response)) *PendingRequest {
	p.afterCallbacks = append(p.afterCallbacks, fn)

	return p
}

// Throw enables automatic error throwing when the response indicates failure.
// Optional callbacks are invoked with the response and error before returning.
func (p *PendingRequest) Throw(callback ...func(*Response, error)) *PendingRequest {
	p.throwOnFailure = true
	p.throwCallbacks = append(p.throwCallbacks, callback...)

	return p
}

// ThrowIf enables automatic error throwing when the given condition is true.
func (p *PendingRequest) ThrowIf(condition bool) *PendingRequest {
	if condition {
		p.throwOnFailure = true
	}

	return p
}

// ThrowUnless enables automatic error throwing when the given condition is
// false.
func (p *PendingRequest) ThrowUnless(condition bool) *PendingRequest {
	return p.ThrowIf(!condition)
}

// Stub sets a per-request stub callback. The stub is checked before the
// factory's fakes.
func (p *PendingRequest) Stub(callback StubCallback) *PendingRequest {
	p.stub = callback

	return p
}

// PreventStrayRequests causes the request to error when neither the per-request
// stub nor the factory fake handles the request.
func (p *PendingRequest) PreventStrayRequests() *PendingRequest {
	p.preventStray = true

	return p
}

// WithAttributes sets arbitrary attributes on the request context.
func (p *PendingRequest) WithAttributes(attrs map[string]any) *PendingRequest {
	p.attributes = attrs

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

	// Merge global query parameters.
	if len(p.queryParams) > 0 {
		fullURL = p.appendQuery(fullURL, p.queryParams)
	}

	body, contentType, err := p.encodeBody(data)

	if err != nil {
		return nil, err
	}

	if contentType != "" && p.headers.Get("Content-Type") == "" {
		p.headers.Set("Content-Type", contentType)
	}

	var resp *Response

	var lastErr error

	// Build the transport once before the retry loop.
	transport := p.buildTransport()

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

		// Set attributes on request context.
		for k, v := range p.attributes {
			req = req.WithContext(context.WithValue(req.Context(), k, v))
		}

		// Run before-sending callbacks.
		for _, fn := range p.beforeCallbacks {
			fn(req)
		}

		// Dispatch RequestSending event.
		if p.factory != nil && p.factory.dispatcher != nil {
			p.factory.dispatcher.Dispatch(RequestSending{Request: req})
		}

		// Check per-request stub first.
		if p.stub != nil {
			if raw := p.stub(req); raw != nil {
				resp = NewResponse(raw)
			} else if p.preventStray {
				return nil, &ConnectionError{URL: fullURL, Err: ErrStrayRequest}
			}
		}

		// Check if factory has a fake/stub.
		if resp == nil && p.factory != nil && p.factory.isFaking() {
			resp, lastErr = p.factory.handleFake(req, body)
		} else if resp == nil {
			// Attach httptrace for handler stats.
			var dnsStart, connectStart, tlsStart, reqStart time.Time
			stats := make(map[string]any)
			reqStart = time.Now()

			trace := &httptrace.ClientTrace{
				DNSStart: func(info httptrace.DNSStartInfo) {
					dnsStart = time.Now()
				},
				DNSDone: func(info httptrace.DNSDoneInfo) {
					if !dnsStart.IsZero() {
						stats["dns_ms"] = float64(time.Since(dnsStart).Milliseconds())
					}
				},
				ConnectStart: func(network, addr string) {
					connectStart = time.Now()
				},
				ConnectDone: func(network, addr string, err error) {
					if !connectStart.IsZero() {
						stats["connect_ms"] = float64(time.Since(connectStart).Milliseconds())
					}
				},
				TLSHandshakeStart: func() {
					tlsStart = time.Now()
				},
				TLSHandshakeDone: func(state tls.ConnectionState, err error) {
					if !tlsStart.IsZero() {
						stats["tls_ms"] = float64(time.Since(tlsStart).Milliseconds())
					}
				},
			}

			traceReq := req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

			// Execute through middleware chain.
			rawResp, err := transport(traceReq)

			stats["total_ms"] = float64(time.Since(reqStart).Milliseconds())

			if err != nil {
				lastErr = &ConnectionError{URL: fullURL, Err: err}

				// Dispatch ConnectionFailed event.
				if p.factory != nil && p.factory.dispatcher != nil {
					p.factory.dispatcher.Dispatch(ConnectionFailed{Request: req, Err: err})
				}

				if p.retryWhen != nil && !p.retryWhen(lastErr, nil) {
					return nil, lastErr
				}

				continue
			}

			resp = NewResponse(rawResp)
			resp.SetStats(stats)
		}

		// Write response body to sink if configured.
		if resp != nil && p.sinkWriter != nil {
			_, _ = p.sinkWriter.Write(resp.Bytes())
		}

		// Run after-response callbacks.
		if resp != nil {
			for _, fn := range p.afterCallbacks {
				fn(resp)
			}

			// Dispatch ResponseReceived event.
			if p.factory != nil && p.factory.dispatcher != nil {
				p.factory.dispatcher.Dispatch(ResponseReceived{Request: req, Response: resp})
			}
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

		// Auto-throw on failure if configured.
		if p.throwOnFailure && resp != nil && resp.Failed() {
			throwErr := &RequestError{Response: resp}

			for _, fn := range p.throwCallbacks {
				fn(resp, throwErr)
			}

			return resp, throwErr
		}

		return resp, nil
	}

	return resp, lastErr
}

func (p *PendingRequest) buildURL(requestURL string) string {
	// Replace URL template parameters.
	for k, v := range p.urlParams {
		requestURL = strings.ReplaceAll(requestURL, "{"+k+"}", v)
	}

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
	if data == nil && len(p.attachments) == 0 {
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
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Write form fields if provided.
	if data != nil {
		fields, ok := data.(map[string]string)

		if !ok {
			return nil, "", fmt.Errorf("client: multipart body must be map[string]string")
		}

		for k, v := range fields {
			_ = writer.WriteField(k, v)
		}
	}

	// Write file attachments.
	for _, att := range p.attachments {
		part, err := writer.CreateFormFile(att.Name, att.Filename)

		if err != nil {
			return nil, "", err
		}

		if _, err := io.Copy(part, att.Contents); err != nil {
			return nil, "", err
		}
	}

	writer.Close()

	return buf.Bytes(), writer.FormDataContentType(), nil
}

func (p *PendingRequest) buildTransport() RoundTripFunc {
	// Configure custom transport when connect timeout or TLS skip is needed.
	if p.connectTimeout > 0 || p.skipTLSVerify {
		transport := &http.Transport{}

		if p.connectTimeout > 0 {
			transport.DialContext = (&net.Dialer{
				Timeout: p.connectTimeout,
			}).DialContext
		}

		if p.skipTLSVerify {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec // intentional user opt-in
			}
		}

		p.httpClient.Transport = transport
	}

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
