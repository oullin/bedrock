// Package wayfinder generates fully-typed, importable TypeScript functions
// for your Go routes.
//
// It is a Go library for Wayfinder package, producing identical
// TypeScript output so the same vitest test suite can validate generated files.
//
// Usage:
//
//	routes := wayfinder.FromRouteCollection(router.GetRoutes(), wayfinder.AdapterOptions{})
//	err := wayfinder.Generate(routes, wayfinder.Options{
//	    Path:     "resources/js",
//	    WithForm: true,
//	})
//
// The generator writes three output directories:
//
//   - {Path}/actions/   — per-controller TypeScript files, one function per route
//   - {Path}/routes/    — named-route helpers grouped by name prefix
//   - {Path}/wayfinder/ — the runtime TypeScript utility (index.ts)
package wayfinder
