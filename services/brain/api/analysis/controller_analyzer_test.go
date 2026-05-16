package analysis

import (
	"testing"

	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/parser"
)

func TestControllerAnalyzer_LinksNamedHandler(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "routes.go", `package fixture

type router struct{}

func (r *router) Get(uri string, fn func()) {}

func ShowHome() {}

func wire(r *router) {
	r.Get("/", ShowHome)
}
`)
	proj, err := parser.Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	g := graph.NewGraph("fixture")
	ctx := &Context{Project: proj, Graph: g}
	if err := (RouteAnalyzer{}).Analyze(ctx); err != nil {
		t.Fatal(err)
	}
	if err := (ControllerAnalyzer{}).Analyze(ctx); err != nil {
		t.Fatal(err)
	}
	if g.Node("action:showhome") == nil {
		t.Errorf("expected action node, got: %v", nodeIDs(g))
	}
	if !hasEdge(g, "route:get:/", "action:showhome", graph.EdgeTypeHandlesBy) {
		t.Errorf("expected route→action edge, got edges: %v", g.Edges)
	}
}

func nodeIDs(g *graph.Graph) []string {
	out := make([]string, 0, len(g.Nodes))
	for _, n := range g.Nodes {
		out = append(out, n.ID)
	}
	return out
}

func hasEdge(g *graph.Graph, src, dst string, t graph.EdgeType) bool {
	for _, e := range g.Edges {
		if e.Source == src && e.Target == dst && e.Type == t {
			return true
		}
	}
	return false
}
