package tools_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/boost/mcp/tools"
)

func TestSearchDocsRequiresQueries(t *testing.T) {
	t.Parallel()

	tool := &tools.SearchDocs{}
	resp, err := tool.Handle(tools.McpRequest{Args: map[string]any{}})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if !resp.IsError {
		t.Error("expected error when queries arg missing")
	}
}

func TestSearchDocsCallsAPI(t *testing.T) {
	t.Parallel()

	apiResp := map[string]any{"results": []any{"doc1", "doc2"}}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, want POST", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(apiResp)
	}))

	defer srv.Close()

	tool := &tools.SearchDocs{
		APIUrl:     srv.URL,
		HTTPClient: srv.Client(),
	}

	resp, err := tool.Handle(tools.McpRequest{
		Args: map[string]any{"queries": []any{"how to use boost"}},
	})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if resp.IsError {
		t.Errorf("Handle returned error: %v", resp.Content)
	}
}

func TestSearchDocsPassesPackageFilter(t *testing.T) {
	t.Parallel()

	var capturedBody map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&capturedBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))

	defer srv.Close()

	tool := &tools.SearchDocs{APIUrl: srv.URL, HTTPClient: srv.Client()}
	_, _ = tool.Handle(tools.McpRequest{
		Args: map[string]any{
			"queries":  []any{"search term"},
			"packages": []any{"mypackage"},
		},
	})

	if pkgs, ok := capturedBody["packages"]; !ok || pkgs == nil {
		t.Error("expected 'packages' key in API request body")
	}
}

func TestSearchDocsIsReadOnly(t *testing.T) {
	t.Parallel()

	if !(&tools.SearchDocs{}).IsReadOnly() {
		t.Error("SearchDocs should be read-only")
	}
}
