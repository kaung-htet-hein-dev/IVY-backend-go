package entity

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      User      `gorm:"foreignKey:UserID"`
	ServiceID uuid.UUID `gorm:"type:uuid;not null"`
	Service   Service   `gorm:"foreignKey:ServiceID"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Status    string    `gorm:"type:varchar(20);check:status IN ('PENDING', 'CONFIRMED', 'CANCELLED', 'COMPLETED')"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
