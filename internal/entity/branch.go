package entity

import (
	"github.com/google/uuid"
)

type Branch struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Location    string    `json:"location" gorm:"type:text;not null"`
	Longitude   string    `json:"longitude" gorm:"type:varchar(50)"`
	Latitude    string    `json:"latitude" gorm:"type:varchar(50)"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(20)"`
	Services    []Service `json:"services" gorm:"foreignKey:BranchID"`
	Bookings    []Booking `json:"bookings" gorm:"foreignKey:BranchID"`
}
