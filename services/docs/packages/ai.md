# ai

Unified API for interacting with AI providers (OpenAI, Anthropic, Gemini, and more).

## Overview

The `ai` package provides a single, expressive interface for multiple AI providers
following the manager + provider + gateway pattern used throughout Bedrock. It
mirrors the Laravel AI package with full functional parity adapted to Go.

**Module:** `github.com/bedrock/packages/ai`

```bash
go get github.com/bedrock/packages/ai@latest
```

## Supported Providers

| Provider  | Text | Image | Audio |
| --------- | ---- | ----- | ----- |
| OpenAI    | ✓    | ✓     | ✓     |
| Anthropic | ✓    |       |       |
| Gemini    | ✓    | ✓     |       |

## Architecture

| Layer      | Responsibility                                          |
| ---------- | ------------------------------------------------------- |
| `Manager`  | Resolves and caches provider instances by `Lab` enum   |
| `Provider` | Implements capability interfaces (Text, Image, Audio…) |
| `Gateway`  | Thin HTTP/SDK adapter used by the provider             |

## Coming Soon

Full documentation is in progress.
