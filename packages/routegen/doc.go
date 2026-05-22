// Package routegen generates fully-typed, importable TypeScript functions
// for your Go routes.
//
// It is a Go library for RouteGen package, producing identical
// TypeScript output so the same vitest test suite can validate generated files.
//
// Usage:
//
//	routes := routegen.FromRouteCollection(router.GetRoutes(), routegen.AdapterOptions{})
//	err := routegen.Generate(routes, routegen.Options{
//	    Path:     "resources/js",
//	    WithForm: true,
//	})
//
// The generator writes three output directories:
//
//   - {Path}/actions/   — per-controller TypeScript files, one function per route
//   - {Path}/routes/    — named-route helpers grouped by name prefix
//   - {Path}/routegen/ — the runtime TypeScript utility (index.ts)
package routegen
