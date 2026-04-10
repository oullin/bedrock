package routing

// ResourceHandlers holds handler functions for standard CRUD routes.
// Nil fields are skipped during route registration.
type ResourceHandlers struct {
	Index   HandlerFunc // GET    /resource
	Show    HandlerFunc // GET    /resource/{id}
	Store   HandlerFunc // POST   /resource
	Update  HandlerFunc // PUT    /resource/{id}
	Destroy HandlerFunc // DELETE /resource/{id}
}
