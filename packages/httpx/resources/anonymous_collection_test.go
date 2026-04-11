package resources_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/resources"
)

func TestAnonymousCollectionToSlice(t *testing.T) {
	t.Parallel()

	items := []map[string]any{
		{"id": 1, "name": "Taylor"},
		{"id": 2, "name": "Otwell"},
	}

	collection := resources.NewAnonymousCollection(items)
	result := collection.ToSlice()

	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}

	if result[0]["name"] != "Taylor" {
		t.Fatalf("expected Taylor, got %v", result[0]["name"])
	}
}

func TestAnonymousCollectionToJSON(t *testing.T) {
	t.Parallel()

	items := []map[string]any{{"id": 1, "name": "Taylor"}}
	collection := resources.NewAnonymousCollection(items)

	b, err := collection.ToJSON()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var data map[string]any
	json.Unmarshal(b, &data)

	arr, ok := data["data"].([]any)

	if !ok {
		t.Fatal("expected data wrapper array")
	}

	if len(arr) != 1 {
		t.Fatalf("expected 1 item, got %d", len(arr))
	}
}

func TestAnonymousCollectionWithoutWrapping(t *testing.T) {
	t.Parallel()

	items := []map[string]any{{"id": 1, "name": "Taylor"}}
	collection := resources.NewAnonymousCollection(items)
	collection.WrapKey = ""

	b, _ := collection.ToJSON()

	var data []any
	json.Unmarshal(b, &data)

	if len(data) != 1 {
		t.Fatalf("expected 1 item without wrapper, got %d", len(data))
	}
}

func TestAnonymousCollectionWithAdditionalData(t *testing.T) {
	t.Parallel()

	items := []map[string]any{{"id": 1, "name": "Taylor"}}
	collection := resources.NewAnonymousCollection(items)
	collection.With = map[string]any{"total": 1}

	b, _ := collection.ToJSON()

	var data map[string]any
	json.Unmarshal(b, &data)

	total, ok := data["total"].(float64)

	if !ok || total != 1 {
		t.Fatal("expected total: 1")
	}
}

func TestAnonymousCollectionResponse(t *testing.T) {
	t.Parallel()

	items := []map[string]any{{"id": 1, "name": "Taylor"}}
	collection := resources.NewAnonymousCollection(items)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	err := collection.Response(rec, req, http.StatusOK)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAnonymousCollectionEmpty(t *testing.T) {
	t.Parallel()

	collection := resources.NewAnonymousCollection(nil)
	result := collection.ToSlice()

	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %d items", len(result))
	}
}

func TestAnonymousCollectionResolvesSentinels(t *testing.T) {
	t.Parallel()

	items := []map[string]any{
		{
			"id":     1,
			"name":   "Taylor",
			"secret": resources.MissingValue{},
		},
	}

	collection := resources.NewAnonymousCollection(items)
	result := collection.ToSlice()

	if _, exists := result[0]["secret"]; exists {
		t.Fatal("MissingValue should have been omitted")
	}

	if result[0]["name"] != "Taylor" {
		t.Fatalf("expected Taylor, got %v", result[0]["name"])
	}
}
