package analysis

import (
	"testing"

	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/parser"
)

func TestMiddlewareAnalyzer_DetectsStandardSignature(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "mw.go", `package fixture

import "net/http"

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func Plain(_ string) string { return "" }
`)
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	g := graph.NewGraph("fixture")

	if err := (MiddlewareAnalyzer{}).Analyze(&Context{Project: proj, Graph: g}); err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	mwIDs := map[string]bool{}

	for _, n := range g.Nodes {
		if n.Type == graph.NodeTypeMiddleware {
			mwIDs[n.ID] = true
		}
	}

	if !anyContains(mwIDs, "Logger") {
		t.Errorf("expected a middleware node for Logger; got %v", mwIDs)
	}

	if anyContains(mwIDs, "Plain") {
		t.Errorf("Plain should not be classified as middleware; got %v", mwIDs)
	}
}

func TestMiddlewareAnalyzer_DetectsConstructorSignature(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "ctor.go", `package fixture

import "net/http"

// WithTimeout returns middleware — the classic constructor shape.
func WithTimeout(d int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = d
			next.ServeHTTP(w, r)
		})
	}
}
`)
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	g := graph.NewGraph("fixture")
	_ = (MiddlewareAnalyzer{}).Analyze(&Context{Project: proj, Graph: g})

	found := false

	for _, n := range g.Nodes {
		if n.Type == graph.NodeTypeMiddleware && n.Label == "WithTimeout" {
			found = true

			break
		}
	}

	if !found {
		t.Errorf("expected middleware for constructor WithTimeout; nodes: %v", g.Nodes)
	}
}

func TestMiddlewareAnalyzer_IgnoresNonMiddlewareFunctions(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "noise.go", `package fixture

func Add(a, b int) int { return a + b }

type T struct{}

func (T) Method() {}
`)
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	g := graph.NewGraph("fixture")
	_ = (MiddlewareAnalyzer{}).Analyze(&Context{Project: proj, Graph: g})

	for _, n := range g.Nodes {
		if n.Type == graph.NodeTypeMiddleware {
			t.Errorf("unexpected middleware node %q in noise-only fixture", n.ID)
		}
	}
}

func anyContains(set map[string]bool, sub string) bool {
	for k := range set {
		if contains(k, sub) {
			return true
		}
	}

	return false
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && indexOf(s, sub) >= 0
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}

	return -1
}
