package application

import (
	"time"

	auth "github.com/franciscoluna/envoy/server/internal/api/auth/domain"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
	ID        auth.UserID `json:"id"`
	Email     string      `json:"email"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewUserResponse(user *auth.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
