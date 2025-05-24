package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time
	Email       string  `gorm:"type:varchar(255);unique;not null"`
	Password    string  `gorm:"type:varchar(255);not null"`
	PhoneNumber *string `gorm:"type:varchar(20)"`
	Role        *string `gorm:"type:varchar(20);default:USER;check:role IN ('USER', 'ADMIN')"`
}
