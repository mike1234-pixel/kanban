package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"kanban/internal/models"
	"kanban/internal/store"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	Store    store.Store
	Validate *validator.Validate
}

// #region Create User
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAt, err := h.Store.SaveUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, store.ErrEmailTaken) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.CreatedAt = createdAt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// #endregion

// #region List Users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Store.GetUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// #endregion

// #region Get User
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := h.Store.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.ID == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// #endregion

// #region Update User
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	existingUser, err := h.Store.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existingUser.ID == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var payload struct {
		Name  string `json:"name" validate:"required,min=1,max=100"`
		Email string `json:"email" validate:"required,email,max=254"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingUser.Name = payload.Name
	existingUser.Email = payload.Email

	updated, err := h.Store.UpdateUser(r.Context(), existingUser)
	if err != nil {
		if errors.Is(err, store.ErrEmailTaken) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !updated {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingUser)
}

//#endregion

// #region Delete User
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deleted, err := h.Store.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//#endregion
