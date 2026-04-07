package hashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argonHash struct {
	algorithm string
	version   int
	memory    uint32
	time      uint32
	threads   uint8
	salt      []byte
	sum       []byte
}

type argonHasher struct {
	algorithm       string
	memory          uint32
	time            uint32
	threads         uint8
	verifyAlgorithm bool
	algorithmErr    error
}

// Argon2i hashes passwords with argon2i.
type Argon2i struct {
	argonHasher
}

// Argon2id hashes passwords with argon2id.
type Argon2id struct {
	argonHasher
}

const (
	argonVersion  = 19
	argonSaltSize = 16
	argonKeySize  = 32
)

var argonRandomReader io.Reader = rand.Reader

// NewArgon2i creates a new argon2i hasher.
func NewArgon2i(cfg ArgonConfig) *Argon2i {
	return &Argon2i{
		argonHasher: newArgonHasher("argon2i", cfg, errArgon2iAlgorithm),
	}
}

// NewArgon2id creates a new argon2id hasher.
func NewArgon2id(cfg ArgonConfig) *Argon2id {
	return &Argon2id{
		argonHasher: newArgonHasher("argon2id", cfg, errArgon2idAlgorithm),
	}
}

func newArgonHasher(algorithm string, cfg ArgonConfig, algorithmErr error) argonHasher {
	cfg = normalizeArgonConfig(cfg)

	return argonHasher{
		algorithm:       algorithm,
		memory:          uint32(cfg.Memory),
		time:            uint32(cfg.Time),
		threads:         uint8(cfg.Threads),
		verifyAlgorithm: cfg.Verify,
		algorithmErr:    algorithmErr,
	}
}

// Info returns the detected metadata for a hashed value.
func (a argonHasher) Info(hashedValue string) Info {
	return parseInfo(hashedValue)
}

// Make hashes a plaintext value.
func (a argonHasher) Make(value string, options map[string]any) (string, error) {
	params, err := a.params(options)

	if err != nil {
		return "", err
	}

	salt := make([]byte, argonSaltSize)

	if _, err := io.ReadFull(argonRandomReader, salt); err != nil {
		return "", fmt.Errorf("hashing: read salt: %w", err)
	}

	sum := a.derive([]byte(value), salt, params)

	return encodeArgonHash(argonHash{
		algorithm: a.algorithm,
		version:   argonVersion,
		memory:    params.memory,
		time:      params.time,
		threads:   params.threads,
		salt:      salt,
		sum:       sum,
	}), nil
}

// Check reports whether value matches the supplied hash.
func (a argonHasher) Check(value string, hashedValue string, _ map[string]any) (bool, error) {
	if hashedValue == "" {
		return false, nil
	}

	parsed, err := parseArgonHash(hashedValue)

	if err != nil {
		if a.verifyAlgorithm {
			return false, a.algorithmErr
		}

		return false, nil
	}

	if parsed.algorithm != a.algorithm {
		if a.verifyAlgorithm {
			return false, a.algorithmErr
		}

		return false, nil
	}

	sum := a.derive([]byte(value), parsed.salt, parsed)

	if subtle.ConstantTimeCompare(sum, parsed.sum) == 1 {
		return true, nil
	}

	return false, nil
}

// NeedsRehash reports whether the hash parameters differ from the configured ones.
func (a argonHasher) NeedsRehash(hashedValue string, options map[string]any) bool {
	parsed, err := parseArgonHash(hashedValue)

	if err != nil {
		return true
	}

	if parsed.algorithm != a.algorithm || parsed.version != argonVersion {
		return true
	}

	params, err := a.params(options)

	if err != nil {
		return true
	}

	return parsed.memory != params.memory || parsed.time != params.time || parsed.threads != params.threads
}

// VerifyConfiguration reports whether the hash parameters are within the configured limits.
func (a argonHasher) VerifyConfiguration(hashedValue string) bool {
	parsed, err := parseArgonHash(hashedValue)

	if err != nil {
		return false
	}

	if parsed.algorithm != a.algorithm || parsed.version != argonVersion {
		return false
	}

	return parsed.memory <= a.memory && parsed.time <= a.time && parsed.threads <= a.threads
}

func (a argonHasher) derive(value []byte, salt []byte, params argonHash) []byte {
	if a.algorithm == "argon2id" {
		return argon2.IDKey(value, salt, params.time, params.memory, params.threads, argonKeySize)
	}

	return argon2.Key(value, salt, params.time, params.memory, params.threads, argonKeySize)
}

