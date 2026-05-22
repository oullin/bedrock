// Package session provides HTTP session management.
// It defines a Store with flash data, CSRF tokens, and lifecycle management,
// backed by swappable Handler implementations (array, file, database, cache,
// cookie, null, and encrypting).
package session
