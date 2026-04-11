package httpx

import (
	"strings"
)

// ContentType returns the request Content-Type without parameters.
func (r *Request) ContentType() string {
	ct := r.raw.Header.Get("Content-Type")

	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = ct[:idx]
	}

	return strings.TrimSpace(ct)
}

// IsJSON returns true when the Content-Type indicates JSON.
func (r *Request) IsJSON() bool {
	return matchesType(r.ContentType(), "application/json", "application/*+json")
}

// ExpectsJSON returns true when the client expects a JSON response. This is
// true when the Accept header includes JSON or the X-Requested-With header is
// XMLHttpRequest (XHR / fetch).
func (r *Request) ExpectsJSON() bool {
	return r.WantsJSON() || r.Ajax()
}

// WantsJSON returns true when the Accept header prefers JSON over HTML.
func (r *Request) WantsJSON() bool {
	accept := r.raw.Header.Get("Accept")

	if accept == "" {
		return false
	}

	types := parseAccept(accept)

	for _, t := range types {
		if matchesType(t, "text/html", "application/xhtml+xml") {
			return false
		}

		if matchesType(t, "application/json", "application/*+json") {
			return true
		}
	}

	return false
}

// WantsMarkdown returns true when the Accept header prefers Markdown.
func (r *Request) WantsMarkdown() bool {
	return r.Prefers("text/markdown", "text/html") == "text/markdown"
}

// Ajax returns true when the request was sent via XMLHttpRequest.
func (r *Request) Ajax() bool {
	return r.raw.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// Pjax returns true when the request was sent via Pjax.
func (r *Request) Pjax() bool {
	return r.raw.Header.Get("X-PJAX") != ""
}

// Prefetch returns true when the request is a prefetch or prerender.
func (r *Request) Prefetch() bool {
	purpose := strings.ToLower(r.raw.Header.Get("Purpose"))
	fetchPurpose := strings.ToLower(r.raw.Header.Get("Sec-Purpose"))

	return purpose == "prefetch" || fetchPurpose == "prefetch" ||
		r.raw.Header.Get("X-Purpose") == "preview" ||
		r.raw.Header.Get("X-moz") == "prefetch"
}

// Accepts returns the first content type from the given list that the client
// accepts, or an empty string if none match.
func (r *Request) Accepts(contentTypes ...string) string {
	accept := parseAccept(r.raw.Header.Get("Accept"))

	if len(accept) == 0 {
		// When no Accept header is present, accept anything.
		if len(contentTypes) > 0 {
			return contentTypes[0]
		}

		return ""
	}

	for _, a := range accept {
		if a == "*/*" || a == "*" {
			if len(contentTypes) > 0 {
				return contentTypes[0]
			}

			return ""
		}

		for _, ct := range contentTypes {
			if matchesType(a, ct) {
				return ct
			}
		}
	}

	return ""
}

// Prefers returns the most preferred content type from the given list based on
// the Accept header's quality ordering.
func (r *Request) Prefers(contentTypes ...string) string {
	return r.Accepts(contentTypes...)
}

// AcceptsAnyContentType returns true when the client has a wildcard accept or
// no Accept header at all.
func (r *Request) AcceptsAnyContentType() bool {
	accept := r.raw.Header.Get("Accept")

	if accept == "" {
		return true
	}

	types := parseAccept(accept)

	return len(types) > 0 && (types[0] == "*/*" || types[0] == "*")
}

// AcceptsJSON returns true when the client accepts JSON responses.
func (r *Request) AcceptsJSON() bool {
	return r.Accepts("application/json") != ""
}

// AcceptsHTML returns true when the client accepts HTML responses.
func (r *Request) AcceptsHTML() bool {
	return r.Accepts("text/html") != ""
}

// AcceptsMarkdown returns true when the client accepts Markdown responses.
func (r *Request) AcceptsMarkdown() bool {
	return r.Accepts("text/markdown") != ""
}

// Format returns a default format based on what the client accepts.
func (r *Request) Format(fallback string) string {
	mapping := map[string]string{
		"text/html":        "html",
		"application/json": "json",
		"application/xml":  "xml",
		"text/plain":       "txt",
	}

	for ct, format := range mapping {
		if r.Accepts(ct) != "" {
			return format
		}
	}

	return fallback
}

// parseAccept splits an Accept header into ordered media types (quality-sorted
// in descending order). Parameters are stripped.
func parseAccept(header string) []string {
	if header == "" {
		return nil
	}

	type entry struct {
		mediaType string
		quality   float64
		order     int
	}

	parts := strings.Split(header, ",")
	entries := make([]entry, 0, len(parts))

	for i, p := range parts {
		p = strings.TrimSpace(p)

		if p == "" {
			continue
		}

		var mediaType string

		q := 1.0
		semicolons := strings.Split(p, ";")
		mediaType = strings.TrimSpace(semicolons[0])

		for _, param := range semicolons[1:] {
			param = strings.TrimSpace(param)

			if strings.HasPrefix(param, "q=") {
				if parsed, err := parseFloat(param[2:]); err == nil {
					q = parsed
				}
			}
		}

		entries = append(entries, entry{mediaType: mediaType, quality: q, order: i})
	}

	// Stable sort by quality descending, then by original order ascending.
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0; j-- {
			if entries[j].quality > entries[j-1].quality ||
				(entries[j].quality == entries[j-1].quality && entries[j].order < entries[j-1].order) {
				entries[j], entries[j-1] = entries[j-1], entries[j]
			}
		}
	}

	result := make([]string, len(entries))

	for i, e := range entries {
		result[i] = e.mediaType
	}

	return result
}

// matchesType checks whether a media type matches any of the given patterns.
// Supports subtype wildcards (e.g., "application/*").
func matchesType(actual string, patterns ...string) bool {
	actual = strings.TrimSpace(strings.ToLower(actual))

	for _, pattern := range patterns {
		pattern = strings.TrimSpace(strings.ToLower(pattern))

		if actual == pattern {
			return true
		}

		// Handle wildcard patterns like "application/*+json".
		if strings.Contains(pattern, "*") {
			patParts := strings.SplitN(pattern, "/", 2)
			actParts := strings.SplitN(actual, "/", 2)

			if len(patParts) != 2 || len(actParts) != 2 {
				continue
			}

			if patParts[0] != actParts[0] && patParts[0] != "*" {
				continue
			}

			subPattern := patParts[1]

			if subPattern == "*" {
				return true
			}

			// Handle "+json" style suffixes.
			if strings.HasPrefix(subPattern, "*+") {
				suffix := subPattern[1:] // "+json"

				if strings.HasSuffix(actParts[1], suffix) {
					return true
				}
			}
		}
	}

	return false
}

// parseFloat is a minimal float parser for quality values.
func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)

	var result float64

	var decimal float64

	var afterDot bool

	divisor := 1.0

	for _, ch := range s {
		if ch == '.' {
			afterDot = true

			continue
		}

		if ch < '0' || ch > '9' {
			break
		}

		digit := float64(ch - '0')

		if afterDot {
			divisor *= 10
			decimal += digit / divisor
		} else {
			result = result*10 + digit
		}
	}

	return result + decimal, nil
}
