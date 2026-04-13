package hashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	contract "github.com/bedrock/packages/contracts/hashing"
	"golang.org/x/crypto/argon2"
)

const (
	argonDefaultMemory  uint32 = 1024
	argonDefaultTime    uint32 = 2
	argonDefaultThreads uint8  = 2
	argonSaltLen               = 16
	argonHashLen        uint32 = 32
)

type keyFunc func(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte

// ArgonHasher hashes values using the argon2i algorithm.
type ArgonHasher struct {
	memory    uint32
	time      uint32
	threads   uint8
	verify    bool
	algorithm string
	derive    keyFunc
}

var _ contract.Hasher = (*ArgonHasher)(nil)

// NewArgonHasher creates an ArgonHasher (argon2i). Recognised option keys:
// "memory" (uint32/int), "time" (uint32/int), "threads" (uint8/int), "verify" (bool).
func NewArgonHasher(opts ...map[string]any) *ArgonHasher {
	h := &ArgonHasher{
		memory:    argonDefaultMemory,
		time:      argonDefaultTime,
		threads:   argonDefaultThreads,
		verify:    false,
		algorithm: "argon2i",
		derive:    argon2.Key,
	}

	if len(opts) > 0 {
		applyArgonOpts(h, opts[0])
	}

	return h
}

// Info returns metadata about an argon2 hashed value.
func (h *ArgonHasher) Info(hashedValue string) (contract.HashInfo, error) {
	params, err := parseArgonHash(hashedValue)
	if err != nil {
		return contract.HashInfo{}, err
	}

	return contract.HashInfo{
		Algorithm: params.algorithm,
		Options: map[string]any{
			"memory":  params.memory,
			"time":    params.time,
			"threads": params.threads,
		},
	}, nil
}

// Make hashes a plaintext value.
func (h *ArgonHasher) Make(value string, options ...map[string]any) (string, error) {
	memory, time, threads := h.params(options...)

	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("hashing: %w", err)
	}

	hash := h.derive([]byte(value), salt, time, memory, threads, argonHashLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		h.algorithm, argon2.Version, memory, time, threads, b64Salt, b64Hash)

	return encoded, nil
}

// Check verifies a plaintext value against an argon2 hash.
func (h *ArgonHasher) Check(value string, hashedValue string, options ...map[string]any) (bool, error) {
	if len(hashedValue) == 0 {
		return false, nil
	}

	if h.verify && !h.isUsingCorrectAlgorithm(hashedValue) {
		return false, ErrAlgorithmMismatch
	}

	params, err := parseArgonHash(hashedValue)
	if err != nil {
		return false, nil
	}

	hash := h.derive([]byte(value), params.salt, params.time, params.memory, params.threads, uint32(len(params.hash)))

	if subtle.ConstantTimeCompare(hash, params.hash) != 1 {
		return false, nil
	}

	return true, nil
}

// NeedsRehash reports whether the hashed value was created with different options.
func (h *ArgonHasher) NeedsRehash(hashedValue string, options ...map[string]any) (bool, error) {
	params, err := parseArgonHash(hashedValue)
	if err != nil {
		return true, nil
	}

	memory, time, threads := h.params(options...)

	return params.memory != memory || params.time != time || params.threads != threads, nil
}

// VerifyConfiguration checks whether the hash was produced with acceptable settings.
func (h *ArgonHasher) VerifyConfiguration(hashedValue string) bool {
	params, err := parseArgonHash(hashedValue)
	if err != nil {
		return false
	}

	return params.memory <= h.memory && params.time <= h.time && params.threads <= h.threads
}

// SetMemory updates the memory cost.
func (h *ArgonHasher) SetMemory(memory uint32) { h.memory = memory }

// SetTime updates the time cost.
func (h *ArgonHasher) SetTime(time uint32) { h.time = time }

// SetThreads updates the parallelism.
func (h *ArgonHasher) SetThreads(threads uint8) { h.threads = threads }

// Memory returns the configured memory cost.
func (h *ArgonHasher) Memory() uint32 { return h.memory }

// Time returns the configured time cost.
func (h *ArgonHasher) Time() uint32 { return h.time }

// Threads returns the configured thread count.
func (h *ArgonHasher) Threads() uint8 { return h.threads }

func (h *ArgonHasher) isUsingCorrectAlgorithm(hashedValue string) bool {
	return strings.HasPrefix(hashedValue, "$"+h.algorithm+"$")
}

func (h *ArgonHasher) params(options ...map[string]any) (uint32, uint32, uint8) {
	memory, time, threads := h.memory, h.time, h.threads

	if len(options) > 0 {
		if v, ok := toUint32(options[0]["memory"]); ok {
			memory = v
		}
		if v, ok := toUint32(options[0]["time"]); ok {
			time = v
		}
		if v, ok := toUint8(options[0]["threads"]); ok {
			threads = v
		}
	}

	return memory, time, threads
}

type argonParams struct {
	algorithm string
	memory    uint32
	time      uint32
	threads   uint8
	salt      []byte
	hash      []byte
}

func parseArgonHash(encoded string) (argonParams, error) {
	// Expected format: $argon2i$v=19$m=M,t=T,p=P$salt$hash
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return argonParams{}, ErrInvalidHash
	}

	algo := parts[1]
	if algo != "argon2i" && algo != "argon2id" {
		return argonParams{}, ErrInvalidHash
	}

	var memory, time uint32
	var threads uint8

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return argonParams{}, ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return argonParams{}, ErrInvalidHash
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return argonParams{}, ErrInvalidHash
	}

	return argonParams{
		algorithm: algo,
		memory:    memory,
		time:      time,
		threads:   threads,
		salt:      salt,
		hash:      hash,
	}, nil
}

func applyArgonOpts(h *ArgonHasher, opts map[string]any) {
	if v, ok := toUint32(opts["memory"]); ok {
		h.memory = v
	}
	if v, ok := toUint32(opts["time"]); ok {
		h.time = v
	}
	if v, ok := toUint8(opts["threads"]); ok {
		h.threads = v
	}
	if v, ok := opts["verify"].(bool); ok {
		h.verify = v
	}
}

func toUint32(v any) (uint32, bool) {
	switch val := v.(type) {
	case uint32:
		return val, true
	case int:
		return uint32(val), true
	case int64:
		return uint32(val), true
	default:
		return 0, false
	}
}

func toUint8(v any) (uint8, bool) {
	switch val := v.(type) {
	case uint8:
		return val, true
	case int:
		return uint8(val), true
	case int64:
		return uint8(val), true
	default:
		return 0, false
	}
}
