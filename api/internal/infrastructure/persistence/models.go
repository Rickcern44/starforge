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

type League struct {
	Base

	Name     string         `gorm:"not null;uniqueIndex:ux_league_name"`
	IsActive bool           `gorm:"default:true"`
	Members  []LeagueMember `gorm:"foreignKey:LeagueID"`
	Games    []Game         `gorm:"foreignKey:LeagueID"`
}

type LeagueMember struct {
	LeagueID string    `gorm:"type:uuid;primaryKey"`
	UserID   string    `gorm:"type:uuid;primaryKey"`
	Role     string    `gorm:"type:varchar(50);not null"`
	JoinedAt time.Time `gorm:"autoCreateTime"`
	League   League    `gorm:"foreignKey:LeagueID;constraint:OnDelete:CASCADE"`
}
type Game struct {
	Base

	LeagueID    string    `gorm:"type:uuid;not null;index:idx_game_league_start,priority:1"`
	StartTime   time.Time `gorm:"not null;index:idx_game_league_start,priority:2"`
	Location    string
	CostInCents int              `gorm:"not null;check:cost_in_cents >= 0"`
	IsCanceled  bool             `gorm:"default:false"`
	Attendance  []GameAttendance `gorm:"foreignKey:GameID"`
	Payments    []GamePayment    `gorm:"foreignKey:GameID"`
	League      League           `gorm:"foreignKey:LeagueID;constraint:OnDelete:CASCADE"`
}

type GameAttendance struct {
	GameID    string    `gorm:"type:uuid;primaryKey"`
	UserID    string    `gorm:"type:uuid;primaryKey"`
	CheckedIn bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Game      Game      `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
}

type GamePayment struct {
	Base

	GameID      string `gorm:"type:uuid;not null;index"`
	UserID      string `gorm:"type:uuid;not null;index"`
	AmountCents int    `gorm:"not null;check:amount_cents > 0"`
	Method      string `gorm:"type:varchar(50);not null"`
	Status      string `gorm:"type:varchar(50);not null"`
	PaidAt      *time.Time
	ConfirmedBy *string `gorm:"type:uuid"`
	Game        Game    `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
}
