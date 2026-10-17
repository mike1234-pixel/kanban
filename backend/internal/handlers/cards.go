package handlers

import (
	"encoding/json"
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
	var card models.Card

	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(card); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Store.SaveCard(card); err != nil {
		http.Error(w, "Failed to save card", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

// #endregion

// #region List Cards
func (h *CardHandler) ListCards(w http.ResponseWriter, r *http.Request) {
	cards, err := h.Store.GetCards()
	if err != nil {
		http.Error(w, "Failed to fetch cards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

// #endregion

// #region Get Card
func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	card, err := h.Store.GetCardByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch card", http.StatusInternalServerError)
		return
	}
	if card.ID == "" {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

// #endregion

// #region Update Card
func (h *CardHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var card models.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := h.Validate.Struct(card); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	success, err := h.Store.UpdateCard(id, card)
	if err != nil {
		http.Error(w, "Failed to update card", http.StatusInternalServerError)
		return
	}
	if !success {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	updatedCard, err := h.Store.GetCardByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch updated card", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCard)
}

//# endregion

// #region Delete Card
func (h *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	success, err := h.Store.DeleteCard(id)
	if err != nil {
		http.Error(w, "Failed to delete card", http.StatusInternalServerError)
		return
	}
	if !success {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//# endregion

// #region Move Card
func (h *CardHandler) MoveCard(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ColumnID string `json:"column_id" validate:"required,uuid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if err := h.Validate.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	moved, err := h.Store.MoveCard(r.PathValue("id"), request.ColumnID)
	if err != nil {
		http.Error(w, "Failed to move card", http.StatusInternalServerError)
		return
	}
	if !moved {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//#endregion
