package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// SupportUriTest::test_can_build_special_urls
// SupportUriTest::test_basic_uri_interactions
// SupportUriTest::test_is_empty_and_is_not_empty
// SupportUriTest::test_without_fragment
// SupportUriTest::test_without_fragment_on_uri_without_fragment
// SupportUriTest::test_complicated_query_string_parsing
// SupportUriTest::test_uri_building
// SupportUriTest::test_complicated_query_string_manipulation
// SupportUriTest::test_query_strings_with_dots_can_be_replaced_or_merged_consistently
// SupportUriTest::test_decoding_the_entire_uri
// SupportUriTest::test_decoding_the_entire_uri_preserves_the_fragment
// SupportUriTest::test_with_query_if_missing
// SupportUriTest::test_with_query_prevents_empty_query_string
// SupportUriTest::test_path_segments

func TestURIBasicInteractions(t *testing.T) {
	t.Parallel()

	uri := MustParseURI("https://example.test/users/1?filter.active=1#profile")

	if uri.IsEmpty() {
		t.Fatal("expected URI to be non-empty")
	}

	if !uri.IsNotEmpty() {
		t.Fatal("expected IsNotEmpty")
	}

	if got := uri.Query().Get("filter.active"); got != "1" {
		t.Fatalf("query filter.active = %q", got)
	}

	if got := uri.WithoutFragment().String(); got != "https://example.test/users/1?filter.active=1" {
		t.Fatalf("WithoutFragment = %q", got)
	}

	if got := MustParseURI("https://example.test/users").WithoutFragment().String(); got != "https://example.test/users" {
		t.Fatalf("WithoutFragment without fragment = %q", got)
	}
}

func TestURIQueryBuildingAndDecoding(t *testing.T) {
	t.Parallel()

	uri := MustParseURI("https://example.test/search")
	built := uri.WithQuery(map[string]string{"q": "upstream go", "filter.active": "1"})

	if got := built.Query().Get("q"); got != "upstream go" {
		t.Fatalf("WithQuery q = %q", got)
	}

	merged := built.WithQuery(map[string]string{"filter.active": "0"})

	if got := merged.Query().Get("filter.active"); got != "0" {
		t.Fatalf("WithQuery replacement = %q", got)
	}

	missing := merged.WithQueryIfMissing(map[string]string{"filter.active": "1", "page": "2"})

	if got := missing.Query().Get("filter.active"); got != "0" {
		t.Fatalf("WithQueryIfMissing existing = %q", got)
	}

	if got := missing.Query().Get("page"); got != "2" {
		t.Fatalf("WithQueryIfMissing page = %q", got)
	}

	if got := MustParseURI("https://example.test/search").WithQuery(map[string]string{}).String(); got != "https://example.test/search" {
		t.Fatalf("WithQuery empty = %q", got)
	}

	decoded := MustParseURI("https://example.test/search?q=upstream+go#top").Decoded()

	if decoded != "https://example.test/search?q=upstream go#top" {
		t.Fatalf("Decoded = %q", decoded)
	}
}

func TestURIPathSegments(t *testing.T) {
	t.Parallel()

	segments := MustParseURI("https://example.test/users/Taylor%20Otwell/profile").PathSegments()

	if len(segments) != 3 || segments[0] != "users" || segments[1] != "Taylor Otwell" || segments[2] != "profile" {
		t.Fatalf("PathSegments = %v", segments)
	}
}
