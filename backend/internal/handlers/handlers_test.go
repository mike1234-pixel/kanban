package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kanban/internal/models"
	"kanban/internal/store"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// testStore overrides only the operations exercised by a test. Embedding the
// interface keeps the fake focused; an unexpected store call will still fail
// immediately with a nil-interface panic.
type testStore struct {
	store.Store
	saveBoard  func(models.Board) error
	getBoard   func(id, ownerID string) (models.Board, error)
	getColumns func(boardID, ownerID string) ([]models.Column, error)
	moveCard   func(id, columnID, ownerID string) (bool, error)
	saveUser   func(models.User) error
	getUser    func(string) (models.User, error)
	updateUser func(models.User) (bool, error)

	getUserByEmail func(string) (models.User, error)
	createSession  func(tokenHash, userID string) error
	getSessionUser func(tokenHash string) (models.User, error)
	deleteSession  func(tokenHash string) error

	updateBoard func(models.Board) (bool, error)

	saveCard func(card models.Card, ownerID string) (bool, error)

	saveColumn   func(col models.Column, ownerID string) (bool, error)
	getColumn    func(id, ownerID string) (models.Column, error)
	updateColumn func(col models.Column, ownerID string) (bool, error)
}

func (s testStore) SaveCard(_ context.Context, card models.Card, ownerID string) (bool, error) {
	return s.saveCard(card, ownerID)
}

func (s testStore) SaveColumn(_ context.Context, col models.Column, ownerID string) (bool, error) {
	return s.saveColumn(col, ownerID)
}

func (s testStore) GetColumnByID(_ context.Context, id, ownerID string) (models.Column, error) {
	return s.getColumn(id, ownerID)
}

func (s testStore) UpdateColumn(_ context.Context, col models.Column, ownerID string) (bool, error) {
	return s.updateColumn(col, ownerID)
}

func (s testStore) UpdateBoard(_ context.Context, board models.Board) (bool, error) {
	return s.updateBoard(board)
}

// testUser is the logged-in user for handlers behind RequireAuth.
var testUser = models.User{ID: "9a1b2c3d-0000-4000-8000-000000000001", Name: "Ada", Email: "ada@example.com"}

// asUser returns a copy of r whose context holds testUser, as RequireAuth would.
func asUser(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userKey, testUser))
}

func (s testStore) DeleteSession(_ context.Context, tokenHash string) error {
	return s.deleteSession(tokenHash)
}

func (s testStore) GetSessionUser(_ context.Context, tokenHash string) (models.User, error) {
	return s.getSessionUser(tokenHash)
}

func (s testStore) GetUserByEmail(_ context.Context, email string) (models.User, error) {
	return s.getUserByEmail(email)
}

func (s testStore) CreateSession(_ context.Context, tokenHash, userID string, _ time.Time) error {
	return s.createSession(tokenHash, userID)
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

func (s testStore) GetBoardByID(_ context.Context, id, ownerID string) (models.Board, error) {
	return s.getBoard(id, ownerID)
}

func (s testStore) GetColumnsWithCards(_ context.Context, boardID, ownerID string) ([]models.Column, error) {
	return s.getColumns(boardID, ownerID)
}

func (s testStore) MoveCard(_ context.Context, id, columnID, ownerID string) (bool, error) {
	return s.moveCard(id, columnID, ownerID)
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

		handler.CreateBoard(response, asUser(request))

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

		handler.CreateBoard(response, asUser(request))

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

		handler.CreateBoard(response, asUser(request))

		if called {
			t.Fatal("store was called for an invalid board")
		}
		assertError(t, response, 400, codeValidationFailed)
	})

	t.Run("returns conflict when the board ID already exists", func(t *testing.T) {
		handler := BoardHandler{
			Store:    testStore{saveBoard: func(models.Board) error { return store.ErrBoardExists }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/boards", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","title":"Sprint planning"}`))
		response := httptest.NewRecorder()

		handler.CreateBoard(response, asUser(request))

		assertError(t, response, 409, codeBoardExists)
	})

	t.Run("returns an internal error when persistence fails", func(t *testing.T) {
		handler := BoardHandler{
			Store:    testStore{saveBoard: func(models.Board) error { return errors.New("database unavailable") }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/boards", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","title":"Sprint planning"}`))
		response := httptest.NewRecorder()

		handler.CreateBoard(response, asUser(request))

		if strings.Contains(response.Body.String(), "database unavailable") {
			t.Fatalf("response leaked internal error: %s", response.Body.String())
		}
		assertError(t, response, 500, codeInternalError)
	})
}

