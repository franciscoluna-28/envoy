package repositories

import (
	"server/internal/models"
	"server/pkg/database"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	dbUser := &database.User{
		Name:  user.Name,
		Email: user.Email,
	}
	return r.db.Create(dbUser).Error
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var dbUser database.User
	err := r.db.First(&dbUser, id).Error
	if err != nil {
		return nil, err
	}
	
	return &models.User{
		ID:    dbUser.ID,
		Name:  dbUser.Name,
		Email: dbUser.Email,
	}, nil
}

func (r *userRepository) GetAll() ([]*models.User, error) {
	var dbUsers []database.User
	err := r.db.Find(&dbUsers).Error
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = &models.User{
			ID:    dbUser.ID,
			Name:  dbUser.Name,
			Email: dbUser.Email,
		}
	}

	return users, nil
}

func (r *userRepository) Update(user *models.User) error {
	dbUser := &database.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return r.db.Save(dbUser).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&database.User{}, id).Error
}
