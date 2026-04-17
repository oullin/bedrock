package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/config"
)

func TestDefaultHead(t *testing.T) {
	head := config.DefaultHead()

	if head.Lang != "en" {
		t.Errorf("Lang = %q, want %q", head.Lang, "en")
	}

	if strings.TrimSpace(head.Title) != "" {
		t.Errorf("Title = %q, want empty", head.Title)
	}

	// Check robots meta tag.
	found := false

	for _, tag := range head.Meta {
		if tag.Name == "robots" && tag.Content == "index, follow" {
			found = true
		}
	}

	if !found {
		t.Error("expected robots meta tag with 'index, follow'")
	}

	// Check og:type meta tag.
	found = false

	for _, tag := range head.Meta {
		if tag.Property == "og:type" && tag.Content == "website" {
			found = true
		}
	}

	if !found {
		t.Error("expected og:type meta tag with 'website'")
	}
}

func TestDefaultHead_EnvOverride(t *testing.T) {
	t.Setenv("INERTIA_SEO_TITLE", "Env Title")

	head := config.DefaultHead()

	if head.Title != "Env Title" {
		t.Errorf("Title = %q, want %q", head.Title, "Env Title")
	}
}

func TestLoadHead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "seo.yml")

	content := `
title: "My App"
meta:
  - name: "description"
    content: "My app description"
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	head, err := config.LoadHead(path)

	if err != nil {
		t.Fatal(err)
	}

	if head.Title != "My App" {
		t.Errorf("Title = %q, want %q", head.Title, "My App")
	}

	// Default lang should apply.
	if head.Lang != "en" {
		t.Errorf("Lang = %q, want %q (default)", head.Lang, "en")
	}

	found := false

	for _, tag := range head.Meta {
		if tag.Name == "robots" && tag.Content == "index, follow" {
			found = true
		}
	}

	if !found {
		t.Error("expected default robots tag to be preserved")
	}
}

func TestLoadHead_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := config.LoadHead("/nonexistent/seo.yml")

	if err == nil {
		t.Error("expected error for missing file")
	}
}
