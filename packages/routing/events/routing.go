package events

// Routing is dispatched at the start of [routing.Router.Dispatch], before any
// route matching occurs. Listeners may inspect the incoming request.
//
// Mirrors @bedrock\Routing\Events\Routing.
type Routing struct {
	Request any // httpx.Request
}
