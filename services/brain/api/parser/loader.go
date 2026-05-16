// Package parser loads Go source for analysis. It wraps go/packages so each
// analyzer sees a uniform view of the target service: an absolute root, a
// shared FileSet, and per-package syntax trees.
package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

// LoadedPackage pairs a packages.Package with its syntax files for AST walks.
type LoadedPackage struct {
	*packages.Package
}

// Project is the full result of loading a target service.
type Project struct {
	Root     string
	FileSet  *token.FileSet
	Packages []*packages.Package
}

// Load resolves `root` to an absolute path and loads every Go package under
// it via `./...`. Syntax, type info, imports, and dependencies are populated
// so analyzers can resolve selector expressions.
func Load(root string) (*Project, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve target %q: %w", root, err)
	}

	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedModule,
		Dir:  abs,
		Fset: fset,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("packages.Load: %w", err)
	}
	return &Project{Root: abs, FileSet: fset, Packages: pkgs}, nil
}

// EachFile iterates over every parsed AST file in the project alongside its
// owning package. Files outside the project root (vendored/transitive) are
// skipped so analyzers stay focused on the target service.
func (p *Project) EachFile(fn func(pkg *packages.Package, file *ast.File, filename string) error) error {
	for _, pkg := range p.Packages {
		for i, file := range pkg.Syntax {
			if i >= len(pkg.CompiledGoFiles) {
				continue
			}
			name := pkg.CompiledGoFiles[i]
			rel, err := filepath.Rel(p.Root, name)
			if err != nil || rel == "" || rel[0] == '.' && len(rel) > 1 && rel[1] == '.' {
				continue
			}
			if err := fn(pkg, file, name); err != nil {
				return err
			}
		}
	}
	return nil
}

// Position returns a human-readable "file:line" for a node in the project.
func (p *Project) Position(n ast.Node) string {
	if n == nil {
		return ""
	}
	pos := p.FileSet.Position(n.Pos())
	if rel, err := filepath.Rel(p.Root, pos.Filename); err == nil {
		return fmt.Sprintf("%s:%d", rel, pos.Line)
	}
	return fmt.Sprintf("%s:%d", pos.Filename, pos.Line)
}
