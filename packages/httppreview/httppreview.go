package httppreview

import "net/http"

// MessageProvider is the minimal interface for checking if validation produced
// any error messages. Both validation.Validator and validation.MessageBag
// satisfy this interface.
type MessageProvider interface {
	IsEmpty() bool
}

// AfterValidationHook returns a closure that should be registered as an "after"
// validation callback. When called, it checks whether validation produced no
// errors (messages are empty) AND the request carries a
// HTTPPreview-Validate-Only header. If both conditions are met it panics with
// [SuccessResponse] to short-circuit the handler chain. The middleware recovers
// this panic and writes a 204 response.
//
// If validation failed or the request has no HTTPPreview-Validate-Only header,
// the closure is a no-op.
//
// Ref: @bedrock/code-0218
//
//	hook := httppreview.AfterValidationHook(r)
//	// Register as after-validation callback:
//	hook(validator.Errors())
func AfterValidationHook(r *http.Request) func(MessageProvider) {
	return func(messages MessageProvider) {
		if messages.IsEmpty() && r.Header.Get("HTTPPreview-Validate-Only") != "" {
			panic(SuccessResponse{})
		}
	}
}
