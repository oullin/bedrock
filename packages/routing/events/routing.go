package events

// Routing is dispatched at the start of [routing.Router.Dispatch], before any
// route matching occurs. Listeners may inspect the incoming request.
//
// Ref: @bedrock/code-0305
type Routing struct {
	Request any // httpx.Request
}
