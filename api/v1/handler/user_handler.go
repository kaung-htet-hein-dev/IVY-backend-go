package handler

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/usecase"
	"KaungHtetHein116/IVY-backend/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	RegisterUser(c echo.Context, user *request.UserRegisterRequest) error
	LoginUser(c echo.Context, user *request.UserLoginRequest) error
	GetMe(c echo.Context) error
	Logout(c echo.Context) error
	GetAllUsers(c echo.Context) error
	UpdateUser(c echo.Context, req *request.UserUpdateRequest) error
	ClerkWebhook(c echo.Context, req request.ClerkWebhookRequest) error
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func (h *userHandler) ClerkWebhook(c echo.Context, req request.ClerkWebhookRequest) error {
	// Handle Clerk webhook events here
	err := h.userUsecase.HandleClerkWebhook(c.Request().Context(), req)

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to process Clerk webhook", err)
	}
	// This is a placeholder implementation
	// You can add logic to handle different Clerk events like user creation, deletion, etc.
	return transport.NewApiSuccessResponse(c, http.StatusOK, "Clerk webhook received successfully", nil)
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) RegisterUser(c echo.Context, user *request.UserRegisterRequest) error {
	err := h.userUsecase.RegisterUser(user)
	if err != nil {

		if err == utils.ErrDuplicateEntry {
			return transport.NewApiErrorResponse(c, http.StatusConflict, "User already exists", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to register user", err)
	}

	// If registration is successful, return login user
	return h.LoginUser(c, &request.UserLoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (h *userHandler) LoginUser(c echo.Context, user *request.UserLoginRequest) error {
	token, err := h.userUsecase.LoginUser(user)
	if err != nil {
		if err == utils.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "User not found", nil)
		}

		if err == utils.ErrInvalidCredentials {
			return transport.NewApiErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", nil)
		}

		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to login user", err)
	}

	// Set JWT token as cookie
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true // Prevents JavaScript access
	cookie.Secure = true   // Only sent over HTTPS
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)

	return transport.NewApiSuccessResponse(c, http.StatusOK, "User logged in successfully",
		echo.Map{"token": token})
}

func (h *userHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(string)

	if userID == "" {
		return transport.NewApiErrorResponse(c, http.StatusUnauthorized, "Unauthorized access", nil)
	}

	userData, err := h.userUsecase.GetMe(userID)

	if err != nil {
		if err == utils.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "User not found", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to get user data", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "User data retrieved successfully", userData)
}

func (h *userHandler) Logout(c echo.Context) error {
	// Create a cookie that expires immediately
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(-1 * time.Hour) // Set expiration in the past
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Successfully logged out", nil)
}

func (h *userHandler) GetAllUsers(c echo.Context) error {
	filter := params.NewUserQueryParams()
	err := c.Bind(filter)

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
	}

	users, err := h.userUsecase.GetAllUsers(c.Request().Context(), filter)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Users retrieved successfully", users)
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
		currentUser, err := h.userUsecase.GetMe(currentUserID)
		if err != nil {
			return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to verify user permissions", err)
		}

		// Only admins can update other users
		if currentUser.Role == nil || *currentUser.Role != "ADMIN" {
			return transport.NewApiErrorResponse(c, http.StatusForbidden, "Only admins can update other users' profiles", nil)
		}
	}

	// Proceed with update
	err := h.userUsecase.UpdateUser(userID, req)
	if err != nil {
		if err == utils.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "User not found", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "User updated successfully", nil)
}
