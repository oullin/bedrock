package graph

import (
	"encoding/json"
	"testing"
)

func TestAddNodeDedupesByID(t *testing.T) {
	g := NewGraph("demo")
	a := g.AddNode(NewNode("n1", NodeTypeRoute, "GET /"))
	b := g.AddNode(NewNode("n1", NodeTypeRoute, "GET /"))
	if a != b {
		t.Fatal("expected duplicate AddNode to return same pointer")
	}
	if g.Meta.NodeCount != 1 {
		t.Fatalf("node count = %d, want 1", g.Meta.NodeCount)
	}
}

func TestAddEdgeDedupesByTuple(t *testing.T) {
	g := NewGraph("demo")
	g.AddNode(NewNode("a", NodeTypeRoute, ""))
	g.AddNode(NewNode("b", NodeTypeController, ""))
	g.AddEdge(&Edge{Source: "a", Target: "b", Type: EdgeTypeHandlesBy, Label: "handles"})
	g.AddEdge(&Edge{Source: "a", Target: "b", Type: EdgeTypeHandlesBy, Label: "handles"})
	if g.Meta.EdgeCount != 1 {
		t.Fatalf("edge count = %d, want 1", g.Meta.EdgeCount)
	}
}

func TestNodeJSONFieldOrder(t *testing.T) {
	n := NewNode("r:get:/", NodeTypeRoute, "GET /")
	n.Set("method", "GET").Set("uri", "/")
	out, err := json.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"id":"r:get:/","type":"route","label":"GET /","data":{"method":"GET","uri":"/"}}`
	if string(out) != want {
		t.Fatalf("json mismatch\n got: %s\nwant: %s", out, want)
	}
}
