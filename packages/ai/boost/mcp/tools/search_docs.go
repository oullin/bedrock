package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SearchDocs performs semantic search over the Boost documentation API.
type SearchDocs struct {
	// APIUrl overrides the default Boost docs search endpoint.
	APIUrl string
	// HTTPClient allows injecting a custom http.Client for testing.
	HTTPClient *http.Client
}

// defaultDocsAPIURL is intentionally empty so that the tool requires an
// explicit endpoint via SearchDocs.APIUrl. Callers point this at whichever
// docs search service they want to query.
const defaultDocsAPIURL = ""

func (t *SearchDocs) Name() string     { return "search_docs" }
func (t *SearchDocs) IsReadOnly() bool { return true }

func (t *SearchDocs) Description() string {
	return "Search up-to-date documentation for Go packages and frameworks " +
		"via the Boost documentation API."
}

func (t *SearchDocs) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"queries": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "string"},
				"description": "One or more search queries.",
			},
			"packages": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "string"},
				"description": "Optional list of package names to filter results.",
			},
			"token_limit": map[string]any{
				"type":        "integer",
				"description": "Maximum token budget for results.",
			},
		},
		"required": []string{"queries"},
	}
}

// Handle calls the Boost docs API and returns the search results.
func (t *SearchDocs) Handle(req McpRequest) (McpResponse, error) {
	queries := extractStringSlice(req.Args["queries"])

	if len(queries) == 0 {
		return ErrorResponse("search_docs: queries argument is required"), nil
	}

	payload := map[string]any{"queries": queries}

	if pkgs := extractStringSlice(req.Args["packages"]); len(pkgs) > 0 {
		payload["packages"] = pkgs
	}

	if limit, ok := req.Args["token_limit"].(float64); ok {
		payload["token_limit"] = int(limit)
	}

	body, err := json.Marshal(payload)

	if err != nil {
		return ErrorResponse(fmt.Sprintf("search_docs: marshal: %v", err)), nil
	}

	apiURL := t.APIUrl

	if apiURL == "" {
		apiURL = defaultDocsAPIURL
	}

	if apiURL == "" {
		return ErrorResponse("search_docs: APIUrl is required (no default endpoint configured)"), nil
	}

	client := t.HTTPClient

	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	resp, err := client.Post(apiURL, "application/json", bytes.NewReader(body)) //nolint:gosec

	if err != nil {
		return ErrorResponse(fmt.Sprintf("search_docs: request: %v", err)), nil
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return ErrorResponse(fmt.Sprintf("search_docs: read response: %v", err)), nil
	}

	var result any

	if err := json.Unmarshal(respBody, &result); err != nil {
		return TextResponse(string(respBody)), nil
	}

	return OkResponse(result), nil
}

// extractStringSlice converts an any value to []string.
func extractStringSlice(v any) []string {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case []string:
		return val
	case []any:
		result := make([]string, 0, len(val))

		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}

		return result
	case string:
		return []string{val}
	}

	return nil
}
