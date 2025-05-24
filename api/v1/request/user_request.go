package request

type AuthCredentials struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type UserRegisterRequest struct {
	AuthCredentials

	Name        string  `json:"name" validate:"required,min=3,max=255"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,max=20"`
	Role        *string `json:"role" validate:"omitempty,oneof=USER ADMIN"`
}

type UserLoginRequest struct {
	AuthCredentials
}
