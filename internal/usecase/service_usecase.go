package usecase

import (
	"context"
	"encoding/json"

	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/internal/repository"

	"github.com/google/uuid"
)

type ServiceUsecase interface {
	CreateService(ctx context.Context, req *request.CreateServiceRequest) (*entity.Service, error)
	GetServiceByID(ctx context.Context, id uuid.UUID) (*entity.Service, error)
	GetAllServices(ctx context.Context) ([]entity.Service, error)
	UpdateService(ctx context.Context, id uuid.UUID, req *request.UpdateServiceRequest) (*entity.Service, error)
	DeleteService(ctx context.Context, id uuid.UUID) error
}

type serviceUsecase struct {
	repo repository.ServiceRepository
}

func NewServiceUsecase(repo repository.ServiceRepository) ServiceUsecase {
	return &serviceUsecase{repo: repo}
}

func (u *serviceUsecase) CreateService(ctx context.Context, req *request.CreateServiceRequest) (*entity.Service, error) {
	// Convert []uuid.UUID to []entity.Branch
	var branches []entity.Branch
	for _, branchID := range req.BranchIDs {

		branches = append(branches, entity.Branch{ID: branchID})
	}

	service := &entity.Service{
		ID:             uuid.New(),
		Name:           req.Name,
		Description:    req.Description,
		DurationMinute: req.DurationMinute,
		Price:          req.Price,
		CategoryID:     req.CategoryID,
		Branches:       branches,
		Image:          req.Image,
		IsActive:       req.IsActive,
	}

	err := u.repo.Create(ctx, service)
	return service, err
}

func (u *serviceUsecase) GetServiceByID(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *serviceUsecase) GetAllServices(ctx context.Context) ([]entity.Service, error) {
	return u.repo.GetAll(ctx)
}

func (u *serviceUsecase) UpdateService(ctx context.Context, id uuid.UUID, req *request.UpdateServiceRequest) (*entity.Service, error) {
	// Check if service exists
	_, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create updates map with only provided fields
	updates := make(map[string]interface{})

	// Get raw JSON data from request to check which fields were actually included
	reqJSON := make(map[string]interface{})
	jsonBytes, _ := json.Marshal(req)
	json.Unmarshal(jsonBytes, &reqJSON)

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.DurationMinute > 0 {
		updates["duration_minute"] = req.DurationMinute
	}
	if req.Price >= 0 { // Allow zero price
		updates["price"] = req.Price
	}
	if req.CategoryID != uuid.Nil {
		updates["category_id"] = req.CategoryID
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}

	// Only include is_active if it was in the request
	if _, ok := reqJSON["is_active"]; ok {
		updates["is_active"] = req.IsActive
	}

	// Convert branch IDs to Branch entities if provided
	var branches []entity.Branch
	if len(req.BranchIDs) > 0 {
		branches = make([]entity.Branch, len(req.BranchIDs))
		for i, branchID := range req.BranchIDs {
			branches[i] = entity.Branch{ID: branchID}
		}
	}

	// Update the service with only the provided fields
	err = u.repo.Update(ctx, id, updates, branches)
	if err != nil {
		return nil, err
	}

	// Fetch and return the updated service
	return u.repo.GetByID(ctx, id)
}

func (u *serviceUsecase) DeleteService(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
