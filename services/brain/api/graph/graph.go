package graph

import (
	"time"
)

// Meta is the summary block emitted at the top of every graph JSON file.
type Meta struct {
	Project    string    `json:"project"`
	AnalyzedAt time.Time `json:"analyzedAt"`
	NodeCount  int       `json:"nodeCount"`
	EdgeCount  int       `json:"edgeCount"`
}

// Graph is the full output of an analyzer pipeline: a meta block plus nodes
// and edges, with an in-memory index for fast lookups during construction.
//
// The JSON shape matches laravel-brain's Graph.php so the Vue SPA can render
// either output unchanged.
type Graph struct {
	Meta  Meta    `json:"meta"`
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`

	nodeIndex        map[string]*Node
	directedEdgeKeys map[string]struct{}
}

// NewGraph returns an empty graph stamped with the given project label.
func NewGraph(project string) *Graph {
	return &Graph{
		Meta:             Meta{Project: project, AnalyzedAt: time.Now().UTC()},
		Nodes:            []*Node{},
		Edges:            []*Edge{},
		nodeIndex:        map[string]*Node{},
		directedEdgeKeys: map[string]struct{}{},
	}
}

// AddNode inserts a node if one with the same ID does not already exist.
// Returns the node already in the graph (new or existing) so callers can chain.
func (g *Graph) AddNode(n *Node) *Node {
	if existing, ok := g.nodeIndex[n.ID]; ok {
		return existing
	}
	g.Nodes = append(g.Nodes, n)
	g.nodeIndex[n.ID] = n
	g.Meta.NodeCount = len(g.Nodes)
	return n
}

// Node returns a node by ID or nil.
func (g *Graph) Node(id string) *Node {
	return g.nodeIndex[id]
}

// AddEdge appends a directed edge. Duplicate (source,target,type,label) tuples
// are coalesced — laravel-brain does the same to keep the graph readable.
func (g *Graph) AddEdge(e *Edge) *Edge {
	key := e.Source + "->" + e.Target + "|" + string(e.Type) + "|" + e.Label
	if _, dup := g.directedEdgeKeys[key]; dup {
		return e
	}
	if e.ID == "" {
		e.ID = "e" + key
	}
	g.directedEdgeKeys[key] = struct{}{}
	g.Edges = append(g.Edges, e)
	g.Meta.EdgeCount = len(g.Edges)
	return e
}

// Stamp refreshes the analyzed-at timestamp and counts before serialisation.
func (g *Graph) Stamp() {
	g.Meta.AnalyzedAt = time.Now().UTC()
	g.Meta.NodeCount = len(g.Nodes)
	g.Meta.EdgeCount = len(g.Edges)
}
