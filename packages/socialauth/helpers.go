package socialauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// formBody encodes a string map as an application/x-www-form-urlencoded body.
func formBody(params map[string]string) *strings.Reader {
	vals := url.Values{}

	for k, v := range params {
		vals.Set(k, v)
	}

	return strings.NewReader(vals.Encode())
}

// decodeJSONBody reads an HTTP response body and decodes it as JSON.
func decodeJSONBody(resp *http.Response) (map[string]any, error) {
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var result map[string]any

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("socialauth: cannot decode response: %w", err)
	}

	return result, nil
}
