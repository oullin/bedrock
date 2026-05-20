# Getting Started

Bedrock is a collection of foundational Go packages for building web applications.
Each package is a standalone Go module with its own `go.mod`, so you import only
what you need.

## Requirements

| Tool    | Version |
| ------- | ------- |
| Go      | ≥ 1.24  |
| Node.js | ≥ 22    |
| pnpm    | ≥ 10.33 |

## Installation

Add any Bedrock package directly with `go get`:

```bash
go get github.com/bedrock/packages/auth@latest
go get github.com/bedrock/packages/cache@latest
```

## Project Layout

```
packages/   Go library packages
services/   Internal tooling, docs, storage, and demo applications
```

## Development Setup

Clone the repository, install Node dependencies, and run the test suite:

```bash
git clone https://github.com/gocanto/bedrock.git
cd bedrock
pnpm install
pnpm run test
```

Run the documentation site locally:

```bash
pnpm run dev --filter=@bedrock/docs
```

Build the static site:

```bash
pnpm run build --filter=@bedrock/docs
# Output: services/docs/.vuepress/dist/
```

## Package Index

### Architecture

| Package                          | Purpose                                           |
| -------------------------------- | ------------------------------------------------- |
| [container](/packages/container) | IoC service container and application bootstrap   |
| [config](/packages/config)       | Configuration repository with dot-notation access |
| [contracts](/packages/contracts) | Shared interface definitions for every package    |

### The Basics

| Package                            | Purpose                                           |
| ---------------------------------- | ------------------------------------------------- |
| [routing](/packages/routing)       | HTTP routing                    |
| [httpx](/packages/httpx)           | HTTP utilities, middleware, and testing helpers   |
| [session](/packages/session)       | Session management with multiple storage handlers |
| [cookie](/packages/cookie)         | HTTP cookie handling                              |
| [validation](/packages/validation) | Rule-based input validation (80+ built-in rules)  |

### Security

| Package                            | Purpose                                            |
| ---------------------------------- | -------------------------------------------------- |
| [auth](/packages/auth)             | Authentication, authorization, password management |
| [encryption](/packages/encryption) | AES encryption with CBC and GCM mode support       |
| [hashing](/packages/hashing)       | Password hashing with bcrypt and Argon2            |
| [fortify](/packages/fortify)       | Rate limiting, two-factor auth, auth pipelines     |

### Data & Storage

| Package                            | Purpose                                        |
| ---------------------------------- | ---------------------------------------------- |
| [cache](/packages/cache)           | Caching layer with multiple driver support     |
| [redis](/packages/redis)           | Full Redis command surface, pipelines, pub/sub |
| [filesystem](/packages/filesystem) | Local filesystem operations                    |
| [pagination](/packages/pagination) | Offset-based and cursor-based pagination       |

### Events & Jobs

| Package                        | Purpose                                          |
| ------------------------------ | ------------------------------------------------ |
| [events](/packages/events)     | Event dispatching and listener management        |
| [bus](/packages/bus)           | Command and event bus with pipeline support      |
| [queue](/packages/queue)       | Background job processing with pluggable drivers |
| [horizon](/packages/horizon)   | Queue monitoring snapshots and metrics           |
| [pipeline](/packages/pipeline) | Middleware-style pipe-and-filter chains          |

### Communication

| Package                                  | Purpose                             |
| ---------------------------------------- | ----------------------------------- |
| [mailx](/packages/mailx)                 | Driver-based email sending          |
| [notifications](/packages/notifications) | Multi-channel notification delivery |

### AI & Integrations

| Package                        | Purpose                                            |
| ------------------------------ | -------------------------------------------------- |
| [ai/sdk](/packages/ai/sdk)     | Unified AI provider API, agents, and RAG tools     |
| [ai/mcp](/packages/ai/mcp)     | Model Context Protocol (MCP) server implementation |
| [ai/boost](/packages/ai/boost) | IDE coding-assistant and agent integration layer   |

### Support & Utilities

| Package                                  | Purpose                                          |
| ---------------------------------------- | ------------------------------------------------ |
| [collection](/packages/collection)       | Fluent collection helpers for slices and maps    |
| [support](/packages/support)             | Helpers, Fluent, Optional, MessageBag, strings   |
| [log](/packages/log)                     | Driver-based structured logging with channels    |
| [translation](/packages/translation)     | Localisation and i18n with CLDR pluralisation    |
| [concurrency](/packages/concurrency)     | Concurrent task execution with pluggable drivers |
| [conditionable](/packages/conditionable) | Conditional method execution with a fluent proxy |
| [jsonx](/packages/jsonx)                 | Fluent JSON Schema builder                       |

### Developer Tools

| Package                          | Purpose                                  |
| -------------------------------- | ---------------------------------------- |
| [prompts](/packages/prompts)     | Interactive terminal prompt components   |
| [pail](/packages/pail)           | Log tail parsing and filtering           |
| [envoy](/packages/envoy)         | Remote task planning and command running |
| [telescope](/packages/telescope) | Application introspection and watchers   |

### Products

| Package                          | Purpose                                              |
| -------------------------------- | ---------------------------------------------------- |
| [inception](/packages/inception) | Unified auth scaffold — Fortify + Jetstream combined |
| [jetstream](/packages/jetstream) | Team and organization management                     |
| [spark](/packages/spark)         | Subscription billing, checkout, and entitlements     |

## AI Assisted Development

Bedrock is built from the ground up to be AI-friendly. By providing clear
interfaces, predictable patterns, and explicit dependencies, Bedrock makes it
easy for AI agents to understand and contribute to your codebase.

Check out the [ai/boost](/packages/ai/boost) package to see how to integrate
Bedrock with your IDE coding assistants and agents.

## Concept Guides

Cross-cutting topics that span multiple packages:

| Guide                                        | What it covers                                          |
| -------------------------------------------- | ------------------------------------------------------- |
| [Request Lifecycle](/architecture/lifecycle) | How an HTTP request flows through a Bedrock application |
| [Testing](/concepts/testing)                 | Built-in test doubles and testing patterns              |
| [Middleware](/basics/middleware)             | Global, per-route, and controller-scoped middleware     |
| [Controllers](/basics/controllers)           | Grouping handlers into types with shared middleware     |
| [URL Generation](/basics/url-generation)     | Named-route URLs, signed URLs, redirects                |
| [CSRF Protection](/basics/csrf)              | Tokens, headers, excluding routes                       |

## Next Steps

Once your app starts up, you'll want to know how it's actually wired
together. The Architecture Concepts guides walk through that:

| Guide                                                       | What it answers                                              |
| ----------------------------------------------------------- | ------------------------------------------------------------ |
| [Directory Structure](/getting-started/directory-structure) | Where to put your code (entry point, routes, models, config) |
| [Application Bootstrap](/architecture/application)          | How `NewApplication` + `StandardProviders` brings up the app |
| [Service Container](/architecture/service-container)        | How to bind services and resolve them in handlers            |
| [Service Providers](/architecture/service-providers)        | How to write your own provider                               |
| [Drivers](/architecture/drivers)                            | How to swap cache/queue/log backends and add custom ones     |
| [Facades](/architecture/facades)                            | The shortcut layer over the container                        |
| [Configuration](/architecture/configuration)                | How options flow into the provider stack                     |
