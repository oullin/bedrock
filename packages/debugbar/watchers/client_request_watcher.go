package watchers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bedrock/packages/debugbar"
)

// ClientRequestWatcher monitors outbound HTTP client requests and records them
// as DebugBar entries. It mirrors Upstream's ClientRequestWatcher class.
//
// Options:
//   - "size_limit" (int): max response body bytes to store (default 64 KB).
type ClientRequestWatcher struct {
	debugbar.BaseWatcher
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

func NewClientRequestWatcher(t *debugbar.DebugBar, options map[string]any) *ClientRequestWatcher {
	w := &ClientRequestWatcher{}
	w.SetDebugBar(t)
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

	if payload, ok := maskedPayload.(map[string]any); ok {
		maskedPayload = maskMap(payload, scope.HiddenRequestParameters())
	}

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

	entry := debugbar.NewEntry(debugbar.EntryTypeClientRequest, content)

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

	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		values, err := url.ParseQuery(string(body))

		if err == nil {
			m := make(map[string]any, len(values))

			for k, v := range values {
				if len(v) == 1 {
					m[k] = v[0]
				} else {
					m[k] = v
				}
			}

			return m
		}
	}

	if strings.Contains(contentType, "multipart/form-data") {
		if payload, ok := parseMultipartBody(body, contentType); ok {
			return payload
		}
	}

	return string(body)
}

func parseMultipartBody(body []byte, contentType string) (map[string]any, bool) {
	_, params, err := mime.ParseMediaType(contentType)

	if err != nil {
		return nil, false
	}

	boundary := params["boundary"]

	if boundary == "" {
		return nil, false
	}

	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	form, err := reader.ReadForm(int64(len(body)))

	if err != nil {
		return nil, false
	}

	defer form.RemoveAll() //nolint:errcheck

	payload := make(map[string]any, len(form.Value)+len(form.File))

	for key, values := range form.Value {
		if len(values) == 1 {
			payload[key] = values[0]
		} else {
			payload[key] = values
		}
	}

	for key, files := range form.File {
		filePayload := make([]map[string]any, 0, len(files))

		for _, header := range files {
			size := header.Size

			if size == 0 {
				if file, err := header.Open(); err == nil {
					n, _ := io.Copy(io.Discard, file)
					_ = file.Close()
					size = n
				}
			}

			filePayload = append(filePayload, map[string]any{
				"name": header.Filename,
				"size": size,
			})
		}

		if len(filePayload) == 1 {
			payload[key] = filePayload[0]
		} else {
			payload[key] = filePayload
		}
	}

	return payload, true
}
