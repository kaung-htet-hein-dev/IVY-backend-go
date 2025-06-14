package repository

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BranchRepository interface {
	Create(ctx context.Context, branch *entity.Branch) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Branch, error)
	GetAll(ctx context.Context, filter *params.BranchQueryParams) ([]entity.Branch, *transport.PaginationResponse, error)
	Update(ctx context.Context, id uuid.UUID, updates interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	BuildQuery(ctx context.Context, params *params.BranchQueryParams, preloads ...string) *gorm.DB
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) Create(ctx context.Context, branch *entity.Branch) error {
	return r.db.WithContext(ctx).Create(branch).Error
}

func (r *branchRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Branch, error) {
	var branch entity.Branch
	err := r.db.WithContext(ctx).First(&branch, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) GetAll(ctx context.Context, params *params.BranchQueryParams) ([]entity.Branch, *transport.PaginationResponse, error) {
	var branches []entity.Branch

	query := r.BuildQuery(ctx, params)

	if params.ServiceID != "" {
		serviceIDUUID, err := uuid.Parse(params.ServiceID)
		if err != nil {
			return nil, nil, err
		}

		err = query.
			Joins("JOIN branch_service ON branches.id = branch_service.branch_id").
			Where("branch_service.service_id = ?", serviceIDUUID).
			Find(&branches).Error

		if err != nil {
			return nil, nil, err
		}
		return branches, nil, nil
	}

	err := query.Find(&branches).Error
	if err != nil {
		return nil, nil, err
	}

	// Calculate pagination using the reusable utility
	pagination, err := utils.CountAndPaginate(ctx, r.db, &entity.Branch{}, params.Limit, params.Offset)
	if err != nil {
		return nil, nil, err
	}

	return branches, pagination, nil
}

func (r *branchRepository) Update(ctx context.Context, id uuid.UUID, updates interface{}) error {
	return r.db.WithContext(ctx).Model(&entity.Branch{}).Where("id = ?", id).Updates(updates).Error
}

func (r *branchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.Branch{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *branchRepository) BuildQuery(ctx context.Context, params *params.BranchQueryParams, preloads ...string) *gorm.DB {
	builder := utils.NewQueryBuilder(r.db, ctx)

	// Apply string filters
	builder.ApplyStringFilters(map[string]string{
		"location":     params.Location,
		"name":         params.Name,
		"phone_number": params.PhoneNumber,
		"service_id":   params.ServiceID,
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