func TestGetBoard(t *testing.T) {
	handler := BoardHandler{
		Store: testStore{getBoard: func(id, ownerID string) (models.Board, error) {
			if id != "board-id" || ownerID != testUser.ID {
				t.Fatalf("requested id = %q, owner = %q; want board-id and the logged-in user", id, ownerID)
			}
			return models.Board{}, nil // "not found" and "not yours" look the same
		}},
		Validate: testValidator(),
	}
	request := httptest.NewRequest("GET", "/boards/board-id", nil)
	request.SetPathValue("id", "board-id")
	response := httptest.NewRecorder()

	handler.GetBoard(response, asUser(request))

	assertError(t, response, 404, codeBoardNotFound)
}

func TestListColumnsRequiresBoardID(t *testing.T) {
	handler := ColumnHandler{Store: testStore{}, Validate: testValidator()}
	request := httptest.NewRequest("GET", "/columns", nil)
	response := httptest.NewRecorder()

	handler.ListColumns(response, asUser(request))

	assertError(t, response, 400, codeBoardIDRequired)
}

func TestMoveCard(t *testing.T) {
	t.Run("moves a card", func(t *testing.T) {
		called := false
		handler := CardHandler{
			Store: testStore{moveCard: func(id, columnID, ownerID string) (bool, error) {
				called = true
				if id != "card-id" || columnID != "a5fc3aed-6138-4e57-a8ec-3059d343eb83" || ownerID != testUser.ID {
					t.Fatalf("move arguments = (%q, %q, %q), unexpected", id, columnID, ownerID)
				}
				return true, nil
			}},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/cards/card-id/move", strings.NewReader(`{"column_id":"a5fc3aed-6138-4e57-a8ec-3059d343eb83"}`))
		request.SetPathValue("id", "card-id")
		response := httptest.NewRecorder()

		handler.MoveCard(response, asUser(request))

		if response.Code != 204 || !called {
			t.Fatalf("status = %d, store called = %v; want 204 and true", response.Code, called)
		}
	})

	t.Run("rejects an invalid destination column", func(t *testing.T) {
		called := false
		handler := CardHandler{
			Store:    testStore{moveCard: func(string, string, string) (bool, error) { called = true; return true, nil }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/cards/card-id/move", strings.NewReader(`{"column_id":"invalid"}`))
		request.SetPathValue("id", "card-id")
		response := httptest.NewRecorder()

		handler.MoveCard(response, asUser(request))

		if called {
			t.Fatal("store was called for an invalid destination column")
		}
		assertError(t, response, 400, codeValidationFailed)
	})

	t.Run("returns not found when the card does not exist", func(t *testing.T) {
		handler := CardHandler{
			Store:    testStore{moveCard: func(string, string, string) (bool, error) { return false, nil }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/cards/card-id/move", strings.NewReader(`{"column_id":"a5fc3aed-6138-4e57-a8ec-3059d343eb83"}`))
		request.SetPathValue("id", "card-id")
		response := httptest.NewRecorder()

		handler.MoveCard(response, asUser(request))

		assertError(t, response, 404, codeCardNotFound)
	})

	t.Run("returns column not found when the destination isn't the user's", func(t *testing.T) {
		handler := CardHandler{
			Store:    testStore{moveCard: func(string, string, string) (bool, error) { return false, store.ErrColumnNotFound }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/cards/card-id/move", strings.NewReader(`{"column_id":"a5fc3aed-6138-4e57-a8ec-3059d343eb83"}`))
		request.SetPathValue("id", "card-id")
		response := httptest.NewRecorder()

		handler.MoveCard(response, asUser(request))

		assertError(t, response, 404, codeColumnNotFound)
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
				if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("correct-horse")); err != nil {
					t.Fatalf("saved hash doesn't match the password: %v", err)
				}
				return nil
			}},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/users", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","name":"Ada","email":"ada@example.com","password":"correct-horse"}`))
		response := httptest.NewRecorder()

		handler.CreateUser(response, request)

		if response.Code != 201 || !called {
			t.Fatalf("status = %d, store called = %v; want 201 and true", response.Code, called)
		}
		if body := response.Body.String(); strings.Contains(body, "password") || strings.Contains(body, "$2a$") {
			t.Fatalf("response leaks the password or its hash: %s", body)
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
		request := httptest.NewRequest("POST", "/users", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","name":"Ada","email":"ada@example.com","password":"correct-horse"}`))
		response := httptest.NewRecorder()

		handler.CreateUser(response, request)

		assertError(t, response, 409, codeUserEmailTaken)
	})

	t.Run("rejects a password shorter than 8 characters", func(t *testing.T) {
		called := false
		handler := UserHandler{
			Store:    testStore{saveUser: func(models.User) error { called = true; return nil }},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("POST", "/users", strings.NewReader(`{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","name":"Ada","email":"ada@example.com","password":"short"}`))
		response := httptest.NewRecorder()

		handler.CreateUser(response, request)

		if called {
			t.Fatal("store was called with a too-short password")
		}
		assertError(t, response, 400, codeValidationFailed)
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

func TestLogin(t *testing.T) {
	const userID = "4c85d3f2-e917-4a3f-929a-11b02499a33f"
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-horse"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	ada := models.User{ID: userID, Name: "Ada", Email: "ada@example.com", PasswordHash: string(hash)}

	login := func(store testStore, body string) *httptest.ResponseRecorder {
		handler := AuthHandler{Store: store, Validate: testValidator()}
		request := httptest.NewRequest("POST", "/login", strings.NewReader(body))
		response := httptest.NewRecorder()
		handler.Login(response, request)
		return response
	}

	t.Run("sets an HttpOnly session cookie and stores only its hash", func(t *testing.T) {
		var storedHash, storedUserID string
		response := login(testStore{
			getUserByEmail: func(string) (models.User, error) { return ada, nil },
			createSession: func(tokenHash, userID string) error {
				storedHash, storedUserID = tokenHash, userID
				return nil
			},
		}, `{"email":"ada@example.com","password":"correct-horse"}`)

		if response.Code != 200 {
			t.Fatalf("status = %d, want 200", response.Code)
		}
		cookies := response.Result().Cookies()
		if len(cookies) != 1 || cookies[0].Name != sessionCookieName {
			t.Fatalf("cookies = %v, want one %q cookie", cookies, sessionCookieName)
		}
		cookie := cookies[0]
		if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode {
			t.Errorf("cookie flags: HttpOnly=%v Secure=%v SameSite=%v; want true, true, Lax", cookie.HttpOnly, cookie.Secure, cookie.SameSite)
		}
		if storedHash == cookie.Value {
			t.Error("the store received the raw session token; it should only get the hash")
		}
		if storedHash != hashToken(cookie.Value) {
			t.Error("stored hash doesn't match the cookie's token")
		}
		if storedUserID != userID {
			t.Errorf("session user = %q, want %q", storedUserID, userID)
		}
		if strings.Contains(response.Body.String(), "$2a$") {
			t.Error("response leaks the password hash")
		}
	})

	t.Run("rejects a wrong password", func(t *testing.T) {
		response := login(testStore{
			getUserByEmail: func(string) (models.User, error) { return ada, nil },
			createSession: func(string, string) error {
				t.Fatal("session created for a wrong password")
				return nil
			},
		}, `{"email":"ada@example.com","password":"wrong-password"}`)

		assertError(t, response, 401, codeInvalidCredentials)
	})

	t.Run("rejects an unknown email with the same error", func(t *testing.T) {
		response := login(testStore{
			getUserByEmail: func(string) (models.User, error) { return models.User{}, nil },
		}, `{"email":"nobody@example.com","password":"correct-horse"}`)

		assertError(t, response, 401, codeInvalidCredentials)
	})
}

func TestRequireAuth(t *testing.T) {
	ada := models.User{ID: "4c85d3f2-e917-4a3f-929a-11b02499a33f", Name: "Ada", Email: "ada@example.com"}

	// next records whether it ran, and which user it found in the context.
	var nextCalled bool
	var gotUser models.User
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		gotUser, _ = userFromContext(r.Context())
	})

	run := func(store testStore, cookie *http.Cookie) *httptest.ResponseRecorder {
		nextCalled, gotUser = false, models.User{}
		handler := AuthHandler{Store: store}
		request := httptest.NewRequest("GET", "/me", nil)
		if cookie != nil {
			request.AddCookie(cookie)
		}
		response := httptest.NewRecorder()
		handler.RequireAuth(next).ServeHTTP(response, request)
		return response
	}

	t.Run("rejects a request without a session cookie", func(t *testing.T) {
		response := run(testStore{}, nil)

		assertError(t, response, 401, codeUnauthenticated)
		if nextCalled {
			t.Error("next handler ran without a session")
		}
	})

	t.Run("rejects an unknown or expired session", func(t *testing.T) {
		response := run(testStore{
			getSessionUser: func(string) (models.User, error) { return models.User{}, nil },
		}, &http.Cookie{Name: sessionCookieName, Value: "not-a-real-token"})

		assertError(t, response, 401, codeUnauthenticated)
		if nextCalled {
			t.Error("next handler ran with an invalid session")
		}
	})

	t.Run("passes a valid session through with the user in the context", func(t *testing.T) {
		var lookedUp string
		response := run(testStore{
			getSessionUser: func(tokenHash string) (models.User, error) {
				lookedUp = tokenHash
				return ada, nil
			},
		}, &http.Cookie{Name: sessionCookieName, Value: "real-token"})

		if response.Code != 200 || !nextCalled {
			t.Fatalf("status = %d, next called = %v; want 200 and true", response.Code, nextCalled)
		}
		if lookedUp != hashToken("real-token") {
			t.Error("the store wasn't given the hash of the cookie's token")
		}
		if gotUser.ID != ada.ID {
			t.Errorf("user in context = %q, want %q", gotUser.ID, ada.ID)
		}
	})
}

func TestMeWithoutRequireAuthIsAServerError(t *testing.T) {
	handler := AuthHandler{}
	response := httptest.NewRecorder()

	handler.Me(response, httptest.NewRequest("GET", "/me", nil))

	assertError(t, response, 500, codeInternalError)
}

func TestLogout(t *testing.T) {
	logout := func(store testStore, cookie *http.Cookie) *httptest.ResponseRecorder {
		handler := AuthHandler{Store: store}
		request := httptest.NewRequest("POST", "/logout", nil)
		if cookie != nil {
			request.AddCookie(cookie)
		}
		response := httptest.NewRecorder()
		handler.Logout(response, request)
		return response
	}

	// assertCookieCleared checks the response tells the browser to delete the session cookie.
	assertCookieCleared := func(t *testing.T, response *httptest.ResponseRecorder) {
		t.Helper()
		cookies := response.Result().Cookies()
		if len(cookies) != 1 || cookies[0].Name != sessionCookieName {
			t.Fatalf("cookies = %v, want one %q cookie", cookies, sessionCookieName)
		}
		if cookies[0].Value != "" || cookies[0].MaxAge >= 0 {
			t.Errorf("cookie value = %q, MaxAge = %d; want empty and deleted", cookies[0].Value, cookies[0].MaxAge)
		}
	}

	t.Run("deletes the session and clears the cookie", func(t *testing.T) {
		var deleted string
		response := logout(testStore{
			deleteSession: func(tokenHash string) error { deleted = tokenHash; return nil },
		}, &http.Cookie{Name: sessionCookieName, Value: "real-token"})

		if response.Code != 204 {
			t.Fatalf("status = %d, want 204", response.Code)
		}
		if deleted != hashToken("real-token") {
			t.Error("the store wasn't asked to delete the hash of the cookie's token")
		}
		assertCookieCleared(t, response)
	})

	t.Run("succeeds without a cookie", func(t *testing.T) {
		response := logout(testStore{
			deleteSession: func(string) error {
				t.Fatal("DeleteSession called without a cookie")
				return nil
			},
		}, nil)

		if response.Code != 204 {
			t.Fatalf("status = %d, want 204", response.Code)
		}
		assertCookieCleared(t, response)
	})

	t.Run("returns 500 if the session can't be deleted", func(t *testing.T) {
		response := logout(testStore{
			deleteSession: func(string) error { return errors.New("database unavailable") },
		}, &http.Cookie{Name: sessionCookieName, Value: "real-token"})

		assertError(t, response, 500, codeInternalError)
	})
}

func TestBoardOwnership(t *testing.T) {
	const otherUserID = "9a1b2c3d-0000-4000-8000-000000000002"

	t.Run("create sets the owner from the session, not the body", func(t *testing.T) {
		var saved models.Board
		handler := BoardHandler{
			Store:    testStore{saveBoard: func(board models.Board) error { saved = board; return nil }},
			Validate: testValidator(),
		}
		body := `{"id":"4c85d3f2-e917-4a3f-929a-11b02499a33f","title":"Hack","owner_id":"` + otherUserID + `"}`
		request := httptest.NewRequest("POST", "/boards", strings.NewReader(body))
		response := httptest.NewRecorder()

		handler.CreateBoard(response, asUser(request))

		if response.Code != 201 {
			t.Fatalf("status = %d, want 201", response.Code)
		}
		if saved.OwnerID != testUser.ID {
			t.Errorf("owner = %q, want the logged-in user %q", saved.OwnerID, testUser.ID)
		}
		if strings.Contains(response.Body.String(), "owner_id") {
			t.Errorf("response includes owner_id: %s", response.Body.String())
		}
	})

	t.Run("update filters by the logged-in owner", func(t *testing.T) {
		var updated models.Board
		handler := BoardHandler{
			Store: testStore{
				getBoard: func(id, _ string) (models.Board, error) {
					return models.Board{ID: id, Title: "Old title"}, nil // like the real store: OwnerID not loaded
				},
				updateBoard: func(board models.Board) (bool, error) { updated = board; return true, nil },
			},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("PUT", "/boards/board-id", strings.NewReader(`{"title":"New title"}`))
		request.SetPathValue("id", "board-id")
		response := httptest.NewRecorder()

		handler.UpdateBoard(response, asUser(request))

		if response.Code != 200 {
			t.Fatalf("status = %d, want 200", response.Code)
		}
		if updated.OwnerID != testUser.ID {
			t.Errorf("UpdateBoard owner = %q, want %q (an empty owner matches no rows)", updated.OwnerID, testUser.ID)
		}
	})

	t.Run("a handler without RequireAuth is a server error", func(t *testing.T) {
		handler := BoardHandler{Store: testStore{}, Validate: testValidator()}
		response := httptest.NewRecorder()

		handler.ListBoards(response, httptest.NewRequest("GET", "/boards", nil)) // no asUser

		assertError(t, response, 500, codeInternalError)
	})
}

func TestColumnOwnership(t *testing.T) {
	const columnJSON = `{"id":"7d1e2f3a-0000-4000-8000-0000000000c1","board_id":"7d1e2f3a-0000-4000-8000-0000000000b1","title":"To do","position":0}`

	t.Run("create passes the logged-in user to the store", func(t *testing.T) {
		var gotOwner string
		handler := ColumnHandler{
			Store:    testStore{saveColumn: func(_ models.Column, ownerID string) (bool, error) { gotOwner = ownerID; return true, nil }},
			Validate: testValidator(),
		}
		response := httptest.NewRecorder()

		handler.CreateColumn(response, asUser(httptest.NewRequest("POST", "/columns", strings.NewReader(columnJSON))))

		if response.Code != 201 || gotOwner != testUser.ID {
			t.Fatalf("status = %d, owner = %q; want 201 and %q", response.Code, gotOwner, testUser.ID)
		}
	})

	t.Run("create on someone else's (or a missing) board is 404", func(t *testing.T) {
		handler := ColumnHandler{
			Store:    testStore{saveColumn: func(models.Column, string) (bool, error) { return false, nil }},
			Validate: testValidator(),
		}
		response := httptest.NewRecorder()

		handler.CreateColumn(response, asUser(httptest.NewRequest("POST", "/columns", strings.NewReader(columnJSON))))

		assertError(t, response, 404, codeBoardNotFound)
	})

	t.Run("list passes the logged-in user to the store", func(t *testing.T) {
		var gotOwner string
		handler := ColumnHandler{
			Store: testStore{getColumns: func(_, ownerID string) ([]models.Column, error) {
				gotOwner = ownerID
				return []models.Column{}, nil
			}},
			Validate: testValidator(),
		}
		response := httptest.NewRecorder()

		handler.ListColumns(response, asUser(httptest.NewRequest("GET", "/columns?board_id=b1", nil)))

		if response.Code != 200 || gotOwner != testUser.ID {
			t.Fatalf("status = %d, owner = %q; want 200 and %q", response.Code, gotOwner, testUser.ID)
		}
	})

	t.Run("update that changes no rows is 404, not 200", func(t *testing.T) {
		handler := ColumnHandler{
			Store: testStore{
				getColumn:    func(id, _ string) (models.Column, error) { return models.Column{ID: id}, nil },
				updateColumn: func(models.Column, string) (bool, error) { return false, nil },
			},
			Validate: testValidator(),
		}
		request := httptest.NewRequest("PUT", "/columns/c1", strings.NewReader(`{"title":"Done","position":1}`))
		request.SetPathValue("id", "c1")
		response := httptest.NewRecorder()

		handler.UpdateColumn(response, asUser(request))

		assertError(t, response, 404, codeColumnNotFound)
	})
}

func TestCreateCardOnSomeoneElsesColumnIs404(t *testing.T) {
	var gotOwner string
	handler := CardHandler{
		Store: testStore{saveCard: func(_ models.Card, ownerID string) (bool, error) {
			gotOwner = ownerID
			return false, nil // the column isn't this user's
		}},
		Validate: testValidator(),
	}
	body := `{"id":"5e6f7a8b-0000-4000-8000-0000000000d1","column_id":"7d1e2f3a-0000-4000-8000-0000000000c1","title":"Write tests","order":0}`
	response := httptest.NewRecorder()

	handler.CreateCard(response, asUser(httptest.NewRequest("POST", "/cards", strings.NewReader(body))))

	if gotOwner != testUser.ID {
		t.Errorf("owner = %q, want %q", gotOwner, testUser.ID)
	}
	assertError(t, response, 404, codeColumnNotFound)
}
