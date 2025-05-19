package handlers

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"go_project/internal/services"
	"net/http"
)

type CardHandler struct {
	service *services.CardService
}

func NewCardHandler(service *services.CardService) *CardHandler {
	return &CardHandler{service: service}
}

// CreateCard обрабатывает POST /cards (выпуск новой карты).
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	var req struct {
		AccountID string `json:"account_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	cardNum, expiry, cvv, err := h.service.CreateCard(userID, req.AccountID)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else if errors.Is(err, services.ErrAccountNotFound) {
			http.Error(w, "Account not found", http.StatusNotFound)
		} else {
			logrus.Error("Failed to create card: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	response := map[string]string{
		"card_number": cardNum,
		"expiry":      expiry,
		"cvv":         cvv,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
