package httpx_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx"
)

func TestInputFromQueryString(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?name=Taylor&age=30", nil)
	req := httpx.NewRequest(raw)

	if req.Input("name") != "Taylor" {
		t.Fatalf("expected Taylor, got %s", req.Input("name"))
	}

	if req.Input("age") != "30" {
		t.Fatalf("expected 30, got %s", req.Input("age"))
	}

	if req.Input("missing", "default") != "default" {
		t.Fatal("expected fallback for missing key")
	}
}

func TestInputFromJSONBody(t *testing.T) {
	t.Parallel()

	body := `{"name":"Taylor","age":30}`
	raw := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	raw.Header.Set("Content-Type", "application/json")
	req := httpx.NewRequest(raw)

	if req.Input("name") != "Taylor" {
		t.Fatalf("expected Taylor, got %s", req.Input("name"))
	}

	if req.Integer("age") != 30 {
		t.Fatalf("expected 30, got %d", req.Integer("age"))
	}
}

func TestInputFromFormBody(t *testing.T) {
	t.Parallel()

	form := url.Values{"name": {"Taylor"}, "email": {"taylor@example.com"}}
	raw := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	raw.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req := httpx.NewRequest(raw)

	if req.Input("name") != "Taylor" {
		t.Fatalf("expected Taylor, got %s", req.Input("name"))
	}

	if req.Post("email") != "taylor@example.com" {
		t.Fatalf("expected taylor@example.com, got %s", req.Post("email"))
	}
}

func TestInputBodyOverridesQuery(t *testing.T) {
	t.Parallel()

	body := `{"name":"Body"}`
	raw := httptest.NewRequest(http.MethodPost, "/?name=Query", strings.NewReader(body))
	raw.Header.Set("Content-Type", "application/json")
	req := httpx.NewRequest(raw)

	if req.Input("name") != "Body" {
		t.Fatalf("expected body value to override query, got %s", req.Input("name"))
	}
}

func TestQuery(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?page=2", nil)
	req := httpx.NewRequest(raw)

	if req.Query("page") != "2" {
		t.Fatalf("expected 2, got %s", req.Query("page"))
	}

	if req.Query("missing", "1") != "1" {
		t.Fatal("expected fallback for missing query key")
	}
}

func TestBoolean(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value    string
		expected bool
	}{
		{"1", true},
		{"true", true},
		{"on", true},
		{"yes", true},
		{"0", false},
		{"false", false},
		{"no", false},
		{"", false},
	}

	for _, tt := range tests {
		raw := httptest.NewRequest(http.MethodGet, "/?active="+tt.value, nil)
		req := httpx.NewRequest(raw)

		if req.Boolean("active") != tt.expected {
			t.Errorf("Boolean(%q) = %v, want %v", tt.value, req.Boolean("active"), tt.expected)
		}
	}
}

func TestInteger(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?count=42", nil)
	req := httpx.NewRequest(raw)

	if req.Integer("count") != 42 {
		t.Fatalf("expected 42, got %d", req.Integer("count"))
	}

	if req.Integer("missing") != 0 {
		t.Fatal("expected 0 for missing integer")
	}

	if req.Integer("missing", 10) != 10 {
		t.Fatal("expected fallback 10 for missing integer")
	}
}

func TestFloat(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?price=19.99", nil)
	req := httpx.NewRequest(raw)

	if req.Float("price") != 19.99 {
		t.Fatalf("expected 19.99, got %f", req.Float("price"))
	}

	if req.Float("missing", 0.0) != 0.0 {
		t.Fatal("expected fallback for missing float")
	}
}

func TestDate(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?date=2024-01-15", nil)
	req := httpx.NewRequest(raw)

	d, err := req.Date("date", "2006-01-02")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if d.Year() != 2024 || d.Month() != 1 || d.Day() != 15 {
		t.Fatalf("unexpected date: %v", d)
	}
}

