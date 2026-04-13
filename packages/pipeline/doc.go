// Package pipeline provides a middleware-style processing chain.
// It allows sending a value through a series of pipes, where each
// pipe can inspect, transform, or short-circuit the chain.
//
// The package also provides a Hub for managing named pipeline
// configurations.
package pipeline
