package analysis

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/parser"
)

func TestRouteAnalyzer_DetectsVerbsAndPaths(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "routes.go", `package fixture

type router struct{}

func (r *router) Get(uri string, fn func()) {}
func (r *router) Post(uri string, fn func()) {}

func wire(r *router) {
	r.Get("/", func() {})
	r.Get("/health", func() {})
	r.Post("/login", func() {})
}
`)
	proj, err := parser.Load(dir)
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	g := graph.NewGraph("fixture")
	if err := (RouteAnalyzer{}).Analyze(&Context{Project: proj, Graph: g}); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	gotByID := map[string]*graph.Node{}
	for _, n := range g.Nodes {
		gotByID[n.ID] = n
	}
	wantIDs := []string{"route:get:/", "route:get:/health", "route:post:/login"}
	for _, id := range wantIDs {
		if _, ok := gotByID[id]; !ok {
			t.Errorf("missing node %s; got: %v", id, gotByID)
		}
	}
	if len(g.Nodes) != len(wantIDs) {
		t.Errorf("node count = %d, want %d", len(g.Nodes), len(wantIDs))
	}
}

func TestRouteAnalyzer_IgnoresNonURIFirstArg(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "noisy.go", `package fixture

type cache struct{}

func (c *cache) Get(key string) string { return "" }

func use(c *cache) {
	_ = c.Get("user:42") // not a route — no leading slash
}
`)
	proj, err := parser.Load(dir)
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	g := graph.NewGraph("fixture")
	if err := (RouteAnalyzer{}).Analyze(&Context{Project: proj, Graph: g}); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if len(g.Nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d: %v", len(g.Nodes), g.Nodes)
	}
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
