// Package auth provides HTTP authentication and authorization.
// It defines a Manager that creates named guards (session, token, request) and
// user providers (ORM, database). Access control is handled by the Gate in the
// access sub-package. Password resets are handled by the passwords sub-package.
package auth
