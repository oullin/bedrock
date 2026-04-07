package hashing

import (
	"encoding/base64"
	"io"
	"math"
	"strings"
	"testing"

	configpkg "github.com/bedrock/packages/config"
)

type errReader struct {
	err error
}

type bytesReader []byte

func TestBcryptHasher(t *testing.T) {
	hasher := NewBcrypt(BcryptConfig{Rounds: 10})

	hashedValue, err := hasher.Make("secret", nil)

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	if !strings.HasPrefix(hashedValue, "$2") {
		t.Fatalf("unexpected bcrypt prefix: %q", hashedValue)
	}

	info := hasher.Info(hashedValue)

	if got, want := info.Algorithm, DriverBcrypt; got != want {
		t.Fatalf("unexpected algorithm: got %q want %q", got, want)
	}

	if got, want := info.Options["rounds"], 10; got != want {
		t.Fatalf("unexpected rounds: got %d want %d", got, want)
	}

	matched, err := hasher.Check("secret", hashedValue, nil)

	if err != nil || !matched {
		t.Fatalf("expected bcrypt check to match, matched=%v err=%v", matched, err)
	}

	matched, err = hasher.Check("wrong", hashedValue, nil)

	if err != nil || matched {
		t.Fatalf("expected bcrypt mismatch, matched=%v err=%v", matched, err)
	}

	matched, err = hasher.Check("secret", "", nil)

	if err != nil || matched {
		t.Fatalf("expected empty hash check to fail, matched=%v err=%v", matched, err)
	}

	matched, err = hasher.Check("secret", "not-a-hash", nil)

	if err != nil || matched {
		t.Fatalf("expected malformed hash check to fail, matched=%v err=%v", matched, err)
	}

	if hasher.NeedsRehash(hashedValue, nil) {
		t.Fatal("expected bcrypt hash not to need rehash")
	}

	if !hasher.NeedsRehash(hashedValue, map[string]any{"rounds": 12}) {
		t.Fatal("expected bcrypt hash to need rehash for different cost")
	}

	if !hasher.NeedsRehash("not-a-hash", nil) {
		t.Fatal("expected invalid bcrypt hash to need rehash")
	}

	if !hasher.NeedsRehash(hashedValue, map[string]any{"rounds": []string{"bad"}}) {
		t.Fatal("expected invalid options to force bcrypt rehash")
	}

	if !hasher.VerifyConfiguration(hashedValue) {
		t.Fatal("expected bcrypt configuration to verify")
	}

	if NewBcrypt(BcryptConfig{Rounds: 4}).VerifyConfiguration(hashedValue) {
		t.Fatal("expected bcrypt configuration to reject higher cost hash")
	}

	if NewBcrypt(BcryptConfig{}).VerifyConfiguration("not-a-hash") {
		t.Fatal("expected bcrypt configuration to reject invalid hash")
	}

	if _, err := NewBcrypt(BcryptConfig{Rounds: 10, Limit: 4}).Make("secret", nil); err == nil {
		t.Fatal("expected bcrypt limit error")
	}

	if _, err := hasher.Make("secret", map[string]any{"rounds": []string{"bad"}}); err == nil {
		t.Fatal("expected bcrypt option type error")
	}

	if _, err := hasher.Make("secret", map[string]any{"rounds": 100}); err == nil {
		t.Fatal("expected bcrypt invalid cost error")
	}

	argonHash, err := NewArgon2id(ArgonConfig{}).Make("secret", nil)

	if err != nil {
		t.Fatalf("Make argon hash: %v", err)
	}

	matched, err = NewBcrypt(BcryptConfig{}).Check("secret", argonHash, nil)

	if err != nil || matched {
		t.Fatalf("expected bcrypt verify-disabled mismatch, matched=%v err=%v", matched, err)
	}

	if _, err := NewBcrypt(BcryptConfig{Verify: true}).Check("secret", argonHash, nil); err != errBcryptAlgorithm {
		t.Fatalf("unexpected bcrypt verify error: %v", err)
	}
}

