package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"go_project/internal/services"
	"net/http"
)

type CreditHandler struct {
	service *services.CreditService
}

func NewCreditHandler(service *services.CreditService) *CreditHandler {
	return &CreditHandler{service: service}
}

// GetPaymentSchedule обрабатывает GET /credits/{creditId}/schedule (график платежей по кредиту).
func (h *CreditHandler) GetPaymentSchedule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	creditID := chi.URLParam(r, "creditId")
	schedule, err := h.service.GetPaymentSchedule(userID, creditID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, sql.ErrNoRows):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}

// CreateCredit обрабатывает POST /credits (создание кредита).
func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	var req struct {
		AccountID  string  `json:"account_id"`
		Amount     float64 `json:"amount"`
		Interest   float64 `json:"interest_rate"`
		TermMonths int     `json:"term_months"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	credit, err := h.service.CreateCredit(userID, req.AccountID, req.Amount, req.Interest, req.TermMonths)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else if errors.Is(err, services.ErrCreditNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credit)
}
