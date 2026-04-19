// Package ai provides a unified, expressive API for interacting with AI providers
// such as OpenAI, Anthropic, Gemini, and more. It mirrors the laravel/ai (0.x)
// package, offering 100% functional parity adapted idiomatically to Go.
//
// The package follows the manager + provider + gateway layering pattern used
// throughout bedrock:
//
//   - Manager: resolves and caches provider instances by Lab (provider enum)
//   - Provider: implements one or more capability interfaces (Text, Image, Audio…)
//   - Gateway: the thin HTTP/SDK adapter used by the provider
//
// For testing, every capability has a fake gateway that records interactions and
// exposes assert helpers, activated via FakeText(), FakeImage(), FakeAudio(), etc.
package ai
