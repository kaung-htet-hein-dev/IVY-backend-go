package handler

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/usecase"
	"KaungHtetHein116/IVY-backend/utils"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ServiceHandler struct {
	usecase usecase.ServiceUsecase
}

func NewServiceHandler(u usecase.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{usecase: u}
}

func (h *ServiceHandler) CreateService(c echo.Context, req *request.CreateServiceRequest) error {
	service, err := h.usecase.CreateService(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, utils.ErrBranchNotFound) || errors.Is(err, utils.ErrCategoryNotFound) {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Category or Branch not found", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to create service", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusCreated, "Service created successfully", service)
}

func (h *ServiceHandler) GetServiceByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid service ID", err)
	}
	service, err := h.usecase.GetServiceByID(c.Request().Context(), id)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "Service not found", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusOK, "Service retrieved successfully", service)
}

func (h *ServiceHandler) GetAllServices(c echo.Context) error {
	filter := params.NewServiceQueryParams()

	err := c.Bind(filter)

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
	}

	services, pagination, err := h.usecase.GetAllServices(c.Request().Context(), filter)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to get services", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Services retrieved successfully", services, pagination)
}

func (h *ServiceHandler) UpdateService(c echo.Context, req *request.UpdateServiceRequest) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid service ID", err)
	}

	service, err := h.usecase.UpdateService(c.Request().Context(), id, req)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Service not found", err)
		}
		if errors.Is(err, utils.ErrBranchNotFound) || errors.Is(err, utils.ErrCategoryNotFound) {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Category or Branch not found", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to update service", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusOK, "Service updated successfully", service)
}

func (h *ServiceHandler) DeleteService(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid service ID", err)
	}
	err = h.usecase.DeleteService(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Service not found", err)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to delete service", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusNoContent, "Service deleted successfully", nil)
}
