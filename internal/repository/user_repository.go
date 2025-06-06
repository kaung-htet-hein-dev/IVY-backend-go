package repository

import (
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	IsUserExist(email string) (bool, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	GetAllUsers(c context.Context, filter *params.UserQueryParams) ([]*entity.User, error)
	Update(id uuid.UUID, updates interface{}) error
	BuildQuery(ctx context.Context, params *params.UserQueryParams, preloads ...string) *gorm.DB
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

func (r *userRepository) GetAllUsers(c context.Context, filter *params.UserQueryParams) ([]*entity.User, error) {
	var users []*entity.User
	query := r.BuildQuery(c, filter)
	if err := query.Find(&users).Error; err != nil {
		return nil, utils.HandleGormError(err, "user")
	}
	return users, nil
}

func (r *userRepository) Update(id uuid.UUID, updates interface{}) error {
	return r.db.Model(&entity.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *userRepository) BuildQuery(ctx context.Context, params *params.UserQueryParams, preloads ...string) *gorm.DB {
	builder := utils.NewQueryBuilder(r.db, ctx)

	// Apply string filters
	builder.ApplyStringFilters(map[string]string{
		"email":        params.Email,
		"name":         params.Name,
		"role":         params.Role,
		"phone_number": params.PhoneNumber,
	})

	// Apply sorting
	if params.SortBy != "" {
		builder.ApplySorting(params.SortBy, params.SortOrder)
	} else {
		builder.ApplySorting("created_at", "desc")
	}

	// Apply pagination and preloads
	builder.ApplyPagination(params.Limit, params.Offset).
		ApplyPreloads(preloads...)

	return builder.Build()
}
