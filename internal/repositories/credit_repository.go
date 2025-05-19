package repositories

import (
	"database/sql"
	"errors"
	"go_project/internal/models"
	"time"
)

type CreditRepository struct {
	DB *sql.DB
}

func NewCreditRepository(db *sql.DB) *CreditRepository {
	return &CreditRepository{DB: db}
}

func (r *CreditRepository) GetById(creditID string) (*models.Credit, error) {
	var c models.Credit
	row := r.DB.QueryRow(`SELECT id, account_id, amount, interest_rate, term_months, start_date 
								FROM credits WHERE id = $1`, creditID)
	var start time.Time
	if err := row.Scan(&c.ID, &c.AccountID, &c.Amount, &c.InterestRate, &c.TermMonths, &start); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	c.StartDate = start
	return &c, nil
}

func (r *CreditRepository) CreateCredit(accountID string, amount float64, interest float64, term int) (string, error) {
	var creditID string
	query := `INSERT INTO credits (id, account_id, amount, interest_rate, term_months, start_date) 
              VALUES (gen_random_uuid(), $1, $2, $3, $4, now()) 
              RETURNING id`
	err := r.DB.QueryRow(query, accountID, amount, interest, term).Scan(&creditID)
	if err != nil {
		return "", err
	}
	return creditID, nil
}
