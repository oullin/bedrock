package routing

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"strings"
)

type contextKey string

// Context holds request-scoped routing state passed to every handler.
type Context struct {
	Writer  nethttp.ResponseWriter
	Request *nethttp.Request
}

// Param returns a route path parameter value by name.
// It relies on Go 1.22+ r.PathValue.
func (c *Context) Param(name string) string {
	return c.Request.PathValue(name)
}

// Text writes a plain-text response with the given status code.
func (c *Context) Text(status int, body string) error {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(status)
	_, _ = c.Writer.Write([]byte(body))

	return nil
}

// JSON marshals v as JSON and writes it with the given status code.
func (c *Context) JSON(status int, v any) error {
	b, err := json.Marshal(v)

	if err != nil {
		return err
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	_, _ = c.Writer.Write(b)

	return nil
}

// HTML writes an HTML response with the given status code.
func (c *Context) HTML(status int, body string) error {
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(status)
	_, _ = c.Writer.Write([]byte(body))

	return nil
}

// Redirect sends an HTTP redirect response.
func (c *Context) Redirect(url string, status ...int) error {
	code := nethttp.StatusFound

	if len(status) > 0 {
		code = status[0]
	}

	nethttp.Redirect(c.Writer, c.Request, url, code)

	return nil
}

// NoContent sends a 204 No Content response.
func (c *Context) NoContent() error {
	c.Writer.WriteHeader(nethttp.StatusNoContent)

	return nil
}

// Blob writes raw bytes with the given content type and status code.
func (c *Context) Blob(status int, contentType string, data []byte) error {
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.WriteHeader(status)
	_, _ = c.Writer.Write(data)

	return nil
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// SetCookie adds a cookie to the response.
func (c *Context) SetCookie(cookie *nethttp.Cookie) {
	nethttp.SetCookie(c.Writer, cookie)
}

// Query returns the query parameter value for the given key.
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// QueryDefault returns the query parameter value, or the fallback if absent.
func (c *Context) QueryDefault(key, fallback string) string {
	val := c.Request.URL.Query().Get(key)

	if val == "" {
		return fallback
	}

	return val
}

// HasQuery returns true if the query parameter is present.
func (c *Context) HasQuery(key string) bool {
	return c.Request.URL.Query().Has(key)
}

// FormValue returns the form value for the given key.
func (c *Context) FormValue(key string) string {
	return c.Request.FormValue(key)
}

// Method returns the HTTP method of the request.
func (c *Context) Method() string {
	return c.Request.Method
}

// Path returns the URL path of the request.
func (c *Context) Path() string {
	return c.Request.URL.Path
}

// IsMethod returns true if the request method matches the given method.
func (c *Context) IsMethod(method string) bool {
	return strings.EqualFold(c.Request.Method, method)
}

// IP returns the client IP address extracted from RemoteAddr.
func (c *Context) IP() string {
	addr := c.Request.RemoteAddr

	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}

	return addr
}

// SetData stores a value in the request context for later retrieval.
func (c *Context) SetData(key string, value any) {
	ctx := context.WithValue(c.Request.Context(), contextKey(key), value)
	c.Request = c.Request.WithContext(ctx)
}

// GetData retrieves a value stored in the request context.
func (c *Context) GetData(key string) any {
	return c.Request.Context().Value(contextKey(key))
}
