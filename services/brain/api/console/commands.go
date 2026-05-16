// Package console exposes brain's CLI surface for embedding into a host
// service's packages/console Application. The contract is intentionally
// dependency-free: callers pass a `Registrar` that adapts to whatever
// console package they use, and brain hands back three `CommandSpec`s.
//
// A host service wires brain like so:
//
//	import (
//		"github.com/bedrock/packages/console"
//		brainconsole "github.com/bedrock/services/brain/api/console"
//	)
//
//	app := console.NewApplication("my-service")
//	for _, spec := range brainconsole.Commands() {
//		cmd, _ := console.NewCommand(spec.Signature, spec.RunWith(app))
//		app.Add(cmd)
//	}
//
// This keeps the bedrock/console <-> brain seam thin and lets brain ship
// as an importable library without owning the host's CLI shape.
package console

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bedrock/services/brain/api/ai"
	"github.com/bedrock/services/brain/api/analysis"
	"github.com/bedrock/services/brain/api/graph"
)

// CommandSpec is a wire-format description of one brain CLI subcommand.
// A host service's console package maps it onto its own Command type.
type CommandSpec struct {
	Name      string
	Signature string
	Summary   string
	Run       func(ctx context.Context, args Args, out io.Writer) error
}

// Args is the subset of console arguments brain needs. Hosts can build it
// from their own argument abstraction (cobra Flag set, packages/console
// InputInterface, etc).
type Args struct {
	Target  string
	Output  string
	Force   bool
	Format  string
}

// Commands returns the three brain subcommands. Each function takes a
// resolved Args and writes progress to `out`.
func Commands() []CommandSpec {
	return []CommandSpec{
		{
			Name:      "brain:scan",
			Signature: "brain:scan {--target=} {--output=}",
			Summary:   "Scan a Go service and write graph JSON to <target>/storage/brain.",
			Run:       runScan,
		},
		{
			Name:      "brain:export-context",
			Signature: "brain:export-context {--target=} {--format=markdown}",
			Summary:   "Render the deterministic AI context export for a service.",
			Run:       runExportContext,
		},
		{
			Name:      "brain:generate-rules",
			Signature: "brain:generate-rules {--target=} {--force}",
			Summary:   "Write editor rules files (CLAUDE.md, .cursorrules, ...).",
			Run:       runGenerateRules,
		},
	}
}

func runScan(ctx context.Context, a Args, out io.Writer) error {
	target, output := resolvePaths(a)
	g, err := scan(target)
	if err != nil {
		return err
	}
	if err := writeGraph(output, g); err != nil {
		return err
	}
	fmt.Fprintf(out, "brain: wrote %s (%d nodes, %d edges)\n", output, g.Meta.NodeCount, g.Meta.EdgeCount)
	_ = ctx
	return nil
}

func runExportContext(_ context.Context, a Args, out io.Writer) error {
	target, _ := resolvePaths(a)
	g, err := scan(target)
	if err != nil {
		return err
	}
	body := ai.RenderMarkdown(g, ai.ContextOptions{})
	_, err = out.Write([]byte(body))
	return err
}

func runGenerateRules(_ context.Context, a Args, out io.Writer) error {
	target, _ := resolvePaths(a)
	g, err := scan(target)
	if err != nil {
		return err
	}
	body := ai.RenderMarkdown(g, ai.ContextOptions{})
	written, err := ai.GenerateRules(target, body, a.Force)
	if err != nil {
		return err
	}
	for _, p := range written {
		fmt.Fprintf(out, "brain: wrote %s\n", p)
	}
	return nil
}

func scan(target string) (*graph.Graph, error) {
	return analysis.NewDefaultProjectAnalyzer().AnalyzeTarget(target)
}

func resolvePaths(a Args) (target string, output string) {
	target = a.Target
	if target == "" {
		target = "."
	}
	abs, err := filepath.Abs(target)
	if err == nil {
		target = abs
	}
	output = a.Output
	if output == "" {
		output = filepath.Join(target, "storage", "brain")
	}
	return
}

func writeGraph(dir string, g *graph.Graph) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	// Re-use graph package's JSON shape by encoding directly here to keep
	// this package dependency-free of cmd/brain.
	if err := writeJSON(filepath.Join(dir, ".graph-all.json"), g); err != nil {
		return err
	}
	return writeJSON(filepath.Join(dir, ".graph-manifest.json"), graph.Manifest{
		Project:     g.Meta.Project,
		AnalyzedAt:  g.Meta.AnalyzedAt,
		TotalNodes:  g.Meta.NodeCount,
		TotalEdges:  g.Meta.EdgeCount,
		TotalRoutes: countRoutes(g),
		Tabs:        []graph.TabEntry{},
	})
}

func writeJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := newIndentEncoder(f)
	return enc.Encode(v)
}

func countRoutes(g *graph.Graph) int {
	c := 0
	for _, n := range g.Nodes {
		if n.Type == graph.NodeTypeRoute {
			c++
		}
	}
	return c
}
