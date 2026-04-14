package contracts

// ControllerDispatcher mirrors
// Illuminate\Routing\Contracts\ControllerDispatcher.
type ControllerDispatcher interface {
	Dispatch(route any, controller any, method string) (any, error)
	GetMiddleware(controller any, method string) []any
}
