// Package conditionable provides conditional method execution.
// It defines generic When and Unless functions that conditionally apply callbacks
// based on the truthiness of a value, supporting both static values and resolver
// functions. A generic Proxy type offers conditional chaining, adapted for
// Go's static type system.
package conditionable
