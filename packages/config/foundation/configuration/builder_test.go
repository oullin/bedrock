package configuration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestBuilderMergesBaseOverlayAndEnv(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "auth.yml"), "session_lifetime: 24h\ncookies:\n  secure: true\n")
	writeFile(t, filepath.Join(dir, "fortify.yml"), "views: true\nlimiters:\n  login: login\nfeatures:\n  - registration\n")
	writeFile(t, filepath.Join(dir, "local", "auth.yml"), "session_lifetime: 12h\ncookies:\n  secure: false\n")

	repo, err := NewBuilder(dir).
		WithAppEnv("local").
		WithEnv(map[string]string{
			"AUTH_SESSION_LIFETIME": "48h",
			"FORTIFY_VIEWS":         "false",
		}).
		Build(context.Background())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	duration, err := repo.Duration("auth.session_lifetime")
	if err != nil {
		t.Fatalf("Duration: %v", err)
	}
	if duration.String() != "48h0m0s" {
		t.Fatalf("unexpected duration: %v", duration)
	}

	secure, err := repo.Bool("auth.cookies.secure")
	if err != nil {
		t.Fatalf("Bool: %v", err)
	}
	if secure {
		t.Fatal("expected overlay bool to apply")
	}

	views, err := repo.Bool("fortify.views")
	if err != nil {
		t.Fatalf("Bool: %v", err)
	}
	if views {
		t.Fatal("expected env override to disable views")
	}
}

func TestBuilderInvalidYAML(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "auth.yml"), "session_lifetime: [\n")

	_, err := NewBuilder(dir).Build(context.Background())
	if err == nil {
		t.Fatal("expected invalid yaml error")
	}
}

func writeFile(t *testing.T, path string, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}