func (a argonHasher) params(options map[string]any) (argonHash, error) {
	memory, err := intOption(options, "memory", int(a.memory))

	if err != nil {
		return argonHash{}, err
	}

	timeCost, err := intOption(options, "time", int(a.time))

	if err != nil {
		return argonHash{}, err
	}

	threads, err := intOption(options, "threads", int(a.threads))

	if err != nil {
		return argonHash{}, err
	}

	if memory <= 0 {
		return argonHash{}, fmt.Errorf("hashing: memory must be positive")
	}

	if timeCost <= 0 {
		return argonHash{}, fmt.Errorf("hashing: time must be positive")
	}

	if threads <= 0 {
		return argonHash{}, fmt.Errorf("hashing: threads must be positive")
	}

	if threads > 255 {
		return argonHash{}, fmt.Errorf("hashing: threads must be <= 255")
	}

	return argonHash{
		algorithm: a.algorithm,
		version:   argonVersion,
		memory:    uint32(memory),
		time:      uint32(timeCost),
		threads:   uint8(threads),
	}, nil
}

func encodeArgonHash(value argonHash) string {
	return fmt.Sprintf(
		"$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		value.algorithm,
		value.version,
		value.memory,
		value.time,
		value.threads,
		base64.RawStdEncoding.EncodeToString(value.salt),
		base64.RawStdEncoding.EncodeToString(value.sum),
	)
}

func parseArgonHash(hashedValue string) (argonHash, error) {
	parts := strings.Split(hashedValue, "$")

	if len(parts) != 6 || parts[0] != "" {
		return argonHash{}, fmt.Errorf("hashing: invalid argon hash")
	}

	algorithm := parts[1]

	if algorithm != "argon2i" && algorithm != "argon2id" {
		return argonHash{}, fmt.Errorf("hashing: invalid argon algorithm")
	}

	versionText, ok := strings.CutPrefix(parts[2], "v=")

	if !ok {
		return argonHash{}, fmt.Errorf("hashing: invalid argon version")
	}

	version, err := strconv.Atoi(versionText)

	if err != nil {
		return argonHash{}, fmt.Errorf("hashing: invalid argon version")
	}

	params, err := parseArgonParams(parts[3])

	if err != nil {
		return argonHash{}, err
	}

	salt, err := decodeArgonBase64(parts[4])

	if err != nil || len(salt) == 0 {
		return argonHash{}, fmt.Errorf("hashing: invalid argon salt")
	}

	sum, err := decodeArgonBase64(parts[5])

	if err != nil || len(sum) == 0 {
		return argonHash{}, fmt.Errorf("hashing: invalid argon sum")
	}

	params.algorithm = algorithm
	params.version = version
	params.salt = salt
	params.sum = sum

	return params, nil
}

func parseArgonParams(raw string) (argonHash, error) {
	result := argonHash{}
	fields := strings.Split(raw, ",")

	if len(fields) != 3 {
		return argonHash{}, fmt.Errorf("hashing: invalid argon parameters")
	}

	for _, field := range fields {
		key, value, ok := strings.Cut(field, "=")

		if !ok {
			return argonHash{}, fmt.Errorf("hashing: invalid argon parameter")
		}

		parsed, err := strconv.Atoi(value)

		if err != nil || parsed <= 0 {
			return argonHash{}, fmt.Errorf("hashing: invalid argon parameter")
		}

		switch key {
		case "m":
			result.memory = uint32(parsed)
		case "t":
			result.time = uint32(parsed)
		case "p":
			if parsed > 255 {
				return argonHash{}, fmt.Errorf("hashing: invalid argon threads")
			}

			result.threads = uint8(parsed)
		default:
			return argonHash{}, fmt.Errorf("hashing: invalid argon parameter")
		}
	}

	if result.memory == 0 || result.time == 0 || result.threads == 0 {
		return argonHash{}, fmt.Errorf("hashing: invalid argon parameters")
	}

	return result, nil
}

func decodeArgonBase64(value string) ([]byte, error) {
	decoded, err := base64.RawStdEncoding.DecodeString(value)

	if err == nil {
		return decoded, nil
	}

	return base64.StdEncoding.DecodeString(value)
}
