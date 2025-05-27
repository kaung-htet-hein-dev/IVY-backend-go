package entity

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ServiceID uuid.UUID `json:"service_id" gorm:"type:uuid;not null"`
	Service   Service   `json:"service" gorm:"foreignKey:ServiceID"`
	BranchID  uuid.UUID `json:"branch_id" gorm:"type:uuid;not null"`
	Branch    Branch    `json:"branch" gorm:"foreignKey:BranchID"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`
	Status    string    `json:"status" gorm:"type:varchar(20);check:status IN ('PENDING', 'CONFIRMED', 'CANCELLED', 'COMPLETED')"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}
