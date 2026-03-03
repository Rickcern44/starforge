package models

import (
	"errors"
	"time"
)

type League struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`

	Members []LeagueMember `json:"members"`
	Games   []*Game        `json:"games"`
}

func (l *League) AddMember(playerId string, role Role) error {
	for _, m := range l.Members {
		if m.PlayerID == playerId {
			return errors.New("user already in league")
		}
	}

	member := CreateLeagueMember(l.ID, playerId, role)
	l.Members = append(l.Members, *member)

	return nil
}

func (l *League) AddGame(gameId string) error {
	for _, g := range l.Games {
		if g.ID == gameId {
			return errors.New("game already in league")
		}
	}

	l.Games = append(l.Games, &Game{
		ID:          gameId,
		LeagueID:    l.ID,
		StartTime:   time.Now(),
		Location:    "Fairview Wellness Center",
		CostInCents: 700,
		IsCanceled:  false,
		Attendance:  nil,
	})

	return nil
}

type LeagueFinancialSummary struct {
	LeagueID       string `json:"leagueId"`
	TotalCollected int    `json:"totalCollected"` // Sum of all payments
	TotalCharges   int    `json:"totalCharges"`   // Sum of all game charges
	TotalAllocated int    `json:"totalAllocated"` // Sum of all allocations (paid towards charges)
	TotalUnpaid    int    `json:"totalUnpaid"`    // TotalCharges - TotalAllocated
	TotalAvailable int    `json:"totalAvailable"` // TotalCollected - TotalAllocated
}
