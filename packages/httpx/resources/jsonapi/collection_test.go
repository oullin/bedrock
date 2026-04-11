package jsonapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/resources/jsonapi"
)

type includedResource struct {
	typeName string
	id       string
	name     string
}

func TestCollectionToDocument(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{ID: "1", Title: "First", Body: "Body 1"},
		{ID: "2", Title: "Second", Body: "Body 2"},
	}

	collection := jsonapi.NewCollection(posts, func(p Post) *jsonapi.JsonApiResource[Post] {
		return postResource(p)
	})

	req := httptest.NewRequest("GET", "/posts", nil)
	doc := collection.ToDocument(req)

	data, ok := doc["data"].([]map[string]any)

	if !ok {
		t.Fatal("expected data array")
	}

	if len(data) != 2 {
		t.Fatalf("expected 2 items, got %d", len(data))
	}

	if data[0]["type"] != "posts" {
		t.Fatalf("expected type posts, got %v", data[0]["type"])
	}
}

func TestCollectionToJSON(t *testing.T) {
	t.Parallel()

	posts := []Post{{ID: "1", Title: "Hello", Body: "World"}}

	collection := jsonapi.NewCollection(posts, func(p Post) *jsonapi.JsonApiResource[Post] {
		return postResource(p)
	})

	req := httptest.NewRequest("GET", "/posts", nil)
	b, err := collection.ToJSON(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc map[string]any

	json.Unmarshal(b, &doc)

	if _, ok := doc["data"]; !ok {
		t.Fatal("expected data key")
	}
}

func TestCollectionResponse(t *testing.T) {
	t.Parallel()

	posts := []Post{{ID: "1", Title: "Hello", Body: "World"}}

	collection := jsonapi.NewCollection(posts, func(p Post) *jsonapi.JsonApiResource[Post] {
		return postResource(p)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/posts", nil)

	err := collection.Response(rec, req, http.StatusOK)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/vnd.api+json" {
		t.Fatalf("expected application/vnd.api+json, got %s", ct)
	}
}

func (r includedResource) ToResourceObject(_ *http.Request) map[string]any {
	return map[string]any{
		"type":       r.typeName,
		"id":         r.id,
		"attributes": map[string]any{"name": r.name},
	}
}

func TestCollectionWithIncluded(t *testing.T) {
	t.Parallel()

	posts := []Post{{ID: "1", Title: "Hello", Body: "World"}}

	collection := jsonapi.NewCollection(posts, func(p Post) *jsonapi.JsonApiResource[Post] {
		return postResource(p)
	}).WithIncluded(
		includedResource{typeName: "users", id: "42", name: "Taylor"},
	)

	req := httptest.NewRequest("GET", "/posts?include=author", nil)
	doc := collection.ToDocument(req)

	included, ok := doc["included"].([]map[string]any)

	if !ok {
		t.Fatal("expected included array")
	}

	if len(included) != 1 {
		t.Fatalf("expected 1 included, got %d", len(included))
	}

	if included[0]["type"] != "users" {
		t.Fatalf("expected type users, got %v", included[0]["type"])
	}
}

func TestCollectionDeduplicatesIncluded(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{ID: "1", Title: "First", Body: "Body 1"},
		{ID: "2", Title: "Second", Body: "Body 2"},
	}

	author := includedResource{typeName: "users", id: "42", name: "Taylor"}

	collection := jsonapi.NewCollection(posts, func(p Post) *jsonapi.JsonApiResource[Post] {
		return postResource(p)
	}).WithIncluded(author, author) // duplicate

	req := httptest.NewRequest("GET", "/posts", nil)
	doc := collection.ToDocument(req)

	included, ok := doc["included"].([]map[string]any)

	if !ok {
		t.Fatal("expected included array")
	}

	if len(included) != 1 {
		t.Fatalf("expected 1 deduplicated included, got %d", len(included))
	}
}

func TestCollectionWithMeta(t *testing.T) {
	t.Parallel()

	posts := []Post{{ID: "1", Title: "Hello", Body: "World"}}

	collection := jsonapi.NewCollection(posts, func(p Post) *jsonapi.JsonApiResource[Post] {
		return postResource(p)
	}).WithMeta(map[string]any{"total": 100})

	req := httptest.NewRequest("GET", "/posts", nil)
	doc := collection.ToDocument(req)

	meta, ok := doc["meta"].(map[string]any)

	if !ok {
		t.Fatal("expected meta")
	}

	if meta["total"] != 100 {
		t.Fatalf("expected total 100, got %v", meta["total"])
	}
}

func TestCollectionWithLinks(t *testing.T) {
	t.Parallel()

	posts := []Post{{ID: "1", Title: "Hello", Body: "World"}}

	collection := jsonapi.NewCollection(posts, func(p Post) *jsonapi.JsonApiResource[Post] {
		return postResource(p)
	}).WithLinks(map[string]any{
		"self": "/posts",
		"next": "/posts?page[number]=2",
	})

	req := httptest.NewRequest("GET", "/posts", nil)
	doc := collection.ToDocument(req)

	links, ok := doc["links"].(map[string]any)

	if !ok {
		t.Fatal("expected links")
	}

	if links["self"] != "/posts" {
		t.Fatalf("expected self link, got %v", links["self"])
	}
}

func TestCollectionEmpty(t *testing.T) {
	t.Parallel()

	var posts []Post

	collection := jsonapi.NewCollection(posts, func(p Post) *jsonapi.JsonApiResource[Post] {
		return postResource(p)
	})

	req := httptest.NewRequest("GET", "/posts", nil)
	doc := collection.ToDocument(req)

	data, ok := doc["data"].([]map[string]any)

	if !ok {
		t.Fatal("expected data array")
	}

	if len(data) != 0 {
		t.Fatalf("expected empty data, got %d items", len(data))
	}
}
