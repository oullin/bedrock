// Package config provides a configuration repository backed
// by Viper. It stores key-value pairs in a nested map with dot-notation
// access, type-safe getters, and array manipulation helpers (prepend and
// push). Consumers get YAML file and environment variable support out of the
// box through the underlying Viper instance.
package config
