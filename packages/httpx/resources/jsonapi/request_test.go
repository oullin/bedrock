package jsonapi_test

import (
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx/resources/jsonapi"
)

func TestRequestIncludes(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/posts?include=author,comments", nil)
	r := jsonapi.NewRequest(req)

	includes := r.Includes()

	if len(includes) != 2 {
		t.Fatalf("expected 2 includes, got %d", len(includes))
	}

	if includes[0] != "author" || includes[1] != "comments" {
		t.Fatalf("expected [author, comments], got %v", includes)
	}
}

func TestRequestIncludesEmpty(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/posts", nil)
	r := jsonapi.NewRequest(req)

	if includes := r.Includes(); includes != nil {
		t.Fatalf("expected nil, got %v", includes)
	}
}

func TestRequestFields(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/posts?fields[posts]=title,body", nil)
	r := jsonapi.NewRequest(req)

	fields := r.Fields("posts")

	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}

	if fields[0] != "title" || fields[1] != "body" {
		t.Fatalf("expected [title, body], got %v", fields)
	}
}

func TestRequestFieldsAbsent(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/posts", nil)
	r := jsonapi.NewRequest(req)

	if fields := r.Fields("posts"); fields != nil {
		t.Fatalf("expected nil, got %v", fields)
	}
}

func TestRequestSort(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/posts?sort=-created_at,title", nil)
	r := jsonapi.NewRequest(req)

	sort := r.Sort()

	if len(sort) != 2 {
		t.Fatalf("expected 2 sort fields, got %d", len(sort))
	}

	if sort[0].Field != "created_at" || !sort[0].Descending {
		t.Fatalf("expected descending created_at, got %+v", sort[0])
	}

	if sort[1].Field != "title" || sort[1].Descending {
		t.Fatalf("expected ascending title, got %+v", sort[1])
	}
}

func TestRequestSortEmpty(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/posts", nil)
	r := jsonapi.NewRequest(req)

	if sort := r.Sort(); sort != nil {
		t.Fatalf("expected nil, got %v", sort)
	}
}

func TestRequestFilter(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/posts?filter[status]=published", nil)
	r := jsonapi.NewRequest(req)

	if v := r.Filter("status"); v != "published" {
		t.Fatalf("expected published, got %s", v)
	}
}

func TestRequestFilterAbsent(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/posts", nil)
	r := jsonapi.NewRequest(req)

	if v := r.Filter("status"); v != "" {
		t.Fatalf("expected empty, got %s", v)
	}
}

func TestRequestPage(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/posts?page[number]=2&page[size]=10", nil)
	r := jsonapi.NewRequest(req)

	page := r.Page()

	if page["number"] != "2" {
		t.Fatalf("expected page number 2, got %s", page["number"])
	}

	if page["size"] != "10" {
		t.Fatalf("expected page size 10, got %s", page["size"])
	}
}

func TestRequestPageAbsent(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/posts", nil)
	r := jsonapi.NewRequest(req)

	if page := r.Page(); page != nil {
		t.Fatalf("expected nil, got %v", page)
	}
}
