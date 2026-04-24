package handler

import (
	"net/http"

	"github.com/bedrock/packages/httpx"
)

func jsonResponse(w http.ResponseWriter, status int, payload any) {
	_ = httpx.NewJsonResponse(w, payload, status).Send()
}

func noContent(w http.ResponseWriter) {
	_ = httpx.NewResponse(w).NoContent()
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	_ = httpx.NewResponse(w).Status(status).SendString(message + "\n")
}

func validationErrors(w http.ResponseWriter, errors map[string][]string) {
	jsonResponse(w, http.StatusUnprocessableEntity, map[string]any{"errors": errors})
}
