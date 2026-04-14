package translation

import "errors"

var (
	// ErrMalformedJSON is returned by FileLoader when a JSON translation file
	// cannot be decoded.
	ErrMalformedJSON = errors.New("translation: malformed JSON file")

	// ErrInvalidLocale is returned by Translator.SetLocale when the locale
	// string contains path separators.
	ErrInvalidLocale = errors.New("translation: locale must not contain path separators")
)
