# filesystem

<!-- upstream-docs: filesystem.md#introduction -->
<!-- upstream-docs: filesystem.md#configuration -->
<!-- upstream-docs: filesystem.md#obtaining-disk-instances -->
<!-- upstream-docs: filesystem.md#retrieving-files -->
<!-- upstream-docs: filesystem.md#storing-files -->
<!-- upstream-docs: filesystem.md#deleting-files -->
<!-- upstream-docs: filesystem.md#directories -->

<!-- BEDROCK:HAND -->

## Introduction

The filesystem package gives every Bedrock app a single, type-safe API
for files and directories. Today it ships a local-disk driver; the same
interface accepts cloud-storage adapters when you bring them.

For the cross-cutting picture, see [Drivers](/architecture/drivers).

## Configuration

The filesystem service is bound by `FilesystemServiceProvider`. The
constructor takes no config — the local driver derives its root from the
caller, and you mount additional disks at runtime:

```go
// services/demo/api/bootstrap.go:140
filesystem.NewFilesystemServiceProvider(application.Container),
```

See [`packages/filesystem/filesystem_service_provider.go`](https://github.com/gocanto/bedrock/blob/main/packages/filesystem/filesystem_service_provider.go).

## Basic Usage

```go
fs := container.Resolve[*filesystem.Filesystem]("filesystem")

// Write
err := fs.Put("uploads/avatar-42.png", data)

// Read
raw, err := fs.Get("uploads/avatar-42.png")

// Stream
reader, err := fs.ReadStream("uploads/avatar-42.png")
defer reader.Close()

// Move/Copy/Delete
err = fs.Move("uploads/avatar-42.png", "archive/avatar-42.png")
err = fs.Delete("archive/avatar-42.png")
```

For directory operations, use `MakeDirectory`, `Files`, `Directories`,
and `DeleteDirectory`. See
[`packages/filesystem/filesystem_dir.go`](https://github.com/gocanto/bedrock/blob/main/packages/filesystem/filesystem_dir.go).

## Drivers

Today only the **`local`** driver ships; it covers POSIX-style file
operations with optional advisory locking
([`filesystem.go`](https://github.com/gocanto/bedrock/blob/main/packages/filesystem/filesystem.go)).

Cloud drivers (S3, GCS, Azure Blob) are not in the box; see _Writing
Custom Drivers_ below.

## Writing Custom Drivers

Implement the `Filesystem` interface (see
[`packages/filesystem/filesystem.go`](https://github.com/gocanto/bedrock/blob/main/packages/filesystem/filesystem.go))
and register your driver as a separate "disk" in your application's
bootstrap:

```go
type s3Disk struct { /* ... */ }

func (d *s3Disk) Get(path string) ([]byte, error)              { /* ... */ }
func (d *s3Disk) Put(path string, contents []byte) error       { /* ... */ }
// ... rest of the Filesystem interface

application.Container.Instance("filesystem.s3", newS3Disk(cfg))
```

Then resolve it under that name when handlers need cloud storage:

```go
disk := container.Resolve[filesystem.Filesystem]("filesystem.s3")
disk.Put("backups/2026-01-01.tar.gz", data)
```

## Events

This package does not currently dispatch events. Wrap operations in
your own event-emitting decorator if you need observability.

## See Also

- [Drivers](/architecture/drivers).
- [Service Providers](/architecture/service-providers).
<!-- /BEDROCK:HAND -->

Package filesystem provides local filesystem operations including reading, writing, copying, moving, and deleting files and directories. It also supports file locking, MIME type detection, hashing, and permission management.

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
go get github.com/bedrock/packages/filesystem@latest
```

When working inside this monorepo, use the repository workspace:

```bash
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/filesystem/...
```

## Source Coverage

| Package      | Purpose                                                                                                                                                                                                                          |
| ------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `filesystem` | Package filesystem provides local filesystem operations including reading, writing, copying, moving, and deleting files and directories. It also supports file locking, MIME type detection, hashing, and permission management. |

## Core Concepts

The filesystem reference is organized around the exported Go surface for package `filesystem`. Start from the source coverage and public surface tables to identify the constructors, managers, interfaces, sentinel errors, and helper functions available to callers. Use the package tests as executable wiring examples for collaborators, default behavior, and Upstream parity expectations.

### Public Surface

| Surface                    | Exported API                                                                                                                                                                                                                                                                              |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Types                      | `Filesystem`, `FilesystemServiceProvider`, `LockableFile`                                                                                                                                                                                                                                 |
| Constructors and functions | `AllDirectories`, `AllFiles`, `Append`, `Basename`, `Chmod`, `CleanDirectory`, `Close`, `Copy`, `CopyDirectory`, `Delete`, `DeleteDirectories`, `DeleteDirectory`, `Directories`, `Dirname`, `EnsureDirectoryExists`, `ExclusiveLock`, `Exists`, `Extension`, `Files`, `Get`, and 38 more |
| Variables                  | `ErrHashAlgorithm`, `ErrLockFailed`, `ErrNotDirectory`, `ErrNotFile`, `ErrNotFound`                                                                                                                                                                                                       |
| Constants                  | None exported from this package root.                                                                                                                                                                                                                                                     |

### Capability Matrix

| Capability                         | Documentation note                                                                                                   |
| ---------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Security-sensitive behavior        | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |
| Serialization or transport formats | Supported by exported API and package tests; use the API reference and parity tests below when wiring this behavior. |

## Usage

Start with the package constructor or manager type when one is exported. Bedrock keeps dependencies explicit, so callers should pass repositories, stores, handlers, dispatchers, clocks, or clients directly instead of relying on global framework state.

```go
package main

import (
    _ "github.com/bedrock/packages/filesystem"
)

func main() {
    // Import the package you use, then wire the exported constructors,
    // managers, stores, handlers, or helpers required by your application.
}
```

Use package tests as executable examples when the exact constructor requires collaborators. The tests under `packages/filesystem` cover the supported creation paths, default values, and Upstream parity behavior.

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
GOWORK=./storage/.cache/go.work go test -count=1 ./packages/filesystem/...
```

Upstream parity is tracked by these tests:

- `packages/filesystem/laravel_inventory_test.go`

## API Reference

### Exported Types

| Type                        | Notes                                                                              |
| --------------------------- | ---------------------------------------------------------------------------------- |
| `Filesystem`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `FilesystemServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LockableFile`              | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Functions

| Function                       | Notes                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------- |
| `AllDirectories`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `AllFiles`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Append`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Basename`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Chmod`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CleanDirectory`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Close`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Copy`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `CopyDirectory`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Delete`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteDirectories`            | Source-backed public surface. See the Go package for exact signature and behavior. |
| `DeleteDirectory`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Directories`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Dirname`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `EnsureDirectoryExists`        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ExclusiveLock`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Exists`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Extension`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Files`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Get`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Glob`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `GuessExtension`               | Source-backed public surface. See the Go package for exact signature and behavior. |
| `HasSameHash`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Hash`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsDirectory`                  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsEmptyDirectory`             | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsFile`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsReadable`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `IsWritable`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `JSON`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `LastModified`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Lines`                        | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Link`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MakeDirectory`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MimeType`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Missing`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Move`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `MoveDirectory`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Name`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `New`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewFilesystemServiceProvider` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `NewLockableFile`              | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Path`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Prepend`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Provides`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Put`                          | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Read`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Register`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `RelativeLink`                 | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Replace`                      | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ReplaceInFile`                | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SharedGet`                    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `SharedLock`                   | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Size`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Truncate`                     | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Type`                         | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Unlock`                       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `Write`                        | Source-backed public surface. See the Go package for exact signature and behavior. |

### Exported Errors, Variables, and Constants

| Name               | Notes                                                                              |
| ------------------ | ---------------------------------------------------------------------------------- |
| `ErrHashAlgorithm` | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrLockFailed`    | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNotDirectory`  | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNotFile`       | Source-backed public surface. See the Go package for exact signature and behavior. |
| `ErrNotFound`      | Source-backed public surface. See the Go package for exact signature and behavior. |

## Upstream Parity Notes

This page should stay aligned with the official Upstream 13.x documentation for the corresponding feature while keeping the Go API explicit. If Bedrock implements a Upstream feature, document the user-facing behavior, the Go entry points, supported drivers, emitted events, error behavior, and the tests that prove parity. If a Upstream feature is PHP-only, record the exclusion in `services/compliance/docs-status.yml` instead of inventing a Go API.
