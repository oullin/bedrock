package responses

import (
	"encoding/json"

	"github.com/bedrock/packages/ai/sdk/data"
)

// StructuredAgentResponse extends AgentResponse with parsed structured output.
type StructuredAgentResponse struct {
	AgentResponse
	Data map[string]any
}

// NewStructuredAgentResponse constructs a StructuredAgentResponse.
func NewStructuredAgentResponse(invocationID string, data_ map[string]any, text string, usage data.Usage, meta data.Meta) *StructuredAgentResponse {
	return &StructuredAgentResponse{
		AgentResponse: *NewAgentResponse(invocationID, text, usage, meta),
		Data:          data_,
	}
}

// ToMap returns the structured data as a map.
func (r *StructuredAgentResponse) ToMap() map[string]any { return r.Data }

// MarshalJSON serialises the structured data.
func (r *StructuredAgentResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Data)
}
