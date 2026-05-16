package analysis

import (
	"go/ast"
	"sort"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// ValidationAnalyzer captures rule sets passed to packages/validation:
//
//   validation.NewFactory().Validate(input, map[string]any{ ... }, nil, nil)
//   v.Make(input, map[string]any{ ... })
//
// We find the *ast.CompositeLit rule map (any arg position) and emit one
// validation_request node per call site, storing the rule keys.
type ValidationAnalyzer struct{}

func (ValidationAnalyzer) Name() string { return "validation" }

func (ValidationAnalyzer) Analyze(ctx *Context) error {
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
			if sel.Sel.Name != "Validate" && sel.Sel.Name != "Make" {
				return true
			}
			rulesMap := findRuleMap(call.Args)
			if rulesMap == nil {
				return true
			}
			fields := mapKeysOf(rulesMap)
			if len(fields) == 0 {
				return true
			}
			pos := ctx.Project.Position(call)
			id := "validation:" + pos
			node := graph.NewNode(id, graph.NodeTypeValidationRequest, "validation@"+pos).
				Set("file", pos).
				Set("fields", fields)
			_ = pkg
			ctx.Graph.AddNode(node)
			return true
		})
		return nil
	})
}

// findRuleMap picks the most-populated map literal among the args. When
// validation.Validate is called as Validate(input, rules, ...) both args are
// map literals — the rule set is the one with the most string-keyed entries
// (rule keys are always string literals; input keys often are too, which is
// why we tie-break on string-keyed count rather than total count).
func findRuleMap(args []ast.Expr) *ast.CompositeLit {
	var best *ast.CompositeLit
	bestScore := -1
	for _, a := range args {
		cl, ok := a.(*ast.CompositeLit)
		if !ok {
			continue
		}
		if _, isMap := cl.Type.(*ast.MapType); !isMap {
			continue
		}
		score := stringKeyCount(cl)
		if score > bestScore {
			best = cl
			bestScore = score
		}
	}
	if bestScore <= 0 {
		return nil
	}
	return best
}

func stringKeyCount(cl *ast.CompositeLit) int {
	c := 0
	for _, e := range cl.Elts {
		kv, ok := e.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		if _, ok := stringLit(kv.Key); ok {
			c++
		}
	}
	return c
}

func mapKeysOf(cl *ast.CompositeLit) []string {
	out := make([]string, 0, len(cl.Elts))
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		if k, ok := stringLit(kv.Key); ok {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}
