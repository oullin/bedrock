// Package precognition is a 1:1 Go port of laravel/framework 13.x
// src/@bedrock/Foundation/Http/Middleware/HandlePrecognitiveRequests and
// src/@bedrock/Foundation/Precognition.
//
// It provides middleware and utilities for handling precognitive HTTP requests
// — live, real-time form validation without duplicating backend validation
// rules in frontend code.
//
// When a precognitive request arrives (Precognition: true header) the
// middleware executes route middleware and resolves controller dependencies
// (triggering validation) but does NOT execute the controller method. If
// validation passes it returns 204 No Content with a Precognition-Success
// header; if validation fails the 422 response with errors is forwarded.
//
// Quick start:
//
//	mw := precognition.New()
//	handler := mw.Wrap(myHandler)
package precognition
