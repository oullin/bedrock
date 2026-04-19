// Package boost provides a Go port of upstream/boost — an IDE coding-assistant
// integration layer. It includes agent detection and MCP configuration writing
// for nine AI coding agents (Cursor, Claude Code, Copilot, Codex, Gemini, Amp,
// Junie, Kiro, OpenCode), guidelines and skills file management, nine MCP server
// tools (application info, database introspection, log reading, docs search), and
// install utilities for writing config/guideline/skill files to disk.
//
// This package is a 1:1 Go port of upstream/boost adapted to Go idioms.
// PHP-specific concepts are adapted: the PHP runtime path becomes GoBinaryPath,
// CLI becomes EntryPointPath, config() calls use config.Repository, and
// the abstract Agent class becomes the CodingAgent interface (named to avoid
// collision with contracts/ai.Agent, which models LLM agents).
package boost
