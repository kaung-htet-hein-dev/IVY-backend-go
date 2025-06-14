package repository

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	GetAll(ctx context.Context, params *params.CategoryQueryParams) ([]entity.Category, *transport.PaginationResponse, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	BuildQuery(ctx context.Context, params *params.CategoryQueryParams, preloads ...string) *gorm.DB
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	var existing entity.Category
	err := r.db.WithContext(ctx).Where("name = ?", category.Name).First(&existing).Error
	if err == nil {
		return gorm.ErrDuplicatedKey
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	var category entity.Category
	err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetAll(ctx context.Context, params *params.CategoryQueryParams) ([]entity.Category, *transport.PaginationResponse, error) {
	var categories []entity.Category
	query := r.BuildQuery(ctx, params)
	err := query.Find(&categories).Error
	if err != nil {
		return nil, nil, err
	}
	// Calculate pagination using the reusable utility
	pagination, err := utils.CountAndPaginate(ctx, r.db, &entity.Category{}, params.Limit, params.Offset)
	if err != nil {
		return nil, nil, err
	}
	return categories, pagination, nil
}

func (r *categoryRepository) BuildQuery(ctx context.Context, params *params.CategoryQueryParams, preloads ...string) *gorm.DB {
	builder := utils.NewQueryBuilder(r.db, ctx)

	builder.ApplyPagination(params.Limit, params.Offset)
	if params.SortBy != "" {
		builder.ApplySorting(params.SortBy, params.SortOrder)
	} else {
		builder.ApplySorting("updated_at", "desc")
	}

	// Apply string filters
	builder.ApplyStringFilters(map[string]string{
		"name": params.Name,
	})

	return builder.Build()
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Updates(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.Category{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
