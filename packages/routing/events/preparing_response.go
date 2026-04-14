package events

// PreparingResponse is dispatched immediately before a response is converted
// (prepared) for sending. Mirrors Illuminate\Routing\Events\PreparingResponse.
type PreparingResponse struct {
	Request  any // httpx.Request
	Response any // any value the route handler returned
}
