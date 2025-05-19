package repositories

import "database/sql"

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) GetUserStats(userID string) (totalSent float64, countSent int, totalReceived float64, countReceived int, err error) {
	query := `
	SELECT 
	COALESCE((SELECT SUM(amount) FROM transactions t JOIN accounts a ON t.from_account = a.id WHERE a.user_id = $1), 0), 
	(SELECT count(*) FROM transactions t JOIN accounts a ON t.from_account = a.id WHERE a.user_id = $1), 
	COALESCE((SELECT SUM(amount) FROM transactions t JOIN accounts b ON t.to_account = b.id WHERE b.user_id = $1), 0),
    (SELECT COUNT(*) FROM transactions t JOIN accounts b ON t.to_account = b.id WHERE b.user_id = $1)`
	row := r.DB.QueryRow(query, userID)
	err = row.Scan(&totalSent, &countSent, &totalReceived, &countReceived)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return
}
