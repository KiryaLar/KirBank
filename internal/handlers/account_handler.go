package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"go_project/internal/services"
	"net/http"
)

type AccountHandler struct {
	service *services.AccountService
}

func NewAccountHandler(service *services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

// CreateAccount обрабатывает POST /accounts (создание нового счета).
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	account, err := h.service.CreateAccount(userID)
	if err != nil {
		logrus.Error("Failed to create account: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetBalance обрабатывает GET /accounts/{accountId}/balance (просмотр баланса).
func (h *AccountHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	accountID := chi.URLParam(r, "accountId")

	balance, err := h.service.GetBalance(userID, accountID)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

// PredictBalance обрабатывает GET /accounts/{accountId}/predict (прогноз баланса).
func (h *AccountHandler) PredictBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	accountID := chi.URLParam(r, "accountId")
	predicted, rate, err := h.service.PredictBalance(userID, accountID)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else if errors.Is(err, services.ErrAccountNotFound) {
			http.Error(w, "Account not found", http.StatusNotFound)
		} else {
			logrus.Error("Failed to predict the balance: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	resp := map[string]interface{}{
		"predicted_balance": predicted,
		"key_rate":          rate,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
