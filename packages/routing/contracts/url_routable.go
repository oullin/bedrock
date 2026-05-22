package contracts

// Ref: @bedrock/code-0198
// User-defined model types implement this interface to participate in
// implicit and scoped route-model binding. The Go port routes binding
// resolution exclusively through this interface, with no Orm dependency.
//
//   - GetRouteKey:     value used in URL generation (e.g. "id" or "slug")
//   - GetRouteKeyName: column or field name the binding resolves against
//   - ResolveRouteBinding: look up an instance by value/field
//   - ResolveChildRouteBinding: scoped lookup off a parent instance
//
// Returning (nil, nil) from the resolve methods signals "not found" without
// an error, matching the upstream nullable model lookup semantics.
type UrlRoutable interface {
	GetRouteKey() any
	GetRouteKeyName() string
	ResolveRouteBinding(value, field string) (any, error)
	ResolveChildRouteBinding(childType, value, field string) (any, error)
}
