package usecase

import (
	"context"

	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/internal/repository"

	"github.com/google/uuid"
)

type CategoryUsecase interface {
	CreateCategory(ctx context.Context, req *request.CreateCategoryRequest) (*entity.Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	GetAllCategories(ctx context.Context) ([]entity.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req *request.UpdateCategoryRequest) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type categoryUsecase struct {
	repo repository.CategoryRepository
}

func NewCategoryUsecase(repo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{repo: repo}
}

func (u *categoryUsecase) CreateCategory(ctx context.Context, req *request.CreateCategoryRequest) (*entity.Category, error) {
	category := &entity.Category{
		ID:   uuid.New(),
		Name: req.Name,
	}
	err := u.repo.Create(ctx, category)
	return category, err
}

func (u *categoryUsecase) GetCategoryByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *categoryUsecase) GetAllCategories(ctx context.Context) ([]entity.Category, error) {
	return u.repo.GetAll(ctx)
}

func (u *categoryUsecase) UpdateCategory(ctx context.Context, id uuid.UUID, req *request.UpdateCategoryRequest) (*entity.Category, error) {
	category, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	err = u.repo.Update(ctx, category)
	return category, err
}

func (u *categoryUsecase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
