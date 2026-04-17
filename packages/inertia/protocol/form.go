package protocol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const maxJSONDepth = 32

// ParseForm parses the request body based on its Content-Type.
// It handles application/json (sent by Inertia.js v3 for non-file forms),
// multipart/form-data, and application/x-www-form-urlencoded.
// After calling ParseForm, r.FormValue() works regardless of content type.
func ParseForm(r *http.Request) error {
	ct := r.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "application/json"):
		return parseJSONForm(r)
	case strings.HasPrefix(ct, "multipart/"):
		return r.ParseMultipartForm(32 << 20)
	default:
		return r.ParseForm()
	}
}

func parseJSONForm(r *http.Request) error {
	const maxBodySize = 32 << 20 // 32 MiB — matches multipart limit

	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)

	var raw map[string]any

	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		return err
	}

	values := make(url.Values)

	if err := flattenJSON("", raw, values, 0); err != nil {
		return err
	}

	r.PostForm = values
	r.Form = make(url.Values)

	for k, v := range values {
		r.Form[k] = v
	}

	if r.URL != nil {
		qs, err := url.ParseQuery(r.URL.RawQuery)

		if err != nil {
			return err
		}

		for k, v := range qs {
			r.Form[k] = append(r.Form[k], v...)
		}
	}

	return nil
}

func flattenJSON(prefix string, data map[string]any, out url.Values, depth int) error {
	if depth > maxJSONDepth {
		return fmt.Errorf("JSON nesting exceeds maximum depth of %d", maxJSONDepth)
	}

	for key, val := range data {
		fullKey := key

		if strings.TrimSpace(prefix) != "" {
			fullKey = prefix + "." + key
		}

		switch v := val.(type) {
		case map[string]any:
			if err := flattenJSON(fullKey, v, out, depth+1); err != nil {
				return err
			}
		case []any:
			if err := flattenArray(fullKey, v, out, depth+1); err != nil {
				return err
			}
		default:
			out.Set(fullKey, toFormValue(val))
		}
	}

	return nil
}

func flattenArray(prefix string, items []any, out url.Values, depth int) error {
	if depth > maxJSONDepth {
		return fmt.Errorf("JSON nesting exceeds maximum depth of %d", maxJSONDepth)
	}

	for i, item := range items {
		arrKey := fmt.Sprintf("%s[%d]", prefix, i)

		switch v := item.(type) {
		case map[string]any:
			if err := flattenJSON(arrKey, v, out, depth+1); err != nil {
				return err
			}
		case []any:
			if err := flattenArray(arrKey, v, out, depth+1); err != nil {
				return err
			}
		default:
			out.Set(arrKey, toFormValue(item))
		}
	}

	return nil
}

func toFormValue(v any) string {
	switch val := v.(type) {
	case bool:
		if val {
			return "1"
		}

		return "0"
	case nil:
		return ""
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}

		return fmt.Sprintf("%g", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
