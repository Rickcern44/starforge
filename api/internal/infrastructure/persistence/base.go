package persistence

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string         `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time      `gorm:"type:timestamp;" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"type:timestamp;" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
