package persistence

import "time"

type GameAttendance struct {
	GameID          string    `gorm:"type:uuid;primaryKey"`
	UserID          string    `gorm:"type:uuid;primaryKey"`
	CheckedIn       bool      `gorm:"default:false"`
	Status          int       `gorm:"type:smallint"`
	CheckInComments string    `json:"checkInComments"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	Game            Game      `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
}
