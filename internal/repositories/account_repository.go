package repositories

import (
	"database/sql"
	"errors"
	"go_project/internal/models"
)

type AccountRepository struct {
	DB *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) CreateAccount(userID string) (string, error) {
	var accountID string
	query := `INSERT INTO accounts (id, user_id, balance) 
              VALUES (gen_random_uuid(), $1, 0) RETURNING id`
	err := r.DB.QueryRow(query, userID).Scan(&accountID)
	if err != nil {
		return "", err
	}
	return accountID, nil
}

func (r *AccountRepository) GetByID(accountID string) (*models.Account, error) {
	var acc models.Account
	row := r.DB.QueryRow(`SELECT id, user_id, accounts.balance 
                                FROM accounts where id = $1`, accountID)
	if err := row.Scan(&acc.ID, &acc.UserID, &acc.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &acc, nil
}

func (r *AccountRepository) GetBalance(accountID string) (float64, error) {
	var balance float64
	query := `SELECT balance FROM accounts where id = $1`
	err := r.DB.QueryRow(query, accountID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *AccountRepository) AddBalance(accountID string, amount float64) error {
	_, err := r.DB.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, accountID)
	return err
}
