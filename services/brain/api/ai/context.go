// Package ai produces the deterministic AI context export and editor-rules
// files brain ships. Per the plan, the heavy lifting belongs in
// packages/ai/sdk (token budgeting) and packages/ai/boost (TARGETS table
// for rules). Until those dependencies land in services/brain/api/go.mod
// the package owns its own minimal implementations of both — the public
// surface here is the seam where the swap happens.
package ai

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bedrock/services/brain/api/graph"
)

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

// EstimateTokens returns a coarse OpenAI-style token estimate (4 chars/tok).
// packages/ai/sdk will replace this with a proper tokenizer.
func EstimateTokens(s string) int {
	if s == "" {
		return 0
	}

	return (len(s) + 3) / 4
}
