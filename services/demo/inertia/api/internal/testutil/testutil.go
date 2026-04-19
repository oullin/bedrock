package testutil

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

const TestTemplate = `<!DOCTYPE html><html><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`

var csrfMetaRe = regexp.MustCompile(
	`<meta\s[^>]*\bname="csrf-token"[^>]*\bcontent="([^"]*)"` +
		`|` +
		`<meta\s[^>]*\bcontent="([^"]*)"[^>]*\bname="csrf-token"`,
)

func FindCookie(t testing.TB, w *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()

	for _, c := range w.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}

	t.Fatalf("cookie %q not found", name)

	return nil
}

func FindCSRFMetaToken(t testing.TB, body string) string {
	t.Helper()

	m := csrfMetaRe.FindStringSubmatch(body)

	if m == nil {
		t.Fatal("csrf meta tag not found")
	}

	if strings.TrimSpace(m[1]) != "" {
		return m[1]
	}

	return m[2]
}
