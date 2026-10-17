package handlers

import (
	"encoding/json"
	"net/http"

	"kanban/internal/models"
	"kanban/internal/store"

	"github.com/go-playground/validator/v10"
)

type BoardHandler struct {
	Store    store.Store
	Validate *validator.Validate
}

// #region Create Board
func (h *BoardHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var board models.Board

	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(board); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Store.SaveBoard(r.Context(), board); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(board)
}

// #endregion

// #region List Boards
func (h *BoardHandler) ListBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := h.Store.GetBoards(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

// #endregion

// #region Get Board
func (h *BoardHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	board, err := h.Store.GetBoardByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if board.ID == "" {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

// #endregion

// #region Update Board
func (h *BoardHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Board ID is required", http.StatusBadRequest)
		return
	}

	existingBoard, err := h.Store.GetBoardByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existingBoard.ID == "" {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	var payload struct {
		Title string `json:"title" validate:"required,min=1,max=100"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingBoard.Title = payload.Title

	updated, err := h.Store.UpdateBoard(r.Context(), existingBoard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !updated {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingBoard)
}

//#endregion

// #region Delete Board
func (h *BoardHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deleted, err := h.Store.DeleteBoard(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//#endregion
