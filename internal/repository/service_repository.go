package repository

import (
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *entity.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error)
	GetAll(ctx context.Context) ([]entity.Service, error)
	Update(ctx context.Context, service *entity.Service) error
	Delete(ctx context.Context, id uuid.UUID) error
	CheckBranchCategoryExist(ctx context.Context, service *entity.Service) error
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

func (r *serviceRepository) GetAll(ctx context.Context) ([]entity.Service, error) {
	var services []entity.Service
	err := r.db.WithContext(ctx).
		Preload("Category").
		Find(&services).Error
	return services, err
}

func (r *serviceRepository) Update(ctx context.Context, service *entity.Service) error {
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
	// Update the service
	if err := tx.Save(service).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Update the branches associated with the service
	if err := tx.Model(service).Association("Branches").Replace(service.Branches); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *serviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.Service{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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
