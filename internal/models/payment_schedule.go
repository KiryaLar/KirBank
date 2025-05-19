package models

import "time"

//	type PaymentSchedule struct {
//		ID              uint64     `json:"id"`
//		CreditID        uint64     `json:"credit_id"`
//		PaymentDate     time.Time  `json:"payment_date"`
//		PaymentAmount   float64    `json:"payment_amount"`
//		PrincipalAmount float64    `json:"principal_amount"`
//		InterestAmount  float64    `json:"interest_amount"`
//		Status          string     `json:"status"`  // pending, paid, overdue
//		PaidAt          *time.Time `json:"paid_at"` // Может быть NULL
//		CreatedAt       time.Time  `json:"created_at"`
//		UpdatedAt       time.Time  `json:"updated_at"`
//	}
type PaymentSchedule struct {
	ID       string     `json:"-"`
	CreditID string     `json:"-"`
	DueDate  time.Time  `json:"due_date"`
	Amount   float64    `json:"amount"`
	Paid     bool       `json:"paid"`
	PaidDate *time.Time `json:"paid_date,omitempty"`
	Penalty  float64    `json:"penalty"`
}
