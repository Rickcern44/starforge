package models

import (
	"errors"
	"time"
)

type League struct {
	ID       string
	Name     string
	IsActive bool

	Members []LeagueMember
	Games   []*Game
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
