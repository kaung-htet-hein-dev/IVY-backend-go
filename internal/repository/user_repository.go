package repository

import (
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *entity.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return utils.HandleGormError(err, "user")
	}
	return nil
}