func TestArgonHashers(t *testing.T) {
	restoreRandom := stubArgonRandomReader(t, bytesReader(make([]byte, 64)))

	defer restoreRandom()

	argon2i := NewArgon2i(ArgonConfig{Memory: 2048, Time: 3, Threads: 2})
	hashI, err := argon2i.Make("secret", nil)

	if err != nil {
		t.Fatalf("Argon2i Make: %v", err)
	}

	if !strings.HasPrefix(hashI, "$argon2i$") {
		t.Fatalf("unexpected argon2i prefix: %q", hashI)
	}

	info := argon2i.Info(hashI)

	if got, want := info.Algorithm, DriverArgon; got != want {
		t.Fatalf("unexpected argon2i algorithm: got %q want %q", got, want)
	}

	if info.Options["memory"] != 2048 || info.Options["time"] != 3 || info.Options["threads"] != 2 {
		t.Fatalf("unexpected argon2i options: %#v", info.Options)
	}

	matched, err := argon2i.Check("secret", hashI, nil)

	if err != nil || !matched {
		t.Fatalf("expected argon2i match, matched=%v err=%v", matched, err)
	}

	matched, err = argon2i.Check("wrong", hashI, nil)

	if err != nil || matched {
		t.Fatalf("expected argon2i mismatch, matched=%v err=%v", matched, err)
	}

	matched, err = argon2i.Check("secret", "", nil)

	if err != nil || matched {
		t.Fatalf("expected empty argon2i hash to fail, matched=%v err=%v", matched, err)
	}

	matched, err = argon2i.Check("secret", "not-a-hash", nil)

	if err != nil || matched {
		t.Fatalf("expected malformed argon2i hash to fail, matched=%v err=%v", matched, err)
	}

	if _, err := NewArgon2i(ArgonConfig{Verify: true}).Check("secret", "not-a-hash", nil); err != errArgon2iAlgorithm {
		t.Fatalf("unexpected argon2i malformed verify error: %v", err)
	}

	argon2id := NewArgon2id(ArgonConfig{Memory: 1024, Time: 2, Threads: 2})
	hashID, err := argon2id.Make("secret", nil)

	if err != nil {
		t.Fatalf("Argon2id Make: %v", err)
	}

	if !strings.HasPrefix(hashID, "$argon2id$") {
		t.Fatalf("unexpected argon2id prefix: %q", hashID)
	}

	matched, err = argon2id.Check("secret", hashID, nil)

	if err != nil || !matched {
		t.Fatalf("expected argon2id match, matched=%v err=%v", matched, err)
	}

	matched, err = argon2id.Check("secret", "", nil)

	if err != nil || matched {
		t.Fatalf("expected empty argon2id hash to fail, matched=%v err=%v", matched, err)
	}

	matched, err = argon2id.Check("secret", hashI, nil)

	if err != nil || matched {
		t.Fatalf("expected verify-disabled argon2id mismatch, matched=%v err=%v", matched, err)
	}

	if _, err := NewArgon2id(ArgonConfig{Verify: true}).Check("secret", hashI, nil); err != errArgon2idAlgorithm {
		t.Fatalf("unexpected argon2id verify error: %v", err)
	}

	if argon2i.NeedsRehash(hashI, nil) {
		t.Fatal("expected argon2i hash not to need rehash")
	}

	if !argon2i.NeedsRehash(hashI, map[string]any{"memory": 4096}) {
		t.Fatal("expected argon2i hash to need rehash for different memory")
	}

	if !argon2i.NeedsRehash("not-a-hash", nil) {
		t.Fatal("expected invalid argon2i hash to need rehash")
	}

	if !argon2i.NeedsRehash(hashI, map[string]any{"memory": []string{"bad"}}) {
		t.Fatal("expected invalid argon2i options to force rehash")
	}

	versionMismatch := strings.Replace(hashI, "v=19", "v=18", 1)

	if !argon2i.NeedsRehash(versionMismatch, nil) {
		t.Fatal("expected argon2i version mismatch to need rehash")
	}

	if !argon2i.NeedsRehash(hashID, nil) {
		t.Fatal("expected argon2i algorithm mismatch to need rehash")
	}

	if !argon2i.VerifyConfiguration(hashI) {
		t.Fatal("expected argon2i configuration to verify")
	}

	if NewArgon2i(ArgonConfig{Memory: 1024, Time: 2, Threads: 1}).VerifyConfiguration(hashI) {
		t.Fatal("expected argon2i configuration to reject larger parameters")
	}

	if NewArgon2i(ArgonConfig{}).VerifyConfiguration("not-a-hash") {
		t.Fatal("expected invalid argon2i hash to fail verification")
	}

	if NewArgon2i(ArgonConfig{}).VerifyConfiguration(versionMismatch) {
		t.Fatal("expected argon2i version mismatch to fail verification")
	}

	if NewArgon2i(ArgonConfig{}).VerifyConfiguration(hashID) {
		t.Fatal("expected argon2i algorithm mismatch to fail verification")
	}

	if _, err := argon2i.Make("secret", map[string]any{"memory": []string{"bad"}}); err == nil {
		t.Fatal("expected argon memory option type error")
	}

	if _, err := argon2i.Make("secret", map[string]any{"time": []string{"bad"}}); err == nil {
		t.Fatal("expected argon time option type error")
	}

	if _, err := argon2i.Make("secret", map[string]any{"threads": []string{"bad"}}); err == nil {
		t.Fatal("expected argon threads option type error")
	}

	if _, err := argon2i.Make("secret", map[string]any{"memory": 0}); err == nil {
		t.Fatal("expected argon memory validation error")
	}

	if _, err := argon2i.Make("secret", map[string]any{"time": 0}); err == nil {
		t.Fatal("expected argon time validation error")
	}

	if _, err := argon2i.Make("secret", map[string]any{"threads": 0}); err == nil {
		t.Fatal("expected argon threads validation error")
	}

	if _, err := argon2i.Make("secret", map[string]any{"threads": 256}); err == nil {
		t.Fatal("expected argon threads >255 validation error")
	}

	restoreRandom = stubArgonRandomReader(t, errReader{err: io.ErrUnexpectedEOF})

	defer restoreRandom()

	if _, err := NewArgon2i(ArgonConfig{}).Make("secret", nil); err == nil || !strings.Contains(err.Error(), "hashing: read salt") {
		t.Fatalf("expected argon salt read error, got %v", err)
	}
}

