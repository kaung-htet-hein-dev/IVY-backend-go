package request

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type UserRegisterRequest struct {
	UserLoginRequest

	Name        string  `json:"name" validate:"required,min=3,max=255"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,max=20"`
	Role        *string `json:"role" validate:"omitempty,oneof=USER ADMIN"`
}

type UserUpdateRequest struct {
	Role        *string `json:"role" validate:"omitempty,oneof=USER ADMIN"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,max=20"`
	Gender      *string `json:"gender" validate:"omitempty,max=20"`
	Birthday    *string `json:"birthday" validate:"omitempty"`
}
