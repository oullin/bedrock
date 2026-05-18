// Package ai produces the deterministic AI context export and editor-rules
// files brain ships. Token budgeting is delegated to packages/ai/sdk and the
// editor-rules TARGETS table is derived from packages/ai/boost; this package
// is the thin seam between brain and those shared implementations.
package ai

import (
	"fmt"
	"sort"
	"strings"

	aisdk "github.com/bedrock/packages/ai/sdk"
	"github.com/bedrock/services/brain/api/graph"
)

// EstimateTokens returns a coarse OpenAI-style token estimate. The seam is
// kept as a package-level variable so brain can swap in a different tokenizer
// from packages/ai/sdk without touching call sites.
var EstimateTokens = aisdk.EstimateTokens

// ContextOptions tunes the export shape.
type ContextOptions struct {
	// MaxNodes truncates the per-type listings.
	MaxNodes int
	// IncludeData controls whether each node's `data` bag is rendered.
	IncludeData bool
}

// RenderMarkdown returns a deterministic Markdown summary of the graph.
// Sections sorted alphabetically by node type so the same scan produces
// byte-identical output between runs.
func RenderMarkdown(g *graph.Graph, opts ContextOptions) string {
	if opts.MaxNodes == 0 {
		opts.MaxNodes = 50
	}

	groups := map[graph.NodeType][]*graph.Node{}

	for _, n := range g.Nodes {
		groups[n.Type] = append(groups[n.Type], n)
	}

	types := make([]graph.NodeType, 0, len(groups))

	for t := range groups {
		types = append(types, t)
	}

	sort.Slice(types, func(i, j int) bool { return types[i] < types[j] })

	var b strings.Builder

	fmt.Fprintf(&b, "# %s — brain context\n\n", g.Meta.Project)
	fmt.Fprintf(&b, "Scanned %d nodes / %d edges at %s.\n\n",
		g.Meta.NodeCount, g.Meta.EdgeCount, g.Meta.AnalyzedAt.Format("2006-01-02T15:04:05Z07:00"))

	for _, t := range types {
		ns := groups[t]

		sort.Slice(ns, func(i, j int) bool { return ns[i].ID < ns[j].ID })

		fmt.Fprintf(&b, "## %s (%d)\n\n", t, len(ns))
		limit := opts.MaxNodes

		if len(ns) < limit {
			limit = len(ns)
		}

		for _, n := range ns[:limit] {
			fmt.Fprintf(&b, "- **%s** — `%s`\n", n.Label, n.ID)

			if opts.IncludeData {
				keys := make([]string, 0, len(n.Data))

				for k := range n.Data {
					keys = append(keys, k)
				}

				sort.Strings(keys)

				for _, k := range keys {
					fmt.Fprintf(&b, "  - %s: %v\n", k, n.Data[k])
				}
			}
		}

		if len(ns) > opts.MaxNodes {
			fmt.Fprintf(&b, "- _…and %d more (omitted; raise --max-nodes)_\n", len(ns)-opts.MaxNodes)
		}

		b.WriteString("\n")
	}

	return b.String()
}

