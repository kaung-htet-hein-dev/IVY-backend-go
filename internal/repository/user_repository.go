package repository

import (
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	IsUserExist(email string) (bool, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	GetAllUsers() ([]*entity.User, error)
	Update(id uuid.UUID, updates interface{}) error
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

func (r *userRepository) GetUserByID(id string) (*entity.User, error) {
	var user entity.User
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, utils.HandleGormError(err, "user")
	}

	if err := r.db.Omit("Password", "DeletedAt").First(&user, "id = ?", uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrRecordNotFound
		}
		return nil, utils.HandleGormError(err, "user")
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]*entity.User, error) {
	var users []*entity.User
	if err := r.db.Omit("Password", "DeletedAt").Find(&users).Error; err != nil {
		return nil, utils.HandleGormError(err, "user")
	}
	return users, nil
}

func (r *userRepository) Update(id uuid.UUID, updates interface{}) error {
	return r.db.Model(&entity.User{}).Where("id = ?", id).Updates(updates).Error
}
