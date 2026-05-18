package analysis

import (
	"testing"

	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/parser"
)

func TestPhase4Analyzers(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "fac/facade.go", `package fac
// Stand-in for github.com/bedrock/packages/facades/cache.
`)
	mustWrite(t, dir, "wire.go", `package fixture

import _ "example.com/fixture/fac"

type dispatcher struct{}

func (d *dispatcher) Dispatch(ctx int, e any) {}
func (d *dispatcher) Listen(name string, fn func()) {}

type Mailer struct{}

type SendWelcomeJob struct{ User string }

func (SendWelcomeJob) QueueDisplayName() string { return "send-welcome" }

type console struct{}
func (console) NewCommand(sig string, fn func()) {}

type broadcaster struct{}
func (broadcaster) Channel(name string, fn func()) {}

type container struct{}
func (c *container) Singleton(name string, fn func() any) {}
func (c *container) Bind(name string, fn func() any, shared bool) {}

func wire(d *dispatcher, c *container, cli console, b broadcaster) {
	d.Dispatch(0, &Mailer{})
	d.Listen("OrderShipped", func() {})

	c.Singleton("mailer", func() any { return &Mailer{} })
	c.Bind("repo", func() any { return nil }, false)

	cli.NewCommand("brain:scan {--watch}", func() {})
	b.Channel("orders.{order}", func() {})
}
`)
	// Simulate the facade import path the analyzer recognises.
	mustWrite(t, dir, "facadeuse.go", `package fixture

// Trick the facade scanner via aliased import comment:
//   "github.com/bedrock/packages/facades/cache"
import _ "example.com/fixture/fac"
`)
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatal(err)
	}

	g := graph.NewGraph("fixture")
	ctx := &Context{Project: proj, Graph: g}

	for _, a := range []Analyzer{
		EventAnalyzer{}, JobAnalyzer{}, ConsoleAnalyzer{},
		ChannelAnalyzer{}, ContainerBindingAnalyzer{},
	} {
		if err := a.Analyze(ctx); err != nil {
			t.Fatalf("%s: %v", a.Name(), err)
		}
	}

	wantIDs := []string{
		"event:Mailer",
		"event:OrderShipped",
		"job:example.com/fixture.SendWelcomeJob",
		"command:brain:scan",
		"channel:orders.{order}",
		"service_provider:mailer",
		"service_provider:repo",
	}

	for _, id := range wantIDs {
		if g.Node(id) == nil {
			t.Errorf("missing node %s; have: %v", id, nodeIDs(g))
		}
	}
}

func TestFacadeNameExtraction(t *testing.T) {
	cases := map[string]string{
		"github.com/bedrock/packages/facades/auth":      "auth",
		"github.com/bedrock/packages/facades/cache":     "cache",
		"github.com/bedrock/packages/facades/log/extra": "log",
		"github.com/other/something":                    "",
	}

	for in, want := range cases {
		got, _ := facadeName(in)

		if got != want {
			t.Errorf("facadeName(%q) = %q, want %q", in, got, want)
		}
	}
}
