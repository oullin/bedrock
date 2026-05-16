package analysis

import (
	"go/ast"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// EventAnalyzer matches calls into packages/events:
//
//   dispatcher.Dispatch(ctx, event)   → event node + edge from caller
//   dispatcher.Listen("Name", fn)     → listener edge into event node
//   dispatcher.Subscribe(subscriber)  → listener node tagged subscriber
//
// Event identity is the type of the dispatched value or the string literal
// passed to Listen — whichever is available.
type EventAnalyzer struct{}

func (EventAnalyzer) Name() string { return "event" }

func (EventAnalyzer) Analyze(ctx *Context) error {
	return ctx.Project.EachFile(func(pkg *packages.Package, file *ast.File, _ string) error {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			switch sel.Sel.Name {
			case "Dispatch":
				recordDispatch(ctx, pkg, call)
			case "Listen":
				recordListen(ctx, pkg, call)
			case "Subscribe":
				recordSubscribe(ctx, pkg, call)
			}
			return true
		})
		return nil
	})
}

func recordDispatch(ctx *Context, pkg *packages.Package, call *ast.CallExpr) {
	if len(call.Args) < 2 {
		return
	}
	label := exprLabel(pkg, call.Args[len(call.Args)-1])
	if label == "" {
		return
	}
	id := "event:" + label
	ctx.Graph.AddNode(
		graph.NewNode(id, graph.NodeTypeEvent, label).Set("file", ctx.Project.Position(call)),
	)
}

func recordListen(ctx *Context, pkg *packages.Package, call *ast.CallExpr) {
	if len(call.Args) < 1 {
		return
	}
	name, ok := stringLit(call.Args[0])
	if !ok {
		name = exprLabel(pkg, call.Args[0])
	}
	if name == "" {
		return
	}
	eventID := "event:" + name
	ctx.Graph.AddNode(graph.NewNode(eventID, graph.NodeTypeEvent, name))
}

func recordSubscribe(ctx *Context, pkg *packages.Package, call *ast.CallExpr) {
	if len(call.Args) < 1 {
		return
	}
	label := exprLabel(pkg, call.Args[0])
	if label == "" {
		return
	}
	ctx.Graph.AddNode(
		graph.NewNode("event_subscriber:"+label, graph.NodeTypeEvent, label).
			Set("kind", "subscriber").
			Set("file", ctx.Project.Position(call)),
	)
}

// exprLabel returns a human label for an expression — its type name when
// possible, else its textual identifier.
func exprLabel(pkg *packages.Package, e ast.Expr) string {
	if pkg != nil && pkg.TypesInfo != nil {
		if tv, ok := pkg.TypesInfo.Types[e]; ok && tv.Type != nil {
			return unwrapType(tv.Type)
		}
	}
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return v.Sel.Name
	case *ast.CompositeLit:
		if id, ok := v.Type.(*ast.Ident); ok {
			return id.Name
		}
		if sel, ok := v.Type.(*ast.SelectorExpr); ok {
			return sel.Sel.Name
		}
	case *ast.UnaryExpr:
		return exprLabel(pkg, v.X)
	}
	return ""
}
