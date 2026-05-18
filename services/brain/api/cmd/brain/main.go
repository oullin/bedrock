// Brain is the static-analysis + viewer tool for bedrock services.
//
// Phase 1 only wires the skeleton: an empty graph + manifest are written so
// downstream phases can plug analyzers and the SPA can be developed against a
// real (if empty) output shape.
package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/bedrock/packages/filesystem"
	"github.com/bedrock/services/brain/api/analysis"
	"github.com/bedrock/services/brain/api/graph"
	brainhttp "github.com/bedrock/services/brain/api/http"
	"github.com/bedrock/services/brain/api/watch"
)

//go:embed resources/views/app.html
var spaHTML []byte

type runOpts struct {
	target string
	output string
	quiet  bool
}

const (
	manifestFile = ".graph-manifest.json"
	allGraphFile = ".graph-all.json"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "brain:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		printUsage(stdout)

		return nil
	}

	sub, rest := args[0], args[1:]

	switch sub {
	case "scan":
		return runScan(rest, stdout, stderr)
	case "serve":
		return runServe(rest, stdout, stderr)
	case "-h", "--help", "help":
		printUsage(stdout)

		return nil
	case "-v", "--version", "version":
		fmt.Fprintln(stdout, "brain 0.0.0 (skeleton)")

		return nil
	default:
		return fmt.Errorf("unknown subcommand %q (try 'brain help')", sub)
	}
}

func runScan(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(stderr)
	target := fs.String("target", ".", "path to the Go service to scan")
	output := fs.String("output", "", "directory for graph JSON (default: <target>/storage/brain)")
	quiet := fs.Bool("quiet", false, "suppress progress output")
	watchMode := fs.Bool("watch", false, "rescan whenever a .go file changes")

	if err := fs.Parse(args); err != nil {
		return err
	}

	opts := runOpts{target: *target, output: *output, quiet: *quiet}

	absTarget, err := filepath.Abs(opts.target)

	if err != nil {
		return fmt.Errorf("resolve target: %w", err)
	}

	outDir := opts.output

	if outDir == "" {
		outDir = filepath.Join(absTarget, "storage", "brain")
	}

	var g *graph.Graph
	analyzer := analysis.NewDefaultProjectAnalyzer()

	if scanned, err := analyzer.AnalyzeTarget(absTarget); err == nil {
		g = scanned
	} else {
		fmt.Fprintf(stderr, "brain: scan warning: %v\n", err)
		g = graph.NewGraph(filepath.Base(absTarget))
		g.Stamp()
	}

	if err := writeJSON(filepath.Join(outDir, allGraphFile), g); err != nil {
		return err
	}

	m := graph.Manifest{
		Project:     g.Meta.Project,
		AnalyzedAt:  g.Meta.AnalyzedAt,
		TotalRoutes: countRoutes(g),
		TotalNodes:  g.Meta.NodeCount,
		TotalEdges:  g.Meta.EdgeCount,
		Tabs:        []graph.TabEntry{},
	}

	if err := writeJSON(filepath.Join(outDir, manifestFile), m); err != nil {
		return err
	}

	if !opts.quiet {
		fmt.Fprintf(stdout, "brain: wrote %s (%d nodes, %d edges)\n",
			outDir, g.Meta.NodeCount, g.Meta.EdgeCount)
	}

	if !*watchMode {
		return nil
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer cancel()

	fmt.Fprintln(stdout, "brain: watching for changes (Ctrl-C to stop)…")

	return watch.Run(ctx, watch.Options{Roots: []string{absTarget}}, func() {
		if scanned, err := analyzer.AnalyzeTarget(absTarget); err == nil {
			_ = writeJSON(filepath.Join(outDir, allGraphFile), scanned)
			_ = writeJSON(filepath.Join(outDir, manifestFile), graph.Manifest{
				Project: scanned.Meta.Project, AnalyzedAt: scanned.Meta.AnalyzedAt,
				TotalRoutes: countRoutes(scanned),
				TotalNodes:  scanned.Meta.NodeCount, TotalEdges: scanned.Meta.EdgeCount,
				Tabs: []graph.TabEntry{},
			})

			if !opts.quiet {
				fmt.Fprintf(stdout, "brain: rescan → %d nodes, %d edges\n",
					scanned.Meta.NodeCount, scanned.Meta.EdgeCount)
			}
		}
	})
}

func runServe(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	target := fs.String("target", ".", "path to the Go service to view")
	addr := fs.String("addr", ":8080", "address to listen on")
	assets := fs.String("assets", "", "directory containing the built SPA assets (default: <target>/storage/dist/brain)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	absTarget, err := filepath.Abs(*target)

	if err != nil {
		return err
	}

	assetDir := *assets

	if assetDir == "" {
		assetDir = filepath.Join(absTarget, "storage", "dist", "brain")
	}

	srv := brainhttp.NewServer(absTarget, spaHTML, assetDir)

	if err := srv.EnsureScanned(); err != nil {
		fmt.Fprintf(stderr, "brain: initial scan warning: %v\n", err)
	}

	displayAddr := *addr

	if strings.HasPrefix(displayAddr, ":") {
		displayAddr = "localhost" + displayAddr
	}

	fmt.Fprintf(stdout, "brain: serving %s on http://%s/_request_cycle/\n", absTarget, displayAddr)

	return http.ListenAndServe(*addr, srv.Routes())
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

// writeJSON marshals v with 2-space indent and writes via
// packages/filesystem so missing parent dirs are created automatically.
func writeJSON(path string, v any) error {
	body, err := json.MarshalIndent(v, "", "  ")

	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}

	if err := filesystem.New().Put(path, append(body, '\n')); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	return nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, strings.TrimSpace(`
brain - bedrock service analyzer

Usage:
  brain <command> [flags]

Commands:
  scan     Scan a Go service and write graph JSON to <target>/storage/brain
  serve    Run the viewer UI (defaults to :8080)
  version  Print the brain version
  help     Show this message

Run 'brain <command> -h' for command-specific flags.
`))
}
