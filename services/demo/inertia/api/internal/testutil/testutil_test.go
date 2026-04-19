package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockTB struct {
	testing.TB
	failed bool
}

func (m *mockTB) Helper()               {}
func (m *mockTB) Fatalf(string, ...any) { m.failed = true }

func TestFindCSRFMetaToken(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "name before content",
			body: `<meta name="csrf-token" content="abc123">`,
			want: "abc123",
		},
		{
			name: "content before name",
			body: `<meta content="xyz789" name="csrf-token">`,
			want: "xyz789",
		},
		{
			name: "extra attributes between",
			body: `<meta name="csrf-token" id="tok" content="mid456">`,
			want: "mid456",
		},
		{
			name: "self-closing tag",
			body: `<meta name="csrf-token" content="sc1" />`,
			want: "sc1",
		},
		{
			name: "embedded in full page",
			body: `<!DOCTYPE html><html><head><meta name="csrf-token" content="page1"></head><body></body></html>`,
			want: "page1",
		},
		{
			name: "content before name with extra attrs",
			body: `<meta class="x" content="rev2" name="csrf-token">`,
			want: "rev2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindCSRFMetaToken(t, tt.body)

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFindCookie(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	http.SetCookie(w, &http.Cookie{Name: "session", Value: "encrypted"})

	cookie := FindCookie(t, w, "session")

	if cookie.Value != "encrypted" {
		t.Fatalf("FindCookie() value = %q, want %q", cookie.Value, "encrypted")
	}
}

func TestFindCookie_Missing(t *testing.T) {
	t.Parallel()

	m := &mockTB{}
	w := httptest.NewRecorder()

	FindCookie(m, w, "missing")

	if !m.failed {
		t.Fatal("FindCookie() should fail when the cookie is missing")
	}
}
