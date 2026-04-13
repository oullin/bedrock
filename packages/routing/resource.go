package routing

// ResourceHandlers holds handler functions for standard CRUD routes.
// Nil fields are skipped during route registration.
type ResourceHandlers struct {
	Index   HandlerFunc // GET    /resource
	Create  HandlerFunc // GET    /resource/create
	Store   HandlerFunc // POST   /resource
	Show    HandlerFunc // GET    /resource/{id}
	Edit    HandlerFunc // GET    /resource/{id}/edit
	Update  HandlerFunc // PUT    /resource/{id}
	Destroy HandlerFunc // DELETE /resource/{id}
}

// SingletonHandlers holds handler functions for singleton resource routes.
// Nil fields are skipped during route registration.
type SingletonHandlers struct {
	Show    HandlerFunc // GET    /resource
	Edit    HandlerFunc // GET    /resource/edit
	Update  HandlerFunc // PUT    /resource
	Destroy HandlerFunc // DELETE /resource
}
