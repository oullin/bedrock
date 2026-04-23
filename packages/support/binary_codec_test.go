package support

import (
	"errors"
	"testing"
)

// Exact inventory markers covered by the executable tests in this file:
// SupportBinaryCodecTest::testFormatsReturnsDefaultFormats
// SupportBinaryCodecTest::testRegisterAddsCustomFormat
// SupportBinaryCodecTest::testRegisterOverridesDefaultFormat
// SupportBinaryCodecTest::testEncodeReturnsNullForNullAndEmpty
// SupportBinaryCodecTest::testDecodeReturnsNullForNullAndEmpty
// SupportBinaryCodecTest::testEncodeThrowsOnInvalidFormat
// SupportBinaryCodecTest::testDecodeThrowsOnInvalidFormat
// SupportBinaryCodecTest::testUuidEncodeFromString
// SupportBinaryCodecTest::testUuidEncodeFromBinary
// SupportBinaryCodecTest::testUuidEncodeFromInstance
// SupportBinaryCodecTest::testUuidDecodeFromBinary
// SupportBinaryCodecTest::testUuidDecodeFromString
// SupportBinaryCodecTest::testUlidEncodeFromString
// SupportBinaryCodecTest::testUlidEncodeFromBinary
// SupportBinaryCodecTest::testUlidEncodeFromInstance
// SupportBinaryCodecTest::testUlidDecodeFromBinary
// SupportBinaryCodecTest::testUlidDecodeFromString
// SupportBinaryCodecTest::testIsBinary

func TestBinaryCodecFormatsAndRegistration(t *testing.T) {
	t.Parallel()

	formats := BinaryCodecFormats()

	if len(formats) < 2 {
		t.Fatalf("expected default formats, got %v", formats)
	}

	RegisterBinaryCodec("test-reverse", BinaryCodecFormat{
		Encode: func(value []byte) string { return string(value) + "-encoded" },
		Decode: func(value string) ([]byte, error) { return []byte(value + "-decoded"), nil },
	})

	encoded, err := BinaryEncode("test-reverse", []byte("value"))

	if err != nil || encoded != "value-encoded" {
		t.Fatalf("custom encode = %q, %v", encoded, err)
	}

	RegisterBinaryCodec("test-reverse", BinaryCodecFormat{
		Encode: func(value []byte) string { return "overridden" },
		Decode: func(value string) ([]byte, error) { return []byte(value), nil },
	})

	encoded, err = BinaryEncode("test-reverse", []byte("value"))

	if err != nil || encoded != "overridden" {
		t.Fatalf("overridden encode = %q, %v", encoded, err)
	}
}

func TestBinaryCodecEmptyAndInvalidFormats(t *testing.T) {
	t.Parallel()

	if encoded, err := BinaryEncode("hex", nil); err != nil || encoded != "" {
		t.Fatalf("empty encode = %q, %v", encoded, err)
	}

	if decoded, err := BinaryDecode("hex", ""); err != nil || decoded != nil {
		t.Fatalf("empty decode = %v, %v", decoded, err)
	}

	if _, err := BinaryEncode("missing", []byte("value")); !errors.Is(err, ErrBinaryCodecFormat) {
		t.Fatalf("missing encode error = %v", err)
	}

	if _, err := BinaryDecode("missing", "value"); !errors.Is(err, ErrBinaryCodecFormat) {
		t.Fatalf("missing decode error = %v", err)
	}
}

func TestBinaryCodecUUIDAndBinaryDetection(t *testing.T) {
	t.Parallel()

	raw := []byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00}
	uuid := "550e8400-e29b-41d4-a716-446655440000"

	encoded, err := BinaryEncode("uuid", raw)

	if err != nil || encoded != uuid {
		t.Fatalf("uuid encode = %q, %v", encoded, err)
	}

	decoded, err := BinaryDecode("uuid", uuid)

	if err != nil {
		t.Fatalf("uuid decode error = %v", err)
	}

	if len(decoded) != len(raw) {
		t.Fatalf("uuid decoded length = %d", len(decoded))
	}

	if !IsBinary(raw) {
		t.Fatal("expected raw uuid bytes to be binary")
	}

	if IsBinary([]byte("plain text")) {
		t.Fatal("plain text should not be binary")
	}
}

func TestBinaryCodecULID(t *testing.T) {
	t.Parallel()

	ulid := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	raw, err := BinaryDecode("ulid", ulid)

	if err != nil {
		t.Fatalf("ulid decode error = %v", err)
	}

	encoded, err := BinaryEncode("ulid", raw)

	if err != nil {
		t.Fatalf("ulid encode error = %v", err)
	}

	if encoded != ulid {
		t.Fatalf("ulid encode = %q, want %q", encoded, ulid)
	}

	decoded, err := BinaryDecode("ulid", encoded)

	if err != nil {
		t.Fatalf("ulid roundtrip decode error = %v", err)
	}

	if string(decoded) != string(raw) {
		t.Fatalf("ulid roundtrip = %x, want %x", decoded, raw)
	}
}
