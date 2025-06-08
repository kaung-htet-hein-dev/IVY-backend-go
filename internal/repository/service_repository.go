package repository

import (
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"
	"context"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *entity.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error)
	GetAll(ctx context.Context, filter *params.ServiceQueryParams) ([]entity.Service, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}, branches []entity.Branch) error
	Delete(ctx context.Context, id uuid.UUID) error
	CheckBranchCategoryExist(ctx context.Context, service *entity.Service) error
	BuildQuery(ctx context.Context, filter *params.ServiceQueryParams, preload ...string) *gorm.DB
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(ctx context.Context, service *entity.Service) error {
	err := r.CheckBranchCategoryExist(ctx, service)
	if err != nil {
		return err
	}

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the service
	if err := tx.Create(service).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Associate branches with the service
	if err := tx.Model(service).Association("Branches").Replace(service.Branches); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *serviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	var service entity.Service
	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&service, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) GetAll(ctx context.Context, filter *params.ServiceQueryParams) ([]entity.Service, error) {
	var services []entity.Service

	query := r.BuildQuery(ctx, filter, "Category", "Branches")

	// Filter by branch ID if provided
	if filter.BranchID != "" {
		if branchID, err := uuid.Parse(filter.BranchID); err == nil && branchID != uuid.Nil {
			query = query.Joins("JOIN branch_service ON services.id = branch_service.service_id").
				Where("branch_service.branch_id = ?", branchID)
		}
	}

	err := query.
		Find(&services).Error

	return services, err
}

func (r *serviceRepository) BuildQuery(ctx context.Context, filter *params.ServiceQueryParams, preload ...string) *gorm.DB {
	query := utils.NewQueryBuilder(r.db, ctx)

	stringFilters := map[string]string{
		"name": filter.Name,
	}
	if filter.DurationMinute != 0 {
		stringFilters["duration_minute"] = strconv.Itoa(filter.DurationMinute)
	}
	if filter.Price != 0 {
		stringFilters["price"] = strconv.Itoa(filter.Price)
	}
	query.ApplyStringFilters(stringFilters)

	query.ApplyUUIDFilter("category_id", filter.CategoryID)

	query.ApplyPreloads(preload...)

	return query.Build()
}

func (r *serviceRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}, branches []entity.Branch) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get existing service to check if it exists
	var service entity.Service
	if err := tx.First(&service, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Check if category exists if categoryID is being updated
	if categoryID, ok := updates["category_id"].(uuid.UUID); ok {
		var category entity.Category
		if err := tx.First(&category, "id = ?", categoryID).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				return utils.ErrCategoryNotFound
			}
			return err
		}
	}

	// Update the service with the provided fields
	if err := tx.Model(&entity.Service{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update branches if provided
	if branches != nil {
		// Check if all branch IDs exist
		branchIDs := make([]uuid.UUID, len(branches))
		for i, branch := range branches {
			branchIDs[i] = branch.ID
		}

		var count int64
		if err := tx.Model(&entity.Branch{}).Where("id IN ?", branchIDs).Count(&count).Error; err != nil {
			tx.Rollback()
			return err
		}
		if int(count) != len(branchIDs) {
			tx.Rollback()
			return utils.ErrBranchNotFound
		}

		// Update the branches
		if err := tx.Model(&entity.Service{ID: id}).Association("Branches").Replace(branches); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *serviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get the service first
	service := &entity.Service{}
	if err := tx.First(service, "id = ?", id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}

	// Clear the associations in branch_service table
	if err := tx.Model(service).Association("Branches").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Delete the service
	if err := tx.Delete(service).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *serviceRepository) CheckBranchCategoryExist(ctx context.Context, service *entity.Service) error {
	var category entity.Category
	var count int64
	// Check if the category exists
	if err := r.db.WithContext(ctx).
		First(&category, "id = ?", service.CategoryID).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrCategoryNotFound
		}
		return err
	}
	// Check if all branch IDs exist
	branchIDs := make([]uuid.UUID, len(service.Branches))
	for i, branch := range service.Branches {
		branchIDs[i] = branch.ID
	}
	if err := r.db.WithContext(ctx).Model(&entity.Branch{}).
		Where("id IN ?", branchIDs).
		Count(&count).Error; err != nil {
		return err
	}
	if int(count) != len(branchIDs) {
		return utils.ErrBranchNotFound
	}

	return nil
}
