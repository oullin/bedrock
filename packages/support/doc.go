// Package support provides Go ports of Laravel's Illuminate/Support utilities.
// It includes string helpers (Str*, StringBuilder), a dynamic key-value object (Fluent),
// a safe nullable wrapper (Optional[T]), an error message collection (MessageBag),
// global helpers (Blank, Filled, Tap, Value, With, Transform, E, Env, Retry), and
// utility types for probabilistic execution (Lottery), flexible sleeping (Sleep), and
// time-constrained execution (Timebox).
//
// This package is a 1:1 Go port of laravel/framework 13.x Illuminate/Support.
// PHP-specific features (Facades, Macroable, ServiceProvider, Carbon, Arr) are
// excluded or handled by other bedrock packages.
package support
