# filesystem

<!-- upstream-docs: filesystem.md#file-storage -->
<!-- upstream-docs: filesystem.md#retrieving-files -->
<!-- upstream-docs: filesystem.md#storing-files -->
<!-- upstream-docs: filesystem.md#deleting-files -->
<!-- upstream-docs: filesystem.md#directories -->
<!-- upstream-docs: filesystem.md#custom-filesystems -->

Local filesystem operations.

## Overview

The `filesystem` package provides a `Filesystem` type that wraps Go's `os` and
`io/fs` packages with a Upstream-inspired interface. It covers reads, writes, copy,
move, delete, directory management, MIME detection, file hashing, and locked
file access.

**Module:** `github.com/gocanto/bedrock/packages/filesystem`

```bash
go get github.com/gocanto/bedrock/packages/filesystem@latest
```

## Creating a Filesystem

```go
fs := filesystem.New()
```

## Existence Checks

```go
fs.Exists("/path/to/file")   // true if exists
fs.Missing("/path/to/file")  // true if not exists
fs.IsFile("/path/to/file")   // true if regular file
fs.IsDirectory("/path/to/dir") // true if directory
```

## Reading

```go
content, err := fs.Get("/path/to/file")       // string
bytes,   err := fs.GetBytes("/path/to/file")  // []byte
lines,   err := fs.Lines("/path/to/file")     // []string
size,    err := fs.Size("/path/to/file")       // int64 (bytes)
modTime, err := fs.LastModified("/path/to/file") // time.Time
```

## Writing

```go
err := fs.Put("/path/to/file", "contents")
err  = fs.Append("/path/to/file", "\nmore")
err  = fs.Prepend("/path/to/file", "header\n")
err  = fs.PutBytes("/path/to/file", []byte{...})
```

## Permissions

```go
err := fs.Chmod("/path/to/file", 0o644)
```

## Copy, Move & Delete

```go
err := fs.Copy("/src", "/dst")
err  = fs.Move("/old", "/new")
err  = fs.Delete("/path/to/file")
err  = fs.DeleteDirectory("/path/to/dir")
```

## Directory Operations

```go
err := fs.MakeDirectory("/path/to/dir", 0o755, true) // recursive
files, err := fs.Files("/path/to/dir")               // []string
all,   err := fs.AllFiles("/path/to/dir")            // recursive
dirs,  err := fs.Directories("/path/to/dir")
```

## MIME & Hashing

```go
mimeType, err := fs.MimeType("/path/to/file") // "image/png"

hash, err := fs.Hash("/path/to/file")         // SHA256 hex string
hash, err  = fs.HashWithAlgorithm("/path/to/file", "md5")
```

## Locked File Access

Serialize concurrent reads and writes using advisory file locks:

```go
lf, err := filesystem.NewLockableFile("/path/to/file", os.O_RDWR|os.O_CREATE, 0o644)
defer lf.Close()

lf.GetSharedLock()   // shared (read) lock
lf.GetExclusiveLock() // exclusive (write) lock
lf.ReleaseLock()
```
