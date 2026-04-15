package data

// TranscriptionSegment represents a timed segment of transcribed audio.
// Mirrors Upstream\Ai\Responses\Data\TranscriptionSegment.
type TranscriptionSegment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

// ToMap returns a map representation.
func (t TranscriptionSegment) ToMap() map[string]any {
	return map[string]any{"start": t.Start, "end": t.End, "text": t.Text}
}
