# MCP (Model Context Protocol)

<!-- laravel-docs: mcp.md#introduction -->
<!-- laravel-docs: mcp.md#installation -->
<!-- laravel-docs: mcp.md#creating-servers -->

Bedrock provides a complete Go implementation of the MCP server specification.

## Building a Server

```go
srv := mcp.NewServer("my-app", "1.0.0")

srv.AddTool(mcp.NewTool("greet", "Greet a user",
    schema,
    func(ctx context.Context, req *mcp.Request) (*mcp.Response, error) {
        return mcp.Text("Hello!"), nil
    },
))
```
