package mcp_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/mcp"
)

// Port of Upstream\Mcp\Tests\ResourceTest

func TestResourcesListReturnsStaticResourcesOnly(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(staticFileResource())
	srv.AddResource(mcp.NewResourceTemplate("users", "User template", "file://users/{id}", "text/plain",
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("user"), nil
		}))

	resp := sendRaw(t, srv, "resources/list", nil)
	result := mustResult(t, resp)
	resources, _ := result["resources"].([]any)
	if len(resources) != 1 {
		t.Fatalf("expected 1 static resource in list (not templates), got %d", len(resources))
	}
}

func TestResourcesTemplatesListReturnsTemplates(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(staticFileResource())
	srv.AddResource(mcp.NewResourceTemplate("users", "User template", "file://users/{id}", "text/plain",
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("user"), nil
		}))

	resp := sendRaw(t, srv, "resources/templates/list", nil)
	result := mustResult(t, resp)
	templates, _ := result["resourceTemplates"].([]any)
	if len(templates) != 1 {
		t.Fatalf("expected 1 template, got %d", len(templates))
	}
	tmpl := templates[0].(map[string]any)
	if tmpl["uriTemplate"] != "file://users/{id}" {
		t.Fatalf("expected uriTemplate field, got %v", tmpl["uriTemplate"])
	}
}

func TestResourcesReadMatchesStaticURI(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(staticFileResource())

	result := srv.Test(t).ReadResource("file://static/config.json")
	result.AssertOK().AssertSee("config data")
}

func TestResourcesReadMatchesURITemplate(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(mcp.NewResourceTemplate(
		"user-profile", "User profile",
		"file://users/{userId}",
		"application/json",
		func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
			uid := req.URIVars["userId"]
			return mcp.Text("profile for " + uid), nil
		},
	))

	result := srv.Test(t).ReadResource("file://users/42")
	result.AssertOK().AssertSee("profile for 42")
}

func TestResourcesReadURIVariablesExtractedCorrectly(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(mcp.NewResourceTemplate(
		"post", "Blog post",
		"file://users/{userId}/posts/{postId}",
		"text/plain",
		func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
			uid := req.URIVars["userId"]
			pid := req.URIVars["postId"]
			return mcp.Text("user=" + uid + " post=" + pid), nil
		},
	))

	result := srv.Test(t).ReadResource("file://users/7/posts/99")
	result.AssertOK().AssertSee("user=7 post=99")
}

func TestResourcesReadURINotFoundReturnsError(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	result := srv.Test(t).ReadResource("file://does/not/exist")
	result.AssertHasErrors()
}

func TestResourcesReadResponseIncludesURI(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(staticFileResource())

	resp := sendRaw(t, srv, "resources/read", map[string]any{"uri": "file://static/config.json"})
	result := mustResult(t, resp)
	contents, _ := result["contents"].([]any)
	if len(contents) == 0 {
		t.Fatal("expected non-empty contents")
	}
	content := contents[0].(map[string]any)
	if content["uri"] != "file://static/config.json" {
		t.Fatalf("expected uri in content, got %v", content["uri"])
	}
}

func TestResourcesListPagination(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0", mcp.WithPagination(2, 10))
	for i := 0; i < 5; i++ {
		srv.AddResource(indexedResource(i))
	}

	resp := sendRaw(t, srv, "resources/list", nil)
	result := mustResult(t, resp)
	resources, _ := result["resources"].([]any)
	if len(resources) != 2 {
		t.Fatalf("expected 2 resources on first page, got %d", len(resources))
	}
	if result["nextCursor"] == nil {
		t.Fatal("expected nextCursor on first page")
	}
}

// --- helpers ---

func staticFileResource() mcp.Resource {
	return mcp.NewResource(
		"config", "Configuration file",
		"file://static/config.json",
		"application/json",
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("config data"), nil
		},
	)
}

func indexedResource(i int) mcp.Resource {
	uri := "file://res-" + string(rune('a'+i))
	return mcp.NewResource("res-"+string(rune('a'+i)), "Resource", uri, "text/plain",
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("data"), nil
		})
}
