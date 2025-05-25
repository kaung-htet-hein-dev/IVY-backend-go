package usecase

import (
	"context"

	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/internal/repository"

	"github.com/google/uuid"
)

type BranchUsecase interface {
	CreateBranch(ctx context.Context, req *request.CreateBranchRequest) (*entity.Branch, error)
	GetBranchByID(ctx context.Context, id uuid.UUID) (*entity.Branch, error)
	GetAllBranches(ctx context.Context) ([]entity.Branch, error)
	UpdateBranch(ctx context.Context, id uuid.UUID, req *request.UpdateBranchRequest) (*entity.Branch, error)
	DeleteBranch(ctx context.Context, id uuid.UUID) error
}

type branchUsecase struct {
	repo repository.BranchRepository
}

func NewBranchUsecase(repo repository.BranchRepository) BranchUsecase {
	return &branchUsecase{repo: repo}
}

func (u *branchUsecase) CreateBranch(ctx context.Context, req *request.CreateBranchRequest) (*entity.Branch, error) {
	branch := &entity.Branch{
		ID:          uuid.New(),
		Name:        req.Name,
		Location:    req.Location,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		PhoneNumber: req.PhoneNumber,
	}
	err := u.repo.Create(ctx, branch)
	return branch, err
}

func (u *branchUsecase) GetBranchByID(ctx context.Context, id uuid.UUID) (*entity.Branch, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *branchUsecase) GetAllBranches(ctx context.Context) ([]entity.Branch, error) {
	return u.repo.GetAll(ctx)
}

func (u *branchUsecase) UpdateBranch(ctx context.Context, id uuid.UUID, req *request.UpdateBranchRequest) (*entity.Branch, error) {
	branch, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		branch.Name = req.Name
	}
	if req.Location != "" {
		branch.Location = req.Location
	}
	if req.Longitude != "" {
		branch.Longitude = req.Longitude
	}
	if req.Latitude != "" {
		branch.Latitude = req.Latitude
	}
	if req.PhoneNumber != "" {
		branch.PhoneNumber = req.PhoneNumber
	}
	err = u.repo.Update(ctx, branch)
	return branch, err
}

func (u *branchUsecase) DeleteBranch(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
