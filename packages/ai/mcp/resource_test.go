package mcp_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/ai/mcp"
)

// Port of Laravel\Mcp\Tests\ResourceTest

func TestResourcesListReturnsStaticResourcesOnly(t *testing.T) {
	t.Parallel()

	// ListResourcesTest::it_returns_a_valid_populated_list_resources_response
	// ListResourcesTest::it_excludes_resource_templates_from_list
	// ListResourcesTest::it_returns_only_static_resources_when_both_templates_and_static_resources_exist
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

	// ReadResourceTest::it_returns_a_valid_resource_result
	// ReadResourceTest::it_returns_a_valid_resource_result_for_text_resources
	// ReadResourceTest::it_handles_a_text_content_object_returned_from_read
	// ReadResourceTest::it_returns_no_meta_by_default
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(staticFileResource())

	result := srv.Test(t).ReadResource("file://static/config.json")
	result.AssertOK().AssertSee("config data")

	resp := sendRaw(t, srv, "resources/read", map[string]any{"uri": "file://static/config.json"})
	decoded := mustResult(t, resp)
	contents, _ := decoded["contents"].([]any)
	content := contents[0].(map[string]any)

	if content["text"] != "config data" {
		t.Fatalf("expected text content payload, got %#v", content)
	}

	if _, ok := content["_meta"]; ok {
		t.Fatalf("expected no content meta for static text resource, got %#v", content)
	}
}

func TestResourcesReadMatchesURITemplate(t *testing.T) {
	t.Parallel()

	// ReadResourceTest::it_reads_resource_template_by_matching_a_uri_pattern
	// ReadResourceTest::it_returns_the_actual_requested_uri_in_response_not_the_template_pattern
	// ReadResourceTest::it_extracts_variables_from_uri_template_and_passes_to_handler
	// ReadResourceTest::it_sets_uri_on_request_when_reading_resource_templates
	// ReadResourceTest::it_provides_both_uri_and_extracted_variables_in_request_for_templates
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(mcp.NewResourceTemplate(
		"user-profile", "User profile",
		"file://users/{userId}",
		"application/json",
		func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
			uid := req.URIVars["userId"]

			if req.URI != "file://users/42" {
				t.Fatalf("expected request URI to be preserved, got %q", req.URI)
			}

			if req.Get("userId") != "42" {
				t.Fatalf("expected extracted URI variable to be available via request getters, got %#v", req.Get("userId"))
			}

			return mcp.Text("profile for " + uid), nil
		},
	))

	result := srv.Test(t).ReadResource("file://users/42")
	result.AssertOK().AssertSee("profile for 42")
}

func TestResourcesReadURIVariablesExtractedCorrectly(t *testing.T) {
	t.Parallel()

	// ReadResourceTest::it_tries_static_resources_before_template_matching
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

func TestResourcesReadReturnsBlobContentPayload(t *testing.T) {
	t.Parallel()

	// ReadResourceTest::it_returns_a_valid_resource_result_for_blob_resources
	// ReadResourceTest::it_handles_a_blob_content_object_returned_from_read
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(mcp.NewResource(
		"binary", "Binary resource",
		"file://binary/blob",
		"application/octet-stream",
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Blob([]byte("binary data"), "application/octet-stream"), nil
		},
	))

	resp := sendRaw(t, srv, "resources/read", map[string]any{"uri": "file://binary/blob"})
	result := mustResult(t, resp)
	contents, _ := result["contents"].([]any)

	if len(contents) != 1 {
		t.Fatalf("expected one blob content item, got %#v", contents)
	}

	content := contents[0].(map[string]any)

	if content["uri"] != "file://binary/blob" {
		t.Fatalf("expected blob resource uri to be preserved, got %#v", content["uri"])
	}

	if content["mimeType"] != "application/octet-stream" {
		t.Fatalf("expected blob mimeType to be preserved, got %#v", content["mimeType"])
	}

	blob, _ := content["blob"].(string)

	if blob == "" {
		t.Fatal("expected blob payload to be present")
	}
}

func TestResourcesReadURINotFoundReturnsError(t *testing.T) {
	t.Parallel()

	// ReadResourceTest::it_throws_exception_when_resource_is_not_found
	srv := mcp.NewServer("srv", "1.0.0")
	result := srv.Test(t).ReadResource("file://does/not/exist")
	result.AssertHasErrors()
}

func TestResourcesListReturnsEmptyWhenNothingRegistered(t *testing.T) {
	t.Parallel()

	// ListResourcesTest::it_returns_a_valid_empty_list_resources_response
	srv := mcp.NewServer("srv", "1.0.0")

	resp := sendRaw(t, srv, "resources/list", nil)
	result := mustResult(t, resp)
	resources, _ := result["resources"].([]any)

	if len(resources) != 0 {
		t.Fatalf("expected empty resource list, got %d entries", len(resources))
	}
}

func TestResourcesReadMissingURIReturnsError(t *testing.T) {
	t.Parallel()

	// ReadResourceTest::it_throws_error_when_uri_is_missing
	srv := mcp.NewServer("srv", "1.0.0")

	resp := sendRaw(t, srv, "resources/read", map[string]any{})

	if _, ok := resp["error"].(map[string]any); !ok {
		t.Fatalf("expected missing resource uri to return an error, got %#v", resp)
	}
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
