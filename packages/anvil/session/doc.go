// Package session provides a Laravel-inspired HTTP session manager.
//
// A Store manages the session lifecycle (start, save, regenerate,
// invalidate) and exposes an attribute bag with flash data and CSRF
// token support. The Handler interface abstracts the storage backend;
// ArrayHandler provides an in-memory implementation for testing.
//
//	h := session.NewArrayHandler()
//	s := session.New("app_session", h)
//	_ = s.Start(ctx)
//	s.Put("user_id", "42")
//	_ = s.Save(ctx)
package session
