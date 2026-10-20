package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kanban/internal/models"
	"kanban/internal/store"

	"github.com/go-playground/validator/v10"
)

// testStore overrides only the operations exercised by a test. Embedding the
// interface keeps the fake focused; an unexpected store call will still fail
// immediately with a nil-interface panic.
type testStore struct {
	store.Store
	saveBoard  func(models.Board) error
	getBoard   func(string) (models.Board, error)
	getColumns func(string) ([]models.Column, error)
	moveCard   func(string, string) (bool, error)
	saveUser   func(models.User) error
	getUser    func(string) (models.User, error)
	updateUser func(models.User) (bool, error)
}

func (s testStore) SaveUser(_ context.Context, user models.User) (time.Time, error) {
	return testCreatedAt, s.saveUser(user)
}

func (s testStore) GetUserByID(_ context.Context, id string) (models.User, error) {
	return s.getUser(id)
}

func (s testStore) UpdateUser(_ context.Context, user models.User) (bool, error) {
	return s.updateUser(user)
}

func (s testStore) SaveBoard(_ context.Context, board models.Board) (time.Time, error) {
	return testCreatedAt, s.saveBoard(board)
}

// testCreatedAt stands in for the timestamp the database assigns on insert.
var testCreatedAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

func (s testStore) GetBoardByID(_ context.Context, id string) (models.Board, error) {
	return s.getBoard(id)
}

func (s testStore) GetColumnsWithCards(_ context.Context, boardID string) ([]models.Column, error) {
	return s.getColumns(boardID)
}

func (s testStore) MoveCard(_ context.Context, id, columnID string) (bool, error) {
	return s.moveCard(id, columnID)
}

func testValidator() *validator.Validate {
	return validator.New()
}

func TestCreateBoard(t *testing.T) {
	t.Run("creates a valid board", func(t *testing.T) {
		called := false
		handler := BoardHandler{
			Store: testStore{saveBoard: func(board models.Board) error {
				called = true
				if board.Title != "Sprint planning" {
					t.Fatalf("saved title = %q, want Sprint planning", board.Title)
				}
				return nil
			}},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/boards", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","title":"Sprint planning","columns":[]}`))
		response := httptest.NewRecorder()

		handler.CreateBoard(response, request)

		if response.Code != 201 || !called {
			t.Fatalf("status = %d, store called = %v; want 201 and true", response.Code, called)
		}
		if got := response.Header().Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		var created models.Board
		if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if !created.CreatedAt.Equal(testCreatedAt) {
			t.Errorf("created_at = %v, want %v", created.CreatedAt, testCreatedAt)
		}
	})

	t.Run("rejects invalid JSON", func(t *testing.T) {
		handler := BoardHandler{Store: testStore{}, Validate: testValidator()}
		request := httptest.NewRequest("POST", "/boards", strings.NewReader("{"))
		response := httptest.NewRecorder()

		handler.CreateBoard(response, request)

		assertError(t, response, 400, codeInvalidJSON)
	})

	t.Run("rejects invalid board fields", func(t *testing.T) {
		called := false
		handler := BoardHandler{
			Store:    testStore{saveBoard: func(models.Board) error { called = true; return nil }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/boards", strings.NewReader(`{"id":"not-a-uuid","title":"x"}`))
		response := httptest.NewRecorder()

		handler.CreateBoard(response, request)

		if called {
			t.Fatal("store was called for an invalid board")
		}
		assertError(t, response, 400, codeValidationFailed)
	})

	t.Run("returns an internal error when persistence fails", func(t *testing.T) {
		handler := BoardHandler{
			Store:    testStore{saveBoard: func(models.Board) error { return errors.New("database unavailable") }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/boards", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","title":"Sprint planning"}`))
		response := httptest.NewRecorder()

		handler.CreateBoard(response, request)

		if strings.Contains(response.Body.String(), "database unavailable") {
			t.Fatalf("response leaked internal error: %s", response.Body.String())
		}
		assertError(t, response, 500, codeInternalError)
	})
}

func TestGetBoard(t *testing.T) {
	handler := BoardHandler{
		Store: testStore{getBoard: func(id string) (models.Board, error) {
			if id != "board-id" {
				t.Fatalf("requested id = %q, want board-id", id)
			}
			return models.Board{}, nil
		}},
		Validate: testValidator(),
	}
	request := httptest.NewRequest("GET", "/boards/board-id", nil)
	request.SetPathValue("id", "board-id")
	response := httptest.NewRecorder()

	handler.GetBoard(response, request)

	assertError(t, response, 404, codeBoardNotFound)
}

func TestListColumnsRequiresBoardID(t *testing.T) {
	handler := ColumnHandler{Store: testStore{}, Validate: testValidator()}
	request := httptest.NewRequest("GET", "/columns", nil)
	response := httptest.NewRecorder()

	handler.ListColumns(response, request)

	assertError(t, response, 400, codeBoardIDRequired)
}

