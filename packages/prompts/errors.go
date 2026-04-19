package prompts

import "errors"

var (
	// ErrCancelled is returned when the user presses Ctrl+C.
	ErrCancelled = errors.New("prompts: cancelled")

	// ErrNonInteractive is returned when the terminal is non-interactive
	// and no default value is available.
	ErrNonInteractive = errors.New("prompts: non-interactive terminal")

	// ErrValidation is returned when the default value fails validation
	// in non-interactive mode.
	ErrValidation = errors.New("prompts: validation failed")

	// ErrRequired is returned when a required prompt receives an empty value.
	ErrRequired = errors.New("prompts: value required")

	// ErrInvalidOptions is returned when options have an unsupported type.
	ErrInvalidOptions = errors.New("prompts: invalid options type")
)
