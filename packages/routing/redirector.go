package routing

import nethttp "net/http"

// Redirector produces HTTP redirect responses using a UrlGenerator for
// named route resolution.
type Redirector struct {
	generator *UrlGenerator
}

// NewRedirector creates a redirector backed by the given URL generator.
func NewRedirector(generator *UrlGenerator) *Redirector {
	return &Redirector{generator: generator}
}

// To sends a redirect to the given path.
func (rd *Redirector) To(ctx *Context, path string, status ...int) error {
	code := nethttp.StatusFound

	if len(status) > 0 {
		code = status[0]
	}

	nethttp.Redirect(ctx.Writer, ctx.Request, rd.generator.To(path), code)

	return nil
}

// Away sends a redirect to an external URL without validation.
func (rd *Redirector) Away(ctx *Context, rawURL string, status ...int) error {
	code := nethttp.StatusFound

	if len(status) > 0 {
		code = status[0]
	}

	nethttp.Redirect(ctx.Writer, ctx.Request, rawURL, code)

	return nil
}

// Secure sends a redirect using HTTPS.
func (rd *Redirector) Secure(ctx *Context, path string, status ...int) error {
	code := nethttp.StatusFound

	if len(status) > 0 {
		code = status[0]
	}

	nethttp.Redirect(ctx.Writer, ctx.Request, rd.generator.Secure(path), code)

	return nil
}

// Back sends a redirect to the previous page (Referer header) or a fallback.
func (rd *Redirector) Back(ctx *Context, fallback ...string) error {
	rd.generator.SetRequest(ctx.Request)

	nethttp.Redirect(ctx.Writer, ctx.Request, rd.generator.Previous(fallback...), nethttp.StatusFound)

	return nil
}

// Refresh sends a redirect to the current URL.
func (rd *Redirector) Refresh(ctx *Context) error {
	rd.generator.SetRequest(ctx.Request)

	nethttp.Redirect(ctx.Writer, ctx.Request, rd.generator.Current(), nethttp.StatusFound)

	return nil
}

// Route sends a redirect to a named route.
func (rd *Redirector) Route(ctx *Context, name string, params map[string]string, status ...int) error {
	code := nethttp.StatusFound

	if len(status) > 0 {
		code = status[0]
	}

	nethttp.Redirect(ctx.Writer, ctx.Request, rd.generator.Route(name, params), code)

	return nil
}

// GetUrlGenerator returns the underlying URL generator.
func (rd *Redirector) GetUrlGenerator() *UrlGenerator {
	return rd.generator
}
