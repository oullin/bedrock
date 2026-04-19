package telescope_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/telescope"
	"github.com/bedrock/packages/telescope/storage"
	"github.com/bedrock/packages/telescope/watchers"
)

func makeRequestTestScope(t *testing.T, opts map[string]any) (*telescope.Telescope, *storage.InMemoryRepository, *watchers.RequestWatcher) {
	t.Helper()

	repo := storage.NewInMemoryRepository()
	scope := telescope.New(telescope.WithRepository(repo))
	scope.StartRecording()
	w := watchers.NewRequestWatcher(scope, opts)

	return scope, repo, w
}

func doTestRequest(t *testing.T, w *watchers.RequestWatcher, method, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	var bodyReader *strings.Reader

	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	var req *http.Request

	var err error

	if bodyReader != nil {
		req, err = http.NewRequest(method, path, bodyReader)
	} else {
		req, err = http.NewRequest(method, path, nil)
	}

	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rr := httptest.NewRecorder()
	handler := w.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"ok":true}`)) //nolint:errcheck
	}))
	handler.ServeHTTP(rr, req)

	return rr
}

func TestRequestWatcherRecordsBasicRequest(t *testing.T) {
	t.Parallel()

	scope, repo, w := makeRequestTestScope(t, nil)

	doTestRequest(t, w, http.MethodGet, "/users", "", nil)

	if err := scope.Store(testContext()); err != nil {
		t.Fatal(err)
	}

	entries, err := repo.Get(telescope.EntryTypeRequest, telescope.DefaultQueryOptions())

	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 request entry, got %d", len(entries))
	}

	e := entries[0]

	if e.Content["method"] != "GET" {
		t.Fatalf("expected method=GET, got %v", e.Content["method"])
	}

	if e.Content["response_status"] != http.StatusOK {
		t.Fatalf("expected response_status=200, got %v", e.Content["response_status"])
	}
}

func TestRequestWatcherMasksPasswordParameter(t *testing.T) {
	t.Parallel()

	scope, repo, w := makeRequestTestScope(t, nil)

	doTestRequest(t, w, http.MethodPost, "/login",
		`{"email":"alice@example.com","password":"secret"}`,
		map[string]string{"Content-Type": "application/json"},
	)

	if err := scope.Store(testContext()); err != nil {
		t.Fatal(err)
	}

	entries, err := repo.Get(telescope.EntryTypeRequest, telescope.DefaultQueryOptions())

	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	payload, _ := entries[0].Content["payload"].(map[string]any)

	if payload["password"] != "********" {
		t.Fatalf("expected password to be masked, got %v", payload["password"])
	}

	if payload["email"] == "********" {
		t.Fatal("email should not be masked")
	}
}

func TestRequestWatcherMasksAuthorizationHeader(t *testing.T) {
	t.Parallel()

	scope, repo, w := makeRequestTestScope(t, nil)

	doTestRequest(t, w, http.MethodGet, "/api/profile", "", map[string]string{
		"Authorization": "Bearer super-secret-token",
	})

	if err := scope.Store(testContext()); err != nil {
		t.Fatal(err)
	}

	entries, err := repo.Get(telescope.EntryTypeRequest, telescope.DefaultQueryOptions())

	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	headers, _ := entries[0].Content["headers"].(map[string]string)

	if headers["Authorization"] != "********" {
		t.Fatalf("expected Authorization header to be masked, got %q", headers["Authorization"])
	}
}

func TestRequestWatcherIgnoresConfiguredHttpMethods(t *testing.T) {
	t.Parallel()

	scope, repo, w := makeRequestTestScope(t, map[string]any{
		"ignore_http_methods": []string{"OPTIONS"},
	})

	doTestRequest(t, w, http.MethodOptions, "/", "", nil)
	doTestRequest(t, w, http.MethodGet, "/users", "", nil)

	if err := scope.Store(testContext()); err != nil {
		t.Fatal(err)
	}

	entries, err := repo.Get(telescope.EntryTypeRequest, telescope.DefaultQueryOptions())

	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry (GET only), got %d", len(entries))
	}

	if entries[0].Content["method"] != "GET" {
		t.Fatalf("expected GET entry, got %v", entries[0].Content["method"])
	}
}
