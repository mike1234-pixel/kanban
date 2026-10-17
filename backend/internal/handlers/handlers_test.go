package handlers

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

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
}

func (s testStore) SaveBoard(_ context.Context, board models.Board) error {
	return s.saveBoard(board)
}

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
	})

	t.Run("rejects invalid JSON", func(t *testing.T) {
		handler := BoardHandler{Store: testStore{}, Validate: testValidator()}
		request := httptest.NewRequest("POST", "/boards", strings.NewReader("{"))
		response := httptest.NewRecorder()

		handler.CreateBoard(response, request)

		if response.Code != 400 {
			t.Fatalf("status = %d, want 400", response.Code)
		}
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

		if response.Code != 400 || called {
			t.Fatalf("status = %d, store called = %v; want 400 and false", response.Code, called)
		}
	})

	t.Run("returns an internal error when persistence fails", func(t *testing.T) {
		handler := BoardHandler{
			Store:    testStore{saveBoard: func(models.Board) error { return errors.New("database unavailable") }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/boards", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","title":"Sprint planning"}`))
		response := httptest.NewRecorder()

		handler.CreateBoard(response, request)

		if response.Code != 500 {
			t.Fatalf("status = %d, want 500", response.Code)
		}
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

	if response.Code != 404 {
		t.Fatalf("status = %d, want 404", response.Code)
	}
}

func TestListColumnsRequiresBoardID(t *testing.T) {
	handler := ColumnHandler{Store: testStore{}, Validate: testValidator()}
	request := httptest.NewRequest("GET", "/columns", nil)
	response := httptest.NewRecorder()

	handler.ListColumns(response, request)

	if response.Code != 400 {
		t.Fatalf("status = %d, want 400", response.Code)
	}
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

		if response.Code != 400 || called {
			t.Fatalf("status = %d, store called = %v; want 400 and false", response.Code, called)
		}
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

		if response.Code != 404 {
			t.Fatalf("status = %d, want 404", response.Code)
		}
	})
}
