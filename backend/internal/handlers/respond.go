package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Error codes returned in the "code" field of error responses. Clients can
// branch on these, so treat them as part of the API: don't rename them lightly.
const (
	codeInvalidJSON      = "INVALID_JSON"
	codeValidationFailed = "VALIDATION_FAILED"
	codeInternalError    = "INTERNAL_ERROR"

	codeBoardNotFound   = "BOARD_NOT_FOUND"
	codeBoardIDRequired = "BOARD_ID_REQUIRED"

	codeColumnNotFound   = "COLUMN_NOT_FOUND"
	codeColumnIDRequired = "COLUMN_ID_REQUIRED"

	codeCardNotFound = "CARD_NOT_FOUND"

	codeUserNotFound   = "USER_NOT_FOUND"
	codeUserIDRequired = "USER_ID_REQUIRED"
	codeUserEmailTaken = "USER_EMAIL_TAKEN"
)

// errorResponse is the body of every error response:
// {"error": {"code": "...", "message": "..."}}
type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeJSON writes v as a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	// Headers must be set before WriteHeader; anything set afterwards is ignored.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// The status line has already been sent, so the client can't be told; just record it.
		slog.Error("failed to encode response", "err", err)
	}
}

// writeError writes an error response in the API's standard error format.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: errorBody{Code: code, Message: message}})
}

// writeInternalError logs the real error for debugging and sends the client a
// generic 500, so internal details (SQL, driver messages) never leak out.
func writeInternalError(w http.ResponseWriter, r *http.Request, err error) {
	slog.ErrorContext(r.Context(), "request failed", "method", r.Method, "path", r.URL.Path, "err", err)
	writeError(w, http.StatusInternalServerError, codeInternalError, "Something went wrong")
}