func TestOnly(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?a=1&b=2&c=3", nil)
	req := httpx.NewRequest(raw)

	result := req.Only("a", "c")

	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}

	if result["a"] != "1" || result["c"] != "3" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestExcept(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?a=1&b=2&c=3", nil)
	req := httpx.NewRequest(raw)

	result := req.Except("b")

	if _, ok := result["b"]; ok {
		t.Fatal("b should be excluded")
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
}

func TestHas(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?a=1&b=2", nil)
	req := httpx.NewRequest(raw)

	if !req.Has("a", "b") {
		t.Fatal("expected Has to return true for a and b")
	}

	if req.Has("a", "c") {
		t.Fatal("expected Has to return false when c is missing")
	}
}

func TestHasAny(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?a=1", nil)
	req := httpx.NewRequest(raw)

	if !req.HasAny("a", "b") {
		t.Fatal("expected HasAny to return true when a is present")
	}

	if req.HasAny("x", "y") {
		t.Fatal("expected HasAny to return false when none present")
	}
}

func TestFilled(t *testing.T) {
	t.Parallel()

	body := `{"name":"Taylor","empty":""}`
	raw := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	raw.Header.Set("Content-Type", "application/json")
	req := httpx.NewRequest(raw)

	if !req.Filled("name") {
		t.Fatal("name should be filled")
	}

	if req.Filled("empty") {
		t.Fatal("empty should not be filled")
	}
}

func TestMissing(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?a=1", nil)
	req := httpx.NewRequest(raw)

	if req.Missing("a") {
		t.Fatal("a should not be missing")
	}

	if !req.Missing("b") {
		t.Fatal("b should be missing")
	}
}

func TestKeys(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?x=1&y=2", nil)
	req := httpx.NewRequest(raw)

	keys := req.Keys()

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestHeader(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("X-Custom", "hello")
	req := httpx.NewRequest(raw)

	if req.Header("X-Custom") != "hello" {
		t.Fatalf("expected hello, got %s", req.Header("X-Custom"))
	}

	if req.Header("X-Missing", "default") != "default" {
		t.Fatal("expected fallback for missing header")
	}
}

func TestHasHeader(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("X-Custom", "hello")
	req := httpx.NewRequest(raw)

	if !req.HasHeader("X-Custom") {
		t.Fatal("expected HasHeader to return true")
	}

	if req.HasHeader("X-Missing") {
		t.Fatal("expected HasHeader to return false")
	}
}

func TestBearerToken(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Authorization", "Bearer abc123")
	req := httpx.NewRequest(raw)

	if req.BearerToken() != "abc123" {
		t.Fatalf("expected abc123, got %s", req.BearerToken())
	}
}

func TestBearerTokenMissing(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httpx.NewRequest(raw)

	if req.BearerToken() != "" {
		t.Fatal("expected empty bearer token")
	}
}

func TestBearerTokenBasicAuth(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("Authorization", "Basic abc123")
	req := httpx.NewRequest(raw)

	if req.BearerToken() != "" {
		t.Fatal("expected empty bearer token for Basic auth")
	}
}

func TestCookie(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.AddCookie(&http.Cookie{Name: "session", Value: "abc"})
	req := httpx.NewRequest(raw)

	if req.Cookie("session") != "abc" {
		t.Fatalf("expected abc, got %s", req.Cookie("session"))
	}

	if req.Cookie("missing", "default") != "default" {
		t.Fatal("expected fallback for missing cookie")
	}
}

func TestHasCookie(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.AddCookie(&http.Cookie{Name: "session", Value: "abc"})
	req := httpx.NewRequest(raw)

	if !req.HasCookie("session") {
		t.Fatal("expected HasCookie to return true")
	}

	if req.HasCookie("missing") {
		t.Fatal("expected HasCookie to return false")
	}
}

func TestAllMergesQueryAndJSONBody(t *testing.T) {
	t.Parallel()

	body := `{"email":"test@example.com"}`
	raw := httptest.NewRequest(http.MethodPost, "/?name=Taylor", strings.NewReader(body))
	raw.Header.Set("Content-Type", "application/json")
	req := httpx.NewRequest(raw)

	all := req.All()

	if all["name"] != "Taylor" {
		t.Fatalf("expected Taylor from query, got %v", all["name"])
	}

	if all["email"] != "test@example.com" {
		t.Fatalf("expected test@example.com from body, got %v", all["email"])
	}
}

func TestFileUpload(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("avatar", "photo.jpg")

	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	part.Write([]byte("fake image data"))
	writer.Close()

	raw := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	raw.Header.Set("Content-Type", writer.FormDataContentType())
	req := httpx.NewRequest(raw)

	if !req.HasFile("avatar") {
		t.Fatal("expected avatar file to be present")
	}

	file := req.File("avatar")

	if file == nil {
		t.Fatal("expected non-nil file")
	}

	if file.ClientOriginalName() != "photo.jpg" {
		t.Fatalf("expected photo.jpg, got %s", file.ClientOriginalName())
	}

	if file.ClientExtension() != "jpg" {
		t.Fatalf("expected jpg, got %s", file.ClientExtension())
	}
}

func TestFileUploadMissing(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httpx.NewRequest(raw)

	if req.HasFile("avatar") {
		t.Fatal("expected no file")
	}

	if req.File("avatar") != nil {
		t.Fatal("expected nil file")
	}
}

func TestAllFiles(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	p1, _ := writer.CreateFormFile("doc", "a.pdf")
	p1.Write([]byte("pdf data"))

	p2, _ := writer.CreateFormFile("img", "b.png")
	p2.Write([]byte("png data"))

	writer.Close()

	raw := httptest.NewRequest(http.MethodPost, "/", &buf)
	raw.Header.Set("Content-Type", writer.FormDataContentType())
	req := httpx.NewRequest(raw)

	files := req.AllFiles()

	if len(files) != 2 {
		t.Fatalf("expected 2 file groups, got %d", len(files))
	}
}

func TestInputCaching(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"key": "value"})
	raw := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	raw.Header.Set("Content-Type", "application/json")
	req := httpx.NewRequest(raw)

	// Call All() twice; second call should use cache.
	all1 := req.All()
	all2 := req.All()

	if all1["key"] != all2["key"] {
		t.Fatal("cached input should return same values")
	}
}

func TestPostFallback(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httpx.NewRequest(raw)

	if req.Post("name", "fallback") != "fallback" {
		t.Fatal("expected fallback for missing POST value")
	}
}

func TestIntegerNonNumeric(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?count=abc", nil)
	req := httpx.NewRequest(raw)

	if req.Integer("count", 5) != 5 {
		t.Fatal("expected fallback for non-numeric integer")
	}
}

func TestFloatNonNumeric(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?price=abc", nil)
	req := httpx.NewRequest(raw)

	if req.Float("price", 9.99) != 9.99 {
		t.Fatal("expected fallback for non-numeric float")
	}
}
