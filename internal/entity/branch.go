package entity

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Location    string    `json:"location" gorm:"type:varchar(50);not null"`
	Longitude   string    `json:"longitude" gorm:"type:varchar(50)"`
	Latitude    string    `json:"latitude" gorm:"type:varchar(50)"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(20)"`
	Service     []Service `json:"-" gorm:"many2many:branch_service;"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
}
