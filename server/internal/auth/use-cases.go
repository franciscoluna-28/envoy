package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, repo Repository, email, rawPassword string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("auth: %w", ErrInternal)
	}

	user := &User{
		ID:           UserID(uuid.New().String()),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := repo.Create(ctx, *user); err != nil {
		return nil, fmt.Errorf("auth: %w", ErrUserExists)
	}
	return user, nil
}

func LoginUser(ctx context.Context, repo Repository, email, rawPassword string) (*User, error) {
	user, err := repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("auth: %w", ErrInvalidCredentials)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(rawPassword)); err != nil {
		return nil, fmt.Errorf("auth: %w", ErrInvalidCredentials)
	}
	return user, nil
}
