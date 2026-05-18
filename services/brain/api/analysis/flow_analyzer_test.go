package analysis

import (
	"testing"

	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/parser"
)

// FlowAnalyzer adds edges only for callers that already exist as nodes.
// In the fixtures below we pre-seed the action node and the target node
// so we can assert the analyzer drew the expected edge type.

func runFlowOn(t *testing.T, src string, seed func(*graph.Graph)) *graph.Graph {
	t.Helper()
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "flow.go", src)
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	g := graph.NewGraph("fixture")

	if seed != nil {
		seed(g)
	}

	if err := (FlowAnalyzer{}).Analyze(&Context{Project: proj, Graph: g}); err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	return g
}

func TestFlowAnalyzer_DispatchEdgeToKnownEvent(t *testing.T) {
	// FlowAnalyzer.Dispatch derives the label from the argument's *type*
	// (via exprLabel) — so the fixture passes an OrderShipped value, not
	// a string literal.
	src := `package fixture

type OrderShipped struct{ ID int }

type bus struct{}

func (b *bus) Dispatch(e any) {}

type Handler struct{ b *bus }

func (h *Handler) Do() {
	h.b.Dispatch(OrderShipped{})
}
`
	g := runFlowOn(t, src, func(g *graph.Graph) {
		g.AddNode(graph.NewNode("action:handler.do", graph.NodeTypeAction, "Handler.Do"))
		g.AddNode(graph.NewNode("event:OrderShipped", graph.NodeTypeEvent, "OrderShipped"))
	})

	if !hasEdge(g, "action:handler.do", "event:OrderShipped", graph.EdgeTypeDispatches) {
		t.Errorf("expected dispatches edge; edges: %v", g.Edges)
	}
}

func TestFlowAnalyzer_ListenEdgeToKnownEvent(t *testing.T) {
	src := `package fixture

type listener struct{}

func (l *listener) Listen(name string, fn func()) {}

func wire(l *listener) {
	l.Listen("OrderShipped", func() {})
}
`
	g := runFlowOn(t, src, func(g *graph.Graph) {
		g.AddNode(graph.NewNode("action:wire", graph.NodeTypeAction, "wire"))
		g.AddNode(graph.NewNode("event:OrderShipped", graph.NodeTypeEvent, "OrderShipped"))
	})

	if !hasEdge(g, "action:wire", "event:OrderShipped", graph.EdgeTypeListensTo) {
		t.Errorf("expected listens_to edge; edges: %v", g.Edges)
	}
}

func TestFlowAnalyzer_QueryEdgeOnDBMethod(t *testing.T) {
	src := `package fixture

type db struct{}

func (d *db) Query(s string) error { return nil }

type Repo struct{ d *db }

func (r *Repo) Find() {
	r.d.Query("SELECT 1")
}
`
	g := runFlowOn(t, src, func(g *graph.Graph) {
		g.AddNode(graph.NewNode("action:repo.find", graph.NodeTypeAction, "Repo.Find"))
	})

	found := false

	for _, e := range g.Edges {
		if e.Source == "action:repo.find" && e.Type == graph.EdgeTypeQueries && e.Label == "Query" {
			found = true

			break
		}
	}

	if !found {
		t.Errorf("expected queries edge with label 'Query'; edges: %v", g.Edges)
	}
}

func TestFlowAnalyzer_RenderEdgeToKnownInertiaPage(t *testing.T) {
	src := `package fixture

type inertia struct{}

func (i *inertia) Render(page string) {}

type Controller struct{ i *inertia }

func (c *Controller) Show() {
	c.i.Render("Dashboard/Index")
}
`
	g := runFlowOn(t, src, func(g *graph.Graph) {
		g.AddNode(graph.NewNode("action:controller.show", graph.NodeTypeAction, "Controller.Show"))
		g.AddNode(graph.NewNode("inertia_page:Dashboard/Index", graph.NodeTypeInertiaPage, "Dashboard/Index"))
	})

	if !hasEdge(g, "action:controller.show", "inertia_page:Dashboard/Index", graph.EdgeTypeRenders) {
		t.Errorf("expected renders edge; edges: %v", g.Edges)
	}
}

func TestFlowAnalyzer_NoEdgeWhenCallerNodeMissing(t *testing.T) {
	src := `package fixture

type bus struct{}

func (b *bus) Dispatch(name string) {}

func orphan() {
	var b bus
	b.Dispatch("X")
}
`
	g := runFlowOn(t, src, func(g *graph.Graph) {
		g.AddNode(graph.NewNode("event:X", graph.NodeTypeEvent, "X"))
	})
	// orphan() has no pre-existing action node — analyzer should skip it.
	if len(g.Edges) != 0 {
		t.Errorf("expected no edges when caller node missing; got %v", g.Edges)
	}
}

func TestFlowAnalyzer_NotifyAlwaysCreatesNotificationNodeAndEdge(t *testing.T) {
	// Like Dispatch, Notify labels by type via exprLabel.
	src := `package fixture

type WelcomeEmail struct{}

type mailer struct{}

func (m *mailer) Notify(n any) {}

type Service struct{ m *mailer }

func (s *Service) Welcome() {
	s.m.Notify(WelcomeEmail{})
}
`
	g := runFlowOn(t, src, func(g *graph.Graph) {
		g.AddNode(graph.NewNode("action:service.welcome", graph.NodeTypeAction, "Service.Welcome"))
	})

	if !hasEdge(g, "action:service.welcome", "notification:WelcomeEmail", graph.EdgeTypeNotifies) {
		t.Errorf("expected notifies edge; edges: %v", g.Edges)
	}

	if g.Node("notification:WelcomeEmail") == nil {
		t.Errorf("expected notification node to be created on demand")
	}
}
