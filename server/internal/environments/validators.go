package environments

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ConnectionValidator interface {
	ValidateConnection(ctx context.Context, connectionString string) (string, error)
}

type PostgresValidator struct{}

func (p *PostgresValidator) ValidateConnection(ctx context.Context, connectionString string) (string, error) {
	connInfo := DatabaseConnection{ConnectionString: connectionString}
	return ValidateDatabaseConnectionAsMigrator(ctx, connInfo)
}

func NewPostgresValidator() *PostgresValidator {
	return &PostgresValidator{}
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) ValidateConnection(ctx context.Context, connectionString string) (string, error) {
	args := m.Called(ctx, connectionString)
	return args.String(0), args.Error(1)
}

func NewMockValidator() *MockValidator {
	return &MockValidator{}
}
