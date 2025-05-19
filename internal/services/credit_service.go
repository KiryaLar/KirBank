package services

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_project/internal/models"
	"go_project/internal/repositories"
	"math"
	"time"
)

type CreditService struct {
	creditRepo   *repositories.CreditRepository
	accountRepo  *repositories.AccountRepository
	scheduleRepo *repositories.PaymentScheduleRepository
}

func NewCreditService(creditRepo *repositories.CreditRepository, accountRepo *repositories.AccountRepository, scheduleRepo *repositories.PaymentScheduleRepository) *CreditService {
	return &CreditService{creditRepo: creditRepo, accountRepo: accountRepo, scheduleRepo: scheduleRepo}
}

var (
	ErrCreditNotFound = errors.New("credit not found")
)

// GetPaymentSchedule возвращает график платежей по кредиту после проверки прав доступа.
func (s *CreditService) GetPaymentSchedule(userID, creditID string) ([]models.PaymentSchedule, error) {
	cred, err := s.creditRepo.GetById(creditID)
	if err != nil {
		return nil, err
	}
	acc, err := s.accountRepo.GetByID(cred.AccountID)
	if err != nil {
		return nil, err
	}
	if acc.UserID != userID {
		return nil, ErrForbidden
	}
	payment, err := s.scheduleRepo.GetByCreditID(creditID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// ProcessOverduePayments обрабатывает просроченные платежи: автосписание или начисление штрафов.
func (s *CreditService) ProcessOverduePayments() {
	overdueList, err := s.scheduleRepo.ListOverdue(time.Now())
	if err != nil {
		logrus.Error("Failed to list overdue payments", err)
		return
	}
	for _, ps := range overdueList {
		cred, err := s.creditRepo.GetById(ps.CreditID)
		if err != nil {
			logrus.Error("Credit not found for payment ", ps.ID)
			continue
		}
		acc, err := s.accountRepo.GetByID(cred.AccountID)
		if err != nil {
			logrus.Error("Account not found for credit")
		}
		if !ps.Paid {
			if acc.Balance >= ps.Amount {
				// Автосписание платежа
				_, err := s.accountRepo.DB.Exec(`UPDATE accounts SET balance = balance - $1 
   														WHERE id=$2`, ps.Amount, acc.ID)
				if err != nil {
					logrus.Error("Failed to deduct payment for schedule ", ps.ID, ": ", err)
					continue
				}
				err = s.scheduleRepo.MarkAsPaid(ps.ID, time.Now())
				if err != nil {
					logrus.Error("Failed to mark payment as paid ", ps.ID, ": ", err)
				}
				logrus.Infof("Auto-paid credit %s installment %.2f from account %s", cred.ID, ps.Amount, acc.ID)
			} else {
				if ps.Penalty == 0 {
					penaltyAmount := ps.Amount * 0.01
					err = s.scheduleRepo.ApplyPenalty(ps.ID, penaltyAmount)
					if err != nil {
						logrus.Error("Failed to apply penalty for payment ", ps.ID, ": ", err)
					} else {
						logrus.Warnf("Applied penalty %.2f for overdue payment %s", penaltyAmount, ps.ID)
					}
				} else {
					logrus.Warnf("Payment %s still overdue; penalty already applied", ps.ID)
				}
			}
		}
	}
}

// StartOverduePayments запускает фоновой планировщик проверки просроченных платежей каждые 12 часов.
func (s *CreditService) StartOverduePayments() {
	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		for {
			<-ticker.C
			s.ProcessOverduePayments()
		}
	}()
}

func (s *CreditService) CreateCredit(userID string, accountID string, amount float64, interest float64, termMonths int) (*models.Credit, error) {
	acc, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}
	if acc.UserID != userID {
		return nil, ErrForbidden
	}
	// Создать кредит
	creditID, err := s.creditRepo.CreateCredit(accountID, amount, interest, termMonths)
	if err != nil {
		return nil, err
	}
	// Обновляем баланс счёта — добавляем сумму кредита
	err = s.accountRepo.AddBalance(accountID, amount)
	if err != nil {
		return nil, err
	}
	// Создать график платежей (аннуитетный)
	err = s.createPaymentSchedule(creditID, amount, interest, termMonths)
	if err != nil {
		return nil, err
	}
	credit := &models.Credit{
		ID:           creditID,
		AccountID:    accountID,
		Amount:       amount,
		InterestRate: interest,
		TermMonths:   termMonths,
		StartDate:    time.Now(),
	}
	return credit, nil
}

func (s *CreditService) createPaymentSchedule(creditID string, principal float64, annualInterest float64, months int) error {
	monthlyRate := annualInterest / 1200
	annuity := principal * (monthlyRate * math.Pow(1+monthlyRate, float64(months))) / (math.Pow(1+monthlyRate, float64(months)) - 1)
	if math.IsNaN(annuity) || annuity <= 0 {
		return fmt.Errorf("invalid annuity calculation")
	}

	paymentDate := time.Now().AddDate(0, 1, 0)
	tx, err := s.creditRepo.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for i := 0; i < months; i++ {
		_, err = tx.Exec(`INSERT INTO payment_schedules (id, credit_id, due_date, amount, is_paid, penalty) 
							    VALUES (gen_random_uuid(), $1, $2, $3, FALSE, 0)`, creditID, paymentDate, annuity)
		if err != nil {
			return err
		}
		paymentDate = paymentDate.AddDate(0, 1, 0)
	}
	err = tx.Commit()
	return err
}
