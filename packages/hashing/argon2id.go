package hashing

import (
	contract "github.com/bedrock/packages/contracts/hashing"
	"golang.org/x/crypto/argon2"
)

// Argon2IdHasher hashes values using the argon2id algorithm.
type Argon2IdHasher struct {
	ArgonHasher
}

var _ contract.Hasher = (*Argon2IdHasher)(nil)

// NewArgon2IdHasher creates an Argon2IdHasher. Recognised option keys are
// the same as NewArgonHasher: "memory", "time", "threads", "verify".
func NewArgon2IdHasher(opts ...map[string]any) *Argon2IdHasher {
	h := &Argon2IdHasher{
		ArgonHasher: ArgonHasher{
			memory:    argonDefaultMemory,
			time:      argonDefaultTime,
			threads:   argonDefaultThreads,
			verify:    false,
			algorithm: "argon2id",
			derive:    argon2.IDKey,
		},
	}

	if len(opts) > 0 {
		applyArgonOpts(&h.ArgonHasher, opts[0])
	}

	return h
}
