package handler

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/usecase"
	"net/http"

	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	usecase usecase.CategoryUsecase
}

func NewCategoryHandler(u usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: u}
}

func (h *CategoryHandler) CreateCategory(c echo.Context, req *request.CreateCategoryRequest) error {
	category, err := h.usecase.CreateCategory(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return transport.NewApiErrorResponse(c, http.StatusConflict, "Category already exists", err)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to create category", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusCreated, "Category created successfully", category)
}

func (h *CategoryHandler) GetCategoryByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid category ID", err)
	}
	category, err := h.usecase.GetCategoryByID(c.Request().Context(), id)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "Category not found", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	filter := params.NewCategoryQueryParams()
	err := c.Bind(filter)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
	}

	categories, pagination, err := h.usecase.GetAllCategories(c.Request().Context(), filter)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to get categories", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categories, pagination)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context, req *request.UpdateCategoryRequest) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid category ID", err)
	}
	category, err := h.usecase.UpdateCategory(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Category not found", err)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to update category", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusOK, "Category updated successfully", category)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid category ID", err)
	}
	err = h.usecase.DeleteCategory(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Category not found", err)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to delete category", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusNoContent, "Category deleted successfully", nil)
}
