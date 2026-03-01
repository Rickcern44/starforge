package persistence

import "time"

type Game struct {
	Base

	LeagueID    string           `gorm:"type:uuid;not null;index:idx_game_league_start,priority:1" json:"leagueId"`
	StartTime   time.Time        `gorm:"not null;index:idx_game_league_start,priority:2" json:"startTime"`
	Location    string           `json:"location"`
	CostInCents int              `gorm:"not null;check:cost_in_cents >= 0" json:"costInCents"`
	IsCanceled  bool             `gorm:"default:false" json:"isCanceled"`
	Attendance  []GameAttendance `gorm:"foreignKey:GameID" json:"attendance"`
	Charges     []GameCharge     `gorm:"foreignKey:GameID" json:"charges"`
	League      League           `gorm:"foreignKey:LeagueID;constraint:OnDelete:CASCADE" json:"-"`
}
