package request

import "time"

type CreateBranchRequest struct {
	Name        string `json:"name" validate:"required"`
	Location    string `json:"location" validate:"required"`
	Longitude   string `json:"longitude" validate:"required"`
	Latitude    string `json:"latitude" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type UpdateBranchRequest struct {
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Longitude   string    `json:"longitude"`
	Latitude    string    `json:"latitude"`
	PhoneNumber string    `json:"phone_number"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
