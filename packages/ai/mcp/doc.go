// Package mcp provides a complete Go implementation of the Model Context
// Protocol (MCP) server specification. It is a behavioural port of the
// Upstream MCP package, adapted idiomatically to Go.
//
// MCP allows applications to expose themselves as servers that AI agents
// (such as Claude) can discover and call. The protocol is built on JSON-RPC
// 2.0 and supports four primitive types:
//
//   - Tools — callable functions with typed JSON Schema inputs.
//   - Resources — URI-addressed data, both at static URIs and via URI
//     templates (e.g. "file://users/{id}").
//   - Prompts — parameterised AI instruction templates that produce a
//     sequence of role/content messages.
//   - Completions — argument value suggestions for prompts and resources
//     that implement the Completable interface.
//
// # Building a server
//
//	srv := mcp.NewServer("my-app", "1.0.0",
//	    mcp.WithInstructions("You are a helpful assistant."),
//	)
//
//	srv.AddTool(mcp.NewTool("greet", "Greet a user",
//	    map[string]any{
//	        "type": "object",
//	        "properties": map[string]any{
//	            "name": map[string]any{"type": "string"},
//	        },
//	        "required": []string{"name"},
//	    },
//	    func(ctx context.Context, req *mcp.Request) (*mcp.Response, error) {
//	        name, _ := req.Get("name").(string)
//	        return mcp.Text("Hello, " + name + "!"), nil
//	    },
//	))
//
// # Transports
//
// Two transports are provided:
//
//   - HTTP with optional Server-Sent Events streaming: register the server
//     as an http.Handler with http.Handle("/mcp", srv).
//   - Stdio (stdin/stdout line-delimited JSON): call srv.ServeStdio(ctx).
//
// # Testing
//
//	result := srv.Test(t).CallTool("greet", map[string]any{"name": "World"})
//	result.AssertOK().AssertSee("Hello, World!")
//
// # Protocol versions
//
// The server negotiates the protocol version during the initialize handshake.
// Supported versions: 2025-11-25, 2025-06-18, 2025-03-26, 2024-11-05.
package mcp
