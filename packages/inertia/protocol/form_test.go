package protocol

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestFlattenJSON(t *testing.T) {
	t.Run("flattens nested map", func(t *testing.T) {
		data := map[string]any{
			"user": map[string]any{
				"name": "Alice",
				"address": map[string]any{
					"city": "Berlin",
				},
			},
			"active": true,
		}

		out := make(url.Values)

		if err := flattenJSON("", data, out, 0); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tests := map[string]string{
			"user.name":         "Alice",
			"user.address.city": "Berlin",
			"active":            "1",
		}

		for key, want := range tests {
			if got := out.Get(key); got != want {
				t.Errorf("key %q = %q, want %q", key, got, want)
			}
		}
	})

	t.Run("returns error when depth exceeds limit", func(t *testing.T) {
		data := buildDeepMap(maxJSONDepth + 2)

		out := make(url.Values)
		err := flattenJSON("", data, out, 0)

		if err == nil {
			t.Fatal("expected error for deeply nested JSON, got nil")
		}

		want := fmt.Sprintf("JSON nesting exceeds maximum depth of %d", maxJSONDepth)

		if err.Error() != want {
			t.Errorf("error = %q, want %q", err.Error(), want)
		}
	})

	t.Run("succeeds at exactly max depth", func(t *testing.T) {
		data := buildDeepMap(maxJSONDepth + 1)

		out := make(url.Values)

		if err := flattenJSON("", data, out, 0); err != nil {
			t.Fatalf("unexpected error at max depth: %v", err)
		}
	})
}

func buildDeepMap(depth int) map[string]any {
	current := map[string]any{"leaf": "value"}

	for i := depth - 1; i > 0; i-- {
		current = map[string]any{fmt.Sprintf("level%d", i): current}
	}

	return current
}

func TestParseForm_JSON(t *testing.T) {
	t.Parallel()

	body := `{"name":"Alice","age":30,"active":true}`
	req := httptest.NewRequest(http.MethodPost, "/submit?page=1", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	if err := ParseForm(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req.PostFormValue("name"); got != "Alice" {
		t.Errorf("PostForm name = %q, want %q", got, "Alice")
	}

	if got := req.PostFormValue("age"); got != "30" {
		t.Errorf("PostForm age = %q, want %q", got, "30")
	}

	if got := req.PostFormValue("active"); got != "1" {
		t.Errorf("PostForm active = %q, want %q", got, "1")
	}

	if got := req.FormValue("page"); got != "1" {
		t.Errorf("Form page = %q, want %q", got, "1")
	}

	if got := req.FormValue("name"); got != "Alice" {
		t.Errorf("Form name = %q, want %q", got, "Alice")
	}
}

func TestParseForm_JSON_BodyOverridesQuery(t *testing.T) {
	t.Parallel()

	body := `{"name":"body"}`
	req := httptest.NewRequest(http.MethodPost, "/submit?name=query", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	if err := ParseForm(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req.FormValue("name"); got != "body" {
		t.Errorf("FormValue(name) = %q, want %q (body should override query)", got, "body")
	}

	if got := req.PostFormValue("name"); got != "body" {
		t.Errorf("PostFormValue(name) = %q, want %q", got, "body")
	}
}

func TestParseForm_URLEncoded(t *testing.T) {
	t.Parallel()

	body := "name=Bob&email=bob%40example.com"
	req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err := ParseForm(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req.FormValue("name"); got != "Bob" {
		t.Errorf("name = %q, want %q", got, "Bob")
	}

	if got := req.FormValue("email"); got != "bob@example.com" {
		t.Errorf("email = %q, want %q", got, "bob@example.com")
	}
}

func TestParseForm_InvalidJSON(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader("{invalid"))

	req.Header.Set("Content-Type", "application/json")

	if err := ParseForm(req); err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParseJSONForm_NestedObjects(t *testing.T) {
	t.Parallel()

	body := `{"user":{"name":"Alice","address":{"city":"Berlin"}}}`
	req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	if err := ParseForm(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req.PostFormValue("user.name"); got != "Alice" {
		t.Errorf("user.name = %q, want %q", got, "Alice")
	}

	if got := req.PostFormValue("user.address.city"); got != "Berlin" {
		t.Errorf("user.address.city = %q, want %q", got, "Berlin")
	}
}

func TestParseJSONForm_Arrays(t *testing.T) {
	t.Parallel()

	body := `{"tags":["go","rust"],"items":[{"name":"a"}]}`
	req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	if err := ParseForm(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req.PostFormValue("tags[0]"); got != "go" {
		t.Errorf("tags[0] = %q, want %q", got, "go")
	}

	if got := req.PostFormValue("tags[1]"); got != "rust" {
		t.Errorf("tags[1] = %q, want %q", got, "rust")
	}

	if got := req.PostFormValue("items[0].name"); got != "a" {
		t.Errorf("items[0].name = %q, want %q", got, "a")
	}
}

func TestFlattenJSON_NestedArrays(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"matrix": []any{
			[]any{float64(1), float64(2)},
			[]any{float64(3), float64(4)},
		},
	}

	out := make(url.Values)

	if err := flattenJSON("", data, out, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]string{
		"matrix[0][0]": "1",
		"matrix[0][1]": "2",
		"matrix[1][0]": "3",
		"matrix[1][1]": "4",
	}

	for key, exp := range want {
		if got := out.Get(key); got != exp {
			t.Errorf("key %q = %q, want %q", key, got, exp)
		}
	}
}

func TestToFormValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"bool true", true, "1"},
		{"bool false", false, "0"},
		{"nil", nil, ""},
		{"integer float64", float64(42), "42"},
		{"fractional float64", float64(3.14), "3.14"},
		{"string", "hello", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toFormValue(tt.input)

			if got != tt.want {
				t.Errorf("toFormValue(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseJSONForm_NoURL(t *testing.T) {
	t.Parallel()

	body := `{"name":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	req.URL = nil

	if err := ParseForm(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req.PostFormValue("name"); got != "test" {
		t.Errorf("name = %q, want %q", got, "test")
	}
}
