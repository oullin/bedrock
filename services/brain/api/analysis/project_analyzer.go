package analysis

import (
	"fmt"

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

// Analyze runs the pipeline against an already-loaded project. Useful in
// tests that construct synthetic projects without a filesystem.
func (p *ProjectAnalyzer) Analyze(proj *parser.Project) (*graph.Graph, error) {
	g := graph.NewGraph(projectLabel(proj))
	ctx := &Context{Project: proj, Graph: g}
	for _, a := range p.Analyzers {
		if err := a.Analyze(ctx); err != nil {
			return nil, fmt.Errorf("analyzer %s: %w", a.Name(), err)
		}
	}
	g.Stamp()
	return g, nil
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
