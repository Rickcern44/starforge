package persistence

import "time"

type GameAttendance struct {
	GameID          string    `gorm:"type:uuid;primaryKey" json:"gameId"`
	UserID          string    `gorm:"type:uuid;primaryKey" json:"userId"`
	CheckedIn       bool      `gorm:"default:false" json:"checkedIn"`
	Status          int       `gorm:"type:smallint" json:"status"`
	CheckInComments string    `gorm:"type:text" json:"checkInComments"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Game Game `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE" json:"-"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}
