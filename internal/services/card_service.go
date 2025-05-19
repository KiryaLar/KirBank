package services

import (
	"crypto/rand"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_project/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"time"
)

type CardService struct {
	cardRepo    *repositories.CardRepository
	accountRepo *repositories.AccountRepository
	encKey      string
}

func NewCardService(cardRepo *repositories.CardRepository, accountRepo *repositories.AccountRepository, encryptionKey string) *CardService {
	return &CardService{cardRepo, accountRepo, encryptionKey}
}

// CreateCard Генерирует новую карту для указанного счета и возвращает её реквизиты (номер, срок и CVV)
func (s *CardService) CreateCard(userID, accountID string) (string, string, string, error) {
	acc, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return "", "", "", ErrAccountNotFound
	}
	if acc.UserID != userID {
		return "", "", "", ErrForbidden
	}
	// Генерируем 16-значный номер карты
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(16), nil)
	randNum, _ := rand.Int(rand.Reader, max)
	cardNumber := fmt.Sprintf("%016d", randNum)
	// Генерируем срок действия (MM/YY) через 3 года от текущей даты
	now := time.Now()
	month := int(now.Month())
	year := (now.Year() + 3) % 100
	expiry := fmt.Sprintf("%02d/%02d", month, year)
	// Генерируем случайный CVV (3 цифры)
	randCVV, _ := rand.Int(rand.Reader, big.NewInt(1000))
	cvv := fmt.Sprintf("%03d", randCVV.Int64())
	// Хешируем CVV
	cvvHashBytes, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return "", "", "", err
	}
	cvvHash := string(cvvHashBytes)
	// Сохраняем карту в базе
	cardID, err := s.cardRepo.CreateCard(accountID, cardNumber, expiry, cvvHash, s.encKey)
	if err != nil {
		return "", "", "", err
	}
	logrus.Infof("Generated new card %s for account %s", cardID, accountID)
	return cardNumber, expiry, cvv, nil
}
