package application

import (
	"context"
	"fmt"
	"time"

	auth "github.com/franciscoluna/envoy/server/internal/api/auth/domain"
	"github.com/franciscoluna/envoy/server/internal/shared"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterWithEmailUseCase struct {
	userRepo auth.UserRepository
}

type LoginWithEmailUseCase struct {
	userRepo auth.UserRepository
}

func NewRegisterUseCase(r auth.UserRepository) *RegisterWithEmailUseCase {
	return &RegisterWithEmailUseCase{userRepo: r}
}

func NewLoginWithEmailUseCase(r auth.UserRepository) *LoginWithEmailUseCase {
	return &LoginWithEmailUseCase{userRepo: r}
}

func (uc *RegisterWithEmailUseCase) Execute(ctx context.Context, email string, rawPassword string) (*auth.User, error) {
	fmt.Printf("Execute called with email: %s, password length: %d\n", email, len(rawPassword))

	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	if err != nil {
		fmt.Printf("bcrypt error: %v\n", err)
		return nil, shared.NewAppError(500, shared.ErrInternalServer, "failed to process password")
	}

	user := &auth.User{
		ID:           auth.UserID(uuid.New().String()),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	fmt.Printf("About to call userRepo.Create with user: %+v\n", user)
	err = uc.userRepo.Create(ctx, user)

	if err != nil {
		fmt.Printf("userRepo.Create error: %v\n", err)
		return nil, shared.NewAppError(409, shared.ErrConflict, "user already exists")
	}

	fmt.Printf("User created successfully: %+v\n", user)
	return user, nil
}

func (uc *LoginWithEmailUseCase) Execute(ctx context.Context, email string, rawPassword string) (*auth.User, error) {

	user, err := uc.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return nil, shared.NewAppError(401, shared.ErrUnauthorized, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(rawPassword)); err != nil {
		return nil, shared.NewAppError(401, shared.ErrUnauthorized, "invalid credentials")
	}

	return user, nil
}
