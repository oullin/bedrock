package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/services/brain/api/ai"
	"github.com/bedrock/services/brain/api/graph"
)

// fixtureRoutes is a minimal Go service the analyzer can scan into a
// non-empty graph: one router type with two routes.

// newFixtureTarget writes a tiny Go module under a temp dir so
// analysis.AnalyzeTarget can produce a real graph for handler tests.

// serverFixture returns a Server pointed at a freshly-built fixture target.
// It does NOT call EnsureScanned — callers decide whether to scan up-front.

// scannedFixture returns a server with a scan already completed.

// Pre-create CLAUDE.md to trigger conflict.

// Empty body → JSON decode error → 400.

// Valid JSON but URL omitted → 400.

type errSentinel struct{}

const fixtureRoutes = `package fixture

type router struct{}

func (r *router) Get(uri string, fn func())  {}
func (r *router) Post(uri string, fn func()) {}

func wire(r *router) {
	r.Get("/", func() {})
	r.Post("/login", func() {})
}
`

func newFixtureTarget(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	mustWriteFile(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWriteFile(t, dir, "routes.go", fixtureRoutes)

	return dir
}

func mustWriteFile(t *testing.T, dir, name, body string) {
	t.Helper()
	full := filepath.Join(dir, name)

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", name, err)
	}

	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func serverFixture(t *testing.T) *Server {
	t.Helper()

	return NewServer(newFixtureTarget(t), []byte("<!doctype html><body>brain</body>"), "")
}

func scannedFixture(t *testing.T) *Server {
	t.Helper()
	s := serverFixture(t)

	if err := s.EnsureScanned(); err != nil {
		t.Fatalf("EnsureScanned: %v", err)
	}

	return s
}

func TestIsLocalURL(t *testing.T) {
	cases := []struct {
		raw  string
		want bool
	}{
		{"http://localhost/", true},
		{"http://localhost:8080/", true},
		{"http://127.0.0.1/", true},
		{"http://127.0.0.1:9000/api", true},
		{"http://[::1]:8080/", true},
		{"http://example.com/", false},
		{"https://api.example.com/", false},
		{"http://192.168.1.1/", false},
		{"://malformed", false},
	}

	for _, c := range cases {
		if got := isLocalURL(c.raw); got != c.want {
			t.Errorf("isLocalURL(%q) = %v, want %v", c.raw, got, c.want)
		}
	}
}

func TestEscapePathDelegates(t *testing.T) {
	got := EscapePath("foo/bar baz")
	want := url.PathEscape("foo/bar baz")

	if got != want {
		t.Errorf("EscapePath = %q, want %q", got, want)
	}
}

func TestNewServerInitializesFields(t *testing.T) {
	s := NewServer("/x", []byte("html"), "/assets")

	if s.Target != "/x" {
		t.Errorf("Target = %q", s.Target)
	}

	if s.Analyzer == nil {
		t.Error("Analyzer is nil")
	}

	if string(s.HTMLBody) != "html" {
		t.Errorf("HTMLBody = %q", s.HTMLBody)
	}

	if s.AssetDir != "/assets" {
		t.Errorf("AssetDir = %q", s.AssetDir)
	}
}

func TestHandleSPAUsesFallbackWhenBodyEmpty(t *testing.T) {
	s := NewServer("/x", nil, "")
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("GET", "/_brain/", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "SPA build missing") {
		t.Errorf("missing fallback message: %s", rr.Body.String())
	}
}

func TestHandleSPAServesEmbeddedBody(t *testing.T) {
	body := "<html><body>brain ui</body></html>"
	s := NewServer("/x", []byte(body), "")
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("GET", "/_brain/anything", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}

	if rr.Body.String() != body {
		t.Errorf("body = %q, want %q", rr.Body, body)
	}

	if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q", ct)
	}
}

func TestHandleManifestReturnsScannedMeta(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("GET", "/_brain/api/manifest", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rr.Code, rr.Body)
	}

	var m graph.Manifest

	if err := json.NewDecoder(rr.Body).Decode(&m); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if m.TotalRoutes < 2 {
		t.Errorf("TotalRoutes = %d, want >= 2", m.TotalRoutes)
	}

	if m.TotalNodes == 0 {
		t.Errorf("TotalNodes = 0; want non-zero from fixture")
	}
}

func TestHandleGraphReturnsFullGraph(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("GET", "/_brain/api/graph", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}

	var g graph.Graph

	if err := json.NewDecoder(rr.Body).Decode(&g); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(g.Nodes) == 0 {
		t.Errorf("graph has no nodes")
	}
}

func TestHandleSourceRequiresPathParam(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("GET", "/_brain/api/source", nil))

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestHandleSourceRejectsPathTraversal(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/_brain/api/source?path=../../etc/passwd", nil)
	s.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400; body=%s", rr.Code, rr.Body)
	}

	if !strings.Contains(rr.Body.String(), "outside target") {
		t.Errorf("body = %s, want 'outside target' message", rr.Body)
	}
}

