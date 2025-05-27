package usecase

import (
	"context"

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
	service := &entity.Service{
		ID:             uuid.New(),
		Name:           req.Name,
		Description:    req.Description,
		DurationMinute: req.DurationMinute,
		Price:          req.Price,
		CategoryID:     req.CategoryID,
		BranchID:       req.BranchID,
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
	service, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		service.Name = req.Name
	}
	if req.Description != "" {
		service.Description = req.Description
	}
	if req.DurationMinute != 0 {
		service.DurationMinute = req.DurationMinute
	}
	if req.Price != 0 {
		service.Price = req.Price
	}
	if req.CategoryID != uuid.Nil {
		service.CategoryID = req.CategoryID
	}
	if req.BranchID != uuid.Nil {
		service.BranchID = req.BranchID
	}
	if req.Image != "" {
		service.Image = req.Image
	}
	service.IsActive = req.IsActive

	err = u.repo.Update(ctx, service)
	return service, err
}

func (u *serviceUsecase) DeleteService(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
