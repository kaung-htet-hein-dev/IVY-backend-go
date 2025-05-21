package entity

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name           string    `gorm:"type:varchar(255);not null"`
	Description    string    `gorm:"type:text"`
	DurationMinute int       `gorm:"not null"`
	Price          float64   `gorm:"type:decimal(10,2);not null"`
	CategoryID     uuid.UUID `gorm:"type:uuid;not null"`
	Category       Category  `gorm:"foreignKey:CategoryID"`
	BranchID       uuid.UUID `gorm:"type:uuid;not null"`
	Branch         Branch    `gorm:"foreignKey:BranchID"`
	Image          string    `gorm:"type:varchar(255)"`
	IsActive       bool      `gorm:"default:true"`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
