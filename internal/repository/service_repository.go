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
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(ctx context.Context, service *entity.Service) error {
	var category entity.Category
	var branch entity.Branch

	// Check if the category exists
	if err := r.db.WithContext(ctx).First(&category, "id = ?", service.CategoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrCategoryNotFound
		}
	}
	// Check if the branch exists
	if err := r.db.WithContext(ctx).First(&branch, "id = ?", service.BranchID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrBranchNotFound

		}
	}

	return r.db.WithContext(ctx).Create(service).Error
}

func (r *serviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	var service entity.Service
	err := r.db.WithContext(ctx).First(&service, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) GetAll(ctx context.Context) ([]entity.Service, error) {
	var services []entity.Service
	err := r.db.WithContext(ctx).Find(&services).Error
	return services, err
}

func (r *serviceRepository) Update(ctx context.Context, service *entity.Service) error {
	return r.db.WithContext(ctx).Save(service).Error
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
