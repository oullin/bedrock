package resources_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/resources"
)

func TestPaginatedResponse(t *testing.T) {
	t.Parallel()

	users := []User{
		{ID: 1, Name: "Taylor"},
		{ID: 2, Name: "Otwell"},
	}

	collection := resources.NewCollection(users, userMapper)
	collection.WrapKey = "" // paginated response handles wrapping

	meta := resources.PaginationMeta{
		CurrentPage: 1,
		LastPage:    5,
		PerPage:     2,
		Total:       10,
		From:        1,
		To:          2,
		Path:        "/users",
	}

	links := resources.PaginationLinks{
		First: "/users?page=1",
		Last:  "/users?page=5",
		Next:  "/users?page=2",
	}

	paginated := resources.NewPaginatedResponse(collection, meta, links)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	b, err := paginated.ToJSON(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var data map[string]any

	json.Unmarshal(b, &data)

	// Check data array.
	items, ok := data["data"].([]any)

	if !ok || len(items) != 2 {
		t.Fatal("expected 2 items in data")
	}

	// Check meta.
	metaMap, ok := data["meta"].(map[string]any)

	if !ok {
		t.Fatal("expected meta object")
	}

	if metaMap["current_page"].(float64) != 1 {
		t.Fatal("expected current_page 1")
	}

	if metaMap["total"].(float64) != 10 {
		t.Fatal("expected total 10")
	}

	// Check links.
	linksMap, ok := data["links"].(map[string]any)

	if !ok {
		t.Fatal("expected links object")
	}

	if linksMap["next"] != "/users?page=2" {
		t.Fatal("expected next link")
	}
}

func TestPaginatedResponseHTTP(t *testing.T) {
	t.Parallel()

	users := []User{{ID: 1, Name: "Taylor"}}
	collection := resources.NewCollection(users, userMapper)

	meta := resources.PaginationMeta{
		CurrentPage: 1,
		LastPage:    1,
		PerPage:     10,
		Total:       1,
	}

	paginated := resources.NewPaginatedResponse(collection, meta)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	err := paginated.Response(rec, req, http.StatusOK)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatal("expected application/json")
	}
}

func TestResourceResponseDefaultStatus(t *testing.T) {
	t.Parallel()

	resource := resources.NewResource(User{ID: 1, Name: "Taylor"}, userMapper)
	rr := resources.NewResourceResponse(resource)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)

	err := rr.Response(rec, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestResourceResponseRecentlyCreated(t *testing.T) {
	t.Parallel()

	resource := resources.NewResource(User{ID: 1, Name: "Taylor"}, userMapper)
	rr := resources.NewResourceResponse(resource).RecentlyCreated()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", nil)

	err := rr.Response(rec, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestResourceResponseWithHeader(t *testing.T) {
	t.Parallel()

	resource := resources.NewResource(User{ID: 1, Name: "Taylor"}, userMapper)
	rr := resources.NewResourceResponse(resource).
		RecentlyCreated().
		WithHeader("Location", "/users/1")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", nil)

	rr.Response(rec, req)

	if rec.Header().Get("Location") != "/users/1" {
		t.Fatalf("expected Location header, got %q", rec.Header().Get("Location"))
	}
}

func TestResourceResponseWithCallback(t *testing.T) {
	t.Parallel()

	resource := resources.NewResource(User{ID: 1, Name: "Taylor"}, userMapper)

	called := false
	rr := resources.NewResourceResponse(resource).
		WithResponse(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.Header().Set("X-Custom", "value")
		})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)

	rr.Response(rec, req)

	if !called {
		t.Fatal("expected callback to be called")
	}

	if rec.Header().Get("X-Custom") != "value" {
		t.Fatalf("expected X-Custom header, got %q", rec.Header().Get("X-Custom"))
	}
}