func TestArgonParsingHelpers(t *testing.T) {
	restoreRandom := stubArgonRandomReader(t, bytesReader(make([]byte, 64)))

	defer restoreRandom()

	hashI, err := NewArgon2i(ArgonConfig{}).Make("secret", nil)

	if err != nil {
		t.Fatalf("Make argon2i hash: %v", err)
	}

	hashID, err := NewArgon2id(ArgonConfig{}).Make("secret", nil)

	if err != nil {
		t.Fatalf("Make argon2id hash: %v", err)
	}

	if parsed, err := parseArgonHash(hashI); err != nil || parsed.algorithm != "argon2i" {
		t.Fatalf("unexpected parseArgonHash argon2i result: parsed=%#v err=%v", parsed, err)
	}

	if parsed, err := parseArgonHash(hashID); err != nil || parsed.algorithm != "argon2id" {
		t.Fatalf("unexpected parseArgonHash argon2id result: parsed=%#v err=%v", parsed, err)
	}

	for _, value := range []string{
		"not-a-hash",
		"$argon2x$v=19$m=1024,t=2,p=2$c2FsdA$c3Vt",
		"$argon2i$broken$m=1024,t=2,p=2$c2FsdA$c3Vt",
		"$argon2i$v=nope$m=1024,t=2,p=2$c2FsdA$c3Vt",
		"$argon2i$v=19$m=nope,t=2,p=2$c2FsdA$c3Vt",
		"$argon2i$v=19$m=1024,t=2,p=2$$c3Vt",
		"$argon2i$v=19$m=1024,t=2,p=2$c2FsdA$",
	} {
		if _, err := parseArgonHash(value); err == nil {
			t.Fatalf("expected parseArgonHash failure for %q", value)
		}
	}

	for _, raw := range []string{
		"m=1024,t=2",
		"m=1024,t=2,p=2,extra=1",
		"m=1024,t=2,p",
		"m=nope,t=2,p=2",
		"m=1024,t=0,p=2",
		"m=1024,t=2,p=256",
		"x=1,t=2,p=2",
		"m=1024,t=2,p=0",
		"m=1,m=2,p=3",
	} {
		if _, err := parseArgonParams(raw); err == nil {
			t.Fatalf("expected parseArgonParams failure for %q", raw)
		}
	}

	if parsed, err := parseArgonParams("m=1024,t=2,p=2"); err != nil || parsed.memory != 1024 || parsed.time != 2 || parsed.threads != 2 {
		t.Fatalf("unexpected parseArgonParams result: parsed=%#v err=%v", parsed, err)
	}

	rawEncoded := base64.RawStdEncoding.EncodeToString([]byte("salt"))

	if decoded, err := decodeArgonBase64(rawEncoded); err != nil || string(decoded) != "salt" {
		t.Fatalf("unexpected raw base64 decode: %q err=%v", string(decoded), err)
	}

	stdEncoded := base64.StdEncoding.EncodeToString([]byte("salt"))

	if decoded, err := decodeArgonBase64(stdEncoded); err != nil || string(decoded) != "salt" {
		t.Fatalf("unexpected std base64 decode: %q err=%v", string(decoded), err)
	}

	if _, err := decodeArgonBase64("%%%invalid%%%"); err == nil {
		t.Fatal("expected invalid base64 decode error")
	}
}

