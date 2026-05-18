package graph

import (
	"fmt"
	"time"
)

// Meta is the summary block emitted at the top of every graph JSON file.
type Meta struct {
	Project    string    `json:"project"`
	AnalyzedAt time.Time `json:"analyzedAt"`
	NodeCount  int       `json:"nodeCount"`
	EdgeCount  int       `json:"edgeCount"`
}

// edgeKey identifies a directed edge for dedupe purposes. Using a struct
// (rather than a concatenated string) keeps the four fields separate so
// values that happen to contain the old separators ("->" or "|") cannot
// silently collide with a different tuple.
type edgeKey struct {
	Source string
	Target string
	Type   EdgeType
	Label  string
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
	directedEdgeKeys map[edgeKey]struct{}
}

// NewGraph returns an empty graph stamped with the given project label.
func NewGraph(project string) *Graph {
	return &Graph{
		Meta:             Meta{Project: project, AnalyzedAt: time.Now().UTC()},
		Nodes:            []*Node{},
		Edges:            []*Edge{},
		nodeIndex:        map[string]*Node{},
		directedEdgeKeys: map[edgeKey]struct{}{},
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
	k := edgeKey{Source: e.Source, Target: e.Target, Type: e.Type, Label: e.Label}

	if _, dup := g.directedEdgeKeys[k]; dup {
		return e
	}

	if e.ID == "" {
		e.ID = fmt.Sprintf("e%d", len(g.Edges)+1)
	}

	g.directedEdgeKeys[k] = struct{}{}
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
