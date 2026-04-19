package watchers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bedrock/packages/telescope"
)

// ClientRequestWatcher monitors outbound HTTP client requests and records them
// as Telescope entries. It mirrors Laravel's ClientRequestWatcher class.
//
// Options:
//   - "size_limit" (int): max response body bytes to store (default 64 KB).
type ClientRequestWatcher struct {
	telescope.BaseWatcher
}

// NewClientRequestWatcher creates a ClientRequestWatcher with the given options.

// Register is a no-op for ClientRequestWatcher; callers drive it via Record.

// ClientResponse carries the response data from an outbound HTTP request.
type ClientResponse struct {
	StatusCode  int
	Headers     http.Header
	Body        []byte
	ContentType string
}

func NewClientRequestWatcher(t *telescope.Telescope, options map[string]any) *ClientRequestWatcher {
	w := &ClientRequestWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

func (w *ClientRequestWatcher) Register(_ any) error { return nil }

// Record records an outbound HTTP client request/response pair.
func (w *ClientRequestWatcher) Record(
	method, uri string,
	requestHeaders http.Header,
	requestBody []byte,
	response *ClientResponse,
	duration time.Duration,
) {
	sizeLimit := w.clientSizeLimit()
	scope := w.Scope()

	maskedHeaders := maskHeaders(requestHeaders, scope.HiddenRequestHeaders())
	maskedPayload := parseClientBody(requestBody, requestHeaders.Get("Content-Type"))

	content := map[string]any{
		"method":   strings.ToUpper(method),
		"uri":      uri,
		"headers":  maskedHeaders,
		"payload":  maskedPayload,
		"duration": float64(duration.Milliseconds()),
	}

	if response != nil {
		respBody := truncate(response.Body, sizeLimit)
		content["response_status"] = response.StatusCode
		content["response_headers"] = maskHeaders(response.Headers, scope.HiddenResponseParameters())
		content["response"] = formatResponseBody(response.ContentType, respBody)
	}

	entry := telescope.NewEntry(telescope.EntryTypeClientRequest, content)

	scope.RecordClientRequest(entry)
}

// clientSizeLimit returns the configured size limit for client request watcher.
func (w *ClientRequestWatcher) clientSizeLimit() int {
	if v, ok := w.Options["size_limit"]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		}
	}

	return defaultRequestSizeLimit
}

// parseClientBody decodes the request body based on content type.
func parseClientBody(body []byte, contentType string) any {
	if len(body) == 0 {
		return nil
	}

	if strings.Contains(contentType, "application/json") {
		var v any

		if json.Unmarshal(body, &v) == nil {
			return v
		}
	}

	return string(body)
}
