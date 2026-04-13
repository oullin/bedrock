package hashing

// HashInfo contains metadata about a hashed value.
type HashInfo struct {
	Algorithm string
	Options   map[string]any
}

// Hasher hashes and verifies values.
type Hasher interface {
	Info(hashedValue string) (HashInfo, error)
	Make(value string, options ...map[string]any) (string, error)
	Check(value string, hashedValue string, options ...map[string]any) (bool, error)
	NeedsRehash(hashedValue string, options ...map[string]any) (bool, error)
}
