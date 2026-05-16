package analysis

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// ControllerAnalyzer links route registrations to the function that handles
// them. Phase 3 supports three handler shapes:
//
//   - `*ast.Ident` referring to a top-level func in the same package
//   - `*ast.SelectorExpr` referring to a top-level func in another package
//   - `*ast.SelectorExpr` referring to a method on a struct receiver
//     (e.g. `controller.Show`) — emits both a controller node and an action node
//
// Inline `func() { ... }` literals are intentionally skipped: anonymous
// handlers have no controller identity. The FlowExtractor in phase 6 will
// trace into their bodies separately.
type ControllerAnalyzer struct{}

func (ControllerAnalyzer) Name() string { return "controller" }

func (ControllerAnalyzer) Analyze(ctx *Context) error {
	return ctx.Project.EachFile(func(pkg *packages.Package, file *ast.File, _ string) error {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			routeNode := routeNodeForCall(call, ctx)
			if routeNode == nil {
				return true
			}
			handlerArg := handlerArgFor(call)
			if handlerArg == nil {
				return true
			}
			linkHandler(ctx, pkg, routeNode, handlerArg, call)
			return true
		})
		return nil
	})
}

// routeNodeForCall is the inverse of RouteAnalyzer.matchVerbCall — given a
// call, return the route node already in the graph (if any).
func routeNodeForCall(call *ast.CallExpr, ctx *Context) *graph.Node {
	// reuse the same matchers but only look up the resulting ID
	if n := matchVerbCall(call); n != nil {
		return ctx.Graph.Node(n.ID)
	}
	if n := matchWayfinderAdd(call); n != nil {
		return ctx.Graph.Node(n.ID)
	}
	if n := matchWayfinderHandle(call); n != nil {
		return ctx.Graph.Node(n.ID)
	}
	return nil
}

func handlerArgFor(call *ast.CallExpr) ast.Expr {
	if len(call.Args) < 2 {
		return nil
	}
	switch call.Args[1].(type) {
	case *ast.Ident, *ast.SelectorExpr:
		return call.Args[1]
	}
	return nil
}

func linkHandler(ctx *Context, pkg *packages.Package, route *graph.Node, handler ast.Expr, call ast.Node) {
	controllerLabel, actionLabel := describeHandler(pkg, handler)
	if actionLabel == "" {
		return
	}
	actionID := "action:" + strings.ToLower(actionLabel)
	action := graph.NewNode(actionID, graph.NodeTypeAction, actionLabel).
		Set("file", ctx.Project.Position(call))
	if controllerLabel != "" {
		action.Set("controller", controllerLabel)
	}
	ctx.Graph.AddNode(action)
	ctx.Graph.AddEdge(&graph.Edge{
		Source: route.ID, Target: actionID,
		Type:  graph.EdgeTypeHandlesBy,
		Label: "handles",
	})

	if controllerLabel == "" {
		return
	}
	controllerID := "controller:" + strings.ToLower(controllerLabel)
	ctrl := graph.NewNode(controllerID, graph.NodeTypeController, controllerLabel)
	ctx.Graph.AddNode(ctrl)
	ctx.Graph.AddEdge(&graph.Edge{
		Source: controllerID, Target: actionID,
		Type: graph.EdgeTypeCalls, Label: "method",
	})
}

// describeHandler returns (controller, action) labels for a handler ast.Expr.
// If the handler is a free function the controller label is empty.
func describeHandler(pkg *packages.Package, h ast.Expr) (string, string) {
	switch e := h.(type) {
	case *ast.Ident:
		// free function in the same package
		if e.Obj != nil && e.Obj.Decl != nil {
			if fn, ok := e.Obj.Decl.(*ast.FuncDecl); ok && fn.Name != nil {
				return "", fn.Name.Name
			}
		}
		return "", e.Name
	case *ast.SelectorExpr:
		// pkg.Func OR receiver.Method
		if pkg != nil && pkg.TypesInfo != nil {
			if sel := pkg.TypesInfo.Selections[e]; sel != nil {
				switch sel.Kind() {
				case types.MethodVal, types.MethodExpr:
					recv := unwrapType(sel.Recv())
					return recv, e.Sel.Name
				}
			}
		}
		if id, ok := e.X.(*ast.Ident); ok {
			return id.Name, e.Sel.Name
		}
	}
	return "", ""
}

func unwrapType(t types.Type) string {
	if t == nil {
		return ""
	}
	if p, ok := t.(*types.Pointer); ok {
		return unwrapType(p.Elem())
	}
	if n, ok := t.(*types.Named); ok {
		return n.Obj().Name()
	}
	return t.String()
}
