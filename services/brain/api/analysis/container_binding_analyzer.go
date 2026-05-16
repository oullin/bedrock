package analysis

import (
	"go/ast"
	"strings"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// ContainerBindingAnalyzer detects calls into packages/container.Container:
//
//	c.Bind("abstract", factory, shared)
//	c.Singleton("abstract", factory)
//	c.Scoped("abstract", factory)
//	c.BindIf / c.SingletonIf
//
// Each binding produces a service_provider node tagged with its lifetime
// (singleton/scoped/bind). Bedrock-canonical service detection is here
// rather than in a separate ServiceAnalyzer.
type ContainerBindingAnalyzer struct{}

func (ContainerBindingAnalyzer) Name() string { return "container_binding" }

var containerBindMethods = map[string]string{
	"Bind":        "bind",
	"BindIf":      "bind",
	"Singleton":   "singleton",
	"SingletonIf": "singleton",
	"Scoped":      "scoped",
	"ScopedIf":    "scoped",
	"Instance":    "instance",
}

func (ContainerBindingAnalyzer) Analyze(ctx *Context) error {
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

			lifetime, ok := containerBindMethods[sel.Sel.Name]

			if !ok || len(call.Args) < 2 {
				return true
			}

			abstract, ok := stringLit(call.Args[0])

			if !ok || strings.TrimSpace(abstract) == "" {
				return true
			}

			concrete := exprLabel(pkg, call.Args[1])

			if strings.Contains(concrete, "func(") || strings.Contains(concrete, " ") {
				concrete = ""
			}

			id := "service_provider:" + abstract
			node := graph.NewNode(id, graph.NodeTypeServiceProvider, abstract).
				Set("lifetime", lifetime).
				Set("file", ctx.Project.Position(call))

			if concrete != "" {
				node.Set("concrete", concrete)
			}

			ctx.Graph.AddNode(node)

			if concrete != "" {
				serviceID := "service:" + concrete
				ctx.Graph.AddNode(graph.NewNode(serviceID, graph.NodeTypeService, concrete))
				ctx.Graph.AddEdge(&graph.Edge{
					Source: id, Target: serviceID,
					Type: graph.EdgeTypeBindsToImpl, Label: lifetime,
				})
			}

			return true
		})

		return nil
	})
}
