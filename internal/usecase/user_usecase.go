package usecase

import (
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/internal/repository"
	"KaungHtetHein116/IVY-backend/utils"
	"context"
)

type UserUsecase interface {
	GetMe(c context.Context, userID string) (*entity.User, error)
	GetAllUsers(c context.Context, filter *params.UserQueryParams) ([]*entity.User, error)
	UpdateUser(c context.Context, userID string, req *request.UserUpdateRequest) error
	HandleClerkWebhook(c context.Context, req *request.ClerkWebhookRequest) error
	DeleteUser(c context.Context, userID string) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) HandleClerkWebhook(c context.Context, req *request.ClerkWebhookRequest) error {
	event := req.Type
	var email string
	var verified bool
	if len(req.Data.EmailAddresses) > 0 {
		email = req.Data.EmailAddresses[0].EmailAddress
		verified = req.Data.EmailAddresses[0].Verification.Status == "verified"
	}

	var phoneNumber *string
	if len(req.Data.PhoneNumbers) > 0 {
		phoneNumber = &req.Data.PhoneNumbers[0]
	}

	clerkUser := &entity.User{
		ID:          req.Data.ID,
		Email:       email,
		FirstName:   req.Data.FirstName,
		LastName:    req.Data.LastName,
		Verified:    verified,
		PhoneNumber: phoneNumber,
		Gender:      &req.Data.Gender,
		Birthday:    &req.Data.Birthday,
	}

	switch event {
	case "user.created":
		err := u.userRepo.CreateUser(c, clerkUser)

		if err != nil {
			return utils.HandleGormError(err, "clerk user")
		}

	case "user.updated":
		err := u.userRepo.UpdateUser(c, clerkUser)

		if err != nil {
			return utils.HandleGormError(err, "clerk user")
		}

	case "user.deleted":
		err := u.userRepo.DeleteUser(c, clerkUser.ID)
		if err != nil {
			return utils.HandleGormError(err, "clerk user")
		}
	}

	return nil
}

func (u *userUsecase) GetMe(c context.Context, userID string) (*entity.User, error) {
	user, err := u.userRepo.GetUserByID(c, userID)
	if err != nil {
		if err == utils.ErrRecordNotFound {
			return nil, utils.ErrRecordNotFound
		}
		return nil, utils.HandleGormError(err, "user")
	}

	return user, nil
}

func (u *userUsecase) GetAllUsers(c context.Context, filter *params.UserQueryParams) ([]*entity.User, error) {

	users, err := u.userRepo.GetUsers(c, filter)
	if err != nil {
		if err == utils.ErrRecordNotFound {
			return nil, utils.ErrRecordNotFound
		}
		return nil, utils.HandleGormError(err, "users")
	}

	return users, nil
}

func (u *userUsecase) UpdateUser(c context.Context, userID string, req *request.UserUpdateRequest) error {
	err := u.userRepo.UpdateUser(c, &entity.User{
		Role:        req.Role,
		PhoneNumber: req.PhoneNumber,
		Gender:      req.Gender,
		Birthday:    req.Birthday,
	})
	// Get existing user
	return err
}

func (u *userUsecase) DeleteUser(c context.Context, userID string) error {
	// Delete user by ID
	err := u.userRepo.DeleteUser(c, userID)
	if err != nil {
		return utils.HandleGormError(err, "user")
	}
	return nil
}
