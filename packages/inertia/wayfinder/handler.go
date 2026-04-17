package wayfinder

import (
	"log"
	"net/http"
)

// Handler returns an http.Handler that serves the registry's routes
// as JSON. This is useful as a development endpoint for Vite plugins
// or other build tools.
func Handler(reg *Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := reg.ToJSON()

		if err != nil {
			http.Error(w, "failed to serialize routes", http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if _, err = w.Write(data); err != nil {
			log.Printf("wayfinder: failed to write response: %v", err)
		}
	})
}
