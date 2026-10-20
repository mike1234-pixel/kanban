package handlers

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

// assertError checks that a response uses the standard JSON error format with
// the expected status and code.
func assertError(t *testing.T, response *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status = %d, want %d", response.Code, status)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	var body errorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if body.Error.Code != code {
		t.Errorf("error code = %q, want %q", body.Error.Code, code)
	}
	if body.Error.Message == "" {
		t.Error("error message is empty")
	}
}

func TestWriteError(t *testing.T) {
	response := httptest.NewRecorder()

	writeError(response, 404, codeBoardNotFound, "Board not found")

	if got, want := response.Body.String(), `{"error":{"code":"BOARD_NOT_FOUND","message":"Board not found"}}`+"\n"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestWriteInternalErrorHidesDetails(t *testing.T) {
	request := httptest.NewRequest("GET", "/boards", nil)
	response := httptest.NewRecorder()

	writeInternalError(response, request, errors.New("pq: password authentication failed"))

	if strings.Contains(response.Body.String(), "pq:") {
		t.Fatalf("response leaked internal error: %s", response.Body.String())
	}
	assertError(t, response, 500, codeInternalError)
}
