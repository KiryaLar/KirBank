package models

import "time"

//	type Credit struct {
//		ID           uint64    `json:"id"`
//		UserID       uint64    `json:"user_id"`
//		Amount       float64   `json:"amount"`
//		InterestRate float64   `json:"interest_rate"`
//		Term         int       `json:"term"` // Срок в месяцах
//		StartDate    time.Time `json:"start_date"`
//		EndDate      time.Time `json:"end_date"`
//		Status       string    `json:"status"` // active, closed
//		CreatedAt    time.Time `json:"created_at"`
//		UpdatedAt    time.Time `json:"updated_at"`
//	}
type Credit struct {
	ID           string    `json:"id"`
	AccountID    string    `json:"account_id"`
	Amount       float64   `json:"amount"`
	InterestRate float64   `json:"interest_rate"`
	TermMonths   int       `json:"term_months"`
	StartDate    time.Time `json:"start_date"`
}
