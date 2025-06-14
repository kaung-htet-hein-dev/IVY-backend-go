package repository

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	BuildQuery(ctx context.Context, params *params.UserQueryParams, preloads ...string) *gorm.DB
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userID string) error

	GetUserByID(ctx context.Context, userID string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUsers(ctx context.Context, params *params.UserQueryParams, preloads ...string) ([]*entity.User, *transport.PaginationResponse, error)
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if existingUser, err := r.GetUserByEmail(ctx, user.Email); err == nil && existingUser != nil {
		return gorm.ErrDuplicatedKey
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return utils.HandleGormError(err, "user")
	}

	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	if existingUser, _ := r.GetUserByID(ctx, user.ID); existingUser == nil {
		return utils.ErrRecordNotFound
	}

	if err := r.db.WithContext(ctx).Model(&user).Updates(user).Error; err != nil {
		return utils.HandleGormError(err, "user")
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", userID).Delete(&entity.User{}).Error; err != nil {
		return utils.HandleGormError(err, "user")
	}
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	var user = new(entity.User)
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user = new(entity.User)
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUsers(ctx context.Context,
	params *params.UserQueryParams,
	preloads ...string) ([]*entity.User, *transport.PaginationResponse, error) {

	query := r.BuildQuery(ctx, params, preloads...)

	var users []*entity.User

	if err := query.Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, utils.ErrRecordNotFound
		}
		return nil, nil, utils.HandleGormError(err, "users")
	}

	// Calculate pagination using the reusable utility
	pagination, err := utils.CountAndPaginate(ctx, r.db, &entity.User{}, params.Limit, params.Offset)
	if err != nil {
		return nil, nil, err
	}

	return users, pagination, nil
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
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
		builder.ApplySorting("updated_at", "desc")
	}

	// Apply pagination and preloads
	builder.ApplyPagination(params.Limit, params.Offset).
		ApplyPreloads(preloads...)

	return builder.Build()
}
