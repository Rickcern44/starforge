package persistence

import (
	"time"

	"github.com/bouncy/bouncy-api/internal/domain/models"
)

type LeagueMember struct {
	LeagueID     string      `gorm:"type:uuid;primaryKey" json:"leagueId"`
	UserID       string      `gorm:"type:uuid;primaryKey" json:"userId"`
	Role         models.Role `gorm:"type:varchar(50);not null" json:"role"`
	PasswordHash string      `gorm:"type:varchar(255);not null" json:"-"`
	JoinedAt     time.Time   `gorm:"autoCreateTime" json:"joinedAt"`
	League       League      `gorm:"foreignKey:LeagueID;constraint:OnDelete:CASCADE" json:"-"`
}
