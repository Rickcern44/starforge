package persistence

import "time"

type Invitation struct {
	Token       string     `gorm:"primaryKey;type:varchar(255)"`
	Email       string     `gorm:"type:varchar(100);not null"`
	LeagueID    string     `gorm:"type:uuid;not null"`
	InvitedBy   string     `gorm:"type:uuid;not null"`
	SpecialRole string     `gorm:"type:varchar(50)"`
	ExpiresAt   time.Time  `gorm:"not null"`
	UsedAt      *time.Time `gorm:"default:null"`
}
