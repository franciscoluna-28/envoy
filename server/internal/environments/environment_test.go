package environments

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

// There must be a better way to do this tbh. It's fine for now.
func (m *MockRepository) CreateEnvironment(ctx context.Context, env Environment) error {
	args := m.Called(ctx, env)
	return args.Error(0)
}

func (m *MockRepository) CreateEnvironmentDbUser(ctx context.Context, user EnvironmentDbUser) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetAllEnvironmentsByProjectID(ctx context.Context, projectID string) ([]Environment, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]Environment), args.Error(1)
}

func (m *MockRepository) GetEnvironmentByID(ctx context.Context, id string) (*Environment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Environment), args.Error(1)
}

func (m *MockRepository) UpdateEnvironment(ctx context.Context, env Environment) error {
	args := m.Called(ctx, env)
	return args.Error(0)
}

func (m *MockRepository) DeleteEnvironment(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateMigration(ctx context.Context, envID string, sqlContent string) error {
	args := m.Called(ctx, envID, sqlContent)
	return args.Error(0)
}

func (m *MockRepository) GetEnvironmentDbUsers(ctx context.Context, envID string) ([]EnvironmentDbUser, error) {
	args := m.Called(ctx, envID)
	return args.Get(0).([]EnvironmentDbUser), args.Error(1)
}

func (m *MockRepository) CreateEnvironmentMigration(ctx context.Context, migration EnvironmentMigration) error {
	args := m.Called(ctx, migration)
	return args.Error(0)
}

func (m *MockRepository) GetEnvironmentMigrations(ctx context.Context, envID string) ([]EnvironmentMigration, error) {
	args := m.Called(ctx, envID)
	return args.Get(0).([]EnvironmentMigration), args.Error(1)
}

func (m *MockRepository) GetEnvironmentMigrationByID(ctx context.Context, id string) (*EnvironmentMigration, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*EnvironmentMigration), args.Error(1)
}

func (m *MockRepository) UpdateEnvironmentMigration(ctx context.Context, migration EnvironmentMigration) error {
	args := m.Called(ctx, migration)
	return args.Error(0)
}

func (m *MockRepository) FindEnvironmentMigrationByClientId(ctx context.Context, clientId string) error {
	args := m.Called(ctx, clientId)
	return args.Error(0)
}

func TestCreateProjectEnvironment(t *testing.T) {
	tests := []struct {
		name          string
		input         CreateEnvironmentRequest
		mockBehavior  func(m *MockRepository, m2 *MockValidator)
		expectedError bool
	}{
		{
			name: "Success: Valid environment created",
			input: CreateEnvironmentRequest{
				Name:          "Production",
				Type:          "prod",
				ProjectID:     "proj_123",
				ConnectionUrl: "postgres://user:pass@localhost:5432/db",
			},
			mockBehavior: func(m *MockRepository, m2 *MockValidator) {
				// m.on calls a valid repository method
				m.On("CreateEnvironment", mock.Anything, mock.Anything).Return(nil)
				m2.On("ValidateConnection", mock.Anything, mock.Anything).Return("validated", nil)
			},
			expectedError: false,
		}, {
			name: "Failure: Repository returns error",
			input: CreateEnvironmentRequest{
				Name:          "Staging",
				Type:          "staging",
				ProjectID:     "proj_456",
				ConnectionUrl: "postgres://user:pass@localhost:5432/db",
			},
			mockBehavior: func(m *MockRepository, m2 *MockValidator) {
				m.On("CreateEnvironment", mock.Anything, mock.Anything).Return(errors.New("repository error"))
				m2.On("ValidateConnection", mock.Anything, mock.Anything).Return("", errors.New("validation error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockValidator := NewMockValidator()

			tt.mockBehavior(mockRepo, mockValidator)

			masterKey := []byte("this-is-a-32-byte-long-key-!!!!!")
			ctx := context.Background()

			err := CreateProjectEnvironment(ctx, tt.input, masterKey, mockRepo, mockValidator)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				mockValidator.AssertExpectations(t)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
