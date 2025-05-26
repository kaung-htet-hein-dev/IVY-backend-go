package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time  `json:"create_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"update_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time `json:"-"`
	Email       string     `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password    string     `json:"-" gorm:"type:varchar(255);not null"`
	PhoneNumber *string    `json:"phone_number" gorm:"type:varchar(20)"`
	Role        *string    `json:"role" gorm:"type:varchar(20);default:USER;check:role IN ('USER', 'ADMIN')"`
	Bookings    []Booking  `json:"bookings" gorm:"foreignKey:UserID"`
}
