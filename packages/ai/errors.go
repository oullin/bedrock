package ai

import "errors"

var (
	// ErrProviderCapability is returned when a provider does not support the
	// requested capability (e.g. asking ElevenLabs for text generation).
	ErrProviderCapability = errors.New("ai: provider does not support this capability")

	// ErrUnsupportedProvider is returned when the manager cannot find a
	// factory for the requested provider Lab.
	ErrUnsupportedProvider = errors.New("ai: unsupported provider")

	// ErrStrayCall is returned when a fake gateway receives a call that was
	// not expected and PreventStray* mode is active.
	ErrStrayCall = errors.New("ai: unexpected call to faked provider (stray call prevention is active)")

	// ErrNoFakeResponses is returned when a fake gateway exhausts its response
	// queue and repeat mode is off.
	ErrNoFakeResponses = errors.New("ai: fake gateway has no more queued responses")
)
