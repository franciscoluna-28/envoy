package auth

import (
	"context"
	"time"

	"github.com/franciscoluna/envoy/server/internal/shared"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, repo UserRepository, email, rawPassword string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, shared.NewAppError(500, shared.ErrInternalServer, "failed to process password")
	}

	user := &User{
		ID:           UserID(uuid.New().String()),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := repo.Create(ctx, user); err != nil {
		return nil, shared.NewAppError(409, shared.ErrConflict, "user already exists")
	}
	return user, nil
}

func LoginUser(ctx context.Context, repo UserRepository, email, rawPassword string) (*User, error) {
	user, err := repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, shared.NewAppError(401, shared.ErrUnauthorized, "invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(rawPassword)); err != nil {
		return nil, shared.NewAppError(401, shared.ErrUnauthorized, "invalid credentials")
	}
	return user, nil
}