func TestManagerAndHelpers(t *testing.T) {
	manager, err := NewManager(Config{})

	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	if got, want := manager.DefaultDriver(), DriverBcrypt; got != want {
		t.Fatalf("unexpected default driver: got %q want %q", got, want)
	}

	if _, err := manager.Driver(DriverArgon); err != nil {
		t.Fatalf("Driver(argon): %v", err)
	}

	if _, err := manager.Driver(DriverArgon2id); err != nil {
		t.Fatalf("Driver(argon2id): %v", err)
	}

	if _, err := manager.Driver("nope"); err == nil {
		t.Fatal("expected unsupported driver error")
	}

	hashedValue, err := manager.Make("secret", nil)

	if err != nil {
		t.Fatalf("Make: %v", err)
	}

	if !manager.IsHashed(hashedValue) {
		t.Fatal("expected manager to detect supported hash")
	}

	if manager.IsHashed("plain-text") {
		t.Fatal("expected plain text not to be detected as hash")
	}

	if info := manager.Info("plain-text"); info.Algorithm != "" || len(info.Options) != 0 {
		t.Fatalf("unexpected plain-text info: %#v", info)
	}

	matched, err := manager.Check("secret", hashedValue, nil)

	if err != nil || !matched {
		t.Fatalf("expected manager bcrypt match, matched=%v err=%v", matched, err)
	}

	if manager.NeedsRehash(hashedValue, nil) {
		t.Fatal("expected manager bcrypt hash not to need rehash")
	}

	if !manager.VerifyConfiguration(hashedValue) {
		t.Fatal("expected manager bcrypt configuration to verify")
	}

	repo := configpkg.NewRepository(map[string]any{
		"security": map[string]any{
			"encryption": map[string]any{
				"key": "base64:MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=",
			},
			"hashing": map[string]any{
				"driver": "argon2id",
				"argon": map[string]any{
					"memory":  2048,
					"time":    3,
					"threads": 4,
					"verify":  true,
				},
			},
		},
	})
	manager, err = NewManagerFromRepository(repo)

	if err != nil {
		t.Fatalf("NewManagerFromRepository: %v", err)
	}

	if got, want := manager.DefaultDriver(), DriverArgon2id; got != want {
		t.Fatalf("unexpected repository driver: got %q want %q", got, want)
	}

	argonHash, err := manager.Make("secret", nil)

	if err != nil {
		t.Fatalf("Argon manager Make: %v", err)
	}

	if info := manager.Info(argonHash); info.Algorithm != DriverArgon2id {
		t.Fatalf("unexpected manager argon info: %#v", info)
	}

	if !manager.VerifyConfiguration(argonHash) {
		t.Fatal("expected manager argon verification to pass")
	}

	bcryptHash, err := NewBcrypt(BcryptConfig{}).Make("secret", nil)

	if err != nil {
		t.Fatalf("Make bcrypt hash: %v", err)
	}

	if manager.VerifyConfiguration(bcryptHash) {
		t.Fatal("expected argon2id manager verification to reject bcrypt hash")
	}

	if _, err := NewManager(Config{Driver: "nope"}); err == nil {
		t.Fatal("expected NewManager invalid driver error")
	}

	if _, err := NewManagerFromRepository(nil); err == nil {
		t.Fatal("expected NewManagerFromRepository nil error")
	}

	badManager := &Manager{defaultDriver: "nope"}

	if _, err := badManager.Make("secret", nil); err == nil {
		t.Fatal("expected Make invalid driver error")
	}

	if _, err := badManager.Check("secret", hashedValue, nil); err == nil {
		t.Fatal("expected Check invalid driver error")
	}

	if !badManager.NeedsRehash(hashedValue, nil) {
		t.Fatal("expected NeedsRehash invalid driver fallback")
	}

	if badManager.VerifyConfiguration(hashedValue) {
		t.Fatal("expected VerifyConfiguration invalid driver fallback")
	}

	if cloned := cloneOptions(nil); len(cloned) != 0 {
		t.Fatalf("expected empty cloned options, got %#v", cloned)
	}

	options := map[string]int{"rounds": 10}
	cloned := cloneOptions(options)
	cloned["rounds"] = 12

	if options["rounds"] != 10 {
		t.Fatalf("expected cloneOptions to copy map, got %#v", options)
	}

	for _, tc := range []struct {
		name     string
		value    any
		fallback int
		want     int
		wantErr  string
	}{
		{name: "nil options", value: nil, fallback: 3, want: 3},
		{name: "missing key", value: map[string]any{"other": 1}, fallback: 4, want: 4},
		{name: "int", value: map[string]any{"rounds": int(5)}, want: 5},
		{name: "int8", value: map[string]any{"rounds": int8(6)}, want: 6},
		{name: "int16", value: map[string]any{"rounds": int16(7)}, want: 7},
		{name: "int32", value: map[string]any{"rounds": int32(8)}, want: 8},
		{name: "int64", value: map[string]any{"rounds": int64(9)}, want: 9},
		{name: "uint", value: map[string]any{"rounds": uint(10)}, want: 10},
		{name: "uint8", value: map[string]any{"rounds": uint8(11)}, want: 11},
		{name: "uint16", value: map[string]any{"rounds": uint16(12)}, want: 12},
		{name: "uint32", value: map[string]any{"rounds": uint32(13)}, want: 13},
		{name: "uint64", value: map[string]any{"rounds": uint64(14)}, want: 14},
		{name: "float32", value: map[string]any{"rounds": float32(15.9)}, want: 15},
		{name: "float64", value: map[string]any{"rounds": float64(16.9)}, want: 16},
		{name: "string", value: map[string]any{"rounds": "17"}, want: 17},
		{name: "bad string", value: map[string]any{"rounds": "nope"}, wantErr: `hashing: option "rounds" must be an integer`},
		{name: "bad type", value: map[string]any{"rounds": []string{"bad"}}, wantErr: `hashing: option "rounds" must be an integer`},
		{name: "uint overflow", value: map[string]any{"rounds": uint(math.MaxInt) + 1}, wantErr: `hashing: option "rounds" exceeds int range`},
		{name: "uint64 overflow", value: map[string]any{"rounds": uint64(math.MaxInt) + 1}, wantErr: `hashing: option "rounds" exceeds int range`},
	} {
		t.Run("intOption/"+tc.name, func(t *testing.T) {
			var options map[string]any

			if typed, ok := tc.value.(map[string]any); ok {
				options = typed
			}

			got, err := intOption(options, "rounds", tc.fallback)

			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("unexpected intOption error: %v", err)
				}

				return
			}

			if err != nil || got != tc.want {
				t.Fatalf("unexpected intOption result: got=%d err=%v", got, err)
			}
		})
	}

	if got := normalizeDriver("  "); got != DriverBcrypt {
		t.Fatalf("unexpected empty normalizeDriver result: %q", got)
	}

	if got := normalizeDriver(" ArGoN2Id "); got != DriverArgon2id {
		t.Fatalf("unexpected normalizeDriver result: %q", got)
	}

	if got := normalizeBcryptConfig(BcryptConfig{}); got.Rounds != 12 {
		t.Fatalf("unexpected normalizeBcryptConfig result: %#v", got)
	}

	argonCfg := normalizeArgonConfig(ArgonConfig{})

	if argonCfg.Memory != 1024 || argonCfg.Time != 2 || argonCfg.Threads != 2 {
		t.Fatalf("unexpected normalizeArgonConfig defaults: %#v", argonCfg)
	}

	argonCfg = normalizeArgonConfig(ArgonConfig{Threads: -1})

	if argonCfg.Threads != 1 {
		t.Fatalf("expected negative threads normalization, got %#v", argonCfg)
	}

	argonCfg = normalizeArgonConfig(ArgonConfig{Threads: math.MaxUint8 + 1})

	if argonCfg.Threads != math.MaxUint8 {
		t.Fatalf("expected oversized threads normalization, got %#v", argonCfg)
	}

	if cfg := normalizeConfig(Config{}); cfg.Driver != DriverBcrypt || cfg.Bcrypt.Rounds != 12 || cfg.Argon.Memory != 1024 {
		t.Fatalf("unexpected normalizeConfig result: %#v", cfg)
	}
}

