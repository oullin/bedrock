package analysis

import (
	"testing"

	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/parser"
)

func TestModelAnalyzer_DetectsTableNameReceiver(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "user.go", `package fixture

type User struct {
	ID   int
	Name string
}

func (User) TableName() string { return "users" }
`)
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	g := graph.NewGraph("fixture")

	if err := (ModelAnalyzer{}).Analyze(&Context{Project: proj, Graph: g}); err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	if !hasModelNamed(g, "User") {
		t.Errorf("expected model 'User'; nodes: %v", nodeIDs(g))
	}
}

func TestModelAnalyzer_DetectsStructTagged(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "order.go", "package fixture\n\ntype Order struct {\n\tID    int    `db:\"id\"`\n\tTotal int    `gorm:\"column:total\"`\n\tNote  string `sql:\"note\"`\n}\n")
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	g := graph.NewGraph("fixture")
	_ = (ModelAnalyzer{}).Analyze(&Context{Project: proj, Graph: g})

	if !hasModelNamed(g, "Order") {
		t.Errorf("expected model 'Order' from tagged fields; nodes: %v", nodeIDs(g))
	}
}

func TestModelAnalyzer_DetectsByPackagePathHint(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "models/customer/customer.go", `package customer

type Customer struct {
	Email string
}
`)
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	g := graph.NewGraph("fixture")
	_ = (ModelAnalyzer{}).Analyze(&Context{Project: proj, Graph: g})

	if !hasModelNamed(g, "Customer") {
		t.Errorf("expected model 'Customer' from /models/ path hint; nodes: %v", nodeIDs(g))
	}
}

func TestModelAnalyzer_IgnoresPlainStructs(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "util.go", `package fixture

type Result struct {
	Value int
}
`)
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	g := graph.NewGraph("fixture")
	_ = (ModelAnalyzer{}).Analyze(&Context{Project: proj, Graph: g})

	for _, n := range g.Nodes {
		if n.Type == graph.NodeTypeModel {
			t.Errorf("unexpected model node %q on untagged plain struct", n.ID)
		}
	}
}

func hasModelNamed(g *graph.Graph, name string) bool {
	for _, n := range g.Nodes {
		if n.Type == graph.NodeTypeModel && n.Label == name {
			return true
		}
	}

	return false
}
