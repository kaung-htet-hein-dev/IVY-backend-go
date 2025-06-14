package entity

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name           string    `json:"name" gorm:"type:varchar(255);not null"`
	Description    string    `json:"description" gorm:"type:text"`
	DurationMinute int       `json:"duration_minute" gorm:"type:smallint;not null"`
	Price          int       `json:"price" gorm:"type:smallint;not null"`
	CategoryID     uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
	Category       Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Image          string    `json:"image" gorm:"type:varchar(255)"`
	IsActive       bool      `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Branches       []Branch  `json:"branches" gorm:"many2many:branch_service;"`
}
