package contracts

// Ref: @bedrock/code-0294
type ControllerDispatcher interface {
	Dispatch(route any, controller any, method string) (any, error)
	GetMiddleware(controller any, method string) []any
}
