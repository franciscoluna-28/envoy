package usecases

import (
	"context"
	"server/internal/api/domain"
	"server/internal/api/repositories"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterWithEmailUseCase struct {
	userRepo *repositories.UserRepository
}

type LoginWithEmailUseCase struct {
	userRepo *repositories.UserRepository
}

func NewRegisterUseCase(r *repositories.UserRepository) *RegisterWithEmailUseCase {
	return &RegisterWithEmailUseCase{userRepo: r}
}

func NewLoginWithEmailUseCase(r *repositories.UserRepository) *LoginWithEmailUseCase {
	return &LoginWithEmailUseCase{userRepo: r}
}

func (uc *RegisterWithEmailUseCase) Execute(ctx context.Context, email string, rawPassword string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           domain.UserID(uuid.New().String()),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	err = uc.userRepo.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *LoginWithEmailUseCase) Execute(ctx context.Context, email string, rawPassword string) (*domain.User, error) {

	user, err := uc.userRepo.GetByEmail(ctx, email)

	if (err) != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(rawPassword))

	if err != nil {
		return nil, err
	}

	return user, nil
}
