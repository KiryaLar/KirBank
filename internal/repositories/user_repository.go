package repositories

import (
	"database/sql"
	"errors"
	"go_project/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(email, username, passwordHash string) (string, error) {
	var newID string
	query := `INSERT INTO users (id, email, username, password_hash) 
	VALUES (gen_random_uuid(), $1, $2, $3) RETURNING id`
	err := r.DB.QueryRow(query, email, username, passwordHash).Scan(&newID)
	if err != nil {
		return "", err
	}
	return newID, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var u models.User
	row := r.DB.QueryRow(`SELECT id, email, username, password_hash 
	FROM users WHERE email=$1`, email)
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var u models.User
	row := r.DB.QueryRow(`SELECT id, email, username, password_hash 
	FROM users WHERE username=$1`, username)
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &u, nil
}
