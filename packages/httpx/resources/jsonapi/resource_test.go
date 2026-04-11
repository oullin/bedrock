package jsonapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/resources/jsonapi"
)

type Post struct {
	ID     string
	Title  string
	Body   string
	Author string
}

func postResource(p Post) *jsonapi.JsonApiResource[Post] {
	return jsonapi.NewJsonApiResource(
		p,
		func(p Post) string { return "posts" },
		func(p Post) string { return p.ID },
		func(p Post, r *http.Request) map[string]any {
			return map[string]any{
				"title": p.Title,
				"body":  p.Body,
			}
		},
	)
}

func TestResourceObjectBasic(t *testing.T) {
	t.Parallel()

	post := Post{ID: "1", Title: "Hello", Body: "World"}
	resource := postResource(post)

	req := httptest.NewRequest("GET", "/posts/1", nil)
	obj := resource.ToResourceObject(req)

	if obj["type"] != "posts" {
		t.Fatalf("expected type posts, got %v", obj["type"])
	}

	if obj["id"] != "1" {
		t.Fatalf("expected id 1, got %v", obj["id"])
	}

	attrs, ok := obj["attributes"].(map[string]any)

	if !ok {
		t.Fatal("expected attributes map")
	}

	if attrs["title"] != "Hello" {
		t.Fatalf("expected title Hello, got %v", attrs["title"])
	}
}

func TestResourceObjectWithRelationships(t *testing.T) {
	t.Parallel()

	post := Post{ID: "1", Title: "Hello", Body: "World", Author: "42"}
	resource := postResource(post).
		WithRelationships(func(p Post, r *http.Request) []jsonapi.Relation {
			return []jsonapi.Relation{
				jsonapi.HasOne("author", "users", p.Author),
			}
		})

	req := httptest.NewRequest("GET", "/posts/1", nil)
	obj := resource.ToResourceObject(req)

	rels, ok := obj["relationships"].(map[string]any)

	if !ok {
		t.Fatal("expected relationships")
	}

	authorRel, ok := rels["author"].(map[string]any)

	if !ok {
		t.Fatal("expected author relationship")
	}

	data, ok := authorRel["data"].(map[string]any)

	if !ok {
		t.Fatal("expected data linkage")
	}

	if data["id"] != "42" {
		t.Fatalf("expected author id 42, got %v", data["id"])
	}
}

func TestResourceObjectWithLinks(t *testing.T) {
	t.Parallel()

	post := Post{ID: "1", Title: "Hello", Body: "World"}
	resource := postResource(post).
		WithLinks(func(p Post, r *http.Request) map[string]any {
			return map[string]any{"self": "/posts/" + p.ID}
		})

	req := httptest.NewRequest("GET", "/posts/1", nil)
	obj := resource.ToResourceObject(req)

	links, ok := obj["links"].(map[string]any)

	if !ok {
		t.Fatal("expected links")
	}

	if links["self"] != "/posts/1" {
		t.Fatalf("expected self link, got %v", links["self"])
	}
}

func TestResourceObjectWithMeta(t *testing.T) {
	t.Parallel()

	post := Post{ID: "1", Title: "Hello", Body: "World"}
	resource := postResource(post).
		WithMeta(func(p Post, r *http.Request) map[string]any {
			return map[string]any{"version": 2}
		})

	req := httptest.NewRequest("GET", "/posts/1", nil)
	obj := resource.ToResourceObject(req)

	meta, ok := obj["meta"].(map[string]any)

	if !ok {
		t.Fatal("expected meta")
	}

	if meta["version"] != 2 {
		t.Fatalf("expected version 2, got %v", meta["version"])
	}
}

func TestResourceSparseFieldsets(t *testing.T) {
	t.Parallel()

	post := Post{ID: "1", Title: "Hello", Body: "World"}
	resource := postResource(post)

	req := httptest.NewRequest("GET", "/posts/1?fields[posts]=title", nil)
	obj := resource.ToResourceObject(req)

	attrs, ok := obj["attributes"].(map[string]any)

	if !ok {
		t.Fatal("expected attributes")
	}

	if _, exists := attrs["body"]; exists {
		t.Fatal("body should have been filtered by sparse fieldsets")
	}

	if attrs["title"] != "Hello" {
		t.Fatalf("expected title Hello, got %v", attrs["title"])
	}
}

func TestResourceToJSON(t *testing.T) {
	t.Parallel()

	post := Post{ID: "1", Title: "Hello", Body: "World"}
	resource := postResource(post)

	req := httptest.NewRequest("GET", "/posts/1", nil)
	b, err := resource.ToJSON(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc map[string]any

	json.Unmarshal(b, &doc)

	data, ok := doc["data"].(map[string]any)

	if !ok {
		t.Fatal("expected data wrapper")
	}

	if data["type"] != "posts" {
		t.Fatalf("expected type posts, got %v", data["type"])
	}
}

func TestResourceResponse(t *testing.T) {
	t.Parallel()

	post := Post{ID: "1", Title: "Hello", Body: "World"}
	resource := postResource(post)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/posts/1", nil)

	err := resource.Response(rec, req, http.StatusOK)

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

func TestResourceIdentificationError(t *testing.T) {
	t.Parallel()

	err := &jsonapi.ResourceIdentificationError{
		Field:    "type",
		Resource: Post{},
	}

	expected := "jsonapi: unable to resolve resource type for jsonapi_test.Post"

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}
