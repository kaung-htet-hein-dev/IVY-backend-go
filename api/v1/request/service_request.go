package request

import (
	"time"

	"github.com/google/uuid"
)

type CreateServiceRequest struct {
	Name           string      `json:"name" validate:"required"`
	Description    string      `json:"description" validate:"required"`
	DurationMinute int         `json:"duration_minute" validate:"required,number,min=1"`
	Price          int         `json:"price" validate:"required,number,min=0"`
	CategoryID     uuid.UUID   `json:"category_id" validate:"required"`
	Image          string      `json:"image" validate:"required"`
	IsActive       bool        `json:"is_active" validate:"required,boolean"`
	BranchIDs      []uuid.UUID `json:"branch_ids" validate:"required,dive,uuid"`
}

type UpdateServiceRequest struct {
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	DurationMinute int         `json:"duration_minute" validate:"omitempty,number,min=1"`
	Price          int         `json:"price" validate:"omitempty,number,min=0"`
	CategoryID     uuid.UUID   `json:"category_id"`
	BranchID       uuid.UUID   `json:"branch_id"`
	Image          string      `json:"image"`
	IsActive       bool        `json:"is_active" validate:"boolean"`
	BranchIDs      []uuid.UUID `json:"branch_ids" validate:"omitempty,dive,uuid"`
	UpdatedAt      time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}
