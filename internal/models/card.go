package models

//	type Card struct {
//		ID         int64     `json:"id"`
//		AccountID  int64     `json:"account_id"`
//		CardNumber string    `json:"card_number"`
//		ExpiryDate string    `json:"expiry_date"`
//		CVVHash    string    `json:"-"`
//		HMAC       string    `json:"-"`
//		CreatedAt  time.Time `json:"created_at"`
//		UpdatedAt  time.Time `json:"updated_at"`
//	}
type Card struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	// Номер карты и срок действия хранятся в БД в зашифрованном виде
	LastFour string `json:"last_four,omitempty"`
}
