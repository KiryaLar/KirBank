package repositories

import "database/sql"

type CardRepository struct {
	DB *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{DB: db}
}

func (r *CardRepository) CreateCard(accountID, cardNumberPlain, expiryPlain, cvvHash, encryptionKey string) (string, error) {
	var cardID string
	query := `INSERT INTO cards (id, account_id, card_number_encrypted, expiry_encrypted, cvv_hash) 
              VALUES (gen_random_uuid(), $1, 
                      pgp_sym_encrypt($2, $5, 'cipher-algo=aes256'), 
                      pgp_sym_encrypt($3, $5, 'cipher-algo=aes256'), 
                      $4) 
              RETURNING id`
	err := r.DB.QueryRow(query, accountID, cardNumberPlain, expiryPlain, cvvHash, encryptionKey).Scan(&cardID)
	if err != nil {
		return "", err
	}
	return cardID, nil
}
