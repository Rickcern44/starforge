package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID        string
	LeagueID  string
	StartTime time.Time
	Location  string

	CostInCents int
	IsCanceled  bool

	Attendance []GameAttendance
	Charges    []GameCharge
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
