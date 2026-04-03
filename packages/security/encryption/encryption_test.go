package encryption

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"math"
	"strings"
	"testing"

	configpkg "github.com/gollin/packages/config"
)

func TestNewSupportedAndAccessors(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want string
	}{
		{name: "cbc", cfg: Config{Key: []byte("1234567890abcdef"), Cipher: AES128CBC}, want: "aes-128-cbc"},
		{name: "gcm mixed case", cfg: Config{Key: []byte("0123456789abcdef"), Cipher: Cipher("AES-128-GCM")}, want: "aes-128-gcm"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := New(tt.cfg)
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			if e.cipher != tt.want {
				t.Fatalf("unexpected normalized cipher: %q", e.cipher)
			}
			if !reflectBytes(e.GetKey(), tt.cfg.Key) {
				t.Fatalf("unexpected key: %#v", e.GetKey())
			}

			key := e.GetKey()
			key[0] = 'x'
			if e.GetKey()[0] == 'x' {
				t.Fatal("expected GetKey to return a clone")
			}
		})
	}

	e, err := New(Config{
		Key:          []byte("aaaaaaaaaaaaaaaa"),
		PreviousKeys: [][]byte{[]byte("bbbbbbbbbbbbbbbb")},
		Cipher:       AES128CBC,
	})
	if err != nil {
		t.Fatalf("New with previous keys: %v", err)
	}
	previous := e.GetPreviousKeys()
	previous[0][0] = 'x'
	if e.GetPreviousKeys()[0][0] == 'x' {
		t.Fatal("expected GetPreviousKeys to return clones")
	}
	allKeys := e.GetAllKeys()
	allKeys[1][0] = 'y'
	if e.GetAllKeys()[1][0] == 'y' {
		t.Fatal("expected GetAllKeys to return clones")
	}

	for _, tc := range []Config{
		{Key: []byte("short"), Cipher: AES128CBC},
		{Key: []byte("1234567890abcdef"), Cipher: Cipher("AES-256-GCM")},
		{Key: []byte("1234567890abcdef"), PreviousKeys: [][]byte{[]byte("short")}, Cipher: AES128CBC},
		{Key: []byte("1234567890abcdef"), Cipher: Cipher("AES-256-CFB8")},
	} {
		if _, err := New(tc); !errors.Is(err, errUnsupportedCipher) {
			t.Fatalf("expected unsupported cipher error for %#v, got %v", tc, err)
		}
	}

	if !Supported([]byte("aaaaaaaaaaaaaaaa"), AES128CBC) {
		t.Fatal("expected AES-128-CBC key to be supported")
	}
	if !Supported([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), AES256GCM) {
		t.Fatal("expected AES-256-GCM key to be supported")
	}
	if Supported([]byte("short"), AES128CBC) {
		t.Fatal("expected short key to be rejected")
	}
	if Supported([]byte("aaaaaaaaaaaaaaaa"), Cipher("unsupported")) {
		t.Fatal("expected unsupported cipher to be rejected")
	}
}

func TestGenerateKey(t *testing.T) {
	restoreRandom := stubRandomReader(t, bytesReader(make([]byte, 64)))
	defer restoreRandom()

	for _, tc := range []struct {
		cipher Cipher
		size   int
	}{
		{cipher: AES128CBC, size: 16},
		{cipher: AES256CBC, size: 32},
		{cipher: AES128GCM, size: 16},
		{cipher: AES256GCM, size: 32},
	} {
		key, err := GenerateKey(tc.cipher)
		if err != nil {
			t.Fatalf("GenerateKey(%q): %v", tc.cipher, err)
		}
		if len(key) != tc.size {
			t.Fatalf("unexpected generated key length for %q: %d", tc.cipher, len(key))
		}
	}

	if _, err := GenerateKey(Cipher("unsupported")); !errors.Is(err, errUnsupportedCipher) {
		t.Fatalf("expected unsupported cipher error, got %v", err)
	}

	restoreRandom = stubRandomReader(t, errReader{err: io.ErrUnexpectedEOF})
	defer restoreRandom()
	if _, err := GenerateKey(AES128CBC); err == nil || !strings.Contains(err.Error(), "read random key bytes") {
		t.Fatalf("expected random key read error, got %v", err)
	}
}

