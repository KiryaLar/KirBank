package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_project/internal/repositories"
)

type TransactionService struct {
	accountRepo     *repositories.AccountRepository
	transactionRepo *repositories.TransactionRepository
	hmacSecret      string
}

func NewTransactionService(accountRepo *repositories.AccountRepository, transactionRepo *repositories.TransactionRepository, hmacSecret string) *TransactionService {
	return &TransactionService{accountRepo, transactionRepo, hmacSecret}
}

var (
	ErrSourceAccountNotFound      = errors.New("source account not found")
	ErrDestinationAccountNotFound = errors.New("destination account not found")
	ErrInvalidAmount              = errors.New("invalid amount")
	ErrInsufficientFunds          = errors.New("insufficient funds")
)

// Transfer выполняет перевод суммы между счетами с проверками и целостностью данных.
func (s *TransactionService) Transfer(userID, fromAccountID, toAccountID string, amount float64) (string, error) {
	if amount <= 0 {
		return "", ErrInvalidAmount
	}
	fromAcc, err := s.accountRepo.GetByID(fromAccountID)
	if err != nil {
		return "", ErrSourceAccountNotFound
	}
	if fromAcc.UserID != userID {
		return "", ErrForbidden
	}
	_, err = s.accountRepo.GetByID(toAccountID)
	if err != nil {
		return "", ErrDestinationAccountNotFound
	}
	if fromAcc.Balance < amount {
		return "", ErrInsufficientFunds
	}
	tx, err := s.accountRepo.DB.Begin()
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromAccountID)
	if err != nil {
		return "", err
	}
	_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toAccountID)
	if err != nil {
		return "", err
	}
	// Рассчитываем HMAC для записи транзакции
	data := fromAccountID + "|" + toAccountID + "|" + fmt.Sprintf("%.2f", amount)
	mac := hmac.New(sha256.New, []byte(s.hmacSecret))
	mac.Write([]byte(data))
	signature := hex.EncodeToString(mac.Sum(nil))
	// Вставляем запись о транзакции
	var newTxID string
	err = tx.QueryRow(`INSERT INTO transactions (id, from_account, to_account, amount, hmac) 
                             VALUES (gen_random_uuid(), $1, $2, $3, $4) RETURNING id`,
		fromAccountID, toAccountID, amount, signature).Scan(&newTxID)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}
	logrus.Infof("Transfer %s: %.2f from %s to %s", newTxID, amount, fromAccountID, toAccountID)
	return newTxID, nil
}

// GetAnalytics возвращает статистику операций для пользователя.
func (s *TransactionService) GetAnalytics(userID string) (int, float64, int, float64, error) {
	totalSent, countSent, totalReceived, countReceived, err := s.transactionRepo.GetUserStats(userID)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return countSent, totalSent, countReceived, totalReceived, nil
}
