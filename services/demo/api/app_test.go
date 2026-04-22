package api_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/bedrock/services/demo/api"
)

func testOptions(t *testing.T) api.Options {
	t.Helper()

	return api.Options{
		BasePath:    t.TempDir(),
		StoragePath: filepath.Join(t.TempDir(), "storage"),
		DatabaseURL: "sqlite:///:memory:",
		AppKey:      "base64:development-test-key",
	}
}

func TestNewApplication_ConfiguresSkeletonDatabase(t *testing.T) {
	t.Parallel()

	opts := testOptions(t)
	opts.DatabaseURL = "sqlite:///" + filepath.ToSlash(filepath.Join(opts.StoragePath, "database.sqlite"))

	application := api.NewApplication(opts)

	raw, err := application.Make("demo.sql")

	if err != nil {
		t.Fatalf("Make(demo.sql) error = %v", err)
	}

	db, ok := raw.(*sql.DB)

	if !ok {
		t.Fatalf("demo.sql binding = %T, want *sql.DB", raw)
	}

	for _, table := range []string{"users", "password_reset_tokens", "sessions", "cache", "cache_locks", "jobs", "job_batches", "failed_jobs"} {
		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&count)

		if err != nil {
			t.Fatalf("query table %s: %v", table, err)
		}

		if count != 1 {
			t.Fatalf("expected table %s to exist", table)
		}
	}
}

func TestNewApplication_SeedsDeterministicUser(t *testing.T) {
	t.Parallel()

	application := api.NewApplication(testOptions(t))

	raw, err := application.Make("demo.sql")

	if err != nil {
		t.Fatalf("Make(demo.sql) error = %v", err)
	}

	db := raw.(*sql.DB)

	var name, email string

	err = db.QueryRow(`SELECT name, email FROM users WHERE email = ?`, "taylor@example.com").Scan(&name, &email)

	if err != nil {
		t.Fatalf("seeded user query error = %v", err)
	}

	if name != "Taylor Otwell" || email != "taylor@example.com" {
		t.Fatalf("seeded user = (%q, %q)", name, email)
	}
}

func TestNewHandler_SkeletonRoutes(t *testing.T) {
	t.Parallel()

	handler, err := api.NewHandler(testOptions(t))

	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	t.Run("welcome", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET / status = %d", rec.Code)
		}

		if body := rec.Body.String(); body == "" || body == "404 page not found\n" {
			t.Fatalf("GET / body = %q", body)
		}
	})

	t.Run("health", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/up", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /up status = %d", rec.Code)
		}

		if rec.Body.String() != "OK\n" {
			t.Fatalf("GET /up body = %q", rec.Body.String())
		}
	})

	t.Run("lottery", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/lottery", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /lottery status = %d", rec.Code)
		}

		if rec.Body.String() != "winner\n" {
			t.Fatalf("GET /lottery body = %q", rec.Body.String())
		}
	})

	t.Run("missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/missing", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /missing status = %d", rec.Code)
		}
	})
}
