package analysis

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// ModelAnalyzer discovers struct types that look like persistence models.
// A struct qualifies when ANY of the following hold:
//
//   - It declares a `TableName() string` method (the bedrock/database ORM hint).
//   - One or more fields carry a `db:`, `gorm:`, `bedrock:` or `sql:` tag.
//   - The struct lives in a package whose import path contains
//     `/models/`, `/entities/`, or `/domain/`.
//
// We deliberately accept a wide net here so the graph shows realistic
// coverage on services that follow any of the common Go ORM conventions.
type ModelAnalyzer struct{}

func (ModelAnalyzer) Name() string { return "model" }

func (ModelAnalyzer) Analyze(ctx *Context) error {
	// First pass: collect names of types that have a TableName receiver.
	hasTableName := map[string]bool{}
	_ = ctx.Project.EachFile(func(_ *packages.Package, file *ast.File, _ string) error {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
				continue
			}
			if fn.Name == nil || fn.Name.Name != "TableName" {
				continue
			}
			hasTableName[recvTypeName(fn.Recv.List[0].Type)] = true
		}
		return nil
	})

	return ctx.Project.EachFile(func(pkg *packages.Package, file *ast.File, _ string) error {
		dirHint := pathHintsModel(pkg)
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				if !qualifiesAsModel(ts.Name.Name, st, dirHint, hasTableName) {
					continue
				}
				id := "model:" + pkg.PkgPath + "." + ts.Name.Name
				node := graph.NewNode(id, graph.NodeTypeModel, ts.Name.Name).
					Set("package", pkg.PkgPath).
					Set("file", ctx.Project.Position(ts))
				ctx.Graph.AddNode(node)
			}
		}
		return nil
	})
}

func qualifiesAsModel(name string, st *ast.StructType, dirHint bool, hasTableName map[string]bool) bool {
	if hasTableName[name] {
		return true
	}
	if structHasORMTag(st) {
		return true
	}
	return dirHint
}

func structHasORMTag(st *ast.StructType) bool {
	if st.Fields == nil {
		return false
	}
	for _, field := range st.Fields.List {
		if field.Tag == nil {
			continue
		}
		raw := field.Tag.Value
		if len(raw) < 2 {
			continue
		}
		tag := reflect.StructTag(raw[1 : len(raw)-1])
		for _, key := range []string{"db", "gorm", "bedrock", "sql"} {
			if _, ok := tag.Lookup(key); ok {
				return true
			}
		}
	}
	return false
}

func pathHintsModel(pkg *packages.Package) bool {
	if pkg == nil {
		return false
	}
	p := pkg.PkgPath
	return strings.Contains(p, "/models/") ||
		strings.Contains(p, "/entities/") ||
		strings.Contains(p, "/domain/")
}

func recvTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return recvTypeName(t.X)
	}
	return ""
}
