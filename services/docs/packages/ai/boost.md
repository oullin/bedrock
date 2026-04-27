# boost

<!-- upstream-docs: boost.md#introduction -->
<!-- upstream-docs: boost.md#features -->
<!-- upstream-docs: boost.md#documentation-api -->
<!-- upstream-docs: boost.md#ai-guidelines -->

Package boost provides a Go port of upstream/boost — an IDE coding-assistant integration layer. It includes agent detection and MCP configuration writing for nine AI coding agents (Cursor, Claude Code, Copilot, Codex, Gemini, Amp, Junie, Kiro, OpenCode), guidelines and skills file management, nine MCP server tools (application info, database introspection, log reading, docs search), and install utilities for writing config/guideline/skill files to disk.

<div class="docs-callout docs-callout-upstream">
  <strong>Upstream baseline.</strong>
  This page follows the Upstream 13.x documentation structure for the matching feature area, then rewrites the examples and edge cases for Bedrock's Go packages.
</div>

<div class="docs-callout docs-callout-go">
  <strong>Go adaptation.</strong>
  Bedrock replaces Upstream facades, service container magic, PHP traits, and CLI commands with explicit Go constructors, interfaces, structs, context propagation, and ordinary package tests.
</div>

## Installation

Install this module directly in applications that consume packages independently:

```bash
go get github.com/bedrock/packages/ai/boost@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/ai/boost/...
```

## Source Coverage

| Package               | Purpose                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `boost`               | Package boost provides a Go port of upstream/boost — an IDE coding-assistant integration layer. It includes agent detection and MCP configuration writing for nine AI coding agents (Cursor, Claude Code, Copilot, Codex, Gemini, Amp, Junie, Kiro, OpenCode), guidelines and skills file management, nine MCP server tools (application info, database introspection, log reading, docs search), and install utilities for writing config/guideline/skill files to disk. |
| `agents`              | Public agents API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                               |
| `guidelines`          | Public guidelines API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                           |
| `install`             | Public install API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                              |
| `internal/boosterr`   | Public internal/boosterr API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                    |
| `internal/jsonconfig` | Public internal/jsonconfig API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `internal/platform`   | Public internal/platform API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                    |
| `mcp`                 | Public mcp API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `mcp/tools`           | Public mcp/tools API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `skills`              | Public skills API surface for this module.                                                                                                                                                                                                                                                                                                                                                                                                                               |

## Core Concepts

The boost reference is organized around the exported Go surface for package `boost`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Upstream parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                                                                              |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `AgentOptions`, `AgentsDetector`, `Amp`, `ApplicationInfo`, `BaseAgent`, `BoostServiceProvider`, `BrowserLogs`, `ClaudeCode`, `Codex`, `CodingAgent`, `Content`, `Copilot`, `Cursor`, `DatabaseConnections`, `DatabaseQuery`, `DatabaseSchema`, `Gemini`, `GetAbsoluteUrl`, `Guideline`, `GuidelineAssist`, and 30 more                   |
| Constructors and functions | `AddFrontmatter`, `AppPath`, `ArtisanCommand`, `AvailableTools`, `BasePath`, `BinCommand`, `BoostSkills`, `Call`, `ClearCache`, `Compose`, `ComposeGuidelines`, `ComposerCommand`, `Current`, `CurrentPlatform`, `CustomGuidelinePath`, `CustomPath`, `DefaultMcpConfig`, `Description`, `DetectInProject`, `DetectOnSystem`, and 94 more |
| Variables                  | `ErrAgentAlreadyRegistered`, `ErrMcpInstallFailed`, `ErrNoMcpConfigPath`                                                                                                                                                                                                                                                                  |
| Constants                  | `Darwin`, `Linux`, `McpStrategyFile`, `McpStrategyNone`, `McpStrategyShell`, `PlatformDarwin`, `PlatformLinux`, `PlatformWindows`, `Windows`                                                                                                                                                                                              |

### Capability Matrix

| Capability                         | Documentation note                                                                                                   |
| ---------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Drivers and managers               | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Database-backed persistence        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/ai/boost"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/ai/boost` cover the supported creation paths, default values, and Upstream parity behavior.

## Configuration

Upstream documents many features through configuration files. Bedrock documents the equivalent behavior through Go options and constructor arguments:

| Upstream shape     | Bedrock shape                                            |
| ----------------- | -------------------------------------------------------- |
| Config file keys  | Typed config structs, options, or constructor parameters |
| Facade defaults   | Explicit manager/default-driver setup                    |
| Service providers | Go service-provider structs or direct application wiring |
| Runtime helpers   | Package functions and interfaces                         |

Prefer narrow interfaces at package boundaries. When a package exposes a manager, register drivers or providers at startup, set the default once, and resolve named instances per request or job.

## Advanced Features

The package reference should be read through these Upstream parity lenses:

| Area              | Documentation coverage                                                                  |
| ----------------- | --------------------------------------------------------------------------------------- |
| Drivers/providers | Available implementations, default selection, custom registration, and failure behavior |
| Events            | Emitted structs, dispatcher hooks, listener timing, transaction or queue interaction    |
| Errors            | Exported sentinel errors, wrapping, and `errors.Is` compatibility                       |
| Context           | Which operations accept `context.Context` and how cancellation/deadlines propagate      |
| Testing           | Fakes, null implementations, assertion helpers, and deterministic clocks/stores         |

