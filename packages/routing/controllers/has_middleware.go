package controllers

// HasMiddleware is the interface a controller implements to declare its
// middleware via the modern non-attribute API.
//
// Ref: @bedrock/code-0299
// declares Middleware() as static; Go has no static methods so the interface
// is a regular instance method that the dispatcher invokes after constructing
// the controller.
type HasMiddleware interface {
	Middleware() []Middleware
}
