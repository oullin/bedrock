package httppreview

import "net/http"

// SuccessResponse is a sentinel value used to signal that precognitive
// validation passed and the handler should short-circuit. The middleware
// recovers this panic and writes a 204 response. This is the Go equivalent
// of the upstream abort(204, headers: ['HTTPPreview-Success' => 'true']).
type SuccessResponse struct{}

// WriteSuccessResponse writes a 204 No Content response with the
// HTTPPreview-Success header. This centralises the httppreview success
// response format.
func WriteSuccessResponse(w http.ResponseWriter) {
	w.Header().Set("HTTPPreview-Success", "true")
	w.WriteHeader(http.StatusNoContent)
}

// AddVaryHeader appends "HTTPPreview" to the Vary response header. Uses Add
// rather than Set to preserve any existing Vary values.
func AddVaryHeader(w http.ResponseWriter) {
	w.Header().Add("Vary", "HTTPPreview")
}

// AddHTTPPreviewHeader sets the HTTPPreview response header to "true",
// indicating that the response was produced by a precognitive request.
func AddHTTPPreviewHeader(w http.ResponseWriter) {
	w.Header().Set("HTTPPreview", "true")
}
