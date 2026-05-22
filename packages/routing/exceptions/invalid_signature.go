package exceptions

import "errors"

// InvalidSignatureException indicates that a signed URL's HMAC did not match
// or that its expiration timestamp has passed.
//
// Ref: @bedrock/code-0308
type InvalidSignatureException struct{ Reason string }

func (e *InvalidSignatureException) Error() string {
	if e.Reason == "" {
		return "the signature is invalid"
	}

	return "the signature is invalid: " + e.Reason
}

// ErrInvalidSignature is the sentinel returned for unwrapping convenience.
var ErrInvalidSignature = errors.New("invalid signature")

func (e *InvalidSignatureException) Unwrap() error { return ErrInvalidSignature }
