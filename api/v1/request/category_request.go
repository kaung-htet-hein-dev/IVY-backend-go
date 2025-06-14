package request

import "time"

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryRequest struct {
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
