# boost

IDE coding-assistant integration for AI agents (Cursor, Claude Code, Copilot, and more).

## Overview

The `boost` package is a Go port of `laravel/boost`. It detects active AI coding
agents, writes MCP configuration files, and manages guidelines and skills for
nine supported agents: Cursor, Claude Code, Copilot, Codex, Gemini, Amp, Junie,
Kiro, and OpenCode.

**Module:** `github.com/bedrock/packages/boost`

```bash
go get github.com/bedrock/packages/boost@latest
```

## MCP Server Tools

| Tool              | Description                                |
| ----------------- | ------------------------------------------ |
| Application info  | Exposes app metadata to agents             |
| Database          | Schema introspection for AI context        |
| Log reader        | Surfaces recent log entries                |
| Docs search       | Searches documentation from the agent     |

## Coming Soon

Full documentation is in progress.
