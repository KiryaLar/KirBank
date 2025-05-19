package repositories

import (
	"database/sql"
	"go_project/internal/models"
	"time"
)

type PaymentScheduleRepository struct {
	DB *sql.DB
}

func NewPaymentScheduleRepository(db *sql.DB) *PaymentScheduleRepository {
	return &PaymentScheduleRepository{DB: db}
}

func (r *PaymentScheduleRepository) GetByCreditID(creditID string) ([]models.PaymentSchedule, error) {
	rows, err := r.DB.Query(`SELECT id, credit_id, due_date, amount, is_paid, penalty 
								    FROM payment_schedules`, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var schedules []models.PaymentSchedule
	for rows.Next() {
		var ps models.PaymentSchedule
		var paidDate sql.NullTime
		if err := rows.Scan(&ps.ID, &ps.CreditID, &ps.DueDate, &ps.Amount, &paidDate, &ps.Penalty); err != nil {
			return nil, err
		}
		if paidDate.Valid {
			pd := paidDate.Time
			ps.PaidDate = &pd
		}
		schedules = append(schedules, ps)
	}
	return schedules, nil
}

func (r *PaymentScheduleRepository) ListOverdue(currentTime time.Time) ([]models.PaymentSchedule, error) {
	rows, err := r.DB.Query(`SELECT id, credit_id, due_date, amount, is_paid, penalty 
									FROM payment_schedules 
									WHERE is_paid = false AND due_date < ?`, currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var overdue []models.PaymentSchedule
	for rows.Next() {
		var ps models.PaymentSchedule
		var paidDate sql.NullTime
		if err := rows.Scan(&ps.ID, &ps.CreditID, &ps.DueDate, &ps.Amount, &paidDate, &ps.Penalty); err != nil {
			return nil, err
		}
		if paidDate.Valid {
			pd := paidDate.Time
			ps.PaidDate = &pd
		}
		overdue = append(overdue, ps)
	}
	return overdue, nil
}

func (r *PaymentScheduleRepository) MarkAsPaid(scheduleID string, paidDate time.Time) error {
	_, err := r.DB.Exec(`UPDATE payment_schedules SET is_paid = true, paid_date = $1 
								WHERE id = $2`, paidDate, scheduleID)
	return err
}

func (r *PaymentScheduleRepository) ApplyPenalty(scheduleID string, penaltyAmount float64) error {
	_, err := r.DB.Exec(`UPDATE payment_schedules SET penalty = penalty + $1 
								WHERE id = $2`, penaltyAmount, scheduleID)
	return err
}
