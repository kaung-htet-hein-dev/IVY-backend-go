package handler

import (
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	RegisterUser(c echo.Context, user *entity.User) error
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) RegisterUser(c echo.Context, user *entity.User) error {
	return nil
}
