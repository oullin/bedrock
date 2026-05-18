package ai

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateRulesWritesAllTargets(t *testing.T) {
	dir := t.TempDir()
	body := "# brain context\n"

	written, err := GenerateRules(dir, body, false)

	if err != nil {
		t.Fatalf("GenerateRules: %v", err)
	}

	if len(written) != len(Targets) {
		t.Fatalf("written count = %d, want %d", len(written), len(Targets))
	}

	for _, target := range Targets {
		full := filepath.Join(dir, target.Path)
		data, err := os.ReadFile(full)

		if err != nil {
			t.Errorf("missing target %s: %v", target.Path, err)

			continue
		}

		if string(data) != body {
			t.Errorf("%s body mismatch: got %q, want %q", target.Path, data, body)
		}
	}
}

func TestGenerateRulesCreatesNestedDirs(t *testing.T) {
	dir := t.TempDir()

	if _, err := GenerateRules(dir, "x", false); err != nil {
		t.Fatalf("GenerateRules: %v", err)
	}

	for _, nested := range []string{".github/copilot-instructions.md", ".junie/guidelines.md"} {
		if _, err := os.Stat(filepath.Join(dir, nested)); err != nil {
			t.Errorf("expected nested file %s: %v", nested, err)
		}
	}
}

func TestGenerateRulesReturnsErrConflictWhenAnyExists(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("old"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	written, err := GenerateRules(dir, "new", false)

	if !errors.Is(err, ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}

	if written != nil {
		t.Errorf("written = %v, want nil on conflict", written)
	}
	// Atomic: no other targets should have been created.
	for _, target := range Targets[1:] {
		if _, err := os.Stat(filepath.Join(dir, target.Path)); err == nil {
			t.Errorf("non-atomic: %s was created despite conflict", target.Path)
		}
	}
	// And the original CLAUDE.md must be untouched.
	body, _ := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))

	if string(body) != "old" {
		t.Errorf("CLAUDE.md overwritten: %q", body)
	}
}

func TestGenerateRulesForceOverwrites(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("old"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	written, err := GenerateRules(dir, "new", true)

	if err != nil {
		t.Fatalf("GenerateRules(force): %v", err)
	}

	if len(written) != len(Targets) {
		t.Errorf("written count = %d, want %d", len(written), len(Targets))
	}

	body, _ := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))

	if string(body) != "new" {
		t.Errorf("force did not overwrite: got %q", body)
	}
}
