package analysis

import (
	"errors"
	"testing"

	"github.com/bedrock/services/brain/api/parser"
)

func TestNewDefaultProjectAnalyzerWiresAllAnalyzers(t *testing.T) {
	p := NewDefaultProjectAnalyzer()

	if p == nil || len(p.Analyzers) == 0 {
		t.Fatal("expected non-empty analyzer pipeline")
	}

	want := map[string]bool{
		"route":              false,
		"controller":         false,
		"middleware":         false,
		"model":              false,
		"event":              false,
		"job":                false,
		"console":            false,
		"channel":            false,
		"facade":             false,
		"container_binding":  false,
		"validation":         false,
		"inertia":            false,
		"flow":               false,
		"security":           false,
	}

	for _, a := range p.Analyzers {
		if _, ok := want[a.Name()]; ok {
			want[a.Name()] = true
		}
	}

	for name, found := range want {
		if !found {
			t.Errorf("missing analyzer in default pipeline: %q", name)
		}
	}
}

func TestProjectAnalyzer_AnalyzeTargetReturnsStampedGraph(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "routes.go", `package fixture

type router struct{}

func (r *router) Get(uri string, fn func()) {}

func wire(r *router) { r.Get("/", func() {}) }
`)

	p := NewDefaultProjectAnalyzer()
	g, err := p.AnalyzeTarget(dir)

	if err != nil {
		t.Fatalf("AnalyzeTarget: %v", err)
	}

	if g.Meta.Project != "" && g.Meta.NodeCount == 0 {
		t.Error("expected NodeCount > 0 after stamp")
	}

	if g.Meta.AnalyzedAt.IsZero() {
		t.Error("AnalyzedAt was not stamped")
	}
}

func TestProjectAnalyzer_AnalyzeTargetSurfacesLoadError(t *testing.T) {
	p := NewDefaultProjectAnalyzer()
	_, err := p.AnalyzeTarget("/this/path/does/not/exist/at/all")

	if err == nil {
		t.Fatal("expected error for nonexistent target")
	}
}

// failingAnalyzer always errors — used to assert pipeline error propagation.
type failingAnalyzer struct{ called bool }

func (failingAnalyzer) Name() string { return "failing" }

func (f *failingAnalyzer) Analyze(_ *Context) error {
	f.called = true

	return errSentinel
}

var errSentinel = errors.New("boom")

func TestProjectAnalyzer_PropagatesAnalyzerError(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "x.go", "package fixture\n")
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	fa := &failingAnalyzer{}
	p := &ProjectAnalyzer{Analyzers: []Analyzer{fa}}
	_, err = p.Analyze(proj)

	if err == nil {
		t.Fatal("expected pipeline error")
	}

	if !fa.called {
		t.Error("failing analyzer was never invoked")
	}
}

func TestProjectAnalyzer_ProjectLabelFromRootPath(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "x.go", "package fixture\n")
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	g, err := (&ProjectAnalyzer{}).Analyze(proj)

	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	if g == nil || g.Meta.Project == "" {
		t.Errorf("expected project label from Root path; got %+v", g.Meta)
	}
}
