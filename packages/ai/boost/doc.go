// Package boost provides a Go port of laravel/boost — an IDE coding-assistant
// integration layer. It includes agent detection and MCP configuration writing
// for nine AI coding agents (Cursor, Claude Code, Copilot, Codex, Gemini, Amp,
// Junie, Kiro, OpenCode), guidelines and skills file management, nine MCP server
// tools (application info, database introspection, log reading, docs search), and
// install utilities for writing config/guideline/skill files to disk.
//
// ## Installation
//
// Install the boost command into your project:
//
// ```bash
// go run github.com/bedrock/packages/ai/boost/cmd/boost@latest install
// ```
//
// This will detect your AI agents and write the necessary MCP configurations,
// guidelines, and skill files to your project root.
//
// ## MCP Server
//
// Boost provides an MCP server that gives AI agents access to your application's
// context, including:
//
//   - Application info and structure
//   - Database introspection (PostgreSQL, MySQL, SQLite)
//   - Log reading (Tail parsing)
//   - Documentation search
//   - Route and model info
//
// ## AI Guidelines and Skills
//
// Boost manages a `.github/linters/guidelines.md` file that provides high-level
// instructions to AI agents about your project's architecture and standards.
//
// Skill files (`SKILL.md`) provide specialized knowledge for specific tasks,
// such as creating new models, writing tests, or deploying the application.
package boost
