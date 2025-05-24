package repository

import (
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	IsUserExist(email string) (bool, error)
	GetUserByEmail(email string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) IsUserExist(email string) (bool, error) {
	var user = new(entity.User)
	if err := r.db.Where("email = ?", email).Select("email").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // User does not exist
		}
		return false, utils.HandleGormError(err, "user")
	}

	return true, nil // User exists
}

func (r *userRepository) CreateUser(user *entity.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return utils.HandleGormError(err, "user")
	}
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrRecordNotFound
		}
		return nil, utils.HandleGormError(err, "user")
	}
	return &user, nil
}
