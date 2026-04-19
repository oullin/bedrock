package prompts

// This file provides package-level documentation for the public API.
// The actual function implementations live in their respective files
// (text.go, password.go, confirm.go, etc.).
//
// All interactive prompts return (value, error). The error is:
//   - ErrCancelled when the user presses Ctrl+C
//   - ErrNonInteractive when the terminal is non-interactive
//   - ErrValidation when validation fails in non-interactive mode
//
// Display functions (Note, Error, Warning, Info, Alert, Intro, Outro,
// Table, Grid, Clear, Title) do not return errors.
//
// Configuration is done via the functional options pattern:
//
//	result, err := prompts.Text("Name?",
//	    prompts.TextWithPlaceholder("Enter your name"),
//	    prompts.TextWithRequired(true),
//	    prompts.TextWithHint("Your full name"),
//	)