func TestHandleSourceReturnsFile(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("GET", "/_brain/api/source?path=routes.go", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rr.Code, rr.Body)
	}

	var got map[string]string

	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got["path"] != "routes.go" || !strings.Contains(got["body"], "package fixture") {
		t.Errorf("payload = %v, want routes.go body", got)
	}
}

func TestHandleScanReturnsManifest(t *testing.T) {
	s := serverFixture(t)
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("POST", "/_brain/api/scan", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rr.Code, rr.Body)
	}

	var m graph.Manifest

	if err := json.NewDecoder(rr.Body).Decode(&m); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if m.TotalNodes == 0 {
		t.Errorf("TotalNodes = 0 after scan")
	}
}

func TestHandleContextMarkdown(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("GET", "/_brain/api/context", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/markdown") {
		t.Errorf("Content-Type = %q, want text/markdown", ct)
	}

	if !strings.HasPrefix(rr.Body.String(), "# ") {
		t.Errorf("body doesn't start with '# ': %s", rr.Body.String()[:min(40, rr.Body.Len())])
	}
}

func TestHandleContextJSON(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("GET", "/_brain/api/context?format=json", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}

	var got map[string]any

	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if _, ok := got["markdown"].(string); !ok {
		t.Errorf("missing markdown field: %v", got)
	}

	if v, ok := got["tokenEstimate"].(float64); !ok || v <= 0 {
		t.Errorf("tokenEstimate = %v, want positive number", got["tokenEstimate"])
	}
}

func TestHandleGenerateRulesWritesAllTargets(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/_brain/api/generate-rules", bytes.NewReader([]byte(`{"force":false}`)))
	req.Header.Set("Content-Length", "16")
	s.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rr.Code, rr.Body)
	}

	var got map[string]any

	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	written, ok := got["written"].([]any)

	if !ok || len(written) != len(ai.Targets) {
		t.Errorf("written = %v, want %d-entry slice", got["written"], len(ai.Targets))
	}

	if _, err := os.Stat(filepath.Join(s.Target, "CLAUDE.md")); err != nil {
		t.Errorf("CLAUDE.md not written: %v", err)
	}
}

func TestHandleGenerateRulesReturns409OnConflict(t *testing.T) {
	s := scannedFixture(t)

	mustWriteFile(t, s.Target, "CLAUDE.md", "existing")

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/_brain/api/generate-rules", bytes.NewReader([]byte(`{"force":false}`)))
	req.Header.Set("Content-Length", "16")
	s.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("status = %d, want 409; body=%s", rr.Code, rr.Body)
	}
}

func TestHandleStressTestEnqueueRejectsRemote(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	body := bytes.NewReader([]byte(`{"url":"https://example.com/"}`))
	s.Routes().ServeHTTP(rr, httptest.NewRequest("POST", "/_brain/api/stress-test", body))

	if rr.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403", rr.Code)
	}
}

func TestHandleStressTestEnqueueRequiresURL(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()

	s.Routes().ServeHTTP(rr, httptest.NewRequest("POST", "/_brain/api/stress-test", bytes.NewReader([]byte(``))))

	if rr.Code != http.StatusBadRequest {
		t.Errorf("empty body status = %d, want 400", rr.Code)
	}

	s.Routes().ServeHTTP(rr2, httptest.NewRequest("POST", "/_brain/api/stress-test", bytes.NewReader([]byte(`{}`))))

	if rr2.Code != http.StatusBadRequest {
		t.Errorf("no-url status = %d, want 400", rr2.Code)
	}
}

func TestHandleStressTestEnqueueLocalhostRuns(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Cleanup(upstream.Close)

	s := scannedFixture(t)
	body := bytes.NewReader([]byte(`{"url":"` + upstream.URL + `","concurrency":2,"requests":3,"timeoutMs":1000}`))
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("POST", "/_brain/api/stress-test", body))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rr.Code, rr.Body)
	}

	var got map[string]any

	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if v, _ := got["total"].(float64); v != 3 {
		t.Errorf("total = %v, want 3", got["total"])
	}
}

func TestHandleStressTestPollReservedForFuture(t *testing.T) {
	s := scannedFixture(t)
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, httptest.NewRequest("GET", "/_brain/api/stress-test/abc123", nil))

	if rr.Code != http.StatusNotImplemented {
		t.Errorf("status = %d, want 501", rr.Code)
	}
}

func TestWriteErrShape(t *testing.T) {
	rr := httptest.NewRecorder()
	writeErr(rr, http.StatusTeapot, errSentinel{})

	if rr.Code != http.StatusTeapot {
		t.Errorf("code = %d", rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q", ct)
	}

	var got map[string]string

	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got["error"] != "sentinel" {
		t.Errorf("body = %v", got)
	}
}

func (errSentinel) Error() string { return "sentinel" }