## Edge Cases

- Do not translate PHP-only behavior literally. If Upstream depends on PHP traits, request globals, Template, CLI, or Orm magic, document the Bedrock Go equivalent instead.
- Preserve error identity when the package exports sentinel errors; callers should be able to use `errors.Is` where the package promises it.
- Treat driver compatibility as observable behavior. Unsupported store/driver combinations should be documented as errors or explicit no-ops, never as silent omissions.
- For I/O paths, document cancellation and timeout behavior whenever the package accepts a `context.Context`.
- For test fakes, document whether assertions inspect recorded calls, stored payloads, emitted events, or rendered output.

## Testing

Run the package tests before changing examples:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/ai/boost/...
```

Upstream parity is tracked by these tests:

- `packages/ai/boost/agents/inventory_parity_test.go`
- `packages/ai/boost/boost_inventory_test.go`
- `packages/ai/boost/guidelines/inventory_parity_test.go`
- `packages/ai/boost/install/inventory_parity_test.go`
- `packages/ai/boost/mcp/inventory_parity_test.go`
- `packages/ai/boost/mcp/tools/inventory_parity_test.go`
- `packages/ai/boost/skills/inventory_parity_test.go`

## API Reference

### Exported Types

| Type                      | Notes                                                                              |
| ------------------------- | ---------------------------------------------------------------------------------- |
| `AgentOptions`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AgentsDetector`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Amp`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ApplicationInfo`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BaseAgent`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BoostServiceProvider`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BrowserLogs`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClaudeCode`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Codex`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CodingAgent`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Content`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Copilot`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Cursor`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseConnections`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseQuery`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DatabaseSchema`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Gemini`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAbsoluteUrl`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Guideline`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GuidelineAssist`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GuidelineComposer`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GuidelineConfig`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GuidelineWriter`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Junie`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Kiro`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LastError`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Manager`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MarkdownFormatter`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpInstallationStrategy` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpRequest`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpResponse`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpTool`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpWriter`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OpenCode`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Platform`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReadLogEntries`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SearchDocs`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Server`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Skill`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SkillComposer`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SkillEntry`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SkillWriter`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsGuidelines`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsGuidelinesPath`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsMcp`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsMcpConfigPath`   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsSkills`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsSkillsPath`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToolExecutor`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ToolRegistry`            | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                         | Notes                                                                              |
| -------------------------------- | ---------------------------------------------------------------------------------- |
| `AddFrontmatter`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AppPath`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ArtisanCommand`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AvailableTools`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BasePath`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BinCommand`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `BoostSkills`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Call`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ClearCache`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Compose`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ComposeGuidelines`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ComposerCommand`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Current`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CurrentPlatform`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CustomGuidelinePath`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CustomPath`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DefaultMcpConfig`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Description`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DetectInProject`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DetectOnSystem`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DiscoverExplicitUserSkills`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DiscoverProjectInstalledAgents` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DiscoverSkillsFromPath`         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DiscoverSystemInstalledAgents`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DisplayName`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EntryPointPath`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrorResponse`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Execute`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExecuteReadOnly`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Executor`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Find`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Frontmatter`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAgents`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetAvailableTools`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GetToolNames`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GoBinaryPath`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Guidelines`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GuidelinesPath`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Handle`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasMcpEnabled`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasPackage`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasSkillsEnabled`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HttpMcpServerConfig`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InstallHttpMcp`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `InstallMcp`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsReadOnly`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsToolAllowed`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JSONContent`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpConfigKey`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpConfigPath`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpInstallationStrategy`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpServerConfig`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Models`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewAmp`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBaseAgent`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewBoostServiceProvider`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewClaudeCode`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCodex`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCopilot`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewCursor`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewDetector`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewExecutor`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGemini`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGuidelineAssist`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGuidelineComposer`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGuidelineConfig`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewGuidelineWriter`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewJunie`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewKiro`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewOpenCode`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewRegistry`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewServer`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewServerWithRegistry`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSkillComposer`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewSkillWriter`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NodePackageManager`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NodePackageManagerCommand`      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NormalizeHeadings`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `OkResponse`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PackageName`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Packages`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseSkill`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ParseSkillFrontmatter`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RegisterAgent`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Registry`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SailBinaryPath`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Schema`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetAllowed`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SetReadOnly`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShellMcpCommand`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ShouldEnforceStrictTypes`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SkillContent`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SkillName`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Skills`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SkillsPath`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `String`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `StripFrontmatter`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SupportsPintAgentFormatter`     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextContent`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TextResponse`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ThirdPartySkills`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `TransformGuidelines`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Trim`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UseAbsolutePathForMcp`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Used`                           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `UserSkills`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WithTimeout`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Write`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WriteEntry`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `WriteFormatted`                 | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name                        | Notes                                                                              |
| --------------------------- | ---------------------------------------------------------------------------------- |
| `Darwin`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrAgentAlreadyRegistered` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrMcpInstallFailed`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNoMcpConfigPath`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Linux`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpStrategyFile`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpStrategyNone`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `McpStrategyShell`          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PlatformDarwin`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PlatformLinux`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `PlatformWindows`           | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Windows`                   | Source-backed public surface. See the Go package for exact signature and behavior. |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
