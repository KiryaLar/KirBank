package models

//	type Account struct {
//		ID            int64     `json:"id"`
//		UserID        uint64    `json:"user_id"`
//		AccountNumber string    `json:"account_number"`
//		Balance       float64   `json:"balance"`
//		Currency      string    `json:"currency"`
//		CreatedAt     time.Time `json:"created_at"`
//		UpdatedAt     time.Time `json:"updated_at"`
//	}
type Account struct {
	ID      string  `json:"id"`
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}
