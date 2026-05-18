package shared

import "net/http"

// ErrorResponse is the standard JSON error envelope for API routes.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteError writes a JSON error body with the given HTTP status.
func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, ErrorResponse{Error: message})
}
