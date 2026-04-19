package mcp_test

import (
	"testing"

	"github.com/bedrock/packages/ai/mcp"
)

// Port of Laravel\Mcp\Tests\UriTemplateTest

func TestUriTemplateStaticURIMatchesExactly(t *testing.T) {
	t.Parallel()

	tmpl, err := mcp.NewUriTemplate("file://static/path")

	if err != nil {
		t.Fatal(err)
	}

	vars, ok := tmpl.Match("file://static/path")

	if !ok {
		t.Fatal("expected match for identical static URI")
	}

	if len(vars) != 0 {
		t.Fatalf("expected no vars, got %v", vars)
	}
}

func TestUriTemplateStaticURIDoesNotMatchDifferent(t *testing.T) {
	t.Parallel()

	tmpl, _ := mcp.NewUriTemplate("file://static/path")
	_, ok := tmpl.Match("file://other/path")

	if ok {
		t.Fatal("expected no match for different URI")
	}
}

func TestUriTemplateSingleVariableExtracted(t *testing.T) {
	t.Parallel()

	tmpl, err := mcp.NewUriTemplate("file://users/{id}")

	if err != nil {
		t.Fatal(err)
	}

	vars, ok := tmpl.Match("file://users/42")

	if !ok {
		t.Fatal("expected match")
	}

	if vars["id"] != "42" {
		t.Fatalf("expected id=42, got %q", vars["id"])
	}
}

func TestUriTemplateMultipleVariablesExtracted(t *testing.T) {
	t.Parallel()

	tmpl, err := mcp.NewUriTemplate("file://users/{userId}/posts/{postId}")

	if err != nil {
		t.Fatal(err)
	}

	vars, ok := tmpl.Match("file://users/7/posts/99")

	if !ok {
		t.Fatal("expected match")
	}

	if vars["userId"] != "7" {
		t.Fatalf("expected userId=7, got %q", vars["userId"])
	}

	if vars["postId"] != "99" {
		t.Fatalf("expected postId=99, got %q", vars["postId"])
	}
}

func TestUriTemplateNoMatchWhenExtraSegments(t *testing.T) {
	t.Parallel()

	tmpl, _ := mcp.NewUriTemplate("file://users/{id}")
	_, ok := tmpl.Match("file://users/42/extra")

	if ok {
		t.Fatal("expected no match when URI has extra segments beyond template")
	}
}

func TestUriTemplateExpandFillsVariables(t *testing.T) {
	t.Parallel()

	tmpl, _ := mcp.NewUriTemplate("file://users/{id}/posts/{postId}")
	result := tmpl.Expand(map[string]string{"id": "5", "postId": "10"})

	if result != "file://users/5/posts/10" {
		t.Fatalf("unexpected expansion: %q", result)
	}
}

func TestUriTemplateExpandLeavesUnknownVarsAsIs(t *testing.T) {
	t.Parallel()

	tmpl, _ := mcp.NewUriTemplate("file://users/{id}")
	result := tmpl.Expand(map[string]string{})

	if result != "file://users/{id}" {
		t.Fatalf("expected template unchanged, got %q", result)
	}
}

func TestUriTemplateEmptyVariableNameIsError(t *testing.T) {
	t.Parallel()

	_, err := mcp.NewUriTemplate("file://users/{}/posts")

	if err == nil {
		t.Fatal("expected error for empty variable name")
	}
}

func TestUriTemplateTemplateReturnsOriginalString(t *testing.T) {
	t.Parallel()

	raw := "file://users/{id}"
	tmpl, _ := mcp.NewUriTemplate(raw)

	if tmpl.Template() != raw {
		t.Fatalf("expected %q, got %q", raw, tmpl.Template())
	}
}

func TestIsTemplateDetectsPlaceholders(t *testing.T) {
	t.Parallel()

	if !mcp.IsTemplate("file://users/{id}") {
		t.Fatal("expected IsTemplate true for URI with placeholder")
	}

	if mcp.IsTemplate("file://users/42") {
		t.Fatal("expected IsTemplate false for static URI")
	}
}
