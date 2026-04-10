package authflows

import (
	"encoding/json"
	"net/http"
	"strings"
)

// RequestInput extracts form or JSON input from a request.
func RequestInput(r *http.Request, keys ...string) map[string]string {
	result := make(map[string]string, len(keys))

	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		var body map[string]any

		if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
			for _, key := range keys {
				if v, ok := body[key]; ok {
					if s, ok := v.(string); ok {
						result[key] = s
					}
				}
			}
		}

		return result
	}

	_ = r.ParseForm()

	for _, key := range keys {
		result[key] = r.FormValue(key)
	}

	return result
}

// RequestIP extracts the client IP from the request.
func RequestIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.SplitN(forwarded, ",", 2)

		return strings.TrimSpace(parts[0])
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	addr := r.RemoteAddr

	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}

	return addr
}

// RequestBool extracts a boolean form value.
func RequestBool(r *http.Request, key string) bool {
	input := RequestInput(r, key)
	v := input[key]

	return v == "true" || v == "1" || v == "on" || v == "yes"
}
