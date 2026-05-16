package analysis

import (
	"go/ast"
	"path"
	"strings"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// InertiaAnalyzer detects packages/inertia render calls. Two shapes:
//
//	inertia.Render(w, r, "Page/Path", props)
//	<handlerRecv>.Render(w, r, "Page/Path", props)
//
// The "Page/Path" string literal becomes an inertia_page node. Each top
// level key in the optional `protocol.Props{...}` composite literal becomes
// an inertia_prop node connected to the page with a `renders` edge.
type InertiaAnalyzer struct{}

func (InertiaAnalyzer) Name() string { return "inertia" }

func (InertiaAnalyzer) Analyze(ctx *Context) error {
	return ctx.Project.EachFile(func(_ *packages.Package, file *ast.File, _ string) error {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)

			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)

			if !ok || sel.Sel.Name != "Render" || len(call.Args) < 3 {
				return true
			}

			pageArg := -1

			var page string

			for i, a := range call.Args {
				if s, ok := stringLit(a); ok && looksLikeInertiaPage(s) {
					page = s
					pageArg = i

					break
				}
			}

			if page == "" {
				return true
			}

			pageID := "inertia_page:" + page
			ctx.Graph.AddNode(
				graph.NewNode(pageID, graph.NodeTypeInertiaPage, page).
					Set("file", ctx.Project.Position(call)),
			)

			if layout := layoutFor(page); layout != "" {
				layoutID := "inertia_layout:" + layout
				ctx.Graph.AddNode(graph.NewNode(layoutID, graph.NodeTypeInertiaLayout, layout))
				ctx.Graph.AddEdge(&graph.Edge{
					Source: layoutID, Target: pageID,
					Type: graph.EdgeTypeRenders, Label: "layout",
				})
			}

			if pageArg+1 < len(call.Args) {
				if cl, ok := call.Args[pageArg+1].(*ast.CompositeLit); ok {
					for _, k := range mapKeysOf(cl) {
						propID := "inertia_prop:" + page + "#" + k
						ctx.Graph.AddNode(
							graph.NewNode(propID, graph.NodeTypeInertiaProp, k).
								Set("page", page),
						)
						ctx.Graph.AddEdge(&graph.Edge{
							Source: pageID, Target: propID,
							Type: graph.EdgeTypeRenders, Label: "prop",
						})
					}
				}
			}

			return true
		})

		return nil
	})
}

// looksLikeInertiaPage accepts strings of the form Foo/Bar or Foo/Bar/Baz —
// PascalCase segments, slash-delimited, no spaces, at least one slash OR
// PascalCase one-segment names.
func looksLikeInertiaPage(s string) bool {
	if s == "" || strings.ContainsAny(s, " \t\n") {
		return false
	}

	if !isPascalSegment(strings.SplitN(s, "/", 2)[0]) {
		return false
	}

	return true
}

func isPascalSegment(s string) bool {
	if s == "" {
		return false
	}

	if s[0] < 'A' || s[0] > 'Z' {
		return false
	}

	return true
}

// layoutFor returns the first path segment if the page has more than one.
// "Auth/Login" → "Auth"; "Dashboard" → "".
func layoutFor(p string) string {
	if !strings.Contains(p, "/") {
		return ""
	}

	return path.Dir(p)
}
