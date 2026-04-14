# Getting Started

Bedrock is a collection of foundational Go packages for building web applications.
Each package is a standalone Go module with its own `go.mod`, so you import only
what you need.

## Requirements

| Tool    | Version |
|---------|---------|
| Go      | ≥ 1.24  |
| Node.js | ≥ 22    |
| pnpm    | ≥ 10.33 |

## Installation

Add any Bedrock package directly with `go get`:

```bash
go get github.com/gocanto/bedrock/packages/auth@latest
go get github.com/gocanto/bedrock/packages/cache@latest
```

## Project Layout

```
packages/   Go library packages
services/   Internal tooling (storage, scripts, docs)
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

| Package                             | Purpose                                             |
|-------------------------------------|-----------------------------------------------------|
| [auth](/packages/auth)              | Authentication, authorization, password management  |
| [bus](/packages/bus)                | Command and event bus with pipeline support         |
| [cache](/packages/cache)            | Caching layer with multiple driver support          |
| [cookie](/packages/cookie)          | HTTP cookie handling                                |
| [fortify](/packages/fortify)        | Rate limiting, two-factor auth, pipelines           |
| [httpx](/packages/httpx)            | HTTP utilities, middleware, and testing helpers     |
| [jetstream](/packages/jetstream)    | Team and organization management                    |
| [queue](/packages/queue)            | Background job processing with pluggable drivers    |
| [routing](/packages/routing)        | HTTP routing (1:1 Laravel port)                     |
| [session](/packages/session)        | Session management with multiple storage handlers   |
| [spark](/packages/spark)            | Subscription billing, checkout, and entitlements    |
