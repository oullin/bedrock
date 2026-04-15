# jsonx

Fluent JSON Schema builder.

## Overview

The `jsonx` package provides type-safe builders for JSON Schema primitives. It
mirrors Upstream's `JsonSchemaTypeFactory` and is primarily used when defining
validation schemas for AI model structured outputs or API contracts.

**Module:** `github.com/gocanto/bedrock/packages/jsonx`

```bash
go get github.com/gocanto/bedrock/packages/jsonx@latest
```

## Schema Types

| Builder      | JSON Schema type |
|--------------|-----------------|
| `Object(…)`  | `"type": "object"` |
| `Array()`    | `"type": "array"` |
| `String()`   | `"type": "string"` |
| `Integer()`  | `"type": "integer"` |
| `Number()`   | `"type": "number"` |
| `Boolean()`  | `"type": "boolean"` |

## Building a Schema

```go
schema := jsonx.Object(map[string]jsonx.SchemaType{
    "name":  jsonx.String().MinLength(1).MaxLength(255),
    "age":   jsonx.Integer().Minimum(0).Maximum(150),
    "email": jsonx.String().Format("email"),
    "score": jsonx.Number().Minimum(0.0).Maximum(100.0),
    "active": jsonx.Boolean(),
    "tags":  jsonx.Array().Items(jsonx.String()),
})
```

## Closure Syntax

Build properties inline using the `Factory`:

```go
schema := jsonx.Object(func(f jsonx.Factory) map[string]jsonx.SchemaType {
    return map[string]jsonx.SchemaType{
        "title":   f.String().MinLength(1),
        "count":   f.Integer().Minimum(1),
        "enabled": f.Boolean(),
    }
})
```

## Serialization

```go
out, err := schema.Serialize() // map[string]any ready for JSON encoding
data, err := json.Marshal(out)
```

## String Constraints

```go
jsonx.String().
    MinLength(3).
    MaxLength(50).
    Pattern(`^[a-z]+$`).
    Format("email")
```

## Numeric Constraints

```go
jsonx.Integer().
    Minimum(0).
    Maximum(100).
    MultipleOf(5)

jsonx.Number().
    ExclusiveMinimum(0.0).
    ExclusiveMaximum(1.0)
```

## Array Constraints

```go
jsonx.Array().
    Items(jsonx.String()).
    MinItems(1).
    MaxItems(10).
    UniqueItems(true)
```

## Object Constraints

```go
jsonx.Object(properties).
    Required("name", "email").
    AdditionalProperties(false)
```
