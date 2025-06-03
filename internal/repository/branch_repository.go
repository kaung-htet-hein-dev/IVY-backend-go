package repository

import (
	"KaungHtetHein116/IVY-backend/internal/entity"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BranchRepository interface {
	Create(ctx context.Context, branch *entity.Branch) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Branch, error)
	GetAll(ctx context.Context, serviceID string) ([]entity.Branch, error)
	Update(ctx context.Context, id uuid.UUID, updates interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
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

func (r *branchRepository) GetAll(ctx context.Context, serviceID string) ([]entity.Branch, error) {
	var branches []entity.Branch

	if serviceID != "" {
		serviceIDUUID, err := uuid.Parse(serviceID)
		if err != nil {
			return nil, err
		}

		err = r.db.WithContext(ctx).
			Joins("JOIN branch_service ON branches.id = branch_service.branch_id").
			Where("branch_service.service_id = ?", serviceIDUUID).
			Find(&branches).Error

		if err != nil {
			return nil, err
		}
		return branches, nil
	}

	err := r.db.WithContext(ctx).Find(&branches).Error
	return branches, err
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
