package auth

import "time"

type UserResponse struct {
	ID        UserID    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func userResponse(u *User) UserResponse {
	return UserResponse{ID: u.ID, Email: u.Email, CreatedAt: u.CreatedAt}
}
