package dto

import (
	"time"

	"github.com/franciscoluna/envoy/server/internal/api/domain"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
	ID        domain.UserID `json:"id"`
	Email     string        `json:"email"`
	CreatedAt time.Time     `json:"created_at"`
}

func NewUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
