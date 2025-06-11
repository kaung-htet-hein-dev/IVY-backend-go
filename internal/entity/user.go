package entity

import (
	"time"
)

type User struct {
	ID        string `json:"id" gorm:"type:varchar(36);primary_key"`
	FirstName string `json:"first_name" gorm:"type:varchar(255)"`
	LastName  string `json:"last_name" gorm:"type:varchar(255)"`
	Email     string `json:"email" gorm:"type:varchar(255);unique;not null"`
	Verified  bool   `json:"verified" gorm:"type:boolean;default:false"`

	// custom fields
	Role        *string `json:"role" gorm:"type:varchar(20);default:USER;check:role IN ('USER', 'ADMIN')"`
	PhoneNumber *string `json:"phone_number" gorm:"type:varchar(20)"`
	Gender      *string `json:"gender" gorm:"type:varchar(20);default:unknown"`
	Birthday    *string `json:"birthday" gorm:"type:date"`

	// auto fields
	CreatedAt *time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
