// Package httppreview provides core functionality for httppreview.
//
// Ref: @bedrock/code-0186
// It provides middleware and utilities for handling precognitive HTTP requests
// — live, real-time form validation without duplicating backend validation
// rules in frontend code.
//
// When a precognitive request arrives (HTTPPreview: true header) the
// middleware executes route middleware and resolves controller dependencies
// (triggering validation) but does NOT execute the controller method. If
// validation passes it returns 204 No Content with a HTTPPreview-Success
// header; if validation fails the 422 response with errors is forwarded.
//
// Quick start:
//
//	mw := httppreview.New()
//	handler := mw.Wrap(myHandler)
package httppreview
