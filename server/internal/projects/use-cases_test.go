package projects

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, project Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id string, userID string) (*Project, error) {
	args := m.Called(ctx, id, userID)
	return args.Get(0).(*Project), args.Error(1)
}

func (m *MockRepository) GetAllByUserID(ctx context.Context, userID string) ([]Project, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]Project), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, project Project, userID string) error {
	args := m.Called(ctx, project, userID)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func GetAllProjectsByUserIDTest(t *testing.T) {
	userID := "user123"

	tests := []struct {
		name         string
		mockBehavior func(*MockRepository)
	}{
		{
			name: "Success",
			mockBehavior: func(m *MockRepository) {
				m.On("GetAllByUserID", context.Background(), userID).Return([]Project{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockBehavior(mockRepo)

			projects, err := GetAllProjectsByUserID(context.Background(), mockRepo, userID)

			if tt.name == "Success" {
				assert.NoError(t, err)
				assert.Len(t, projects, 0)
			} else {
				assert.Error(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}

}

func TestCreateProject(t *testing.T) {
	projectID := uuid.New().String()

	tests := []struct {
		name         string
		input        Project
		mockBehavior func(*MockRepository)
	}{
		{
			name: "Success",
			input: Project{
				ID:   projectID,
				Name: "Test Project",
			},
			mockBehavior: func(m *MockRepository) {
				m.On("Create", context.Background(), Project{
					ID:   projectID,
					Name: "Test Project",
				}).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockBehavior(mockRepo)

			err := CreateProject(context.Background(), mockRepo, tt.input)

			if tt.name == "Success" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
