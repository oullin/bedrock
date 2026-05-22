package events

// ResponsePrepared is dispatched after a response has been prepared and is
// ready to be returned to the client.
// Ref: @bedrock/code-0303
type ResponsePrepared struct {
	Request  any // httpx.Request
	Response any // httpx.Response
}
