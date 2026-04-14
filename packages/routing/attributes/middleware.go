// Package attributes mirrors
// upstream/framework/src/Framework/Routing/Attributes/Controllers.
//
// PHP 8 attributes have no direct Go equivalent. The strict-parity story here
// is "the value object is identical, but you attach it via the
// [controllers.HasMiddleware] interface instead of an attribute on the class".
// The struct itself is exposed for code that wants to construct a middleware
// declaration out-of-band (e.g. when porting a controller mechanically).
package attributes

import "github.com/bedrock/packages/routing/controllers"

// Middleware is a value object describing a controller-attached middleware.
// It is the parity counterpart of #[Middleware] in PHP.
type Middleware = controllers.Middleware
