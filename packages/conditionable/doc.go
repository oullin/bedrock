// Package conditionable provides Laravel-inspired conditional method execution.
// It defines generic When and Unless functions that conditionally apply callbacks
// based on the truthiness of a value, supporting both static values and resolver
// functions. A generic Proxy type offers conditional chaining similar to Laravel's
// HigherOrderWhenProxy, adapted for Go's static type system.
package conditionable
