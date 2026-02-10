package persistence

import (
	"time"

	"github.com/bouncy/bouncy-api/internal/domain/models"
)

type LeagueMember struct {
	LeagueID     string      `gorm:"type:uuid;primaryKey"`
	UserID       string      `gorm:"type:uuid;primaryKey"`
	Role         models.Role `gorm:"type:varchar(50);not null"`
	PasswordHash string      `gorm:"type:varchar(255);not null"`
	JoinedAt     time.Time   `gorm:"autoCreateTime"`
	League       League      `gorm:"foreignKey:LeagueID;constraint:OnDelete:CASCADE"`
}
