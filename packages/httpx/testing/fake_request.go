package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/bedrock/packages/httpx"
)

// NewGetRequest creates a test GET request wrapped in httpx.Request.
func NewGetRequest(path string, query ...map[string]string) *httpx.Request {
	u := path

	if len(query) > 0 {
		vals := url.Values{}

		for k, v := range query[0] {
			vals.Set(k, v)
		}

		u += "?" + vals.Encode()
	}

	return httpx.NewRequest(httptest.NewRequest(http.MethodGet, u, nil))
}

// NewJSONRequest creates a test request with a JSON body.
func NewJSONRequest(method, path string, data any) *httpx.Request {
	b, _ := json.Marshal(data)
	raw := httptest.NewRequest(method, path, strings.NewReader(string(b)))
	raw.Header.Set("Content-Type", "application/json")

	return httpx.NewRequest(raw)
}

// NewFormRequest creates a test request with form-encoded body.
func NewFormRequest(method, path string, data map[string]string) *httpx.Request {
	vals := url.Values{}

	for k, v := range data {
		vals.Set(k, v)
	}

	raw := httptest.NewRequest(method, path, strings.NewReader(vals.Encode()))
	raw.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return httpx.NewRequest(raw)
}

// NewRequestWithHeaders creates a test request with custom headers.
func NewRequestWithHeaders(method, path string, headers map[string]string) *httpx.Request {
	raw := httptest.NewRequest(method, path, nil)

	for k, v := range headers {
		raw.Header.Set(k, v)
	}

	return httpx.NewRequest(raw)
}
