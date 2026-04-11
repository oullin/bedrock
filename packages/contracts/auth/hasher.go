package auth

import "context"

// PasswordHasher hashes and verifies passwords.
type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Check(ctx context.Context, password string, hashedPassword string) (bool, error)
	NeedsRehash(hashedPassword string) bool
}
