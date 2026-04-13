// Package container provides a Laravel-inspired inversion-of-control (IoC)
// container for Go. It supports binding factories and instances by string key,
// singleton and scoped lifecycles, contextual bindings, tagging, extension
// callbacks, lifecycle hooks (before/resolving/after), method invocation with
// dependency injection, and rebinding notifications. The container is safe for
// concurrent use.
package container
