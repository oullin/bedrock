package jsonapi_test

import (
	"testing"

	"github.com/bedrock/packages/httpx/resources/jsonapi"
)

func TestHasOneRelation(t *testing.T) {
	t.Parallel()

	rel := jsonapi.HasOne("author", "users", "42")
	resolver := jsonapi.NewRelationResolver(rel)
	result := resolver.Resolve()

	authorRel, ok := result["author"].(map[string]any)

	if !ok {
		t.Fatal("expected author relationship")
	}

	data, ok := authorRel["data"].(map[string]any)

	if !ok {
		t.Fatal("expected data linkage object")
	}

	if data["type"] != "users" {
		t.Fatalf("expected type users, got %v", data["type"])
	}

	if data["id"] != "42" {
		t.Fatalf("expected id 42, got %v", data["id"])
	}
}

func TestHasManyRelation(t *testing.T) {
	t.Parallel()

	rel := jsonapi.HasMany("comments", "comments", []string{"1", "2", "3"})
	resolver := jsonapi.NewRelationResolver(rel)
	result := resolver.Resolve()

	commentsRel, ok := result["comments"].(map[string]any)

	if !ok {
		t.Fatal("expected comments relationship")
	}

	data, ok := commentsRel["data"].([]map[string]any)

	if !ok {
		t.Fatal("expected data linkage array")
	}

	if len(data) != 3 {
		t.Fatalf("expected 3 items, got %d", len(data))
	}

	if data[0]["type"] != "comments" || data[0]["id"] != "1" {
		t.Fatalf("unexpected first item: %v", data[0])
	}
}

func TestRelationResolverMultiple(t *testing.T) {
	t.Parallel()

	resolver := jsonapi.NewRelationResolver(
		jsonapi.HasOne("author", "users", "1"),
		jsonapi.HasMany("tags", "tags", []string{"a", "b"}),
	)

	result := resolver.Resolve()

	if len(result) != 2 {
		t.Fatalf("expected 2 relationships, got %d", len(result))
	}

	if _, ok := result["author"]; !ok {
		t.Fatal("expected author relationship")
	}

	if _, ok := result["tags"]; !ok {
		t.Fatal("expected tags relationship")
	}
}

func TestRelationResolverEmpty(t *testing.T) {
	t.Parallel()

	resolver := jsonapi.NewRelationResolver()
	result := resolver.Resolve()

	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestHasOneNullRelation(t *testing.T) {
	t.Parallel()

	rel := jsonapi.Relation{
		Name:   "author",
		Type:   "users",
		IDs:    nil,
		IsMany: false,
	}

	resolver := jsonapi.NewRelationResolver(rel)
	result := resolver.Resolve()

	authorRel, ok := result["author"].(map[string]any)

	if !ok {
		t.Fatal("expected author relationship")
	}

	if authorRel["data"] != nil {
		t.Fatal("expected null data for empty to-one relation")
	}
}
