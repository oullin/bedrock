package cookie

import (
	"net/http"
	"time"
)

// SameSite matches the http.SameSite constants for convenient use.
const (
	SameSiteDefault = http.SameSiteDefaultMode
	SameSiteLax     = http.SameSiteLaxMode
	SameSiteStrict  = http.SameSiteStrictMode
	SameSiteNone    = http.SameSiteNoneMode
)

// Options configures cookie attributes.
type Options struct {
	Path     string
	Domain   string
	MaxAge   int // seconds; 0 = session cookie, negative = delete
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
	Raw      bool // do not URL-encode the value
}

// DefaultOptions returns sensible production defaults.
func DefaultOptions() Options {
	return Options{
		Path:     "/",
		HTTPOnly: true,
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
		Secure:   opts.Secure,
		HttpOnly: opts.HTTPOnly,
		SameSite: opts.SameSite,
	}

	if opts.Raw {
		c.Raw = name + "=" + value
	}

	return c
}

// Forever creates a cookie that expires in 400 days (~Laravel's "forever").
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
