# ai/mcp

<!-- upstream-docs: mcp.md#upstream-mcp -->
<!-- upstream-docs: mcp.md#creating-servers -->
<!-- upstream-docs: mcp.md#tools -->
<!-- upstream-docs: mcp.md#resources -->

Model Context Protocol (MCP) server implementation.

## Overview

The `ai/mcp` package provides a complete Go implementation of the MCP server
specification — a Go port of the Upstream MCP package. MCP allows applications
to expose themselves as servers that AI agents (such as Claude) can discover
and call. The protocol is built on JSON-RPC 2.0.

**Module:** `github.com/bedrock/packages/ai/mcp`

```bash
go get github.com/bedrock/packages/ai/mcp@latest
```

## Primitive Types

| Type        | Description                                          |
| ----------- | ---------------------------------------------------- |
| `Tools`     | Callable functions with typed JSON Schema inputs     |
| `Resources` | URI-addressed data (static and template-based URIs)  |
| `Prompts`   | Reusable prompt templates with argument substitution |
| `Sampling`  | Request LLM completions from the connected client    |

## Coming Soon

Full documentation is in progress.
