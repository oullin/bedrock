package analysis

import (
	"go/ast"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// ChannelAnalyzer detects broadcast channel registrations from
// packages/broadcasting. The bedrock convention is:
//
//	broadcaster.Channel("orders.{order}", fn)
//
// The first arg is the channel name pattern.
type ChannelAnalyzer struct{}

func (ChannelAnalyzer) Name() string { return "channel" }

func (ChannelAnalyzer) Analyze(ctx *Context) error {
	return ctx.Project.EachFile(func(_ *packages.Package, file *ast.File, _ string) error {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)

			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)

			if !ok || sel.Sel.Name != "Channel" || len(call.Args) < 2 {
				return true
			}

			name, ok := stringLit(call.Args[0])

			if !ok || name == "" {
				return true
			}

			id := "channel:" + name
			ctx.Graph.AddNode(
				graph.NewNode(id, graph.NodeTypeChannel, name).
					Set("file", ctx.Project.Position(call)),
			)

			return true
		})

		return nil
	})
}
