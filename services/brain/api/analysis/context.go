// Package analysis hosts the static-analysis pipeline. Each Analyzer reads
// the target project's AST and writes nodes/edges into the shared graph.
package analysis

import (
	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/parser"
)

// Context is shared state passed to every analyzer in the pipeline.
type Context struct {
	Project *parser.Project
	Graph   *graph.Graph
}

// Analyzer is one stage of the scan pipeline. Phases are intentionally small
// so the upstream-brain analyzer split (Route, Controller, Model, ...) maps
// 1:1 to a Go type here.
type Analyzer interface {
	Name() string
	Analyze(ctx *Context) error
}
