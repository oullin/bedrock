package analysis

import (
	"go/ast"
	"strings"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// ConsoleAnalyzer detects artisan-style commands registered via
// packages/console: `console.NewCommand("signature", fn)`. The signature's
// leading token becomes the command name (Laravel convention preserved).
type ConsoleAnalyzer struct{}

func (ConsoleAnalyzer) Name() string { return "console" }

func (ConsoleAnalyzer) Analyze(ctx *Context) error {
	return ctx.Project.EachFile(func(_ *packages.Package, file *ast.File, _ string) error {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)

			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)

			if !ok || sel.Sel.Name != "NewCommand" || len(call.Args) < 1 {
				return true
			}

			sig, ok := stringLit(call.Args[0])

			if !ok {
				return true
			}

			name := strings.SplitN(strings.TrimSpace(sig), " ", 2)[0]

			if name == "" {
				return true
			}

			id := "command:" + name
			ctx.Graph.AddNode(
				graph.NewNode(id, graph.NodeTypeCommand, name).
					Set("signature", sig).
					Set("file", ctx.Project.Position(call)),
			)

			return true
		})

		return nil
	})
}
