# prompts

Beautiful, user-friendly terminal UI forms for Go applications.

## Overview

The `prompts` package provides interactive terminal UI components with full
functional parity with Upstream Prompts. All prompts are rendered with ANSI
escape codes, support validation and transformation, and include non-interactive
fallbacks. Terminal I/O is fully injectable for testing via `Fake`.

**Module:** `github.com/bedrock/packages/prompts`

```bash
go get github.com/bedrock/packages/prompts@latest
```

## Available Components

| Component    | Description                             |
| ------------ | --------------------------------------- |
| Text input   | Single-line text prompt with validation |
| Password     | Masked password input                   |
| Select       | Single-choice selection list            |
| Multi-select | Multiple-choice selection list          |
| Search       | Searchable select with live filtering   |
| Spinner      | Progress indicator for async operations |
| Progress bar | Determinate progress display            |
| Table        | Formatted tabular data output           |
| Form builder | Multi-step form with state across steps |

## Coming Soon

Full documentation is in progress.
