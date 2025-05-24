package usecase

import (
	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/internal/repository"
	"KaungHtetHein116/IVY-backend/utils"
	"errors"

	"gorm.io/gorm"
)

type UserUsecase interface {
	RegisterUser(user *request.UserRegisterRequest) error
	LoginUser(user *request.UserLoginRequest) (string, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) RegisterUser(user *request.UserRegisterRequest) error {
	isExist, err := u.userRepo.IsUserExist(user.Email)

	if err != nil {
		return err
	}

	if isExist {
		return utils.ErrDuplicateEntry
	}

	userData := &entity.User{
		Name:        user.Name,
		Email:       user.Email,
		Password:    utils.GenerateHashedPassword(user.Password),
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	}

	return u.userRepo.CreateUser(userData)
}

func (u *userUsecase) LoginUser(user *request.UserLoginRequest) (string, error) {
	// fetch user by email
	userData, err := u.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		// handle not found gracefully
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utils.ErrRecordNotFound
		}
		return "", err
	}

	// check password
	if !utils.IsPasswordCorrect(userData.Password, user.Password) {
		return "", utils.ErrInvalidCredentials
	}

	// generate JWT
	token, err := utils.GenerateJWTToken(uint(userData.ID.ID()))
	if err != nil {
		return "", err
	}

	return token, nil
}
