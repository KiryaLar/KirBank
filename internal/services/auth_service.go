package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"go_project/internal/models"
	"go_project/internal/repositories"
	"go_project/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

var (
	ErrEmailInUse    = errors.New("email is already in use")
	ErrUsernameTaken = errors.New("username is already taken")
)

func (s *AuthService) RegisterUser(email, username, password string) (*models.User, error) {
	if user, _ := s.userRepo.GetByEmail(email); user != nil {
		return nil, ErrEmailInUse
	}
	if user, _ := s.userRepo.GetByUsername(username); user != nil {
		return nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	passwordHash := string(hash)
	userID, err := s.userRepo.CreateUser(email, username, passwordHash)
	if err != nil {
		return nil, err
	}

	welcomeBody := "Welcome, " + username + "! Thank you for registering."
	err = utils.SendEmail(email, "Welcome to Bank", welcomeBody)
	if err != nil {
		logrus.Errorf("Failed to send welcome email to %s: %v", email, err)
	}
	logrus.Infof("Registered new user %s (email: %s)", username, email)

	user := &models.User{ID: userID, Email: email, Username: username}
	return user, nil
}

func (s *AuthService) LoginUser(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenObj.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	logrus.Infof("User %s logged in (email: %s)", user.Username, email)
	return token, nil
}
