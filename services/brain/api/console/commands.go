// Package console builds brain's three commands as native
// packages/console.Command values so a host service can register them
// directly on its own packages/console.Application:
//
//	import (
//		"github.com/bedrock/packages/console"
//		brainconsole "github.com/bedrock/services/brain/api/console"
//	)
//
//	app := console.NewApplication()
//	for _, cmd := range brainconsole.Commands() {
//		app.Add(cmd)
//	}
//
// No ad-hoc Command/Spec abstraction — the brain commands live alongside
// any other artisan-style command in the host's CLI.
package console

import (
	"context"
	"fmt"
	"path/filepath"

	bcons "github.com/bedrock/packages/console"
	"github.com/bedrock/services/brain/api/ai"
	"github.com/bedrock/services/brain/api/analysis"
	"github.com/bedrock/services/brain/api/graph"
)

// Commands returns the three brain commands ready to be registered.
// Errors from NewCommand are unwrapped via must — the signatures are
// constants, so a parse error here is a programmer bug rather than a
// runtime concern.
func Commands() []*bcons.Command {
	return []*bcons.Command{
		must(bcons.NewCommand(
			"brain:scan {--target=.} {--output=}",
			scanRun,
		)),
		must(bcons.NewCommand(
			"brain:export-context {--target=.}",
			exportContextRun,
		)),
		must(bcons.NewCommand(
			"brain:generate-rules {--target=.} {--force}",
			generateRulesRun,
		)),
	}
}

func scanRun(ctx context.Context, in *bcons.Input, out *bcons.Output) error {
	target, output := resolvePaths(in)
	g, err := scan(target)

	if err != nil {
		return err
	}

	if err := writeGraph(output, g); err != nil {
		return err
	}

	out.Writeln(fmt.Sprintf("brain: wrote %s (%d nodes, %d edges)",
		output, g.Meta.NodeCount, g.Meta.EdgeCount))
	_ = ctx

	return nil
}

func exportContextRun(_ context.Context, in *bcons.Input, out *bcons.Output) error {
	target, _ := resolvePaths(in)
	g, err := scan(target)

	if err != nil {
		return err
	}

	out.Write(ai.RenderMarkdown(g, ai.ContextOptions{}))

	return nil
}

func generateRulesRun(_ context.Context, in *bcons.Input, out *bcons.Output) error {
	target, _ := resolvePaths(in)
	force := in.Option("force") == "true"
	g, err := scan(target)

	if err != nil {
		return err
	}

	body := ai.RenderMarkdown(g, ai.ContextOptions{})
	written, err := ai.GenerateRules(target, body, force)

	if err != nil {
		return err
	}

	for _, p := range written {
		out.Writeln("brain: wrote " + p)
	}

	return nil
}

func scan(target string) (*graph.Graph, error) {
	return analysis.NewDefaultProjectAnalyzer().AnalyzeTarget(target)
}

func resolvePaths(in *bcons.Input) (target, output string) {
	target = in.Option("target")

	if target == "" {
		target = "."
	}

	if abs, err := filepath.Abs(target); err == nil {
		target = abs
	}

	output = in.Option("output")

	if output == "" {
		output = filepath.Join(target, "storage", "brain")
	}

	return
}

func writeGraph(dir string, g *graph.Graph) error {
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

func countRoutes(g *graph.Graph) int {
	c := 0

	for _, n := range g.Nodes {
		if n.Type == graph.NodeTypeRoute {
			c++
		}
	}

	return c
}

func must(cmd *bcons.Command, err error) *bcons.Command {
	if err != nil {
		panic("brain console signature: " + err.Error())
	}

	return cmd
}
