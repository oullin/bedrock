package precognition

import "net/http"

// SuccessResponse is a sentinel value used to signal that precognitive
// validation passed and the handler should short-circuit. The middleware
// recovers this panic and writes a 204 response. This is the Go equivalent
// of Laravel's abort(204, headers: ['Precognition-Success' => 'true']).
type SuccessResponse struct{}

// WriteSuccessResponse writes a 204 No Content response with the
// Precognition-Success header. This centralises the precognition success
// response format.
func WriteSuccessResponse(w http.ResponseWriter) {
	w.Header().Set("Precognition-Success", "true")
	w.WriteHeader(http.StatusNoContent)
}

// AddVaryHeader appends "Precognition" to the Vary response header. Uses Add
// rather than Set to preserve any existing Vary values.
func AddVaryHeader(w http.ResponseWriter) {
	w.Header().Add("Vary", "Precognition")
}

// AddPrecognitionHeader sets the Precognition response header to "true",
// indicating that the response was produced by a precognitive request.
func AddPrecognitionHeader(w http.ResponseWriter) {
	w.Header().Set("Precognition", "true")
}
