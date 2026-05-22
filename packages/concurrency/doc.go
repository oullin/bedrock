// Package concurrency provides concurrent task execution.
// It defines a Driver interface with multiple implementations: GoroutineDriver
// for true parallel execution via goroutines, and SyncDriver for sequential
// execution useful in testing. A Manager handles named driver instances with
// lazy initialization and thread-safe access.
package concurrency
