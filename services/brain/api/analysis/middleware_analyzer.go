package analysis

import (
	"go/ast"
	"go/types"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// MiddlewareAnalyzer discovers functions with the canonical Go middleware
// signature `func(http.Handler) http.Handler`. Bedrock's
// packages/inertia/middleware and packages/httpx middlewares follow this
// shape, as does almost every public Go middleware in the wild.
//
// We also accept the closure-style `func(...) func(http.Handler) http.Handler`
// (a constructor returning a middleware) by recognising the return type.
type MiddlewareAnalyzer struct{}

func (MiddlewareAnalyzer) Name() string { return "middleware" }

func (MiddlewareAnalyzer) Analyze(ctx *Context) error {
	return ctx.Project.EachFile(func(pkg *packages.Package, file *ast.File, _ string) error {
		if pkg == nil || pkg.TypesInfo == nil {
			return nil
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)

			if !ok || fn.Name == nil {
				continue
			}

			if !looksLikeMiddleware(pkg, fn) {
				continue
			}

			id := "middleware:" + pkg.PkgPath + "." + fn.Name.Name
			node := graph.NewNode(id, graph.NodeTypeMiddleware, fn.Name.Name).
				Set("package", pkg.PkgPath).
				Set("file", ctx.Project.Position(fn))
			ctx.Graph.AddNode(node)
		}

		return nil
	})
}

func looksLikeMiddleware(pkg *packages.Package, fn *ast.FuncDecl) bool {
	sig := signatureOf(pkg, fn)

	if sig == nil {
		return false
	}

	if isHandlerToHandler(sig) {
		return true
	}
	// constructor: returns a func(http.Handler) http.Handler
	if sig.Results().Len() == 1 {
		if inner, ok := sig.Results().At(0).Type().(*types.Signature); ok {
			return isHandlerToHandler(inner)
		}
	}

	return false
}

func isHandlerToHandler(sig *types.Signature) bool {
	if sig.Params().Len() != 1 || sig.Results().Len() != 1 {
		return false
	}

	return isHTTPHandler(sig.Params().At(0).Type()) && isHTTPHandler(sig.Results().At(0).Type())
}

func isHTTPHandler(t types.Type) bool {
	named, ok := t.(*types.Named)

	if !ok {
		// interfaces from net/http are not named in the usual sense,
		// fall through to string match
		return t.String() == "net/http.Handler"
	}

	obj := named.Obj()

	if obj == nil || obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == "net/http" && obj.Name() == "Handler"
}

func signatureOf(pkg *packages.Package, fn *ast.FuncDecl) *types.Signature {
	if pkg == nil || pkg.TypesInfo == nil {
		return nil
	}

	def, ok := pkg.TypesInfo.Defs[fn.Name]

	if !ok || def == nil {
		return nil
	}

	sig, _ := def.Type().(*types.Signature)

	return sig
}
