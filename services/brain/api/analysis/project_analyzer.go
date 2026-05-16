package analysis

import (
	"context"
	"fmt"

	"github.com/bedrock/packages/pipeline"
	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/parser"
)

// ProjectAnalyzer is the top-level orchestrator. It loads the target project,
// runs the analyzer pipeline, and returns the assembled graph.
//
// The set of analyzers grows with each phase. Tests stub the slice directly.
type ProjectAnalyzer struct {
	Analyzers []Analyzer
}

// NewDefaultProjectAnalyzer returns the analyzer pipeline as wired for the
// current phase. Subsequent phases prepend their analyzer here.
func NewDefaultProjectAnalyzer() *ProjectAnalyzer {
	return &ProjectAnalyzer{
		Analyzers: []Analyzer{
			RouteAnalyzer{},
			ControllerAnalyzer{},
			MiddlewareAnalyzer{},
			ModelAnalyzer{},
			EventAnalyzer{},
			JobAnalyzer{},
			ConsoleAnalyzer{},
			ChannelAnalyzer{},
			FacadeAnalyzer{},
			ContainerBindingAnalyzer{},
			ValidationAnalyzer{},
			InertiaAnalyzer{},
			FlowAnalyzer{},
			&SecurityAnalyzer{},
		},
	}
}

// AnalyzeTarget runs every analyzer against `target` and returns the graph.
func (p *ProjectAnalyzer) AnalyzeTarget(target string) (*graph.Graph, error) {
	proj, err := parser.Load(target)
	if err != nil {
		return nil, err
	}
	return p.Analyze(proj)
}

// Analyze runs the pipeline against an already-loaded project. Each
// analyzer is wrapped in a packages/pipeline.Pipe so the chain can be
// composed, decorated, or short-circuited by external callers (for
// instance, an MCP server might insert a permission check between two
// analyzers without forking this code).
func (p *ProjectAnalyzer) Analyze(proj *parser.Project) (*graph.Graph, error) {
	g := graph.NewGraph(projectLabel(proj))
	pipes := make([]any, 0, len(p.Analyzers))
	for _, a := range p.Analyzers {
		pipes = append(pipes, analyzerAsPipe(a))
	}
	_, err := pipeline.New().
		Send(&Context{Project: proj, Graph: g}).
		Through(pipes...).
		Then(context.Background(), func(v any) (any, error) { return v, nil })
	if err != nil {
		return nil, err
	}
	g.Stamp()
	return g, nil
}

// analyzerAsPipe adapts an Analyzer into a pipeline.Pipe. The pipe runs
// the analyzer against the in-flight Context (the "passable"), then calls
// next to advance the chain.
func analyzerAsPipe(a Analyzer) pipeline.Pipe {
	return func(_ context.Context, passable any, next func(any) (any, error)) (any, error) {
		c, ok := passable.(*Context)
		if !ok {
			return nil, fmt.Errorf("analyzer %s: pipeline passable is %T, want *analysis.Context", a.Name(), passable)
		}
		if err := a.Analyze(c); err != nil {
			return nil, fmt.Errorf("analyzer %s: %w", a.Name(), err)
		}
		return next(c)
	}
}

func projectLabel(p *parser.Project) string {
	if p == nil || p.Root == "" {
		return ""
	}
	// last path segment
	for i := len(p.Root) - 1; i >= 0; i-- {
		if p.Root[i] == '/' {
			return p.Root[i+1:]
		}
	}
	return p.Root
}