// ======================== Additional Upstream compliance tests ========================

// Upstream: testEmptyHashedValueReturnsFalse
func TestEmptyHashedValueReturnsFalse(t *testing.T) {
	bcryptHasher := NewBcrypt(BcryptConfig{Rounds: 4})

	ok, err := bcryptHasher.Check("password", "", nil)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if ok {
		t.Fatal("expected empty hashed value to return false")
	}
}

// Upstream: testNullHashedValueReturnsFalse (Go equivalent: empty string)
func TestCheckAgainstEmptyPasswordHash(t *testing.T) {
	argonHasher := NewArgon2id(ArgonConfig{Memory: 1024, Time: 2, Threads: 2})

	ok, err := argonHasher.Check("password", "", nil)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if ok {
		t.Fatal("expected check against empty hash to return false")
	}
}

// Upstream: testIsHashedWithNonHashedValue
func TestIsHashedWithNonHashedValue(t *testing.T) {
	manager, err := NewManager(Config{Driver: DriverBcrypt, Bcrypt: BcryptConfig{Rounds: 4}})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	if manager.IsHashed("not-a-hash") {
		t.Fatal("expected plain text to not be detected as hashed")
	}
	if manager.IsHashed("") {
		t.Fatal("expected empty string to not be detected as hashed")
	}
	if manager.IsHashed("password123") {
		t.Fatal("expected simple string to not be detected as hashed")
	}

	hash, _ := manager.Make("secret", nil)
	if !manager.IsHashed(hash) {
		t.Fatal("expected actual hash to be detected as hashed")
	}
}

