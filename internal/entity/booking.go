package entity

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ServiceID uuid.UUID `json:"service_id" gorm:"type:uuid;not null"`
	BranchID  uuid.UUID `json:"branch_id" gorm:"type:uuid;not null"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`
	Status    string    `json:"status" gorm:"type:varchar(20);default:PENDING;check:status IN ('PENDING', 'CONFIRMED', 'CANCELLED', 'COMPLETED')"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}
