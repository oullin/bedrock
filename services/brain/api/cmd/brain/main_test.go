package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/services/brain/api/graph"
)

const fixtureGoMod = "module example.com/fixture\n\ngo 1.26\n"

const fixtureRoutes = `package fixture

type router struct{}

func (r *router) Get(uri string, fn func()) {}

func wire(r *router) { r.Get("/", func() {}) }
`

func newFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", fixtureGoMod)
	mustWrite(t, dir, "routes.go", fixtureRoutes)

	return dir
}

func mustWrite(t *testing.T, dir, name, body string) {
	t.Helper()
	full := filepath.Join(dir, name)

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", name, err)
	}

	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestRunNoArgsPrintsUsage(t *testing.T) {
	var out, errBuf bytes.Buffer

	if err := run(nil, &out, &errBuf); err != nil {
		t.Fatalf("run: %v", err)
	}

	if !strings.Contains(out.String(), "brain - bedrock service analyzer") {
		t.Errorf("usage not printed:\n%s", out.String())
	}
}

func TestRunVersionSubcommand(t *testing.T) {
	for _, arg := range []string{"version", "-v", "--version"} {
		var out, errBuf bytes.Buffer

		if err := run([]string{arg}, &out, &errBuf); err != nil {
			t.Fatalf("run(%q): %v", arg, err)
		}

		if !strings.Contains(out.String(), "brain 0.0.0") {
			t.Errorf("run(%q) version output: %q", arg, out.String())
		}
	}
}

func TestRunHelpSubcommand(t *testing.T) {
	for _, arg := range []string{"help", "-h", "--help"} {
		var out, errBuf bytes.Buffer

		if err := run([]string{arg}, &out, &errBuf); err != nil {
			t.Fatalf("run(%q): %v", arg, err)
		}

		if !strings.Contains(out.String(), "brain - bedrock service analyzer") {
			t.Errorf("run(%q) usage missing", arg)
		}
	}
}

func TestRunUnknownSubcommand(t *testing.T) {
	var out, errBuf bytes.Buffer
	err := run([]string{"bogus"}, &out, &errBuf)

	if err == nil {
		t.Fatal("run returned nil error for unknown subcommand")
	}

	if !strings.Contains(err.Error(), "unknown subcommand") {
		t.Errorf("err = %v, want 'unknown subcommand'", err)
	}
}

func TestRunScanWritesManifestToExplicitOutput(t *testing.T) {
	fix := newFixture(t)
	outDir := filepath.Join(t.TempDir(), "out")

	var out, errBuf bytes.Buffer
	err := run([]string{"scan", "-target=" + fix, "-output=" + outDir, "-quiet"}, &out, &errBuf)

	if err != nil {
		t.Fatalf("run scan: %v\nstderr=%s", err, errBuf.String())
	}

	for _, name := range []string{".graph-all.json", ".graph-manifest.json"} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}

	body, err := os.ReadFile(filepath.Join(outDir, ".graph-manifest.json"))

	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}

	var m graph.Manifest

	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}

	if m.TotalNodes == 0 {
		t.Errorf("manifest TotalNodes = 0; expected non-zero")
	}
}

func TestRunScanDefaultsOutputToTargetStorage(t *testing.T) {
	fix := newFixture(t)

	var out, errBuf bytes.Buffer
	err := run([]string{"scan", "-target=" + fix, "-quiet"}, &out, &errBuf)

	if err != nil {
		t.Fatalf("run scan: %v\nstderr=%s", err, errBuf.String())
	}

	defaultOut := filepath.Join(fix, "storage", "brain", ".graph-manifest.json")

	if _, err := os.Stat(defaultOut); err != nil {
		t.Errorf("expected default output at %s: %v", defaultOut, err)
	}
}

func TestRunScanPrintsSummaryWhenNotQuiet(t *testing.T) {
	fix := newFixture(t)
	outDir := filepath.Join(t.TempDir(), "out")

	var out, errBuf bytes.Buffer
	err := run([]string{"scan", "-target=" + fix, "-output=" + outDir}, &out, &errBuf)

	if err != nil {
		t.Fatalf("run scan: %v", err)
	}

	if !strings.Contains(out.String(), "brain: wrote") {
		t.Errorf("missing 'brain: wrote' summary in stdout: %q", out.String())
	}
}

func TestCountRoutes(t *testing.T) {
	g := graph.NewGraph("svc")

	if got := countRoutes(g); got != 0 {
		t.Errorf("empty = %d, want 0", got)
	}

	g.AddNode(graph.NewNode("r:1", graph.NodeTypeRoute, ""))
	g.AddNode(graph.NewNode("r:2", graph.NodeTypeRoute, ""))
	g.AddNode(graph.NewNode("m:1", graph.NodeTypeMiddleware, ""))

	if got := countRoutes(g); got != 2 {
		t.Errorf("got = %d, want 2", got)
	}
}

func TestWriteJSONCreatesNestedDirs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "c.json")

	if err := writeJSON(path, map[string]int{"x": 1}); err != nil {
		t.Fatalf("writeJSON: %v", err)
	}

	body, err := os.ReadFile(path)

	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var got map[string]int

	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got["x"] != 1 {
		t.Errorf("payload = %v, want x=1", got)
	}
}

func TestPrintUsageMentionsSubcommands(t *testing.T) {
	var buf bytes.Buffer

	printUsage(&buf)
	s := buf.String()

	for _, want := range []string{"scan", "serve", "version", "help"} {
		if !strings.Contains(s, want) {
			t.Errorf("usage missing %q:\n%s", want, s)
		}
	}
}
