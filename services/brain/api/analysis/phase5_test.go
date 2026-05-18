package analysis

import (
	"testing"

	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/parser"
)

func TestValidationAndInertiaAnalyzers(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "handlers.go", `package fixture

type w struct{}
type r struct{}

type factory struct{}
func (factory) Validate(input map[string]any, rules map[string]any, _ any, _ any) (map[string]any, error) {
	return input, nil
}
func NewFactory() factory { return factory{} }

type container struct{}
func (container) Render(_ *w, _ *r, page string, props map[string]any) {}

func login() {
	NewFactory().Validate(map[string]any{}, map[string]any{
		"email": "required|email",
		"name":  "required|string",
	}, nil, nil)
}

func dashboard(c container, ww *w, rr *r) {
	c.Render(ww, rr, "Dashboard/Index", map[string]any{
		"user":   "u",
		"counts": 42,
	})
}
`)
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatal(err)
	}

	g := graph.NewGraph("fixture")
	ctx := &Context{Project: proj, Graph: g}

	for _, a := range []Analyzer{ValidationAnalyzer{}, InertiaAnalyzer{}} {
		if err := a.Analyze(ctx); err != nil {
			t.Fatalf("%s: %v", a.Name(), err)
		}
	}

	hasValidation := false

	for _, n := range g.Nodes {
		if n.Type == graph.NodeTypeValidationRequest {
			hasValidation = true
			fields, ok := n.Data["fields"].([]string)

			if !ok || len(fields) != 2 {
				t.Errorf("validation fields = %v, want [email name]", n.Data["fields"])
			}
		}
	}

	if !hasValidation {
		t.Errorf("expected a validation_request node; got: %v", nodeIDs(g))
	}

	if g.Node("inertia_page:Dashboard/Index") == nil {
		t.Errorf("expected inertia_page; got: %v", nodeIDs(g))
	}

	if g.Node("inertia_layout:Dashboard") == nil {
		t.Errorf("expected inertia_layout:Dashboard")
	}

	if g.Node("inertia_prop:Dashboard/Index#user") == nil {
		t.Errorf("expected inertia_prop user")
	}
}
