package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDistPath(t *testing.T) {
	t.Run("finds first candidate", func(t *testing.T) {
		tmp := t.TempDir()
		origDir, _ := os.Getwd()

		os.Chdir(tmp)

		t.Cleanup(func() { os.Chdir(origDir) })

		os.MkdirAll(filepath.Join(tmp, "storage", "dist", "app"), 0o755)

		got, err := resolveDistPath()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != filepath.Clean("storage/dist/app") {
			t.Errorf("got %q, want %q", got, "storage/dist/app")
		}
	})

	t.Run("no candidates returns error", func(t *testing.T) {
		tmp := t.TempDir()
		origDir, _ := os.Getwd()

		os.Chdir(tmp)

		t.Cleanup(func() { os.Chdir(origDir) })

		_, err := resolveDistPath()

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestResolveResourcePath(t *testing.T) {
	t.Run("finds resource file", func(t *testing.T) {
		tmp := t.TempDir()
		origDir, _ := os.Getwd()

		os.Chdir(tmp)

		t.Cleanup(func() { os.Chdir(origDir) })

		dir := filepath.Join(tmp, "resources")

		os.MkdirAll(dir, 0o755)

		os.WriteFile(filepath.Join(dir, "seo.yml"), []byte("title: test"), 0o644)

		got, err := resolveResourcePath("seo.yml")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != filepath.Clean("resources/seo.yml") {
			t.Errorf("got %q, want %q", got, "resources/seo.yml")
		}
	})

	t.Run("no match returns error", func(t *testing.T) {
		tmp := t.TempDir()
		origDir, _ := os.Getwd()

		os.Chdir(tmp)

		t.Cleanup(func() { os.Chdir(origDir) })

		_, err := resolveResourcePath("nonexistent.yml")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("finds in cmd/resources", func(t *testing.T) {
		tmp := t.TempDir()
		origDir, _ := os.Getwd()

		os.Chdir(tmp)

		t.Cleanup(func() { os.Chdir(origDir) })

		dir := filepath.Join(tmp, "cmd", "resources")

		os.MkdirAll(dir, 0o755)

		os.WriteFile(filepath.Join(dir, "csrf.yml"), []byte("secure: true"), 0o644)

		got, err := resolveResourcePath("csrf.yml")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != filepath.Clean("cmd/resources/csrf.yml") {
			t.Errorf("got %q, want %q", got, "cmd/resources/csrf.yml")
		}
	})
}

func TestMustResolveResourcePath(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()

	os.Chdir(tmp)

	t.Cleanup(func() { os.Chdir(origDir) })

	dir := filepath.Join(tmp, "resources")

	os.MkdirAll(dir, 0o755)

	os.WriteFile(filepath.Join(dir, "test.yml"), []byte("ok"), 0o644)

	got := mustResolveResourcePath("test.yml")

	if got != filepath.Clean("resources/test.yml") {
		t.Errorf("got %q, want %q", got, "resources/test.yml")
	}
}
