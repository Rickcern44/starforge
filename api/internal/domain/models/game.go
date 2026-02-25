package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID        string    `json:"id"`
	LeagueID  string    `json:"leagueId"`
	StartTime time.Time `json:"startTime"`
	Location  string    `json:"location"`

	CostInCents int  `json:"costInCents"`
	IsCanceled  bool `json:"isCanceled"`

	Attendance []GameAttendance `json:"attendance"`
	Charges    []GameCharge     `json:"charges"`
}

func CreateGame(leagueId, location string, costInCents int) *Game {
	gameId, _ := uuid.NewV7()

	return &Game{
		ID:          gameId.String(),
		LeagueID:    leagueId,
		StartTime:   time.Now(),
		Location:    location,
		CostInCents: costInCents,
		IsCanceled:  false,
		Attendance:  make([]GameAttendance, 0),
		Charges:     make([]GameCharge, 0),
	}
}
