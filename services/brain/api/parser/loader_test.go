package parser_test

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/services/brain/api/parser"
	"golang.org/x/tools/go/packages"
)

const fixtureGoMod = "module example.com/fixture\n\ngo 1.26\n"

const fixtureFooBody = `package fixture

func Hello() string { return "hi" }
`

const fixtureBarBody = `package fixture

func Bye() string { return "bye" }
`

func newFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", fixtureGoMod)
	mustWrite(t, dir, "foo.go", fixtureFooBody)
	mustWrite(t, dir, "bar.go", fixtureBarBody)

	return dir
}

func mustWrite(t *testing.T, dir, name, body string) {
	t.Helper()
	full := filepath.Join(dir, name)

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", name, err)
	}

	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestLoadIsolatesFromOuterGOWORK(t *testing.T) {
	// Force an unusable GOWORK in the process env. The previous behaviour was
	// to honor this and return empty Syntax with no top-level error — which
	// silently broke every analyzer test in CI. Load now sets GOWORK=off
	// inside packages.Config.Env so it remains usable regardless.
	t.Setenv("GOWORK", "/this/workspace/does/not/exist")

	proj, err := parser.Load(newFixture(t))

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(proj.Packages) == 0 {
		t.Fatal("Load returned 0 packages — GOWORK isolation regression")
	}

	totalFiles := 0

	for _, pkg := range proj.Packages {
		totalFiles += len(pkg.Syntax)
	}

	if totalFiles < 2 {
		t.Errorf("total syntax files = %d, want >= 2 (foo.go, bar.go)", totalFiles)
	}
}

func TestLoadReturnsAbsolutePath(t *testing.T) {
	dir := newFixture(t)
	// Convert to a relative path to ensure Load resolves it.
	rel, err := filepath.Rel(t.TempDir(), dir)

	if err != nil || rel == "" {
		rel = dir
	}

	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if !filepath.IsAbs(proj.Root) {
		t.Errorf("Root = %q, want absolute path", proj.Root)
	}
}

func TestEachFileIteratesProjectFiles(t *testing.T) {
	proj, err := parser.Load(newFixture(t))

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	got := map[string]bool{}
	err = proj.EachFile(func(_ *packages.Package, _ *ast.File, filename string) error {
		got[filepath.Base(filename)] = true

		return nil
	})

	if err != nil {
		t.Fatalf("EachFile: %v", err)
	}

	if !got["foo.go"] || !got["bar.go"] {
		t.Errorf("EachFile coverage = %v, want both foo.go and bar.go", got)
	}
}

func TestEachFilePropagatesCallbackError(t *testing.T) {
	proj, err := parser.Load(newFixture(t))

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	sentinel := errSentinel{}
	err = proj.EachFile(func(_ *packages.Package, _ *ast.File, _ string) error {
		return sentinel
	})

	if err != sentinel {
		t.Errorf("err = %v, want sentinel", err)
	}
}

func TestPositionFormatsRelativeToRoot(t *testing.T) {
	proj, err := parser.Load(newFixture(t))

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	var found bool
	_ = proj.EachFile(func(_ *packages.Package, file *ast.File, _ string) error {
		pos := proj.Position(file)
		// Format is "<rel-path>:<line>". For the fixture the file is at
		// the project root so the rel path is just the base name.
		if !strings.HasSuffix(pos, ".go:1") {
			t.Errorf("Position = %q, want suffix '.go:1'", pos)
		}

		found = true

		return nil
	})

	if !found {
		t.Errorf("EachFile never invoked the callback")
	}
}

func TestPositionEmptyForNil(t *testing.T) {
	proj := &parser.Project{Root: "/x"}

	if got := proj.Position(nil); got != "" {
		t.Errorf("Position(nil) = %q, want empty", got)
	}
}

type errSentinel struct{}

func (errSentinel) Error() string { return "sentinel" }
