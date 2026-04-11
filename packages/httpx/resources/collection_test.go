package resources_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/resources"
)

func TestCollectionToSlice(t *testing.T) {
	t.Parallel()

	users := []User{
		{ID: 1, Name: "Taylor"},
		{ID: 2, Name: "Otwell"},
	}

	collection := resources.NewCollection(users, userMapper)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := collection.ToSlice(req)

	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}

	if result[0]["name"] != "Taylor" {
		t.Fatalf("expected Taylor, got %v", result[0]["name"])
	}
}

func TestCollectionToJSON(t *testing.T) {
	t.Parallel()

	users := []User{{ID: 1, Name: "Taylor"}}
	collection := resources.NewCollection(users, userMapper)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	b, err := collection.ToJSON(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var data map[string]any

	json.Unmarshal(b, &data)

	items, ok := data["data"].([]any)

	if !ok {
		t.Fatal("expected data wrapper array")
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

func TestCollectionWithoutWrapping(t *testing.T) {
	t.Parallel()

	users := []User{{ID: 1, Name: "Taylor"}}
	collection := resources.NewCollection(users, userMapper)
	collection.WrapKey = ""

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	b, _ := collection.ToJSON(req)

	var data []any

	json.Unmarshal(b, &data)

	if len(data) != 1 {
		t.Fatalf("expected 1 item without wrapper, got %d", len(data))
	}
}

func TestCollectionWithAdditionalData(t *testing.T) {
	t.Parallel()

	users := []User{{ID: 1, Name: "Taylor"}}
	collection := resources.NewCollection(users, userMapper)
	collection.With = map[string]any{"total": 1}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	b, _ := collection.ToJSON(req)

	var data map[string]any

	json.Unmarshal(b, &data)

	total, ok := data["total"].(float64)

	if !ok || total != 1 {
		t.Fatal("expected total: 1")
	}
}

func TestCollectionResponse(t *testing.T) {
	t.Parallel()

	users := []User{{ID: 1, Name: "Taylor"}}
	collection := resources.NewCollection(users, userMapper)

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

func TestCollectionEmptySlice(t *testing.T) {
	t.Parallel()

	var users []User
	collection := resources.NewCollection(users, userMapper)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	result := collection.ToSlice(req)

	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %d items", len(result))
	}
}
