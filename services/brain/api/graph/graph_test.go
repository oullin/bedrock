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

func TestAddEdgeStructKeyAvoidsConcatCollision(t *testing.T) {
	g := NewGraph("demo")
	g.AddNode(NewNode("a", NodeTypeRoute, ""))
	g.AddNode(NewNode("b->c", NodeTypeController, ""))
	g.AddNode(NewNode("a->b", NodeTypeRoute, ""))
	g.AddNode(NewNode("c", NodeTypeController, ""))

	// Two distinct edges whose old concatenated key collided as
	// "a->b->c|handles_by|handles"; the second was silently dropped.
	g.AddEdge(&Edge{Source: "a", Target: "b->c", Type: EdgeTypeHandlesBy, Label: "handles"})
	g.AddEdge(&Edge{Source: "a->b", Target: "c", Type: EdgeTypeHandlesBy, Label: "handles"})

	if g.Meta.EdgeCount != 2 {
		t.Fatalf("edge count = %d, want 2 (concat-style collision)", g.Meta.EdgeCount)
	}

	if g.Edges[0].ID == g.Edges[1].ID {
		t.Fatalf("edge IDs collided: %q == %q", g.Edges[0].ID, g.Edges[1].ID)
	}
}

func TestAddEdgeAssignsSequentialIDs(t *testing.T) {
	g := NewGraph("demo")
	g.AddNode(NewNode("a", NodeTypeRoute, ""))
	g.AddNode(NewNode("b", NodeTypeController, ""))
	g.AddNode(NewNode("c", NodeTypeController, ""))

	g.AddEdge(&Edge{Source: "a", Target: "b", Type: EdgeTypeHandlesBy, Label: "x"})
	g.AddEdge(&Edge{Source: "a", Target: "c", Type: EdgeTypeHandlesBy, Label: "y"})

	if g.Edges[0].ID != "e1" || g.Edges[1].ID != "e2" {
		t.Fatalf("edge IDs = %q,%q want e1,e2", g.Edges[0].ID, g.Edges[1].ID)
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
