package support

import (
	"net/url"
	"strings"
)

// URI wraps net/url.URL with Upstream-style immutable helpers.
type URI struct {
	value *url.URL
}

// ParseURI parses a URI string.
func ParseURI(raw string) (URI, error) {
	parsed, err := url.Parse(raw)

	if err != nil {
		return URI{}, err
	}

	return URI{value: parsed}, nil
}

// MustParseURI parses a URI string and panics on invalid input.
func MustParseURI(raw string) URI {
	uri, err := ParseURI(raw)

	if err != nil {
		panic(err)
	}

	return uri
}

// String returns the encoded URI.
func (u URI) String() string {
	if u.value == nil {
		return ""
	}

	return u.value.String()
}

// IsEmpty reports whether the URI is empty.
func (u URI) IsEmpty() bool {
	return u.String() == ""
}

// IsNotEmpty reports whether the URI is not empty.
func (u URI) IsNotEmpty() bool {
	return !u.IsEmpty()
}

// WithoutFragment returns a copy without the fragment.
func (u URI) WithoutFragment() URI {
	next := u.clone()
	next.value.Fragment = ""

	return next
}

// WithQuery returns a copy with query values merged or replaced.
func (u URI) WithQuery(values map[string]string) URI {
	next := u.clone()
	query := next.value.Query()

	for key, value := range values {
		query.Set(key, value)
	}

	next.value.RawQuery = query.Encode()

	return next
}

// WithQueryIfMissing returns a copy with query values only when absent.
func (u URI) WithQueryIfMissing(values map[string]string) URI {
	next := u.clone()
	query := next.value.Query()

	for key, value := range values {
		if _, ok := query[key]; !ok {
			query.Set(key, value)
		}
	}

	next.value.RawQuery = query.Encode()

	return next
}

// Query returns a decoded copy of the query values.
func (u URI) Query() url.Values {
	if u.value == nil {
		return url.Values{}
	}

	return u.value.Query()
}

// PathSegments returns non-empty path segments.
func (u URI) PathSegments() []string {
	if u.value == nil {
		return nil
	}

	parts := strings.Split(strings.Trim(u.value.EscapedPath(), "/"), "/")

	if len(parts) == 1 && parts[0] == "" {
		return nil
	}

	for i, part := range parts {
		decoded, err := url.PathUnescape(part)

		if err == nil {
			parts[i] = decoded
		}
	}

	return parts
}

// Decoded returns the decoded URI string while preserving the fragment.
func (u URI) Decoded() string {
	if u.value == nil {
		return ""
	}

	decoded, err := url.QueryUnescape(u.value.String())

	if err != nil {
		return u.value.String()
	}

	return decoded
}

func (u URI) clone() URI {
	if u.value == nil {
		return URI{value: &url.URL{}}
	}

	copy := *u.value

	return URI{value: &copy}
}
