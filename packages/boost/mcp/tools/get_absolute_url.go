package tools

import (
	"net/url"
	"strings"
)

// GetAbsoluteUrl resolves a relative path or named route to an absolute URL.
// Mirrors Upstream\Boost\Mcp\Tools\GetAbsoluteUrl.
// Tagged IsReadOnly.
type GetAbsoluteUrl struct {
	// BaseURL is the application base URL (e.g. "http://localhost:8080").
	BaseURL string
	// Routes maps named route identifiers to their paths.
	// Optional — only needed when callers pass "route" instead of "path".
	Routes map[string]string
}

func (t *GetAbsoluteUrl) Name() string     { return "get_absolute_url" }
func (t *GetAbsoluteUrl) IsReadOnly() bool { return true }

func (t *GetAbsoluteUrl) Description() string {
	return "Convert a relative URL path or named route to an absolute URL using the application base URL."
}

func (t *GetAbsoluteUrl) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Relative path (e.g. /users/1). Mutually exclusive with route.",
			},
			"route": map[string]any{
				"type":        "string",
				"description": "Named route identifier. Mutually exclusive with path.",
			},
		},
		"required": []string{},
	}
}

// Handle resolves the path or route to an absolute URL.
func (t *GetAbsoluteUrl) Handle(req McpRequest) (McpResponse, error) {
	base := strings.TrimRight(t.BaseURL, "/")

	if base == "" {
		base = "http://localhost"
	}

	var relativePath string

	if routeName, ok := req.Args["route"].(string); ok && routeName != "" {
		if t.Routes != nil {
			if p, exists := t.Routes[routeName]; exists {
				relativePath = p
			}
		}

		if relativePath == "" {
			relativePath = "/" + routeName
		}
	} else if path, ok := req.Args["path"].(string); ok {
		relativePath = path
	}

	absolute, err := url.JoinPath(base, relativePath)

	if err != nil {
		return ErrorResponse("get_absolute_url: " + err.Error()), nil
	}

	return OkResponse(map[string]any{"url": absolute}), nil
}
