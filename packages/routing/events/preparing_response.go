package events

// PreparingResponse is dispatched immediately before a response is converted
// Ref: @bedrock/code-0302
type PreparingResponse struct {
	Request  any // httpx.Request
	Response any // any value the route handler returned
}
