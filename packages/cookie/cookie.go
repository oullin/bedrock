package cookie

import (
	"net/http"
	"time"
)

// SameSite matches the http.SameSite constants for convenient use.

// Options configures cookie attributes. Boolean fields use *bool so that
// callers can distinguish "not set" (nil) from an explicit false, allowing
// defaults to be overridden in either direction.
type Options struct {
	Path     string
	Domain   string
	MaxAge   int // seconds; 0 = session cookie, negative = delete
	Secure   *bool
	HTTPOnly *bool
	SameSite http.SameSite
	Raw      *bool // do not URL-encode the value
}

// BoolPtr returns a pointer to v, useful for setting Options boolean fields.
func BoolPtr(v bool) *bool { return &v }

const (
	SameSiteDefault = http.SameSiteDefaultMode
	SameSiteLax     = http.SameSiteLaxMode
	SameSiteStrict  = http.SameSiteStrictMode
	SameSiteNone    = http.SameSiteNoneMode
)

// DefaultOptions returns sensible production defaults.
func DefaultOptions() Options {
	return Options{
		Path:     "/",
		HTTPOnly: BoolPtr(true),
		SameSite: SameSiteLax,
	}
}

// Make creates an *http.Cookie from the given options.
func Make(name, value string, opts Options) *http.Cookie {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     opts.Path,
		Domain:   opts.Domain,
		MaxAge:   opts.MaxAge,
		SameSite: opts.SameSite,
	}

	if opts.Secure != nil {
		c.Secure = *opts.Secure
	}

	if opts.HTTPOnly != nil {
		c.HttpOnly = *opts.HTTPOnly
	}

	if opts.Raw != nil && *opts.Raw {
		c.Raw = name + "=" + value
	}

	return c
}

// Forever creates a cookie that expires in 400 days (~the upstream "forever").
func Forever(name, value string, opts Options) *http.Cookie {
	opts.MaxAge = int((400 * 24 * time.Hour).Seconds())

	return Make(name, value, opts)
}

// Forget creates a cookie with a negative max-age to instruct the browser
// to delete it.
func Forget(name string, opts Options) *http.Cookie {
	opts.MaxAge = -1

	return Make(name, "", opts)
}
