package events

// ResponsePrepared is dispatched after a response has been prepared and is
// ready to be returned to the client. Mirrors
// Illuminate\Routing\Events\ResponsePrepared.
type ResponsePrepared struct {
	Request  any // httpx.Request
	Response any // httpx.Response
}
