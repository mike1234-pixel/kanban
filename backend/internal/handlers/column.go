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
	var col models.Column
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(col); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Store.SaveColumn(r.Context(), col); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(col)
}

//#endregion

// #region List Columns
func (h *ColumnHandler) ListColumns(w http.ResponseWriter, r *http.Request) {
	boardID := r.URL.Query().Get("board_id")
	if boardID == "" {
		http.Error(w, "board_id query parameter is required", http.StatusBadRequest)
		return
	}

	columns, err := h.Store.GetColumnsWithCards(r.Context(), boardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(columns)
}

//#endregion

// #region Get Column
func (h *ColumnHandler) GetColumn(w http.ResponseWriter, r *http.Request) {
	// Extract 'id' path parameter (Go 1.22+ net/http router)
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Column ID is required", http.StatusBadRequest)
		return
	}

	col, err := h.Store.GetColumnByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if col.ID == "" {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(col)
}

//#endregion

// #region Update Column
func (h *ColumnHandler) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Column ID is required", http.StatusBadRequest)
		return
	}

	// 1. Fetch existing column from database
	existingCol, err := h.Store.GetColumnByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existingCol.ID == "" {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}

	// 2. Decode update request body
	var req struct {
		Title    string `json:"title" validate:"required,min=1,max=100"`
		Position *int   `json:"position" validate:"required,gte=0"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Apply updates to the existing model
	existingCol.Title = req.Title
	existingCol.Position = *req.Position

	// 4. Save back to database
	_, err = h.Store.UpdateColumn(r.Context(), existingCol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Return complete updated column
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingCol)
}

//#endregion

// #region Delete Column
func (h *ColumnHandler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	// Extract 'id' from path parameter (Go 1.22+ net/http pattern routing)
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Column ID is required", http.StatusBadRequest)
		return
	}

	deleted, err := h.Store.DeleteColumn(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !deleted {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//#endregion
