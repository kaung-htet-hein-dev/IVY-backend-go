package entity

import (
	"github.com/google/uuid"
)

type Branch struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Location    string    `gorm:"type:text;not null"`
	Longitude   string    `gorm:"type:varchar(50)"`
	Latitude    string    `gorm:"type:varchar(50)"`
	PhoneNumber string    `gorm:"type:varchar(20)"`
}
