package middleware

import (
	"github.com/bedrock/packages/routing/exceptions"
)

// SignatureValidator is the surface ValidateSignature needs from a request.
//
// Mirrors httpx.Request's HasValidSignatureWhileIgnoring narrowly. The real
// httpx request will satisfy this in M11.
type SignatureValidator interface {
	HasValidSignatureWhileIgnoring(ignore []string, absolute bool) bool
}

// ValidateSignature is the middleware form of [routing.UrlGenerator]'s signed
// URL check. It rejects requests whose signature does not match.
//
// Mirrors @bedrock\Routing\Middleware\ValidateSignature.
type ValidateSignature struct {
	// Ignore lists parameter names that should be skipped when computing the
	// HMAC input. Useful for unrelated tracking parameters.
	Ignore []string
	// Relative flips the validation to relative-URL mode.
	Relative bool
}

// NeverValidate is the global skip list, mirroring the static $neverValidate
// property on the PHP class.
var NeverValidate []string

// Except adds parameters to the global skip list.
func Except(parameters ...string) {
	for _, p := range parameters {
		seen := false

		for _, n := range NeverValidate {
			if n == p {
				seen = true

				break
			}
		}

		if !seen {
			NeverValidate = append(NeverValidate, p)
		}
	}
}

// Handle runs the middleware: returns the next chain output if the signature
// is valid, otherwise returns an [*exceptions.InvalidSignatureException].
func (v *ValidateSignature) Handle(request SignatureValidator, next func(any) any) (any, error) {
	ignore := append([]string(nil), v.Ignore...)
	ignore = append(ignore, NeverValidate...)

	if request.HasValidSignatureWhileIgnoring(ignore, !v.Relative) {
		return next(request), nil
	}

	return nil, &exceptions.InvalidSignatureException{}
}

// Relative returns the middleware-spec string used to register a relative
// validator with optional ignored parameters.
//
// In the upstream framework this is consumed by the middleware-name resolver to produce the
// "ValidateSignature:relative,foo,bar" syntax. The Go form is a thin helper
// for code-generated registrations.
func Relative(ignore ...string) string {
	out := "routing.middleware.ValidateSignature:relative"

	for _, p := range ignore {
		out += "," + p
	}

	return out
}

// Absolute is the absolute-URL counterpart of [Relative].
func Absolute(ignore ...string) string {
	if len(ignore) == 0 {
		return "routing.middleware.ValidateSignature"
	}

	out := "routing.middleware.ValidateSignature:"

	for i, p := range ignore {
		if i > 0 {
			out += ","
		}

		out += p
	}

	return out
}
