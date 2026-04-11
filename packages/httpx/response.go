package httpx

import (
	"encoding/json"
	"net/http"
)

// Response wraps http.ResponseWriter with a fluent API for building HTTP
// responses. Headers and cookies are buffered until Send, JSON or NoContent is
// called, at which point they are flushed to the underlying writer.
type Response struct {
	writer      http.ResponseWriter
	statusCode  int
	headers     http.Header
	cookies     []*http.Cookie
	exception   error
	original    any
	headersSent bool
}

// NewResponse creates a Response backed by the given http.ResponseWriter.
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{
		writer:     w,
		statusCode: http.StatusOK,
		headers:    make(http.Header),
	}
}

// Status sets the HTTP status code.
func (r *Response) Status(code int) *Response {
	r.statusCode = code

	return r
}

// GetStatusCode returns the current status code.
func (r *Response) GetStatusCode() int {
	return r.statusCode
}

// Header sets a single response header.
func (r *Response) Header(key, value string) *Response {
	r.headers.Set(key, value)

	return r
}

// WithHeaders sets multiple response headers at once.
func (r *Response) WithHeaders(headers map[string]string) *Response {
	for k, v := range headers {
		r.headers.Set(k, v)
	}

	return r
}

// WithoutHeader removes a previously set header.
func (r *Response) WithoutHeader(key string) *Response {
	r.headers.Del(key)

	return r
}

// GetHeader returns the value of a buffered response header.
func (r *Response) GetHeader(key string) string {
	return r.headers.Get(key)
}

// Cookie appends a cookie to the response.
func (r *Response) Cookie(c *http.Cookie) *Response {
	r.cookies = append(r.cookies, c)

	return r
}

// WithoutCookie appends a cookie-deletion instruction (MaxAge = -1).
func (r *Response) WithoutCookie(name string, path ...string) *Response {
	p := "/"

	if len(path) > 0 {
		p = path[0]
	}

	r.cookies = append(r.cookies, &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   p,
		MaxAge: -1,
	})

	return r
}

// WithException attaches an error to the response for later inspection.
func (r *Response) WithException(err error) *Response {
	r.exception = err

	return r
}

// GetException returns the attached exception, if any.
func (r *Response) GetException() error {
	return r.exception
}

// SetOriginal stores the original content before it was converted to bytes.
func (r *Response) SetOriginal(original any) *Response {
	r.original = original

	return r
}

// GetOriginal returns the original content.
func (r *Response) GetOriginal() any {
	return r.original
}

// flushHeaders writes buffered headers and cookies to the underlying writer.
// Must be called before WriteHeader.
func (r *Response) flushHeaders() {
	if r.headersSent {
		return
	}

	r.headersSent = true

	dest := r.writer.Header()

	for k, vals := range r.headers {
		for _, v := range vals {
			dest.Add(k, v)
		}
	}

	for _, c := range r.cookies {
		http.SetCookie(r.writer, c)
	}
}

// Send writes the given body and finalises the response.
func (r *Response) Send(body []byte) error {
	r.flushHeaders()
	r.writer.WriteHeader(r.statusCode)
	_, err := r.writer.Write(body)

	return err
}

// SendString is a convenience wrapper around Send for string bodies.
func (r *Response) SendString(body string) error {
	if r.headers.Get("Content-Type") == "" {
		r.headers.Set("Content-Type", "text/plain; charset=utf-8")
	}

	return r.Send([]byte(body))
}

// JSON marshals v as JSON and sends it.
func (r *Response) JSON(v any) error {
	b, err := json.Marshal(v)

	if err != nil {
		return err
	}

	r.original = v

	if r.headers.Get("Content-Type") == "" {
		r.headers.Set("Content-Type", "application/json")
	}

	return r.Send(b)
}

// NoContent sends a response with no body. The default status is 204.
func (r *Response) NoContent(status ...int) error {
	code := http.StatusNoContent

	if len(status) > 0 {
		code = status[0]
	}

	r.statusCode = code
	r.flushHeaders()
	r.writer.WriteHeader(r.statusCode)

	return nil
}

// Writer returns the underlying http.ResponseWriter.
func (r *Response) Writer() http.ResponseWriter {
	return r.writer
}
