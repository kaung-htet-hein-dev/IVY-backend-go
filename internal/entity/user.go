package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Email       string    `gorm:"type:varchar(255);unique;not null"`
	PhoneNumber string    `gorm:"type:varchar(20)"`
	Role        string    `gorm:"type:varchar(20);check:role IN ('USER', 'ADMIN')"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
