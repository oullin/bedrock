// Package enums contains string-typed enumerations for the AI package,
// mirroring PHP backed enums from laravel/ai.
package enums

// Lab identifies an AI provider.
// Mirrors Laravel\Ai\Enums\Lab.
type Lab string

const (
	LabAnthropic  Lab = "anthropic"
	LabAzure      Lab = "azure"
	LabCohere     Lab = "cohere"
	LabDeepSeek   Lab = "deepseek"
	LabElevenLabs Lab = "elevenlabs"
	LabGemini     Lab = "gemini"
	LabGroq       Lab = "groq"
	LabJina       Lab = "jina"
	LabMistral    Lab = "mistral"
	LabOllama     Lab = "ollama"
	LabOpenAI     Lab = "openai"
	LabOpenRouter Lab = "openrouter"
	LabVoyageAI   Lab = "voyageai"
	LabXAI        Lab = "xai"
)

// String returns the raw string value of the Lab.
func (l Lab) String() string { return string(l) }

// Valid reports whether the Lab is one of the known provider values.
func (l Lab) Valid() bool {
	switch l {
	case LabAnthropic, LabAzure, LabCohere, LabDeepSeek, LabElevenLabs,
		LabGemini, LabGroq, LabJina, LabMistral, LabOllama,
		LabOpenAI, LabOpenRouter, LabVoyageAI, LabXAI:
		return true
	}

	return false
}
