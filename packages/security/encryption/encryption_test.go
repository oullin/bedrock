package encryption

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"testing"

	configpkg "github.com/gollin/packages/config"
	"github.com/gollin/packages/security/serialize"
)

func TestEncryptDecryptStringCBC(t *testing.T) {
	t.Parallel()

	e, err := New(Config{Key: []byte("1234567890abcdef"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	encrypted, err := e.EncryptString("foo")
	if err != nil {
		t.Fatalf("EncryptString: %v", err)
	}
	if encrypted == "foo" {
		t.Fatal("expected encrypted payload")
	}

	decrypted, err := e.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString: %v", err)
	}
	if decrypted != "foo" {
		t.Fatalf("unexpected decrypted value: %q", decrypted)
	}
}

func TestEncryptDecryptStringWithPreviousKeys(t *testing.T) {
	t.Parallel()

	previous, err := New(Config{Key: []byte("bbbbbbbbbbbbbbbb"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New previous: %v", err)
	}
	payload, err := previous.EncryptString("foo")
	if err != nil {
		t.Fatalf("EncryptString: %v", err)
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
		t.Fatalf("DecryptString: %v", err)
	}
	if decrypted != "foo" {
		t.Fatalf("unexpected decrypted value: %q", decrypted)
	}
}

func TestEncryptDecryptSerializedValues(t *testing.T) {
	t.Parallel()

	e, err := New(Config{Key: []byte("0123456789abcdef0123456789abcdef"), Cipher: AES256CBC})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	tests := []struct {
		name  string
		value any
		check func(t *testing.T, got any)
	}{
		{
			name:  "array",
			value: map[string]any{"foo": "bar", "count": 2},
			check: func(t *testing.T, got any) {
				t.Helper()
				mapped, ok := got.(map[string]interface{})
				if !ok {
					t.Fatalf("expected map, got %T", got)
				}
				if mapped["foo"] != "bar" {
					t.Fatalf("unexpected map value: %#v", mapped)
				}
				if mapped["count"].(int64) != 2 {
					t.Fatalf("unexpected count: %#v", mapped["count"])
				}
			},
		},
		{
			name:  "nil",
			value: nil,
			check: func(t *testing.T, got any) {
				t.Helper()
				if got != nil {
					t.Fatalf("expected nil, got %#v", got)
				}
			},
		},
		{
			name:  "object",
			value: serialize.PHPObject{ClassName: "User", Properties: map[string]any{"id": 1}},
			check: func(t *testing.T, got any) {
				t.Helper()
				obj, ok := got.(serialize.PHPObject)
				if !ok {
					t.Fatalf("expected PHPObject, got %T", got)
				}
				if obj.ClassName != "User" {
					t.Fatalf("unexpected class name: %s", obj.ClassName)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := e.Encrypt(tt.value)
			if err != nil {
				t.Fatalf("Encrypt: %v", err)
			}

			decrypted, err := e.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt: %v", err)
			}

			tt.check(t, decrypted)
		})
	}
}

func TestEncryptDecryptGCM(t *testing.T) {
	t.Parallel()

	e, err := New(Config{Key: []byte("0123456789abcdef0123456789abcdef"), Cipher: AES256GCM})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	encrypted, err := e.EncryptString("bar")
	if err != nil {
		t.Fatalf("EncryptString: %v", err)
	}

	decrypted, err := e.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString: %v", err)
	}
	if decrypted != "bar" {
		t.Fatalf("unexpected decrypted value: %q", decrypted)
	}

	data := decodePayload(t, encrypted)
	if data.MAC != "" {
		t.Fatalf("expected empty mac for AEAD payload, got %q", data.MAC)
	}
	if data.Tag == "" {
		t.Fatal("expected AEAD tag")
	}
}

func TestEncryptedLengthIsFixed(t *testing.T) {
	t.Parallel()

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

func TestPreviousKeyPayloadFixtureDecrypts(t *testing.T) {
	t.Parallel()

	e, err := New(Config{
		Key:          []byte("aaaaaaaaaaaaaaaa"),
		PreviousKeys: [][]byte{[]byte("bbbbbbbbbbbbbbbb")},
		Cipher:       AES128CBC,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	payload := "eyJpdiI6Ilg0dFM5TVRibEFqZW54c3lQdWJoVVE9PSIsInZhbHVlIjoiRGJpa2p2ZHI3eUs0dUtRakJneUhUUT09IiwibWFjIjoiMjBjZWYxODdhNThhOTk4MTk1NTc0YTE1MDgzODU1OWE0ZmQ4MDc5ZjMxYThkOGM1ZmM1MzlmYzBkYTBjMWI1ZiIsInRhZyI6IiJ9"

	decrypted, err := e.DecryptString(payload)
	if err != nil {
		t.Fatalf("DecryptString: %v", err)
	}
	if decrypted != "foo" {
		t.Fatalf("unexpected decrypted value: %q", decrypted)
	}
}

func TestDecryptFailures(t *testing.T) {
	t.Parallel()

	base, err := New(Config{Key: []byte("aaaaaaaaaaaaaaaa"), Cipher: AES128CBC})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	gcm, err := New(Config{Key: []byte("0123456789abcdef0123456789abcdef"), Cipher: AES256GCM})
	if err != nil {
		t.Fatalf("New GCM: %v", err)
	}

	payload, err := base.EncryptString("foo")
	if err != nil {
		t.Fatalf("EncryptString: %v", err)
	}
	gcmPayload, err := gcm.EncryptString("foo")
	if err != nil {
		t.Fatalf("EncryptString GCM: %v", err)
	}

	tests := []struct {
		name    string
		payload string
		target  *Encrypter
		want    error
		mutate  func(encryptedPayload) string
	}{
		{
			name:    "invalid payload base64",
			payload: payload,
			target:  base,
			want:    errInvalidPayload,
			mutate: func(_ encryptedPayload) string {
				return "%%%invalid%%%"
			},
		},
		{
			name:    "invalid mac",
			payload: payload,
			target:  func() *Encrypter { e, _ := New(Config{Key: []byte("bbbbbbbbbbbbbbbb"), Cipher: AES128CBC}); return e }(),
			want:    errInvalidMAC,
			mutate: func(p encryptedPayload) string {
				p.MAC = strings.Repeat("0", len(p.MAC))
				return encodePayload(p)
			},
		},
		{
			name:    "invalid iv",
			payload: payload,
			target:  base,
			want:    errInvalidPayload,
			mutate: func(p encryptedPayload) string {
				p.IV += "A"
				return encodePayload(p)
			},
		},
		{
			name:    "unexpected tag on CBC",
			payload: payload,
			target:  base,
			want:    errUnexpectedTag,
			mutate: func(p encryptedPayload) string {
				p.Tag = base64.StdEncoding.EncodeToString([]byte("set-manually"))
				return encodePayload(p)
			},
		},
		{
			name:    "truncated gcm tag",
			payload: gcmPayload,
			target:  gcm,
			want:    errDecryptFailed,
			mutate: func(p encryptedPayload) string {
				tag, _ := base64.StdEncoding.DecodeString(p.Tag)
				p.Tag = base64.StdEncoding.EncodeToString(tag[:4])
				return encodePayload(p)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var raw string
			if tt.mutate == nil {
				raw = tt.payload
			} else {
				raw = tt.mutate(decodePayload(t, tt.payload))
			}

			_, err := tt.target.DecryptString(raw)
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}

func TestSupportedGenerateKeyAndParseUpstreamKey(t *testing.T) {
	t.Parallel()

	if !Supported([]byte("aaaaaaaaaaaaaaaa"), AES128CBC) {
		t.Fatal("expected AES-128-CBC key to be supported")
	}
	if Supported([]byte("short"), AES128CBC) {
		t.Fatal("expected short key to be rejected")
	}

	generated, err := GenerateKey(AES256GCM)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	if len(generated) != 32 {
		t.Fatalf("unexpected generated key length: %d", len(generated))
	}

	parsed, err := ParseUpstreamKey("base64:YWJjZGVmZ2hpamtsbW5vcA==")
	if err != nil {
		t.Fatalf("ParseUpstreamKey: %v", err)
	}
	if string(parsed) != "abcdefghijklmnop" {
		t.Fatalf("unexpected parsed key: %q", string(parsed))
	}
}

func TestAppearsEncrypted(t *testing.T) {
	t.Parallel()

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
	if AppearsEncrypted("APP_NAME=Upstream") {
		t.Fatal("expected plain text to be rejected")
	}
}

func TestNewFromRepository(t *testing.T) {
	t.Parallel()

	repo := configpkg.NewRepository(map[string]any{
		"security": map[string]any{
			"encryption": map[string]any{
				"key":    "aaaaaaaaaaaaaaaa",
				"cipher": "aes-128-cbc",
			},
		},
	})

	encrypter, err := NewFromRepository(repo)
	if err != nil {
		t.Fatalf("NewFromRepository: %v", err)
	}
	if got := string(encrypter.GetKey()); got != "aaaaaaaaaaaaaaaa" {
		t.Fatalf("unexpected key: %q", got)
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
