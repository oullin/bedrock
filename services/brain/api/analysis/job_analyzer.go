package analysis

import (
	"go/ast"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// JobAnalyzer detects queued jobs via:
//
//   - calls to `queue.CreatePayloadFor(_, _, job, _, _)` — the canonical
//     bedrock dispatch shape; the 3rd arg is the job value.
//   - types declaring a `QueueDisplayName() string` method (the bedrock
//     `Namer` contract) — surfaces jobs even when dispatch happens via a
//     facade or wrapper.
type JobAnalyzer struct{}

func (JobAnalyzer) Name() string { return "job" }

func (JobAnalyzer) Analyze(ctx *Context) error {
	// 1) jobs by type marker — QueueDisplayName()
	_ = ctx.Project.EachFile(func(pkg *packages.Package, file *ast.File, _ string) error {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
				continue
			}
			if fn.Name == nil || fn.Name.Name != "QueueDisplayName" {
				continue
			}
			recv := recvTypeName(fn.Recv.List[0].Type)
			if recv == "" {
				continue
			}
			id := "job:" + pkg.PkgPath + "." + recv
			ctx.Graph.AddNode(
				graph.NewNode(id, graph.NodeTypeJob, recv).
					Set("package", pkg.PkgPath).
					Set("file", ctx.Project.Position(fn)).
					Set("source", "namer"),
			)
		}
		return nil
	})

	// 2) jobs by dispatch site
	return ctx.Project.EachFile(func(pkg *packages.Package, file *ast.File, _ string) error {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "CreatePayloadFor" || len(call.Args) < 3 {
				return true
			}
			label := exprLabel(pkg, call.Args[2])
			if label == "" {
				return true
			}
			id := "job:" + label
			ctx.Graph.AddNode(
				graph.NewNode(id, graph.NodeTypeJob, label).
					Set("file", ctx.Project.Position(call)).
					Set("source", "dispatch"),
			)
			return true
		})
		return nil
	})
}
