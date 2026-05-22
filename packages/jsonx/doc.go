// Package jsonx provides a fluent builder API for constructing JSON Schema
// Ref: @bedrock/code-0224
// package, offering type-safe builders for all JSON Schema primitive types
// (string, integer, number, boolean, array, object) with support for
// validation constraints, nullable types, required fields, and recursive
// schema composition.
//
// Usage:
//
//	schema := jsonx.Object(map[string]jsonx.SchemaType{
//	    "name": jsonx.String().Required(),
//	    "age":  jsonx.Integer().Min(0),
//	}).Title("User")
//
//	m := schema.ToMap()   // map[string]any
//	s := schema.String()  // JSON string
package jsonx
