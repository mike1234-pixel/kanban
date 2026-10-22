package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"kanban/internal/models"
	"kanban/internal/store"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Store    store.Store
	Validate *validator.Validate
}

// #region Create User
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Embedding models.User reuses its fields and validation tags; Password exists only on the request.
	var payload struct {
		models.User
		Password string `json:"password" validate:"required,min=8,max=72"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		writeError(w, http.StatusBadRequest, codeValidationFailed, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		// max=72 counts characters, but bcrypt's limit is 72 bytes, and some characters (e.g. é, emoji) take several.
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			writeError(w, http.StatusBadRequest, codeValidationFailed, "Password must be at most 72 bytes")
			return
		}
		writeInternalError(w, r, err)
		return
	}

	user := payload.User
	user.PasswordHash = string(hash)

	createdAt, err := h.Store.SaveUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, store.ErrEmailTaken) {
			writeError(w, http.StatusConflict, codeUserEmailTaken, "Email already in use")
			return
		}
		writeInternalError(w, r, err)
		return
	}
	user.CreatedAt = createdAt

	writeJSON(w, http.StatusCreated, user)
}

// #endregion

// #region List Users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Store.GetUsers(r.Context())
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

// #endregion

// #region Get User
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := h.Store.GetUserByID(r.Context(), id)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	if user.ID == "" {
		writeError(w, http.StatusNotFound, codeUserNotFound, "User not found")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// #endregion

// #region Update User
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, codeUserIDRequired, "User ID is required")
		return
	}

	existingUser, err := h.Store.GetUserByID(r.Context(), id)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if existingUser.ID == "" {
		writeError(w, http.StatusNotFound, codeUserNotFound, "User not found")
		return
	}

	var payload struct {
		Name  string `json:"name" validate:"required,min=1,max=100"`
		Email string `json:"email" validate:"required,email,max=254"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		writeError(w, http.StatusBadRequest, codeValidationFailed, err.Error())
		return
	}

	existingUser.Name = payload.Name
	existingUser.Email = payload.Email

	updated, err := h.Store.UpdateUser(r.Context(), existingUser)
	if err != nil {
		if errors.Is(err, store.ErrEmailTaken) {
			writeError(w, http.StatusConflict, codeUserEmailTaken, "Email already in use")
			return
		}
		writeInternalError(w, r, err)
		return
	}
	if !updated {
		writeError(w, http.StatusNotFound, codeUserNotFound, "User not found")
		return
	}

	writeJSON(w, http.StatusOK, existingUser)
}

//#endregion

// #region Delete User
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deleted, err := h.Store.DeleteUser(r.Context(), id)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if !deleted {
		writeError(w, http.StatusNotFound, codeUserNotFound, "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//#endregion