func TestMoveCard(t *testing.T) {
	t.Run("moves a card", func(t *testing.T) {
		called := false
		handler := CardHandler{
			Store: testStore{moveCard: func(id, columnID string) (bool, error) {
				called = true
				if id != "card-id" || columnID != "a5fc3aed-6138-4e57-a8ec-3059d343eb83" {
					t.Fatalf("move arguments = (%q, %q), unexpected", id, columnID)
				}
				return true, nil
			}},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/cards/card-id/move", strings.NewReader(`{"column_id":"a5fc3aed-6138-4e57-a8ec-3059d343eb83"}`))
		request.SetPathValue("id", "card-id")
		response := httptest.NewRecorder()

		handler.MoveCard(response, request)

		if response.Code != 204 || !called {
			t.Fatalf("status = %d, store called = %v; want 204 and true", response.Code, called)
		}
	})

	t.Run("rejects an invalid destination column", func(t *testing.T) {
		called := false
		handler := CardHandler{
			Store:    testStore{moveCard: func(string, string) (bool, error) { called = true; return true, nil }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/cards/card-id/move", strings.NewReader(`{"column_id":"invalid"}`))
		request.SetPathValue("id", "card-id")
		response := httptest.NewRecorder()

		handler.MoveCard(response, request)

		if called {
			t.Fatal("store was called for an invalid destination column")
		}
		assertError(t, response, 400, codeValidationFailed)
	})

	t.Run("returns not found when the card does not exist", func(t *testing.T) {
		handler := CardHandler{
			Store:    testStore{moveCard: func(string, string) (bool, error) { return false, nil }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/cards/card-id/move", strings.NewReader(`{"column_id":"a5fc3aed-6138-4e57-a8ec-3059d343eb83"}`))
		request.SetPathValue("id", "card-id")
		response := httptest.NewRecorder()

		handler.MoveCard(response, request)

		assertError(t, response, 404, codeCardNotFound)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("creates a valid user", func(t *testing.T) {
		called := false
		handler := UserHandler{
			Store: testStore{saveUser: func(user models.User) error {
				called = true
				if user.Email != "ada@example.com" {
					t.Fatalf("saved email = %q, want ada@example.com", user.Email)
				}
				return nil
			}},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/users", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","name":"Ada","email":"ada@example.com"}`))
		response := httptest.NewRecorder()

		handler.CreateUser(response, request)

		if response.Code != 201 || !called {
			t.Fatalf("status = %d, store called = %v; want 201 and true", response.Code, called)
		}
		var created models.User
		if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if !created.CreatedAt.Equal(testCreatedAt) {
			t.Errorf("created_at = %v, want %v", created.CreatedAt, testCreatedAt)
		}
	})

	t.Run("rejects an invalid email", func(t *testing.T) {
		called := false
		handler := UserHandler{
			Store:    testStore{saveUser: func(models.User) error { called = true; return nil }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/users", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","name":"Ada","email":"not-an-email"}`))
		response := httptest.NewRecorder()

		handler.CreateUser(response, request)

		if called {
			t.Fatal("store was called for an invalid user")
		}
		assertError(t, response, 400, codeValidationFailed)
	})

	t.Run("returns conflict when the email is taken", func(t *testing.T) {
		handler := UserHandler{
			Store:    testStore{saveUser: func(models.User) error { return store.ErrEmailTaken }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/users", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","name":"Ada","email":"ada@example.com"}`))
		response := httptest.NewRecorder()

		handler.CreateUser(response, request)

		assertError(t, response, 409, codeUserEmailTaken)
	})
}

func TestUpdateUser(t *testing.T) {
	existing := func(string) (models.User, error) {
		return models.User{ID: "user-id", Name: "Ada", Email: "ada@example.com"}, nil
	}

	t.Run("updates name and email", func(t *testing.T) {
		handler := UserHandler{
			Store: testStore{getUser: existing, updateUser: func(user models.User) (bool, error) {
				if user.ID != "user-id" || user.Name != "Ada Lovelace" || user.Email != "ada@lovelace.dev" {
					t.Fatalf("updated user = %+v, unexpected", user)
				}
				return true, nil
			}},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("PUT", "/users/user-id", strings.NewReader(`{"name":"Ada Lovelace","email":"ada@lovelace.dev"}`))
		request.SetPathValue("id", "user-id")
		response := httptest.NewRecorder()

		handler.UpdateUser(response, request)

		if response.Code != 200 {
			t.Fatalf("status = %d, want 200", response.Code)
		}
	})

	t.Run("returns not found for a missing user", func(t *testing.T) {
		handler := UserHandler{
			Store:    testStore{getUser: func(string) (models.User, error) { return models.User{}, nil }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("PUT", "/users/user-id", strings.NewReader(`{"name":"Ada","email":"ada@example.com"}`))
		request.SetPathValue("id", "user-id")
		response := httptest.NewRecorder()

		handler.UpdateUser(response, request)

		assertError(t, response, 404, codeUserNotFound)
	})

	t.Run("returns conflict when the email is taken", func(t *testing.T) {
		handler := UserHandler{
			Store: testStore{getUser: existing, updateUser: func(models.User) (bool, error) {
				return false, store.ErrEmailTaken
			}},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("PUT", "/users/user-id", strings.NewReader(`{"name":"Ada","email":"taken@example.com"}`))
		request.SetPathValue("id", "user-id")
		response := httptest.NewRecorder()

		handler.UpdateUser(response, request)

		assertError(t, response, 409, codeUserEmailTaken)
	})
}
