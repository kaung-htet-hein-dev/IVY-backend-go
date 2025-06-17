package handler

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/usecase"
	"KaungHtetHein116/IVY-backend/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	GetMe(c echo.Context) error
	GetAllUsers(c echo.Context) error
	UpdateUser(c echo.Context, req *request.UserUpdateRequest) error
	ClerkWebhook(c echo.Context) error
	GetUserByID(c echo.Context) error
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) ClerkWebhook(c echo.Context) error {
	var req *request.ClerkWebhookRequest
	if err := c.Bind(&req); err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
	}

	err := h.userUsecase.HandleClerkWebhook(c.Request().Context(), req)

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to process Clerk webhook", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Clerk webhook received successfully", nil)
}

func (h *userHandler) GetUserByID(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "User ID is required", nil)
	}

	userData, err := h.userUsecase.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		if err == utils.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "User not found", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to get user data", err)
	}
	return transport.NewApiSuccessResponse(c, http.StatusOK, "User data retrieved successfully", userData)
}

func (h *userHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(string)

	if userID == "" {
		return transport.NewApiErrorResponse(c, http.StatusUnauthorized, "Unauthorized access", nil)
	}

	userData, err := h.userUsecase.GetMe(c.Request().Context(), userID)

	if err != nil {
		if err == utils.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "User not found", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to get user data", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "User data retrieved successfully", userData)
}

func (h *userHandler) GetAllUsers(c echo.Context) error {
	filter := params.NewUserQueryParams()
	err := c.Bind(filter)

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
	}

	users, pagination, err := h.userUsecase.GetAllUsers(c.Request().Context(), filter)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Users retrieved successfully", users, pagination)
}

func (h *userHandler) UpdateUser(c echo.Context, req *request.UserUpdateRequest) error {
	// Get the ID of the user to update
	userID := c.Param("id")
	if userID == "" {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "User ID is required", nil)
	}

	// Get the current user's ID from the context
	currentUserID := c.Get("user_id").(string)

	// If trying to update someone else's profile
	if currentUserID != userID {
		// Check if current user is admin
		currentUser, err := h.userUsecase.GetMe(c.Request().Context(), currentUserID)
		if err != nil {
			return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to verify user permissions", err)
		}

		// Only admins can update other users
		if currentUser.Role == nil || *currentUser.Role != "ADMIN" {
			return transport.NewApiErrorResponse(c, http.StatusForbidden, "Only admins can update other users' profiles", nil)
		}
	}

	// Proceed with update
	err := h.userUsecase.UpdateUser(c.Request().Context(), userID, req)
	if err != nil {
		if err == utils.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "User not found", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "User updated successfully", nil)
}
