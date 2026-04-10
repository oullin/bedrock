package routing

import (
	"encoding/json"
	nethttp "net/http"
)

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
