// Package support provides Go ports of the upstream @bedrock/Support utilities.
// It includes array and dot-notation helpers (Arr*, Dot* internals), a dynamic
// key-value object (Fluent), a safe nullable wrapper (Optional[T]), an error
// message collection (MessageBag), global helpers (Blank, Filled, Tap, Value,
// With, Transform, E, Env, Retry), flexible sleeping (Sleep), and
// time-constrained execution (Timebox).
//
// This package is a 1:1 Go port of upstream framework 13.x @bedrock/Support.
// String utilities and Lottery were split into sibling packages:
// github.com/bedrock/packages/str and github.com/bedrock/packages/lottery.
// PHP-specific features (Facades, Macroable, ServiceProvider, Carbon) are
// excluded or handled by other Bedrock packages.
package support