func TestEncryptDecryptStringAndJSONValues(t *testing.T) {
	cbc, err := New(Config{Key: []byte("1234567890abcdef"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New CBC: %v", err)
	}
	gcm, err := New(Config{Key: []byte("0123456789abcdef0123456789abcdef"), Cipher: AES256GCM})
	if err != nil {
		t.Fatalf("New GCM: %v", err)
	}

	for _, tc := range []struct {
		name  string
		value any
		check func(t *testing.T, got any)
	}{
		{name: "map", value: map[string]any{"foo": "bar", "count": 2}, check: func(t *testing.T, got any) {
			mapped, ok := got.(map[string]any)
			if !ok || mapped["foo"] != "bar" || mapped["count"].(float64) != 2 {
				t.Fatalf("unexpected map value: %#v", got)
			}
		}},
		{name: "slice", value: []any{"foo", true, 2}, check: func(t *testing.T, got any) {
			values, ok := got.([]any)
			if !ok || len(values) != 3 || values[0] != "foo" || values[1] != true || values[2].(float64) != 2 {
				t.Fatalf("unexpected slice value: %#v", got)
			}
		}},
		{name: "string", value: "hello", check: func(t *testing.T, got any) {
			if got != "hello" {
				t.Fatalf("unexpected string value: %#v", got)
			}
		}},
		{name: "bool", value: true, check: func(t *testing.T, got any) {
			if got != true {
				t.Fatalf("unexpected bool value: %#v", got)
			}
		}},
		{name: "nil", value: nil, check: func(t *testing.T, got any) {
			if got != nil {
				t.Fatalf("unexpected nil value: %#v", got)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := cbc.Encrypt(tc.value)
			if err != nil {
				t.Fatalf("Encrypt: %v", err)
			}
			decrypted, err := cbc.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt: %v", err)
			}
			tc.check(t, decrypted)
		})
	}

	for _, tc := range []struct {
		name      string
		encrypter *Encrypter
		value     string
	}{
		{name: "cbc", encrypter: cbc, value: ""},
		{name: "cbc long", encrypter: cbc, value: strings.Repeat("a", 1000)},
		{name: "gcm", encrypter: gcm, value: "bar"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := tc.encrypter.EncryptString(tc.value)
			if err != nil {
				t.Fatalf("EncryptString: %v", err)
			}
			if encrypted == tc.value {
				t.Fatal("expected encrypted payload")
			}
			decrypted, err := tc.encrypter.DecryptString(encrypted)
			if err != nil {
				t.Fatalf("DecryptString: %v", err)
			}
			if decrypted != tc.value {
				t.Fatalf("unexpected decrypted value: %q", decrypted)
			}
		})
	}

	gcmPayload, err := gcm.EncryptString("bar")
	if err != nil {
		t.Fatalf("EncryptString GCM: %v", err)
	}
	gcmData := decodePayload(t, gcmPayload)
	if gcmData.MAC != "" || gcmData.Tag == "" {
		t.Fatalf("unexpected GCM payload: %#v", gcmData)
	}

	cbcPayload, err := cbc.EncryptString("foo")
	if err != nil {
		t.Fatalf("EncryptString CBC: %v", err)
	}
	cbcData := decodePayload(t, cbcPayload)
	if cbcData.MAC == "" || cbcData.Tag != "" {
		t.Fatalf("unexpected CBC payload: %#v", cbcData)
	}
}

func TestEncryptedLengthIsFixed(t *testing.T) {
	e, err := New(Config{Key: []byte("aaaaaaaaaaaaaaaa"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	minLen, maxLen := math.MaxInt, 0
	for i := 0; i < 50; i++ {
		encrypted, err := e.EncryptString("foo")
		if err != nil {
			t.Fatalf("EncryptString: %v", err)
		}
		if len(encrypted) < minLen {
			minLen = len(encrypted)
		}
		if len(encrypted) > maxLen {
			maxLen = len(encrypted)
		}
	}

	if minLen != maxLen {
		t.Fatalf("expected fixed payload length, got min=%d max=%d", minLen, maxLen)
	}
}

func TestEncryptFailures(t *testing.T) {
	e, err := New(Config{Key: []byte("aaaaaaaaaaaaaaaa"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if _, err := e.Encrypt(make(chan int)); err == nil || !strings.Contains(err.Error(), "serialise value") {
		t.Fatalf("expected JSON serialise failure, got %v", err)
	}

	restoreRandom := stubRandomReader(t, errReader{err: io.ErrUnexpectedEOF})
	defer restoreRandom()
	if _, err := e.EncryptString("foo"); err == nil || !strings.Contains(err.Error(), "read random iv bytes") {
		t.Fatalf("expected random IV read failure, got %v", err)
	}

	restoreRandom = stubRandomReader(t, bytesReader(make([]byte, 64)))
	defer restoreRandom()
	restoreMarshal := stubMarshalPayloadJSON(t, func(encryptedPayload) ([]byte, error) {
		return nil, errors.New("boom")
	})
	defer restoreMarshal()
	if _, err := e.EncryptString("foo"); !errors.Is(err, errEncryptFailed) {
		t.Fatalf("expected payload marshal failure to map to encrypt error, got %v", err)
	}
}

func TestDecryptFailuresAndPayloadValidation(t *testing.T) {
	cbc, err := New(Config{Key: []byte("aaaaaaaaaaaaaaaa"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New CBC: %v", err)
	}
	gcm, err := New(Config{Key: []byte("0123456789abcdef0123456789abcdef"), Cipher: AES256GCM})
	if err != nil {
		t.Fatalf("New GCM: %v", err)
	}

	cbcPayload, err := cbc.EncryptString("foo")
	if err != nil {
		t.Fatalf("EncryptString CBC: %v", err)
	}
	gcmPayload, err := gcm.EncryptString("foo")
	if err != nil {
		t.Fatalf("EncryptString GCM: %v", err)
	}

	legacyPayload, err := cbc.encryptBytes([]byte(`a:2:{s:3:"foo";s:3:"bar";}`))
	if err != nil {
		t.Fatalf("encryptBytes legacy payload: %v", err)
	}
	if _, err := cbc.Decrypt(legacyPayload); err == nil || !strings.Contains(err.Error(), "unserialise value") {
		t.Fatalf("expected legacy payload JSON decode failure, got %v", err)
	}
	if _, err := cbc.Decrypt("%%%invalid%%%"); !errors.Is(err, errInvalidPayload) {
		t.Fatalf("expected Decrypt to return invalid payload error, got %v", err)
	}

	badJSON := base64.StdEncoding.EncodeToString([]byte("{not-json"))
	if _, err := cbc.DecryptString("%%%invalid%%%"); !errors.Is(err, errInvalidPayload) {
		t.Fatalf("expected invalid payload error, got %v", err)
	}
	if _, err := cbc.DecryptString(badJSON); !errors.Is(err, errInvalidPayload) {
		t.Fatalf("expected invalid JSON payload error, got %v", err)
	}

	tampered := []struct {
		name string
		raw  string
		want error
	}{
		{name: "wrong mac", raw: encodePayload(func() encryptedPayload {
			p := decodePayload(t, cbcPayload)
			p.MAC = strings.Repeat("0", len(p.MAC))
			return p
		}()), want: errInvalidMAC},
		{name: "iv length mismatch", raw: encodePayload(func() encryptedPayload {
			p := decodePayload(t, cbcPayload)
			p.IV += "A"
			return p
		}()), want: errInvalidPayload},
		{name: "unexpected tag", raw: encodePayload(func() encryptedPayload {
			p := decodePayload(t, cbcPayload)
			p.Tag = base64.StdEncoding.EncodeToString([]byte("set-manually"))
			return p
		}()), want: errUnexpectedTag},
		{name: "bad ciphertext base64", raw: encodePayload(func() encryptedPayload {
			p := decodePayload(t, cbcPayload)
			p.Value = "%%%invalid%%%"
			return p
		}()), want: errDecryptFailed},
		{name: "bad tag base64", raw: encodePayload(func() encryptedPayload {
			p := decodePayload(t, gcmPayload)
			p.Tag = "%%%invalid%%%"
			return p
		}()), want: errDecryptFailed},
		{name: "truncated gcm tag", raw: encodePayload(func() encryptedPayload {
			p := decodePayload(t, gcmPayload)
			tag, _ := base64.StdEncoding.DecodeString(p.Tag)
			p.Tag = base64.StdEncoding.EncodeToString(tag[:4])
			return p
		}()), want: errDecryptFailed},
		{name: "modified gcm tag", raw: encodePayload(func() encryptedPayload {
			p := decodePayload(t, gcmPayload)
			tag, _ := base64.StdEncoding.DecodeString(p.Tag)
			tag[0] ^= 0x01
			p.Tag = base64.StdEncoding.EncodeToString(tag)
			return p
		}()), want: errDecryptFailed},
	}

	for _, tc := range tampered {
		t.Run(tc.name, func(t *testing.T) {
			target := cbc
			if strings.Contains(tc.name, "gcm") || strings.Contains(tc.name, "tag") && !strings.Contains(tc.name, "unexpected") {
				target = gcm
			}
			if _, err := target.DecryptString(tc.raw); !errors.Is(err, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, err)
			}
		})
	}

	otherCBC, err := New(Config{Key: []byte("bbbbbbbbbbbbbbbb"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New wrong-key CBC: %v", err)
	}
	if _, err := otherCBC.DecryptString(cbcPayload); !errors.Is(err, errInvalidMAC) {
		t.Fatalf("expected invalid MAC for wrong key, got %v", err)
	}

	if _, err := cbc.DecryptString(encodePayload(encryptedPayload{
		IV:    base64.StdEncoding.EncodeToString([]byte("short")),
		Value: base64.StdEncoding.EncodeToString([]byte("value")),
		MAC:   "",
	})); !errors.Is(err, errInvalidPayload) {
		t.Fatalf("expected invalid short IV payload, got %v", err)
	}
}

func TestPreviousKeyDecryptionAndFixture(t *testing.T) {
	previous, err := New(Config{Key: []byte("bbbbbbbbbbbbbbbb"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New previous: %v", err)
	}
	payload, err := previous.EncryptString("foo")
	if err != nil {
		t.Fatalf("EncryptString previous: %v", err)
	}

	current, err := New(Config{
		Key:          []byte("aaaaaaaaaaaaaaaa"),
		PreviousKeys: [][]byte{[]byte("bbbbbbbbbbbbbbbb")},
		Cipher:       AES128CBC,
	})
	if err != nil {
		t.Fatalf("New current: %v", err)
	}

	decrypted, err := current.DecryptString(payload)
	if err != nil {
		t.Fatalf("DecryptString previous key: %v", err)
	}
	if decrypted != "foo" {
		t.Fatalf("unexpected previous-key decrypted value: %q", decrypted)
	}

	fixture := "eyJpdiI6Ilg0dFM5TVRibEFqZW54c3lQdWJoVVE9PSIsInZhbHVlIjoiRGJpa2p2ZHI3eUs0dUtRakJneUhUUT09IiwibWFjIjoiMjBjZWYxODdhNThhOTk4MTk1NTc0YTE1MDgzODU1OWE0ZmQ4MDc5ZjMxYThkOGM1ZmM1MzlmYzBkYTBjMWI1ZiIsInRhZyI6IiJ9"
	decrypted, err = current.DecryptString(fixture)
	if err != nil {
		t.Fatalf("DecryptString fixture: %v", err)
	}
	if decrypted != "foo" {
		t.Fatalf("unexpected fixture decrypted value: %q", decrypted)
	}
}

func TestAppearsEncryptedAndRepositoryLoading(t *testing.T) {
	e, err := New(Config{Key: []byte("aaaaaaaaaaaaaaaa"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	payload, err := e.EncryptString("foo")
	if err != nil {
		t.Fatalf("EncryptString: %v", err)
	}

	if !AppearsEncrypted(payload) {
		t.Fatal("expected payload to look encrypted")
	}
	if !AppearsEncrypted(encodePayload(encryptedPayload{IV: "aXY=", Value: "dmFsdWU=", MAC: "mac"})) {
		t.Fatal("expected payload shape to look encrypted")
	}
	for _, value := range []string{"foo", "APP_NAME=Laravel", "APP_NAME=Laravel\nAPP_ENV=local", "%%%invalid%%%", base64.StdEncoding.EncodeToString([]byte("{"))} {
		if AppearsEncrypted(value) {
			t.Fatalf("expected plain value to be rejected: %q", value)
		}
	}

	repo := configpkg.NewRepository(map[string]any{
		"security": map[string]any{
			"encryption": map[string]any{
				"key":    "aaaaaaaaaaaaaaaa",
				"cipher": "aes-128-cbc",
			},
		},
	})
	fromRepo, err := NewFromRepository(repo)
	if err != nil {
		t.Fatalf("NewFromRepository: %v", err)
	}
	if got := string(fromRepo.GetKey()); got != "aaaaaaaaaaaaaaaa" {
		t.Fatalf("unexpected repository key: %q", got)
	}
	if _, err := NewFromRepository(nil); err == nil {
		t.Fatal("expected NewFromRepository nil error")
	}
}

func TestPayloadAndHelperFunctions(t *testing.T) {
	cbc, err := New(Config{Key: []byte("aaaaaaaaaaaaaaaa"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New CBC: %v", err)
	}
	gcm, err := New(Config{Key: []byte("0123456789abcdef0123456789abcdef"), Cipher: AES256GCM})
	if err != nil {
		t.Fatalf("New GCM: %v", err)
	}

	validCBC := encryptedPayload{
		IV:    base64.StdEncoding.EncodeToString([]byte("1234567890abcdef")),
		Value: base64.StdEncoding.EncodeToString([]byte("value")),
		MAC:   "mac",
	}
	if payload, ok := cbc.validPayload(map[string]any{"iv": validCBC.IV, "value": validCBC.Value, "mac": validCBC.MAC}); !ok || payload.MAC != "mac" {
		t.Fatalf("expected valid CBC payload, got %#v ok=%v", payload, ok)
	}
	if _, ok := cbc.validPayload(nil); ok {
		t.Fatal("expected nil payload to be rejected")
	}
	for _, payload := range []map[string]any{
		{"iv": []string{"bad"}, "value": "", "mac": ""},
		{"iv": validCBC.IV, "value": []string{"bad"}, "mac": ""},
		{"iv": validCBC.IV, "value": "", "mac": []string{"bad"}},
		{"iv": validCBC.IV, "value": "", "mac": "", "tag": []string{"bad"}},
		{"iv": "%%%invalid%%%", "value": "", "mac": ""},
		{"iv": base64.StdEncoding.EncodeToString([]byte("short")), "value": "", "mac": ""},
	} {
		if _, ok := cbc.validPayload(payload); ok {
			t.Fatalf("expected invalid payload to be rejected: %#v", payload)
		}
	}

	if _, err := cbc.getPayload("%%%invalid%%%"); !errors.Is(err, errInvalidPayload) {
		t.Fatalf("expected invalid payload decode error, got %v", err)
	}
	if _, err := cbc.getPayload(base64.StdEncoding.EncodeToString([]byte("{"))); !errors.Is(err, errInvalidPayload) {
		t.Fatalf("expected invalid payload JSON error, got %v", err)
	}

	if err := cbc.ensureTagValid([]byte("unexpected")); !errors.Is(err, errUnexpectedTag) {
		t.Fatalf("expected unexpected tag error, got %v", err)
	}
	if err := gcm.ensureTagValid([]byte("short")); !errors.Is(err, errDecryptFailed) {
		t.Fatalf("expected short AEAD tag error, got %v", err)
	}
	if err := gcm.ensureTagValid(make([]byte, 16)); err != nil {
		t.Fatalf("expected valid AEAD tag, got %v", err)
	}

	if cbc.shouldValidateMAC() != true || gcm.shouldValidateMAC() != false {
		t.Fatal("unexpected MAC validation flags")
	}
	if meta, ok := cipherMeta(AES256GCM); !ok || meta.keySize != 32 || !meta.aead {
		t.Fatalf("unexpected cipher meta: %#v ok=%v", meta, ok)
	}
	if _, ok := cipherMeta(Cipher("unsupported")); ok {
		t.Fatal("expected unsupported cipher meta lookup to fail")
	}
	if got := normalizeCipher(Cipher(" AeS-128-GcM ")); got != "aes-128-gcm" {
		t.Fatalf("unexpected normalized cipher: %q", got)
	}
	if isAEADCipher(string(AES256GCM)) != true || isAEADCipher(string(AES128CBC)) != false {
		t.Fatal("unexpected AEAD detection")
	}
	if ivLength(string(AES256GCM)) != 12 || ivLength(string(AES128CBC)) != 16 {
		t.Fatal("unexpected IV lengths")
	}

	mac := hash("iv", "value", []byte("key"))
	if len(mac) != 64 {
		t.Fatalf("unexpected MAC length: %d", len(mac))
	}
	if !validMACForKey(encryptedPayload{IV: "iv", Value: "value", MAC: mac}, []byte("key")) {
		t.Fatal("expected MAC to validate")
	}
	if validMACForKey(encryptedPayload{IV: "iv", Value: "value", MAC: mac}, []byte("nope")) {
		t.Fatal("expected MAC mismatch")
	}
}

func TestCipherHelpers(t *testing.T) {
	plaintext := []byte("hello world")
	iv := []byte("1234567890abcdef")

	ciphertext, err := encryptCBC([]byte("1234567890abcdef"), iv, plaintext)
	if err != nil {
		t.Fatalf("encryptCBC: %v", err)
	}
	decrypted, err := decryptCBC([]byte("1234567890abcdef"), iv, ciphertext)
	if err != nil {
		t.Fatalf("decryptCBC: %v", err)
	}
	if string(decrypted) != string(plaintext) {
		t.Fatalf("unexpected CBC decrypted value: %q", string(decrypted))
	}
	if _, err := encryptCBC([]byte("short"), iv, plaintext); err == nil {
		t.Fatal("expected encryptCBC invalid key error")
	}
	if _, err := (&Encrypter{key: []byte("short"), cipher: string(AES128CBC)}).encryptBytes(plaintext); !errors.Is(err, errEncryptFailed) {
		t.Fatalf("expected encryptBytes invalid CBC key error, got %v", err)
	}
	if _, err := decryptCBC([]byte("1234567890abcdef"), iv, []byte("short")); !errors.Is(err, errDecryptFailed) {
		t.Fatalf("expected decryptCBC invalid ciphertext error, got %v", err)
	}
	if _, err := decryptCBC([]byte("short"), iv, ciphertext); err == nil {
		t.Fatal("expected decryptCBC invalid key error")
	}

	gcmIV := []byte("123456789012")
	gcmCiphertext, tag, err := encryptGCM(string(AES256GCM), []byte("0123456789abcdef0123456789abcdef"), gcmIV, plaintext)
	if err != nil {
		t.Fatalf("encryptGCM: %v", err)
	}
	decrypted, err = decryptGCM(string(AES256GCM), []byte("0123456789abcdef0123456789abcdef"), gcmIV, gcmCiphertext, tag)
	if err != nil {
		t.Fatalf("decryptGCM: %v", err)
	}
	if string(decrypted) != string(plaintext) {
		t.Fatalf("unexpected GCM decrypted value: %q", string(decrypted))
	}
	if _, _, err := encryptGCM(string(AES256GCM), []byte("short"), gcmIV, plaintext); err == nil {
		t.Fatal("expected encryptGCM invalid key error")
	}
	if _, err := (&Encrypter{key: []byte("short"), cipher: string(AES256GCM)}).encryptBytes(plaintext); !errors.Is(err, errEncryptFailed) {
		t.Fatalf("expected encryptBytes invalid GCM key error, got %v", err)
	}
	if _, err := decryptGCM(string(AES256GCM), []byte("short"), gcmIV, gcmCiphertext, tag); err == nil {
		t.Fatal("expected decryptGCM invalid key error")
	}
}

func TestPaddingAndKeyParsing(t *testing.T) {
	padded := pkcs7Pad([]byte("1234567890abcdef"), 16)
	if len(padded) != 32 {
		t.Fatalf("expected exact-block padding to add a full block, got %d", len(padded))
	}
	unpadded, err := pkcs7Unpad(padded, 16)
	if err != nil {
		t.Fatalf("pkcs7Unpad: %v", err)
	}
	if string(unpadded) != "1234567890abcdef" {
		t.Fatalf("unexpected unpadded value: %q", string(unpadded))
	}
	for _, data := range [][]byte{
		nil,
		[]byte("short"),
		append([]byte("1234567890abcdef"), 0),
		append([]byte("1234567890abcde"), 17),
		append([]byte("1234567890abcd"), []byte{1, 2}...),
	} {
		if _, err := pkcs7Unpad(data, 16); !errors.Is(err, errDecryptFailed) {
			t.Fatalf("expected invalid padding error for %#v, got %v", data, err)
		}
	}

	parsed, err := ParseLaravelKey("base64:YWJjZGVmZ2hpamtsbW5vcA==")
	if err != nil {
		t.Fatalf("ParseLaravelKey base64: %v", err)
	}
	if string(parsed) != "abcdefghijklmnop" {
		t.Fatalf("unexpected parsed key: %q", string(parsed))
	}
	parsed, err = ParseLaravelKey(" raw ")
	if err != nil {
		t.Fatalf("ParseLaravelKey raw: %v", err)
	}
	if string(parsed) != "raw" {
		t.Fatalf("unexpected raw key parse: %q", string(parsed))
	}
	if _, err := ParseLaravelKey(" "); err == nil {
		t.Fatal("expected blank ParseLaravelKey failure")
	}
	if _, err := ParseLaravelKey("base64:not-valid"); err == nil {
		t.Fatal("expected invalid base64 ParseLaravelKey failure")
	}
}

func decodePayload(t *testing.T, raw string) encryptedPayload {
	t.Helper()

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		t.Fatalf("DecodeString: %v", err)
	}
	var payload encryptedPayload
	if err := json.Unmarshal(decoded, &payload); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	return payload
}

func encodePayload(payload encryptedPayload) string {
	encoded, _ := json.Marshal(payload)
	return base64.StdEncoding.EncodeToString(encoded)
}

func stubRandomReader(t *testing.T, reader io.Reader) func() {
	t.Helper()
	previous := randomReader
	randomReader = reader
	return func() {
		randomReader = previous
	}
}

func stubMarshalPayloadJSON(t *testing.T, fn func(encryptedPayload) ([]byte, error)) func() {
	t.Helper()
	previous := marshalPayloadJSON
	marshalPayloadJSON = fn
	return func() {
		marshalPayloadJSON = previous
	}
}

type errReader struct {
	err error
}

func (r errReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

type bytesReader []byte

func (r bytesReader) Read(p []byte) (int, error) {
	n := copy(p, r)
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

func reflectBytes(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
