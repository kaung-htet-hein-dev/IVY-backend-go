package handler

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/usecase"
	"KaungHtetHein116/IVY-backend/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	RegisterUser(c echo.Context, user *request.UserRegisterRequest) error
	LoginUser(c echo.Context, user *request.UserLoginRequest) error
	GetMe(c echo.Context) error
}

type userHandler struct {
	userUsecase usecase.UserUsecase
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
