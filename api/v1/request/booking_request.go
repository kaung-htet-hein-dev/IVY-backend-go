package request

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	ServiceID uuid.UUID `json:"service_id" validate:"required"`
	BranchID  uuid.UUID `json:"branch_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

type UpdateBookingRequest struct {
	Status string `json:"status" validate:"required,oneof=PENDING CONFIRMED CANCELLED COMPLETED"`
}
