package auth

import "context"

// My approach towards hexagonal architecture is simple. I DO NOT like to separate things per file, but rather by logic
// Why creating a file for an interface when we can have them grouped logically? Makes no sense to me.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id UserID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type TokenProvider interface {
	GenerateToken(user *User) (string, error)
	ParseToken(tokenString string) (map[string]any, error)
}
