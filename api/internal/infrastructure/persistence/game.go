package persistence

import "time"

type Game struct {
	Base

	LeagueID    string    `gorm:"type:uuid;not null;index:idx_game_league_start,priority:1"`
	StartTime   time.Time `gorm:"not null;index:idx_game_league_start,priority:2"`
	Location    string
	CostInCents int              `gorm:"not null;check:cost_in_cents >= 0"`
	IsCanceled  bool             `gorm:"default:false"`
	Attendance  []GameAttendance `gorm:"foreignKey:GameID"`
	League      League           `gorm:"foreignKey:LeagueID;constraint:OnDelete:CASCADE"`
}
