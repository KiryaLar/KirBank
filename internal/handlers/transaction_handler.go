package handlers

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"go_project/internal/services"
	"net/http"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Transfer обрабатывает POST /transfer (перевод между счетами).
func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	var req struct {
		FromAccount string  `json:"from_account"`
		ToAccount   string  `json:"to_account"`
		Amount      float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	txID, err := h.service.Transfer(userID, req.FromAccount, req.ToAccount, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, services.ErrInsufficientFunds):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, services.ErrSourceAccountNotFound) || errors.Is(err, services.ErrDestinationAccountNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			{
				logrus.Error("Transfer failed: ", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
			}
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"transaction_id": txID,
		"status":         "success",
	})
}

// Analytics обрабатывает GET /analytics (статистика операций).
func (h *TransactionHandler) Analytics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	countSent, totalSent, countReceived, totalreceived, err := h.service.GetAnalytics(userID)
	if err != nil {
		logrus.Error("Failed to retrieve analytics: ", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"sent": map[string]interface{}{
			"count": countSent,
			"total": totalSent,
		},
		"received": map[string]interface{}{
			"count": countReceived,
			"total": totalreceived,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
