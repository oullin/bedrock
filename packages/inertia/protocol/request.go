package protocol

import (
	"net/http"
	"strings"
)

// IsInertiaRequest reports whether r was initiated by the Inertia.js
// client (i.e. it carries the X-Inertia header).
func IsInertiaRequest(r *http.Request) bool {
	return strings.TrimSpace(r.Header.Get(HeaderInertia)) == "true"
}

// IsHTTPPreviewRequest reports whether r carries the HTTPPreview: true
// header, indicating it is a httppreview validation request.
func IsHTTPPreviewRequest(r *http.Request) bool {
	return strings.TrimSpace(r.Header.Get(HeaderHTTPPreview)) == "true"
}

// ValidateOnly returns the list of field names from the Validate-Only
// header, or nil if the header is absent.
func ValidateOnly(r *http.Request) []string {
	header := strings.TrimSpace(r.Header.Get(HeaderValidateOnly))

	if strings.TrimSpace(header) == "" {
		return nil
	}

	fields := strings.Split(header, ",")
	result := make([]string, 0, len(fields))

	for _, f := range fields {
		if trimmed := strings.TrimSpace(f); strings.TrimSpace(trimmed) != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
