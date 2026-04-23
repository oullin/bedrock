package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func queryInt(r *http.Request, key string, fallback int) int {
	raw := strings.TrimSpace(r.URL.Query().Get(key))

	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)

	if err != nil {
		return fallback
	}

	return value
}

func nonNegativeQueryInt(r *http.Request, key string, fallback int) int {
	value := queryInt(r, key, fallback)

	if value < 0 {
		return 0
	}

	return value
}

func normaliseBatchQuery(query string) string {
	// Laravel Horizon escapes LIKE wildcards so users can search for literal
	// % and _ characters; Bedrock does the same by matching against a trimmed
	// lowercase needle.
	return strings.ToLower(strings.TrimSpace(query))
}

func containsFold(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), needle)
}
