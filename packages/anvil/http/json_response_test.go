package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	bedhttp "github.com/bedrock/packages/anvil/http"
)

// Laravel: testSetAndRetrieveData
func TestJsonResponseSetAndRetrieveData(t *testing.T) {
	t.Parallel()

	jr := bedhttp.NewJsonResponse(map[string]string{"name": "John"}, http.StatusOK)

	data, ok := jr.GetData().(map[string]string)
	if !ok {
		t.Fatal("expected map[string]string data")
	}
	if data["name"] != "John" {
		t.Fatalf("expected name 'John', got %q", data["name"])
	}

	if err := jr.SetData(map[string]string{"name": "Jane"}); err != nil {
		t.Fatalf("SetData: %v", err)
	}

	data2, ok := jr.GetData().(map[string]string)
	if !ok {
		t.Fatal("expected map[string]string after SetData")
	}
	if data2["name"] != "Jane" {
		t.Fatalf("expected name 'Jane', got %q", data2["name"])
	}
}

// Laravel: testGetOriginalContent
func TestJsonResponseGetOriginalContent(t *testing.T) {
	t.Parallel()

	original := map[string]int{"count": 42}
	jr := bedhttp.NewJsonResponse(original, http.StatusOK)

	got, ok := jr.GetOriginalContent().(map[string]int)
	if !ok {
		t.Fatal("expected map[string]int")
	}
	if got["count"] != 42 {
		t.Fatalf("expected count 42, got %d", got["count"])
	}
}

// Laravel: testSetAndRetrieveOptions (indent)
func TestJsonResponseSetAndRetrieveOptions(t *testing.T) {
	t.Parallel()

	jr := bedhttp.NewJsonResponse(map[string]string{"key": "value"}, http.StatusOK)
	jr.SetIndent("  ")

	if got := jr.GetIndent(); got != "  " {
		t.Fatalf("expected indent '  ', got %q", got)
	}

	rec := httptest.NewRecorder()
	if err := jr.WriteTo(rec); err != nil {
		t.Fatalf("WriteTo: %v", err)
	}

	body := rec.Body.String()
	// Pretty-printed JSON should contain newlines.
	if len(body) < 5 {
		t.Fatalf("expected indented JSON, got %q", body)
	}

	var parsed map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if parsed["key"] != "value" {
		t.Fatalf("expected key 'value', got %q", parsed["key"])
	}
}

// Laravel: testSetAndRetrieveDefaultOptions
func TestJsonResponseDefaultOptions(t *testing.T) {
	t.Parallel()

	jr := bedhttp.NewJsonResponse(nil, http.StatusOK)
	if got := jr.GetIndent(); got != "" {
		t.Fatalf("expected empty indent by default, got %q", got)
	}
}

// Laravel: testSetAndRetrieveStatusCode
func TestJsonResponseSetAndRetrieveStatusCode(t *testing.T) {
	t.Parallel()

	jr := bedhttp.NewJsonResponse(nil, http.StatusOK)
	if got := jr.GetStatusCode(); got != http.StatusOK {
		t.Fatalf("expected 200, got %d", got)
	}

	jr.SetStatusCode(http.StatusCreated)
	if got := jr.GetStatusCode(); got != http.StatusCreated {
		t.Fatalf("expected 201, got %d", got)
	}
}

// Laravel: testInvalidArgumentExceptionOnJsonError
func TestJsonResponseErrorOnUnencodableData(t *testing.T) {
	t.Parallel()

	// Channels cannot be JSON-encoded.
	ch := make(chan int)
	err := bedhttp.NewJsonResponse(nil, http.StatusOK).SetData(ch)
	if err == nil {
		t.Fatal("expected error for unencodable data (channel)")
	}

	// Functions cannot be JSON-encoded.
	err = bedhttp.NewJsonResponse(nil, http.StatusOK).SetData(func() {})
	if err == nil {
		t.Fatal("expected error for unencodable data (func)")
	}
}

// Laravel: testFromJsonString
func TestJsonResponseFromString(t *testing.T) {
	t.Parallel()

	raw := `{"name":"John","age":30}`
	jr, err := bedhttp.NewJsonResponseFromString(raw, http.StatusOK)
	if err != nil {
		t.Fatalf("NewJsonResponseFromString: %v", err)
	}

	if jr.GetStatusCode() != http.StatusOK {
		t.Fatalf("expected 200, got %d", jr.GetStatusCode())
	}

	// Original content should be the raw string.
	if got, ok := jr.GetOriginalContent().(string); !ok || got != raw {
		t.Fatalf("expected original content to be raw JSON string, got %v", jr.GetOriginalContent())
	}

	rec := httptest.NewRecorder()
	if err := jr.WriteTo(rec); err != nil {
		t.Fatalf("WriteTo: %v", err)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected 'application/json', got %q", ct)
	}

	// Invalid JSON string should return error.
	_, err = bedhttp.NewJsonResponseFromString("not json", http.StatusOK)
	if err == nil {
		t.Fatal("expected error for invalid JSON string")
	}
}
