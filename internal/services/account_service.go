package services

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"github.com/sirupsen/logrus"
	"go_project/internal/models"
	"go_project/internal/repositories"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AccountService struct {
	accountRepo *repositories.AccountRepository
}

func NewAccountService(accountRepo *repositories.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrForbidden       = errors.New("forbidden")
)

func (s *AccountService) CreateAccount(userID string) (*models.Account, error) {
	accountID, err := s.accountRepo.CreateAccount(userID)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Account %s for user %s was created", accountID, userID)
	return &models.Account{ID: accountID, UserID: userID, Balance: 0}, nil
}

func (s *AccountService) GetBalance(userID, accountID string) (float64, error) {
	acc, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return 0, err
	}
	if acc.UserID != userID {
		return 0, fmt.Errorf("forbidden")
	}
	return acc.Balance, nil
}

func (s *AccountService) PredictBalance(userID, accountID string) (float64, float64, error) {
	acc, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return 0, 0, ErrAccountNotFound
	}
	if acc.UserID != userID {
		return 0, 0, ErrForbidden
	}
	currentBalance := acc.Balance

	toDate := time.Now().Format("2006-01-02") + "T00:00:00"
	fromDate := time.Now().AddDate(-1, 0, 0).Format("2006-01-02") + "T00:00:00"
	soapEnvelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
               xmlns:xsd="http://www.w3.org/2001/XMLSchema"
               xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <KeyRate xmlns="http://web.cbr.ru/">
      <fromDate>%s</fromDate>
      <ToDate>%s</ToDate>
    </KeyRate>
  </soap:Body>
</soap:Envelope>`, fromDate, toDate)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx",
		strings.NewReader(soapEnvelope))
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SoapAction", "http://web.cbr.ru/KeyRate")
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		return 0, 0, err
	}
	var keyRate float64
	elements := doc.FindElements("//Rate")
	if len(elements) > 0 {
		lastRateText := elements[len(elements)-1].Text()
		if val, err := strconv.ParseFloat(lastRateText, 64); err == nil {
			keyRate = val
		}
	}
	if keyRate == 0 {
		logrus.Warn("Key rate not found, assuming 0%")
	}
	predictedBalance := currentBalance * (1 + keyRate)
	return predictedBalance, keyRate, nil
}
