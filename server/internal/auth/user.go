package auth

import (
	"time"
)

type UserID string

type User struct {
	ID           UserID    `json:"id"            db:"id"`
	Email        string    `json:"email"         db:"email"`
	PasswordHash string    `json:"-"             db:"password_hash"`
	CreatedAt    time.Time `json:"created_at"    db:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenProvider interface {
	GenerateToken(userID string) (string, error)
	ParseToken(token string) (map[string]interface{}, error)
}
