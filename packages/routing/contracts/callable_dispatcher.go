// Package contracts mirrors laravel/framework/src/Illuminate/Routing/Contracts.
//
// These interfaces are the binding points the routing service provider wires
// into the container so consumers can swap implementations.
package contracts

// CallableDispatcher mirrors
// Illuminate\Routing\Contracts\CallableDispatcher.
//
// The route argument is typed as `any` to break the import cycle between this
// package and the parent routing package; concrete dispatchers will narrow it
// to *routing.Route.
type CallableDispatcher interface {
	Dispatch(route any, callable any) (any, error)
}
