package analysis

import (
	"go/ast"
	"strings"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// FlowAnalyzer walks every function body and, for each "interesting" call
// inside it, draws an edge from the enclosing function to the target node
// already in the graph. This is the Go analogue of upstream-brain's
// FlowExtractor + MethodTracer + QueryTracer rolled into a single pass.
//
// Pre-existing analyzers populate the target nodes (route, event, job,
// validation, inertia_page, service_provider, ...). FlowAnalyzer only adds
// edges, never new nodes — so it must run last.
//
// Recognised intra-body calls:
//
//	*.Query / Exec / QueryRow / QueryContext / ExecContext / Prepare → query edge
//	*.Dispatch                                                       → dispatches
//	*.Listen / Subscribe                                             → listens_to
//	*.Validate / Make (with map literal)                             → validates
//	*.Render (inertia page string)                                   → renders
//	*.CreatePayloadFor / Push (job)                                  → queues
//	*.Notify                                                         → notifies
type FlowAnalyzer struct{}

func (FlowAnalyzer) Name() string { return "flow" }

func (FlowAnalyzer) Analyze(ctx *Context) error {
	return ctx.Project.EachFile(func(pkg *packages.Package, file *ast.File, _ string) error {
		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)

			if !ok || fn.Body == nil {
				return true
			}

			caller := callerNodeID(fn)

			if caller == "" {
				return true
			}
			// only emit edges if a node with that ID already exists
			if ctx.Graph.Node(caller) == nil {
				return true
			}

			ast.Inspect(fn.Body, func(inner ast.Node) bool {
				call, ok := inner.(*ast.CallExpr)

				if !ok {
					return true
				}

				traceCall(ctx, pkg, caller, call)

				return true
			})

			return true
		})

		return nil
	})
}

// callerNodeID returns the ID of the enclosing function's controller/action
// node, if such a node will have been created upstream.
func callerNodeID(fn *ast.FuncDecl) string {
	if fn.Name == nil {
		return ""
	}

	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return "action:" + strings.ToLower(fn.Name.Name)
	}

	recv := recvTypeName(fn.Recv.List[0].Type)

	return "action:" + strings.ToLower(recv+"."+fn.Name.Name)
}

func traceCall(ctx *Context, pkg *packages.Package, caller string, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)

	if !ok {
		return
	}

	switch sel.Sel.Name {
	case "Query", "Exec", "QueryRow", "QueryContext", "ExecContext", "QueryRowContext", "Prepare":
		ctx.Graph.AddEdge(&graph.Edge{
			Source: caller, Target: caller + ":db",
			Type: graph.EdgeTypeQueries, Label: sel.Sel.Name,
		})
	case "Dispatch":
		if label := exprLabel(pkg, lastArg(call)); label != "" {
			if ctx.Graph.Node("event:"+label) != nil {
				ctx.Graph.AddEdge(&graph.Edge{
					Source: caller, Target: "event:" + label,
					Type: graph.EdgeTypeDispatches, Label: "dispatch",
				})
			}
		}
	case "Listen", "Subscribe":
		if len(call.Args) > 0 {
			if name, ok := stringLit(call.Args[0]); ok && ctx.Graph.Node("event:"+name) != nil {
				ctx.Graph.AddEdge(&graph.Edge{
					Source: caller, Target: "event:" + name,
					Type: graph.EdgeTypeListensTo, Label: sel.Sel.Name,
				})
			}
		}
	case "Validate", "Make":
		if findRuleMap(call.Args) != nil {
			pos := ctx.Project.Position(call)
			target := "validation:" + pos

			if ctx.Graph.Node(target) != nil {
				ctx.Graph.AddEdge(&graph.Edge{
					Source: caller, Target: target,
					Type: graph.EdgeTypeValidates, Label: "validate",
				})
			}
		}
	case "Render":
		for _, a := range call.Args {
			if s, ok := stringLit(a); ok && looksLikeInertiaPage(s) {
				if ctx.Graph.Node("inertia_page:"+s) != nil {
					ctx.Graph.AddEdge(&graph.Edge{
						Source: caller, Target: "inertia_page:" + s,
						Type: graph.EdgeTypeRenders, Label: "render",
					})
				}

				break
			}
		}
	case "CreatePayloadFor":
		if len(call.Args) >= 3 {
			if label := exprLabel(pkg, call.Args[2]); label != "" && ctx.Graph.Node("job:"+label) != nil {
				ctx.Graph.AddEdge(&graph.Edge{
					Source: caller, Target: "job:" + label,
					Type: graph.EdgeTypeQueues, Label: "queue",
				})
			}
		}
	case "Notify":
		if label := exprLabel(pkg, lastArg(call)); label != "" {
			target := "notification:" + label
			ctx.Graph.AddNode(graph.NewNode(target, graph.NodeTypeNotification, label))
			ctx.Graph.AddEdge(&graph.Edge{
				Source: caller, Target: target,
				Type: graph.EdgeTypeNotifies, Label: "notify",
			})
		}
	}
}

func lastArg(call *ast.CallExpr) ast.Expr {
	if len(call.Args) == 0 {
		return nil
	}

	return call.Args[len(call.Args)-1]
}