// Upstream: testBcryptValueTooLong
// In Upstream/PHP, bcrypt silently truncates at 72 bytes.
// In Go, golang.org/x/crypto/bcrypt rejects passwords exceeding 72 bytes with an error.
// This is a behavioral difference — Go enforces the limit strictly.
func TestBcryptValueTooLong(t *testing.T) {
	bcryptHasher := NewBcrypt(BcryptConfig{Rounds: 4})

	longPassword := strings.Repeat("a", 100)
	_, err := bcryptHasher.Make(longPassword, nil)
	if err == nil {
		t.Fatal("expected Go bcrypt to reject passwords longer than 72 bytes")
	}

	// Exactly 72 bytes should work
	maxPassword := strings.Repeat("a", 72)
	hash, err := bcryptHasher.Make(maxPassword, nil)
	if err != nil {
		t.Fatalf("Make with 72-byte password: %v", err)
	}

	ok, err := bcryptHasher.Check(maxPassword, hash, nil)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !ok {
		t.Fatal("expected 72-byte password to verify")
	}
}

// Upstream: testBasicBcryptVerification / testBasicArgon2iVerification / testBasicArgon2idVerification
func TestCrossHasherVerification(t *testing.T) {
	bcryptHasher := NewBcrypt(BcryptConfig{Rounds: 4})
	argon2idHasher := NewArgon2id(ArgonConfig{Memory: 1024, Time: 2, Threads: 2})

	// Hash with bcrypt
	bcryptHash, _ := bcryptHasher.Make("secret", nil)

	// Verify bcrypt hash works with bcrypt
	ok, _ := bcryptHasher.Check("secret", bcryptHash, nil)
	if !ok {
		t.Fatal("expected bcrypt to verify its own hash")
	}

	// Verify bcrypt hash fails with wrong password
	ok, _ = bcryptHasher.Check("wrong", bcryptHash, nil)
	if ok {
		t.Fatal("expected bcrypt to reject wrong password")
	}

	// Hash with argon2id
	argonHash, _ := argon2idHasher.Make("secret", nil)

	// Verify argon hash works with argon
	ok, _ = argon2idHasher.Check("secret", argonHash, nil)
	if !ok {
		t.Fatal("expected argon2id to verify its own hash")
	}

	// Cross-check: bcrypt hasher should fail on argon hash
	ok, _ = bcryptHasher.Check("secret", argonHash, nil)
	if ok {
		t.Fatal("expected bcrypt to not verify argon2id hash")
	}
}

// Upstream: testBasicBcryptNotSupported / testBasicArgon2iNotSupported / testBasicArgon2idNotSupported
// These tests verify that unavailable drivers are handled. In Go, all drivers are always available
// since they're compiled in, so this is a gap documentation.
// GAP: Go doesn't have runtime driver availability issues like PHP extensions.

func stubArgonRandomReader(t *testing.T, reader io.Reader) func() {
	t.Helper()
	previous := argonRandomReader
	argonRandomReader = reader

	return func() {
		argonRandomReader = previous
	}
}

func (r errReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

func (r bytesReader) Read(p []byte) (int, error) {
	n := copy(p, r)

	if n < len(p) {
		return n, io.EOF
	}

	return n, nil
}
