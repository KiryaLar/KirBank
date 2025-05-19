package models

import "time"

//type Transaction struct {
//	ID              uint64    `json:"id"`
//	FromAccountID   *uint64   `json:"from_account_id"` // Может быть NULL для депозитов
//	ToAccountID     *uint64   `json:"to_account_id"`   // Может быть NULL для снятий
//	Amount          float64   `json:"amount"`
//	TransactionType string    `json:"transaction_type"` // transfer, deposit, withdrawal, payment
//	CreatedAt       time.Time `json:"created_at"`
//}

type Transaction struct {
	ID            string    `json:"id"`
	FromAccountID string    `json:"from_account"`
	ToAccountID   string    `json:"to_account"`
	Amount        float64   `json:"amount"`
	Timestamp     time.Time `json:"timestamp"`
	HMAC          string    `json:"-"`
}
