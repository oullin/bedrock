package support

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"sync"
)

var (
	ErrBinaryCodecFormat = errors.New("support: unknown binary codec format")

	binaryCodecMu      sync.RWMutex
	binaryCodecFormats = map[string]BinaryCodecFormat{
		"hex": {
			Encode: hex.EncodeToString,
			Decode: hex.DecodeString,
		},
		"uuid": {
			Encode: encodeUUID,
			Decode: decodeUUID,
		},
		"ulid": {
			Encode: encodeULID,
			Decode: decodeULID,
		},
	}
)

// BinaryCodecFormat contains binary encode/decode functions.
type BinaryCodecFormat struct {
	Encode func([]byte) string
	Decode func(string) ([]byte, error)
}

// BinaryCodecFormats returns registered format names.
func BinaryCodecFormats() []string {
	binaryCodecMu.RLock()
	defer binaryCodecMu.RUnlock()

	formats := make([]string, 0, len(binaryCodecFormats))
	for name := range binaryCodecFormats {
		formats = append(formats, name)
	}

	return formats
}

// RegisterBinaryCodec registers or overrides a binary codec format.
func RegisterBinaryCodec(name string, format BinaryCodecFormat) {
	binaryCodecMu.Lock()
	defer binaryCodecMu.Unlock()

	binaryCodecFormats[name] = format
}

// BinaryEncode encodes bytes using a registered format.
func BinaryEncode(format string, value []byte) (string, error) {
	if len(value) == 0 {
		return "", nil
	}

	binaryCodecMu.RLock()
	codec, ok := binaryCodecFormats[format]
	binaryCodecMu.RUnlock()
	if !ok {
		return "", ErrBinaryCodecFormat
	}

	return codec.Encode(value), nil
}

// BinaryDecode decodes a string using a registered format.
func BinaryDecode(format string, value string) ([]byte, error) {
	if value == "" {
		return nil, nil
	}

	binaryCodecMu.RLock()
	codec, ok := binaryCodecFormats[format]
	binaryCodecMu.RUnlock()
	if !ok {
		return nil, ErrBinaryCodecFormat
	}

	return codec.Decode(value)
}

// IsBinary reports whether value contains non-printable bytes.
func IsBinary(value []byte) bool {
	for _, b := range value {
		if b == 0 || b < 0x20 || b > 0x7e {
			return true
		}
	}

	return false
}

func encodeUUID(value []byte) string {
	if len(value) != 16 {
		return hex.EncodeToString(value)
	}

	encoded := hex.EncodeToString(value)

	return encoded[:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:]
}

func decodeUUID(value string) ([]byte, error) {
	normalized := strings.ReplaceAll(value, "-", "")

	return hex.DecodeString(normalized)
}

const crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

func encodeULID(value []byte) string {
	if len(value) != 16 {
		return hex.EncodeToString(value)
	}

	n := new(big.Int).SetBytes(value)
	base := big.NewInt(32)
	mod := new(big.Int)
	out := make([]byte, 26)

	for i := len(out) - 1; i >= 0; i-- {
		n.DivMod(n, base, mod)
		out[i] = crockfordAlphabet[mod.Int64()]
	}

	return string(out)
}

func decodeULID(value string) ([]byte, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	normalized = strings.Map(func(r rune) rune {
		switch r {
		case 'I', 'L':
			return '1'
		case 'O':
			return '0'
		default:
			return r
		}
	}, normalized)

	if len(normalized) != 26 {
		return hex.DecodeString(normalized)
	}

	n := big.NewInt(0)
	base := big.NewInt(32)

	for _, r := range normalized {
		idx := strings.IndexRune(crockfordAlphabet, r)
		if idx < 0 {
			return nil, hex.InvalidByteError(r)
		}
		n.Mul(n, base)
		n.Add(n, big.NewInt(int64(idx)))
	}

	out := n.Bytes()
	if len(out) > 16 {
		return nil, ErrBinaryCodecFormat
	}

	padded := make([]byte, 16)
	copy(padded[16-len(out):], out)

	return padded, nil
}
