package analysis

import (
	"go/ast"
	"strings"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// httpVerbs is the set of router method names that register a route.
// Names match packages/routing: Get/Post/Put/Patch/Delete/Options/Head.

// RouteAnalyzer discovers HTTP route registrations. Matches:
//
//   - `<recv>.Get("/path", handler)` and the rest of httpVerbs.
//   - `<recv>.Handle("name", handler, mux)` from wayfinder.Registry.
//   - `<recv>.Add("name", "METHOD", "/path")` from wayfinder.Group.
//
// Phase 2 matches by method name and string-literal arity rather than by
// resolved receiver type. Later phases tighten this using type info from
// go/packages once the analyzer-list grows large enough to warrant it.
type RouteAnalyzer struct{}

var httpVerbs = map[string]string{
	"Get":     "GET",
	"Post":    "POST",
	"Put":     "PUT",
	"Patch":   "PATCH",
	"Delete":  "DELETE",
	"Options": "OPTIONS",
	"Head":    "HEAD",
	"Any":     "ANY",
}

func (RouteAnalyzer) Name() string { return "route" }

func (RouteAnalyzer) Analyze(ctx *Context) error {
	return ctx.Project.EachFile(func(_ *packages.Package, file *ast.File, filename string) error {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)

			if !ok {
				return true
			}

			if node := matchVerbCall(call); node != nil {
				node.Set("file", ctx.Project.Position(call))
				ctx.Graph.AddNode(node)

				return true
			}

			if node := matchWayfinderAdd(call); node != nil {
				node.Set("file", ctx.Project.Position(call))
				ctx.Graph.AddNode(node)

				return true
			}

			if node := matchWayfinderHandle(call); node != nil {
				node.Set("file", ctx.Project.Position(call))
				ctx.Graph.AddNode(node)

				return true
			}

			return true
		})
		_ = filename

		return nil
	})
}

// matchVerbCall handles `<recv>.Get("/path", handler[, ...])`.
func matchVerbCall(call *ast.CallExpr) *graph.Node {
	sel, ok := call.Fun.(*ast.SelectorExpr)

	if !ok {
		return nil
	}

	method, ok := httpVerbs[sel.Sel.Name]

	if !ok {
		return nil
	}

	if len(call.Args) < 1 {
		return nil
	}

	path, ok := stringLit(call.Args[0])

	if !ok {
		return nil
	}

	if !looksLikeURIPath(path) {
		return nil
	}

	id := "route:" + strings.ToLower(method) + ":" + path

	return graph.NewNode(id, graph.NodeTypeRoute, method+" "+path).
		Set("method", method).
		Set("uri", path).
		Set("source", "router")
}

// matchWayfinderAdd handles `g.Add("name", "METHOD", "/path"[, ...])`.
func matchWayfinderAdd(call *ast.CallExpr) *graph.Node {
	sel, ok := call.Fun.(*ast.SelectorExpr)

	if !ok || sel.Sel.Name != "Add" || len(call.Args) < 3 {
		return nil
	}

	name, ok1 := stringLit(call.Args[0])
	method, ok2 := stringLit(call.Args[1])
	path, ok3 := stringLit(call.Args[2])

	if !ok1 || !ok2 || !ok3 || !isHTTPMethodLiteral(method) || !looksLikeURIPath(path) {
		return nil
	}

	id := "route:" + strings.ToLower(method) + ":" + path

	return graph.NewNode(id, graph.NodeTypeRoute, strings.ToUpper(method)+" "+path).
		Set("method", strings.ToUpper(method)).
		Set("uri", path).
		Set("name", name).
		Set("source", "wayfinder")
}

// matchWayfinderHandle handles `r.Handle("name", handler, mux)` — a route
// with no inline method/path; we record the name only and the controller
// analyzer (phase 3) will link it to its handler.
func matchWayfinderHandle(call *ast.CallExpr) *graph.Node {
	sel, ok := call.Fun.(*ast.SelectorExpr)

	if !ok || sel.Sel.Name != "Handle" || len(call.Args) < 2 {
		return nil
	}

	name, ok := stringLit(call.Args[0])

	if !ok || !looksLikeRouteName(name) {
		return nil
	}

	id := "route:name:" + name

	return graph.NewNode(id, graph.NodeTypeRoute, name).
		Set("name", name).
		Set("source", "wayfinder.handle")
}

func stringLit(e ast.Expr) (string, bool) {
	bl, ok := e.(*ast.BasicLit)

	if !ok {
		return "", false
	}

	if len(bl.Value) < 2 {
		return "", false
	}

	first, last := bl.Value[0], bl.Value[len(bl.Value)-1]

	if first != '"' && first != '`' {
		return "", false
	}

	if last != first {
		return "", false
	}

	return bl.Value[1 : len(bl.Value)-1], true
}

func looksLikeURIPath(s string) bool {
	return strings.HasPrefix(s, "/")
}

// looksLikeRouteName accepts dotted/slug names like "auth.login" or "home".
func looksLikeRouteName(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '.' || r == '_' || r == '-':
		default:
			return false
		}
	}

	return true
}

func isHTTPMethodLiteral(s string) bool {
	switch strings.ToUpper(s) {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD":
		return true
	}

	return false
}
