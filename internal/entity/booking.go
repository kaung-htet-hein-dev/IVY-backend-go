package entity

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID     string    `json:"user_id" gorm:"type:varchar(36);not null"`
	ServiceID  uuid.UUID `json:"service_id" gorm:"type:uuid;not null"`
	BranchID   uuid.UUID `json:"branch_id" gorm:"type:uuid;not null"`
	BookedDate string    `json:"booked_date" gorm:"type:varchar(20);not null"`
	BookedTime string    `json:"booked_time" gorm:"type:varchar(20);not null"`
	Status     string    `json:"status" gorm:"type:varchar(20);default:PENDING;check:status IN ('PENDING', 'CONFIRMED', 'CANCELLED', 'COMPLETED')"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Service    Service   `json:"service" gorm:"foreignKey:ServiceID"`
	Branch     Branch    `json:"branch" gorm:"foreignKey:BranchID"`
	Note       *string   `json:"note" gorm:"type:text;default:''"`
}
