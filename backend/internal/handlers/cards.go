package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"kanban/internal/models"
	"kanban/internal/store"

	"github.com/go-playground/validator/v10"
)

type CardHandler struct {
	Store    store.Store
	Validate *validator.Validate
}

// #region Create Card
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	var card models.Card

	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(card); err != nil {
		writeError(w, http.StatusBadRequest, codeValidationFailed, err.Error())
		return
	}

	saved, err := h.Store.SaveCard(r.Context(), card, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if !saved {
		writeError(w, http.StatusNotFound, codeColumnNotFound, "Column not found")
		return
	}

	writeJSON(w, http.StatusCreated, card)
}

// #endregion

// #region List Cards
func (h *CardHandler) ListCards(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	cards, err := h.Store.GetCards(r.Context(), user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, cards)
}

// #endregion

// #region Get Card
func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	id := r.PathValue("id")

	card, err := h.Store.GetCardByID(r.Context(), id, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if card.ID == "" {
		writeError(w, http.StatusNotFound, codeCardNotFound, "Card not found")
		return
	}

	writeJSON(w, http.StatusOK, card)
}

// #endregion

// #region Update Card
func (h *CardHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	id := r.PathValue("id")

	var card models.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(card); err != nil {
		writeError(w, http.StatusBadRequest, codeValidationFailed, err.Error())
		return
	}

	success, err := h.Store.UpdateCard(r.Context(), id, card, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if !success {
		writeError(w, http.StatusNotFound, codeCardNotFound, "Card not found")
		return
	}

	updatedCard, err := h.Store.GetCardByID(r.Context(), id, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, updatedCard)
}

//# endregion

// #region Delete Card
func (h *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	id := r.PathValue("id")

	success, err := h.Store.DeleteCard(r.Context(), id, user.ID)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if !success {
		writeError(w, http.StatusNotFound, codeCardNotFound, "Card not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//# endregion

// #region Move Card
func (h *CardHandler) MoveCard(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(w, r)
	if !ok {
		return
	}

	var request struct {
		ColumnID string `json:"column_id" validate:"required,uuid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(request); err != nil {
		writeError(w, http.StatusBadRequest, codeValidationFailed, err.Error())
		return
	}

	moved, err := h.Store.MoveCard(r.Context(), r.PathValue("id"), request.ColumnID, user.ID)
	if errors.Is(err, store.ErrColumnNotFound) {
		writeError(w, http.StatusNotFound, codeColumnNotFound, "Column not found")
		return
	}
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	if !moved {
		writeError(w, http.StatusNotFound, codeCardNotFound, "Card not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//#endregion
