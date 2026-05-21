// Package routing is a 1:1 Go port of laravel/framework 13.x src/@bedrock/Routing.
//
// The directory layout, type names, method names, and behavior mirror the
// upstream PHP package as faithfully as Go allows. PHP traits are realized as
// embedded structs; PHP attributes are realized via the [controllers.HasMiddleware]
// interface; Symfony's RouteCompiler is reimplemented under [compiler] using the
// Go regexp/syntax package (RE2).
//
// The accompanying _test.go files are translations of laravel/framework
// tests/Routing/*.php — each PHP test method maps to a t.Run subtest with the
// same snake_case name so cross-referencing against the upstream test file is
// trivial.
package routing
