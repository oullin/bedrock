package analysis

import (
	"go/ast"
	"strings"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// FacadeAnalyzer flags every import under `…/packages/facades/<name>` and
// records a facade node per distinct facade package. Cheap and complete:
// no AST walking needed beyond the import block.
type FacadeAnalyzer struct{}

func (FacadeAnalyzer) Name() string { return "facade" }

func (FacadeAnalyzer) Analyze(ctx *Context) error {
	return ctx.Project.EachFile(func(_ *packages.Package, file *ast.File, _ string) error {
		for _, imp := range file.Imports {
			path, err := unquote(imp.Path.Value)
			if err != nil {
				continue
			}
			name, ok := facadeName(path)
			if !ok {
				continue
			}
			id := "facade:" + name
			ctx.Graph.AddNode(
				graph.NewNode(id, graph.NodeTypeFacade, name).
					Set("import", path).
					Set("file", ctx.Project.Position(imp)),
			)
		}
		return nil
	})
}

func facadeName(importPath string) (string, bool) {
	const marker = "/packages/facades/"
	i := strings.Index(importPath, marker)
	if i < 0 {
		return "", false
	}
	rest := importPath[i+len(marker):]
	if rest == "" {
		return "", false
	}
	if idx := strings.Index(rest, "/"); idx >= 0 {
		rest = rest[:idx]
	}
	return rest, true
}

func unquote(s string) (string, error) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return s, nil
	}
	return s[1 : len(s)-1], nil
}
