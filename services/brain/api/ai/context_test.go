package ai

import (
	"strings"
	"testing"

	"github.com/bedrock/services/brain/api/graph"
)

func TestEstimateTokens(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"a", 1},
		{"abcd", 1},
		{"abcde", 2},
		{strings.Repeat("x", 100), 25},
	}

	for _, c := range cases {
		if got := EstimateTokens(c.in); got != c.want {
			t.Errorf("EstimateTokens(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestRenderMarkdownHeaderOnly(t *testing.T) {
	g := graph.NewGraph("svc")
	g.Stamp()
	out := RenderMarkdown(g, ContextOptions{})

	if !strings.HasPrefix(out, "# svc — brain context\n\n") {
		t.Fatalf("missing header; got:\n%s", out)
	}

	if !strings.Contains(out, "Scanned 0 nodes / 0 edges") {
		t.Errorf("missing scanned summary; got:\n%s", out)
	}
}

func TestRenderMarkdownSectionsSortedByType(t *testing.T) {
	g := graph.NewGraph("svc")
	g.AddNode(graph.NewNode("r:get:/", graph.NodeTypeRoute, "GET /"))
	g.AddNode(graph.NewNode("c:Auth", graph.NodeTypeController, "Auth"))
	g.AddNode(graph.NewNode("a:Show", graph.NodeTypeAction, "Show"))

	out := RenderMarkdown(g, ContextOptions{})

	// NodeType values are alphabetical strings: "action" < "controller" < "route".
	idxAction := strings.Index(out, "## action")
	idxCtrl := strings.Index(out, "## controller")
	idxRoute := strings.Index(out, "## route")

	if idxAction < 0 || idxCtrl < 0 || idxRoute < 0 {
		t.Fatalf("missing one or more sections; got:\n%s", out)
	}

	if !(idxAction < idxCtrl && idxCtrl < idxRoute) {
		t.Errorf("sections not sorted alphabetically: action=%d controller=%d route=%d",
			idxAction, idxCtrl, idxRoute)
	}
}

func TestRenderMarkdownTruncatesAtMaxNodes(t *testing.T) {
	g := graph.NewGraph("svc")

	for _, id := range []string{"r1", "r2", "r3", "r4", "r5"} {
		g.AddNode(graph.NewNode(id, graph.NodeTypeRoute, id))
	}

	out := RenderMarkdown(g, ContextOptions{MaxNodes: 2})

	if !strings.Contains(out, "r1") || !strings.Contains(out, "r2") {
		t.Errorf("expected first 2 IDs in output:\n%s", out)
	}

	if strings.Contains(out, "r3") || strings.Contains(out, "r4") || strings.Contains(out, "r5") {
		t.Errorf("expected truncation after 2 nodes:\n%s", out)
	}

	if !strings.Contains(out, "…and 3 more") {
		t.Errorf("expected '…and 3 more' footer:\n%s", out)
	}
}

func TestRenderMarkdownIncludesDataWhenRequested(t *testing.T) {
	g := graph.NewGraph("svc")
	g.AddNode(graph.NewNode("r:get:/", graph.NodeTypeRoute, "GET /").
		Set("zeta", "last").
		Set("alpha", "first"))

	out := RenderMarkdown(g, ContextOptions{IncludeData: true})

	if !strings.Contains(out, "- alpha: first") || !strings.Contains(out, "- zeta: last") {
		t.Errorf("missing data entries:\n%s", out)
	}
	// Keys sorted alphabetically — alpha must appear before zeta.
	idxAlpha := strings.Index(out, "- alpha:")
	idxZeta := strings.Index(out, "- zeta:")

	if idxAlpha < 0 || idxZeta < 0 || idxAlpha > idxZeta {
		t.Errorf("data keys not sorted: alpha=%d zeta=%d", idxAlpha, idxZeta)
	}
}

func TestRenderMarkdownDeterministic(t *testing.T) {
	build := func() *graph.Graph {
		g := graph.NewGraph("svc")
		g.AddNode(graph.NewNode("b", graph.NodeTypeRoute, ""))
		g.AddNode(graph.NewNode("a", graph.NodeTypeRoute, ""))
		g.AddNode(graph.NewNode("c", graph.NodeTypeController, ""))

		return g
	}

	a := RenderMarkdown(build(), ContextOptions{})
	b := RenderMarkdown(build(), ContextOptions{})

	if a != b {
		t.Fatalf("non-deterministic output:\n--- a ---\n%s\n--- b ---\n%s", a, b)
	}
}
