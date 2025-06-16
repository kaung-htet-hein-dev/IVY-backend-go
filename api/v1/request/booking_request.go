package request

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	ServiceID  uuid.UUID `json:"service_id" validate:"required"`
	BranchID   uuid.UUID `json:"branch_id" validate:"required"`
	BookedDate string    `json:"booked_date" validate:"required"`
	BookedTime string    `json:"booked_time" validate:"required"`
	Note       *string   `json:"note" validate:"omitempty,max=100"`
}

type UpdateBookingRequest struct {
	Status    string    `json:"status" validate:"required,oneof=PENDING CONFIRMED CANCELLED COMPLETED"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
