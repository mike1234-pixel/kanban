package handlers

import (
	"encoding/json"
	"net/http"

	"kanban/internal/models"
	"kanban/internal/store"

	"github.com/go-playground/validator/v10"
)

type ColumnHandler struct {
	Store    store.Store
	Validate *validator.Validate
}

// #region Create Column
func (h *ColumnHandler) CreateColumn(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	var col models.Column
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(col); err != nil {
		writeError(w, http.StatusBadRequest, codeValidationFailed, err.Error())
		return
	}

	saved, err := h.Store.SaveColumn(r.Context(), col, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if !saved {
		writeError(w, http.StatusNotFound, codeBoardNotFound, "Board not found")
		return
	}

	writeJSON(w, http.StatusCreated, col)
}

//#endregion

// #region List Columns
func (h *ColumnHandler) ListColumns(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	boardID := r.URL.Query().Get("board_id")
	if boardID == "" {
		writeError(w, http.StatusBadRequest, codeBoardIDRequired, "board_id query parameter is required")
		return
	}

	columns, err := h.Store.GetColumnsWithCards(r.Context(), boardID, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, columns)
}

//#endregion

// #region Get Column
func (h *ColumnHandler) GetColumn(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	// Extract 'id' path parameter (Go 1.22+ net/http router)
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, codeColumnIDRequired, "Column ID is required")
		return
	}

	col, err := h.Store.GetColumnByID(r.Context(), id, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	if col.ID == "" {
		writeError(w, http.StatusNotFound, codeColumnNotFound, "Column not found")
		return
	}

	writeJSON(w, http.StatusOK, col)
}

//#endregion

// #region Update Column
func (h *ColumnHandler) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, codeColumnIDRequired, "Column ID is required")
		return
	}

	// 1. Fetch existing column from database
	existingCol, err := h.Store.GetColumnByID(r.Context(), id, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if existingCol.ID == "" {
		writeError(w, http.StatusNotFound, codeColumnNotFound, "Column not found")
		return
	}

	// 2. Decode update request body
	var req struct {
		Title    string `json:"title" validate:"required,min=1,max=100"`
		Position *int   `json:"position" validate:"required,gte=0"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, codeValidationFailed, err.Error())
		return
	}

	// 3. Apply updates to the existing model
	existingCol.Title = req.Title
	existingCol.Position = *req.Position

	// 4. Save back to database
	updated, err := h.Store.UpdateColumn(r.Context(), existingCol, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if !updated {
		// Deleted (or moved out of reach) since step 1
		writeError(w, http.StatusNotFound, codeColumnNotFound, "Column not found")
		return
	}

	// 5. Return complete updated column
	writeJSON(w, http.StatusOK, existingCol)
}

//#endregion

// #region Delete Column
func (h *ColumnHandler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	// Extract 'id' from path parameter (Go 1.22+ net/http pattern routing)
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, codeColumnIDRequired, "Column ID is required")
		return
	}

	deleted, err := h.Store.DeleteColumn(r.Context(), id, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	if !deleted {
		writeError(w, http.StatusNotFound, codeColumnNotFound, "Column not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//#endregion
