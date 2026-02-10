package persistence

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string         `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time      `gorm:"type:timestamp;"`
	UpdatedAt time.Time      `gorm:"type:timestamp;"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
