package handlers

import (
	"encoding/json"
	"errors"
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
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	var board models.Board

	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(board); err != nil {
		writeError(w, http.StatusBadRequest, codeValidationFailed, err.Error())
		return
	}

	board.OwnerID = user.ID

	createdAt, err := h.Store.SaveBoard(r.Context(), board)
	if err != nil {
		if errors.Is(err, store.ErrBoardExists) {
			writeError(w, http.StatusConflict, codeBoardExists, "A board with this ID already exists")
			return
		}
		writeInternalError(w, r, err)
		return
	}
	board.CreatedAt = createdAt

	writeJSON(w, http.StatusCreated, board)
}

// #endregion

// #region List Boards
func (h *BoardHandler) ListBoards(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	boards, err := h.Store.GetBoards(r.Context(), user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, boards)
}

// #endregion

// #region Get Board
func (h *BoardHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	id := r.PathValue("id")

	board, err := h.Store.GetBoardByID(r.Context(), id, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	if board.ID == "" {
		writeError(w, http.StatusNotFound, codeBoardNotFound, "Board not found")
		return
	}

	writeJSON(w, http.StatusOK, board)
}

// #endregion

// #region Update Board
func (h *BoardHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, codeBoardIDRequired, "Board ID is required")
		return
	}

	existingBoard, err := h.Store.GetBoardByID(r.Context(), id, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if existingBoard.ID == "" {
		writeError(w, http.StatusNotFound, codeBoardNotFound, "Board not found")
		return
	}

	var payload struct {
		Title string `json:"title" validate:"required,min=1,max=100"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		writeError(w, http.StatusBadRequest, codeValidationFailed, err.Error())
		return
	}

	existingBoard.Title = payload.Title
	existingBoard.OwnerID = user.ID // GetBoardByID doesn't load it; UpdateBoard filters on it

	updated, err := h.Store.UpdateBoard(r.Context(), existingBoard)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if !updated {
		writeError(w, http.StatusNotFound, codeBoardNotFound, "Board not found")
		return
	}

	writeJSON(w, http.StatusOK, existingBoard)
}

//#endregion

// #region Delete Board
func (h *BoardHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	id := r.PathValue("id")

	deleted, err := h.Store.DeleteBoard(r.Context(), id, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if !deleted {
		writeError(w, http.StatusNotFound, codeBoardNotFound, "Board not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//#endregion
