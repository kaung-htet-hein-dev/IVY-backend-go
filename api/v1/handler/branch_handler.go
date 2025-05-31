package handler

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/usecase"
	"net/http"

	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type BranchHandler struct {
	usecase usecase.BranchUsecase
}

func NewBranchHandler(u usecase.BranchUsecase) *BranchHandler {
	return &BranchHandler{usecase: u}
}

func (h *BranchHandler) CreateBranch(c echo.Context, req *request.CreateBranchRequest) error {
	branch, err := h.usecase.CreateBranch(c.Request().Context(), req)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to create branch", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusCreated, "Branch created successfully", branch)
}

func (h *BranchHandler) GetBranchByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid branch ID", err)
	}
	branch, err := h.usecase.GetBranchByID(c.Request().Context(), id)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "Branch not found", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusOK, "Branch retrieved successfully", branch)
}

func (h *BranchHandler) GetAllBranches(c echo.Context) error {
	serviceID := c.QueryParam("service_id")
	branches, err := h.usecase.GetAllBranches(c.Request().Context(), serviceID)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to get branches", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusOK, "Branches retrieved successfully", branches)
}

func (h *BranchHandler) UpdateBranch(c echo.Context, req *request.UpdateBranchRequest) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid branch ID", err)
	}

	branch, err := h.usecase.UpdateBranch(c.Request().Context(), id, req)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Branch not found", err)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to update branch", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusOK, "Branch updated successfully", branch)
}

func (h *BranchHandler) DeleteBranch(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid branch ID", err)
	}
	err = h.usecase.DeleteBranch(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Branch not found", err)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to delete branch", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusNoContent, "Branch deleted successfully", nil)
}
