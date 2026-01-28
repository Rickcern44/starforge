package persistence

import "time"

type GameAttendance struct {
	GameID    string    `gorm:"type:uuid;primaryKey"`
	UserID    string    `gorm:"type:uuid;primaryKey"`
	CheckedIn bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Game      Game      `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
}
