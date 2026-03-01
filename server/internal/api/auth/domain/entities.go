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
